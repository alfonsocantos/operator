package secureutils

import (
	"encoding/hex"
	"regexp"
	"testing"
)

// --- Helpers for tests ---

func hasAny(s string, chars string) bool {
	for _, c := range s {
		for _, b := range chars {
			if c == b {
				return true
			}
		}
	}
	return false
}

func countUnique(ss []string) int {
	seen := make(map[string]struct{}, len(ss))
	for _, s := range ss {
		seen[s] = struct{}{}
	}
	return len(seen)
}

func isURLSafeBase64(s string) bool {
	// Base64 URL-safe without padding: only A-Z a-z 0-9 - _
	re := regexp.MustCompile(`^[A-Za-z0-9_-]+$`)
	return re.MatchString(s)
}

// --- Tests ---

func TestGeneratePassword_Defaults(t *testing.T) {
	pw, err := GeneratePassword(PasswordOptions{})
	if err != nil {
		t.Fatalf("GeneratePassword (defaults) error: %v", err)
	}
	if len(pw) != 16 {
		t.Fatalf("expected default length 16, got %d", len(pw))
	}

	// Defaults avoid ambiguous: no '0','1','l','I','O'
	if hasAny(pw, "01lIO") {
		t.Fatalf("password contains ambiguous characters: %q", pw)
	}
}

func TestGeneratePassword_RequireEachChosen(t *testing.T) {
	opts := PasswordOptions{
		Length:            24,
		UseLower:          true,
		UseUpper:          true,
		UseDigits:         true,
		UseSymbols:        true,
		AvoidAmbiguous:    true,
		RequireEachChosen: true,
	}
	pw, err := GeneratePassword(opts)
	if err != nil {
		t.Fatalf("GeneratePassword error: %v", err)
	}
	if len(pw) != opts.Length {
		t.Fatalf("expected length %d, got %d", opts.Length, len(pw))
	}

	// With AvoidAmbiguous=true, sets are:
	lower := "abcdefghijkmnopqrstuvwxyz" // no 'l'
	upper := "ABCDEFGHJKLMNPQRSTUVWXYZ"  // no 'I','O'
	digits := "23456789"                 // no '0','1'
	symbols := "!@#$%^&*()-_=+[]{};:,.?/~"

	if !hasAny(pw, lower) {
		t.Fatalf("password missing lowercase: %q", pw)
	}
	if !hasAny(pw, upper) {
		t.Fatalf("password missing uppercase: %q", pw)
	}
	if !hasAny(pw, digits) {
		t.Fatalf("password missing digits: %q", pw)
	}
	if !hasAny(pw, symbols) {
		t.Fatalf("password missing symbols: %q", pw)
	}
	// still ensure no ambiguous chars
	if hasAny(pw, "01lIO") {
		t.Fatalf("password contains ambiguous characters: %q", pw)
	}
}

func TestGeneratePassword_LengthTooShort(t *testing.T) {
	// Require 4 sets but length 3 should fail
	_, err := GeneratePassword(PasswordOptions{
		Length:            3,
		UseLower:          true,
		UseUpper:          true,
		UseDigits:         true,
		UseSymbols:        true,
		RequireEachChosen: true,
	})
	if err == nil {
		t.Fatalf("expected error for length too short when RequireEachChosen=true")
	}
}

func TestGeneratePassword_RandomnessBasic(t *testing.T) {
	// Not a statistical test, just checks low collision likelihood
	N := 50
	samples := make([]string, 0, N)
	for i := 0; i < N; i++ {
		pw, err := GeneratePassword(PasswordOptions{})
		if err != nil {
			t.Fatalf("GeneratePassword error: %v", err)
		}
		samples = append(samples, pw)
	}
	uniq := countUnique(samples)
	if uniq < 45 {
		t.Fatalf("too many duplicates: %d unique out of %d", uniq, N)
	}
}

func TestRandomString_OK_AndErrors(t *testing.T) {
	// OK
	rs, err := RandomString(10, "abcXYZ")
	if err != nil {
		t.Fatalf("RandomString error: %v", err)
	}
	if len(rs) != 10 {
		t.Fatalf("expected length 10, got %d", len(rs))
	}
	if !hasAny(rs, "abcXYZ") {
		t.Fatalf("result should contain only provided charset: %q", rs)
	}
	// Errors
	if _, err := RandomString(5, ""); err == nil {
		t.Fatalf("expected error for empty charset")
	}
	if _, err := RandomString(0, "abc"); err == nil {
		t.Fatalf("expected error for non-positive length")
	}
}

func TestRandomHex(t *testing.T) {
	h, err := RandomHex(16)
	if err != nil {
		t.Fatalf("RandomHex error: %v", err)
	}
	if len(h) != 32 {
		t.Fatalf("expected hex length 32 for 16 bytes, got %d", len(h))
	}
	if _, err := hex.DecodeString(h); err != nil {
		t.Fatalf("invalid hex string: %v", err)
	}
	if _, err := RandomHex(0); err == nil {
		t.Fatalf("expected error for nBytes<=0")
	}
}

func TestRandomBase64(t *testing.T) {
	s, err := RandomBase64(16)
	if err != nil {
		t.Fatalf("RandomBase64 error: %v", err)
	}
	if !isURLSafeBase64(s) {
		t.Fatalf("base64 not URL-safe or contains invalid chars: %q", s)
	}
	if hasAny(s, "=") {
		t.Fatalf("base64 string should be unpadded (no '='): %q", s)
	}
	if _, err := RandomBase64(0); err == nil {
		t.Fatalf("expected error for nBytes<=0")
	}
}

func TestEstimateEntropy_Basics(t *testing.T) {
	total, perChar, size, err := EstimateEntropy(PasswordOptions{})
	if err != nil {
		t.Fatalf("EstimateEntropy error: %v", err)
	}
	// With defaults (AvoidAmbiguous + all sets):
	// lower(25) + upper(24) + digits(8) + symbols(len("!@#$%^&*()-_=+[]{};:,.?/~"))
	symbols := "!@#$%^&*()-_=+[]{};:,.?/~"
	expectedSize := 25 + 24 + 8 + len(symbols)

	if size != expectedSize {
		t.Fatalf("unexpected charset size: got %d, want %d", size, expectedSize)
	}
	if total <= 0 || perChar <= 0 {
		t.Fatalf("entropy values should be positive, got total=%.4f perChar=%.4f", total, perChar)
	}

	// Monotonicity: increasing length increases total bits
	total2, _, _, err := EstimateEntropy(PasswordOptions{Length: 24})
	if err != nil {
		t.Fatalf("EstimateEntropy error (length 24): %v", err)
	}
	if total2 <= total {
		t.Fatalf("entropy should increase with length: %f !< %f", total, total2)
	}
}
