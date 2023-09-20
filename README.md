# iuc-parser
Small application to parse .i6z dossiers

Binaries are available for Windows and Linux under Releases.

## Building

It can be built by cloning the repo and running
```
go build bpr.go structs.go
```
This software was developed on Linux, but crosscompiles to Windows with the prefix
```
GOOS=windows GOARCH=amd64
```
An external .csv containing the keys needed to rename and move extracted files must be sourced manually.
