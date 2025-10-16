package secureutils

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
)

// RandomString generates a random string of length n using the provided charset.
// It uses crypto/rand for cryptographically secure randomness.
func RandomString(n int, charset string) (string, error) {
	if len(charset) == 0 {
		return "", fmt.Errorf("charset cannot be empty")
	}
	if n <= 0 {
		return "", fmt.Errorf("length must be positive")
	}

	result := make([]byte, n)
	max := big.NewInt(int64(len(charset)))

	for i := 0; i < n; i++ {
		num, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("failed to generate random number: %v", err)
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}

// RandomHex generates a random hexadecimal string of length nBytes * 2.
// It is suitable for tokens, hashes, or unique IDs.
func RandomHex(nBytes int) (string, error) {
	if nBytes <= 0 {
		return "", fmt.Errorf("nBytes must be positive")
	}
	buf := make([]byte, nBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("failed to read random bytes: %v", err)
	}
	return hex.EncodeToString(buf), nil
}

// RandomBase64 generates a random Base64 (URL-safe) string without padding.
// It is convenient for API keys or URLs.
func RandomBase64(nBytes int) (string, error) {
	if nBytes <= 0 {
		return "", fmt.Errorf("nBytes must be positive")
	}
	buf := make([]byte, nBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("failed to read random bytes: %v", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// PasswordOptions defines configuration options for GeneratePassword.
// By default, all character sets are used, ambiguous characters are avoided,
// and at least one character from each chosen set is guaranteed.
type PasswordOptions struct {
	Length            int  // Total password length
	UseLower          bool // Include lowercase letters
	UseUpper          bool // Include uppercase letters
	UseDigits         bool // Include digits
	UseSymbols        bool // Include punctuation/symbols
	AvoidAmbiguous    bool // Exclude ambiguous chars (0/O, 1/l/I)
	RequireEachChosen bool // Require at least one char from each chosen set
}

// GeneratePassword creates a secure random password according to the given options.
// Default behavior (if no options are specified):
//   - 16 characters long
//   - Uses upper, lower, digits, and symbols
//   - Avoids ambiguous characters
//   - Guarantees at least one of each chosen type
func GeneratePassword(opts PasswordOptions) (string, error) {
	// Safe defaults for typical usage
	if opts.Length == 0 {
		opts.Length = 16
	}
	if !(opts.UseLower || opts.UseUpper || opts.UseDigits || opts.UseSymbols) {
		opts.UseLower, opts.UseUpper, opts.UseDigits, opts.UseSymbols = true, true, true, true
	}
	if !opts.RequireEachChosen {
		opts.RequireEachChosen = true
	}
	if !opts.AvoidAmbiguous {
		opts.AvoidAmbiguous = true
	}

	lower := "abcdefghijklmnopqrstuvwxyz"
	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	symbols := "!@#$%^&*()-_=+[]{};:,.?/~" // avoids < > | for readability

	if opts.AvoidAmbiguous {
		lower = "abcdefghijkmnopqrstuvwxyz" // no 'l'
		upper = "ABCDEFGHJKLMNPQRSTUVWXYZ"  // no 'I','O'
		digits = "23456789"                 // no '0','1'
	}

	sets := []string{}
	if opts.UseLower && len(lower) > 0 {
		sets = append(sets, lower)
	}
	if opts.UseUpper && len(upper) > 0 {
		sets = append(sets, upper)
	}
	if opts.UseDigits && len(digits) > 0 {
		sets = append(sets, digits)
	}
	if opts.UseSymbols && len(symbols) > 0 {
		sets = append(sets, symbols)
	}
	if len(sets) == 0 {
		return "", fmt.Errorf("no character sets enabled")
	}
	if opts.Length < len(sets) && opts.RequireEachChosen {
		return "", fmt.Errorf("length too short to include all required sets")
	}

	// Build combined charset
	combined := ""
	for _, s := range sets {
		combined += s
	}

	// 1) Ensure one from each set if required
	var password []byte
	if opts.RequireEachChosen {
		for _, s := range sets {
			idx, err := secureRandInt(len(s))
			if err != nil {
				return "", fmt.Errorf("rand index: %w", err)
			}
			password = append(password, s[idx])
		}
	}

	// 2) Fill the rest from combined charset
	for len(password) < opts.Length {
		idx, err := secureRandInt(len(combined))
		if err != nil {
			return "", fmt.Errorf("rand index: %w", err)
		}
		password = append(password, combined[idx])
	}

	// 3) Shuffle to randomize position distribution
	if err := secureShuffle(password); err != nil {
		return "", err
	}

	return string(password), nil
}

// EstimateEntropy returns an approximate entropy estimation for a password
// generated with the given options. It calculates:
//
//	totalBits ≈ Length * log2(|charset|)
//
// This gives an upper bound if RequireEachChosen = true.
func EstimateEntropy(opts PasswordOptions) (totalBits float64, perCharBits float64, charsetSize int, err error) {
	// Apply same defaults as GeneratePassword
	if opts.Length == 0 {
		opts.Length = 16
	}
	if !(opts.UseLower || opts.UseUpper || opts.UseDigits || opts.UseSymbols) {
		opts.UseLower, opts.UseUpper, opts.UseDigits, opts.UseSymbols = true, true, true, true
	}
	if !opts.AvoidAmbiguous {
		opts.AvoidAmbiguous = true
	}

	lower := "abcdefghijklmnopqrstuvwxyz"
	upper := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits := "0123456789"
	symbols := "!@#$%^&*()-_=+[]{};:,.?/~"

	if opts.AvoidAmbiguous {
		lower = "abcdefghijkmnopqrstuvwxyz"
		upper = "ABCDEFGHJKLMNPQRSTUVWXYZ"
		digits = "23456789"
	}

	combined := ""
	if opts.UseLower {
		combined += lower
	}
	if opts.UseUpper {
		combined += upper
	}
	if opts.UseDigits {
		combined += digits
	}
	if opts.UseSymbols {
		combined += symbols
	}

	if len(combined) == 0 {
		return 0, 0, 0, fmt.Errorf("no character sets enabled")
	}
	if opts.Length <= 0 {
		return 0, 0, 0, fmt.Errorf("length must be positive")
	}

	charsetSize = len(combined)
	perCharBits = math.Log2(float64(charsetSize))
	totalBits = float64(opts.Length) * perCharBits

	return totalBits, perCharBits, charsetSize, nil
}

//
// ──────────────────────────────────────────────────────────────────────────────
//  Secure helpers
// ──────────────────────────────────────────────────────────────────────────────
//

// secureRandInt returns a random integer in [0, n) using crypto/rand.
func secureRandInt(n int) (int, error) {
	if n <= 0 {
		return 0, fmt.Errorf("n must be positive")
	}
	max := big.NewInt(int64(n))
	v, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	return int(v.Int64()), nil
}

// secureShuffle shuffles a byte slice in-place using Fisher–Yates algorithm
// with crypto/rand for secure randomization.
func secureShuffle(b []byte) error {
	for i := len(b) - 1; i > 0; i-- {
		j, err := secureRandInt(i + 1)
		if err != nil {
			return err
		}
		b[i], b[j] = b[j], b[i]
	}
	return nil
}
