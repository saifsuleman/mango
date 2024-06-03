package main

import (
	"fmt"
    "io"
	"log"
	"net/http"
)

const (
    TEXT = 0
    IMG = 1
    OK = 2
    ERROR = 3
)

func main() {
	fmt.Println("Hello world!")

    httpCache := NewLoadingCacheBuilder[string, string]().
        MaximumSize(10000).
        ExpirationSeconds(60).
        WithLoad(computeHttpResponse).
        Build()

    for range 10 {
        response, err := httpCache.Get("1.1.1.1")
        if err != nil {
            log.Fatalln(err)
        }
        fmt.Printf("received response: %s\n", response)
    }
}

func computeHttpResponse(id string) (str string, err error) {
    url := fmt.Sprintf("http://check.getipintel.net/check.php?ip=%s&contact=admin@vortexnetwork.net&flags=f", id)
    fmt.Printf("computing result from %s\n", url)
    response, err := http.Get(url)
    if err != nil {
        return
    }

    body, err := io.ReadAll(response.Body)
    if err != nil {
        return
    }

    str = string(body)
    return
}
