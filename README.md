# The Operators Library

## How to use

### Run commands

It only support public ket authentication.

```go

term, _ := NewTerm(PrivatekeyBytes, "my-user", "my-sftp:22")

// send the file
output, _ := term.Run("ls -l")

```


### Upload File using ssh (sftp)

It only support public ket authentication.

```go

term, _ := NewTerm(PrivatekeyBytes, "my-user", "my-sftp:22")

// send the file
size_bytes, _ := term.Upload("file-path", "file-remote-path")

```
If it is successful, `size_bytes` is the size in bytes sent to destination

# Develop the library

##  download dependencies

`$ make deps`



##  run tests

`$ make test`

