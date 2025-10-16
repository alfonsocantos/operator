# The Operators Library

## How to use

### Run commands

It only supports public key authentication.

```go
term, _ := NewTerm(PrivatekeyBytes, "my-user", "my-sftp:22")

// run remote command
output, _ := term.Run("ls -l")
```

### Upload File using ssh (sftp)

It only supports public key authentication.

```go
term, _ := NewTerm(PrivatekeyBytes, "my-user", "my-sftp:22")

// send the file
sizeBytes, _ := term.Upload("file-path", "file-remote-path")
```
If it is successful, `sizeBytes` is the size in bytes sent to destination.

## Secure utilities

The `secureutils` package provides cryptographically secure helpers for generating
random identifiers, tokens, and passwords.

Import it from within the module:

```go
import "operator/secureutils"
```

### Generate random strings

Use `secureutils.RandomString` when you need a random string built from a custom
character set (for example, when creating invitation codes). Pass the desired
length and the characters that can be used:

```go
charset := "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
token, err := secureutils.RandomString(12, charset)
if err != nil {
        // handle error
}
```

### Generate hexadecimal or Base64 tokens

If you only need hexadecimal or Base64-URL output you can rely on the helpers
that already choose a secure character set and size for you:

```go
hexToken, _ := secureutils.RandomHex(16)      // 32 hex characters
b64Token, _ := secureutils.RandomBase64(24)  // URL-safe Base64 without padding
```

### Create strong passwords

`secureutils.GeneratePassword` builds strong passwords with sensible defaults
(16 characters, mixed character classes, no ambiguous characters). You can
customise the behaviour via `secureutils.PasswordOptions`:

```go
password, err := secureutils.GeneratePassword(secureutils.PasswordOptions{
        Length:            20,
        UseLower:          true,
        UseUpper:          true,
        UseDigits:         true,
        UseSymbols:        false, // only alphanumerics
        AvoidAmbiguous:    true,
        RequireEachChosen: true,
})
if err != nil {
        // handle error
}
```

When auditing password strength you can estimate its entropy using
`secureutils.EstimateEntropy`, which returns the number of bits of entropy per
character and in total for the chosen options.

## Develop the library

### Download dependencies

```sh
make deps
```

### Run tests

```sh
make test
```
