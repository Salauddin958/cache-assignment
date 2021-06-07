package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
)

const httpHOST = "http://godoc.org"
const httpsHOST = "https://godoc.org"

var url string

var cachePage = cache.New(5*time.Minute, 10*time.Minute)

func loadWebsite(w http.ResponseWriter, r *http.Request) {
	if r.Host == "localhost:8080" {
		url = httpHOST + r.URL.Path
	} else {
		url = httpsHOST + r.URL.Path
	}
	cachedResponse, found := cachePage.Get(r.URL.Path)
	fmt.Println("url : ", url)
	fmt.Println("HOST :: ", r.Host)
	if found {
		fmt.Println("Loading from cache")
		fmt.Fprintf(w, cachedResponse.(string))
	} else {
		response, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer response.Body.Close()
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		str := string(bodyBytes)
		cachePage.Set(r.URL.Path, str, cache.DefaultExpiration)
		w.Header().Set("Cache-Control", "max-age=60")
		fmt.Fprintf(w, str)
	}
}

func main() {

	http.HandleFunc("/github.com/stretchr/testify/assert", loadWebsite)
	go http.ListenAndServe(":8080", nil)
	err1 := http.ListenAndServeTLS(":8090", "https-server.crt", "https-server.key", nil)
	if err1 != nil {
		log.Fatal(err1)
	}

}
