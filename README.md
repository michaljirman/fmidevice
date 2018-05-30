# fmidevice

A command line program and library which locates all iDevices associated with an Apple ID registered on the iCloud services. 

By default, a web client is used to access location data. A user is notified via email same way as he would be in the case of accessing an iCloud services through a web browser. Additionally, the 2FA (Two Factor Authentication) is triggered on devices associated with an Apple ID.

An alternative way is to use a mobile client through silent mode option. A user won't be notified via email and the 2FA is not triggered. 

## Usage

### CLI 
#### Print usage
```go
./fmidevice --help
iDevice locating tool

Usage:
  fmidevice [command]

Available Commands:
  help        Help about any command
  locate      Locate iDevices for an Apple account

Flags:
  -h, --help   help for fmidevice

Use "fmidevice [command] --help" for more information about a command.
```

```go
./fmidevice locate --help
Locate iDevices for an Apple account

Usage:
  fmidevice locate [flags]

Flags:
      --account string    An Apple account name
  -h, --help              help for locate
      --password string   An Apple password
      --silent            Silent mode allows to access location through a mobile client
```

```go
./fmidevice locate --account <APPLE_ID> --password <APPLE_PASSWORD> --silent
----------
[1] Device: iPad Air 2, Tester’s iPad ... Location: 0.000000, 0.000000
----------
[2] Device: iPhone 3GS, Tester’s iPhone ... Location: 0.000000, 0.000000
----------
[3] Device: iPhone 7, iPhone ... Location: 51.470493, 0.124431
----------
[4] Device: iMac 21.5", iMac ... Location: 0.000000, 0.000000
----------
[5] Device: iPhone 7, Tester’s iPhone ... Location: 0.000000, 0.000000
----------
[6] Device: iMac 27", iMac ... Location: 0.000000, 0.000000
----------
[7] Device: MacBook Pro 13", MacBook Pro ... Location: 51.470482, 0.124411
```

#### File configuration
A configuration file such as (~/.fmidevice.yaml) can be loaded in an application if exists in a home directory. 
In silent mode, it allows to configure an AccountName (Apple ID), a PrsId and an AuthToken
to bypass authentication via an AccountName and a Password.
Normal (default) mode does not use content of a configuration file.

Configuration flags are prioritised over a values from a configuration file.

Example of a content present in a ~/.fmidevice.yaml file:
```go
---
AccountName: tester.test@gmail.com
prsID: 422669977
AuthToken: IAAAAAAABLwIAAAAAFsNsMMRjbG91ZC5.....bwCrDpbTR3MOSXbonheI2OivFvl9JW4FKA~~
```

If a valid configuration file exists then you can run binary without parameters in silent mode.
```go
./fmidevice locate --silent
```

### Library

#### Testing
```go
go test -v ./...
```

### Requirements
#### GO installation (e.g. brew install go)
#### ~/.bash_profile or similar
```go
export GOPATH=/Users/USER/go
export PATH=$GOPATH/bin:$PATH
source ~/.bash_profile
```

### Build & install 
### Install
```go
go install ./cmd/...
```

### Proxy configuration
Proxy server can be configured by setting evironmental variables such as HTTP_PROXY and HTTPS_PROXY.

```go
HTTP_PROXY="http://localhost:8888" HTTPS_PROXY="http://localhost:8888" ./fmidevice locate --account <APPLE_ID> --password <APPLE_PASSWORD> --silent
```