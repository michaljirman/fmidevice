[![Build Status](https://travis-ci.org/michaljirman/coreapi.svg?branch=master)](https://travis-ci.org/michaljirman/coreapi)
# COREAPI

A simple REST API support library written in Go (Golang).


## Usage
```go
var api = coreapi.NewAPI("https://httpbin.org")

func main() {
	router := coreapi.NewRouter()
	router.RegisterFunc(200, func(resp *http.Response) error {
		defer resp.Body.Close()
		content, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Println(string(content))
		return nil
	})
	resource := coreapi.NewResource("/get", "GET", router)
	if err := api.Call(resource, nil, nil, nil); err != nil {
		log.Fatalln(err)
	}
}
```

```go
./coreapi
{"args":{},"headers":{"Accept-Encoding":"gzip","Connection":"close","Host":"httpbin.org","User-Agent":"Go-http-client/1.1"},"origin":"82.37.173.154","url":"https://httpbin.org/get"}
```

## Testing
```go
go test -v ./...
```

## Build & install 
```go
go build ./cmd/...
```
OR
```go
go build -i -o coreapi ./cmd/...
```
OR
```go
go install ./cmd/...
ls ~/go/bin
```


## Requirements
#### GO installation (e.g. brew install go)
#### ~/.bash_profile or similar
```go
export GOPATH=/Users/USER/go
export PATH=$GOPATH/bin:$PATH
source ~/.bash_profile
```

## Build & install 
### Install
```go
go install ./cmd/...
```