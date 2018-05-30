package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/michaljirman/coreapi"
)

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
