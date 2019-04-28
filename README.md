http-server
===========

[![Build Status](https://dev.azure.com/dohrmichael/commons/_apis/build/status/dohr-michael.http-server?branchName=master)](https://dev.azure.com/dohrmichael/commons/_build/latest?definitionId=2&branchName=master)

- Commands (from Magefile)
    - Run unit test : `mage test` 
    - Build locally : `mage build` 
- Instal dependencies
```
go mod download
```
- Run project
```
go run main.go start
```
- Hot Reload
```
go get github.com/codegangsta/gin
gin --appPort 8080 --buildArgs main.go -i run start
```
