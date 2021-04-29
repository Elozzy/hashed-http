package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"sync"
)
var wg sync.WaitGroup

func main() {

	parallel := flag.Int("parallel", 10, "number of parallel requests defaults to 10")
	flag.Parse()
	maxRequest := *parallel

	if maxRequest > 10 {
		log.Fatal("number of parallel requests should be less than or equal to 10")
	}
	concurrent := make(chan struct{}, maxRequest)
	ch := make(chan string)
	wg.Add(len(flag.Args()))
	for _, rawUrl := range flag.Args() {
		go fetchUrl(rawUrl, ch, concurrent)
	}
	for range flag.Args() {
		fmt.Println(<-ch)
		<-concurrent
	}

	wg.Wait()
}


func fetchUrl(rawUrl string, ch chan <- string, concurrent chan <- struct{})  {
	defer wg.Done()
	hostname := parseRawUrl(rawUrl)
	resp, err := http.Get(hostname)
	concurrent <- struct{}{}

	if err != nil {
		ch <- fmt.Sprintf( "fetch: %v\n", err)
		return
	}
	responseBody, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		ch <- fmt.Sprintf("fetch: reading %s: %v\n", hostname, err)
		return
	}

	ch <- fmt.Sprintf("%s %x\n", hostname, md5.Sum(responseBody))

}

func parseRawUrl(rawUrl string) (hostname string) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		log.Fatal(err)
	}
	if u.Scheme == "" && u.Host == ""{
		return "http://" + rawUrl
	}

	return rawUrl
}

