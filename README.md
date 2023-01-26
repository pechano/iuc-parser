# iuc-parser
Small application to parse .i6z dossiers

The program uses a local key which is provided as a .csv file when prompted.
A default key is also available in the source code.

To build, simply run 

> go build bpr.go structs.go

The program was developed on Linux, but crosscompiles to windows with the prefix 
> GOOS=windows GOARCH=amd64

