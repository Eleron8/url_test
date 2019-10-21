package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net/http"
	"os"
	"sync"
)

type Answer struct {
	Url   string
	Count int
}

func main() {
	answers := make(chan Answer)
	end := make(chan struct{})
	go func(answers <-chan Answer, end chan<- struct{}) {
		var total int
		for answer := range answers {
			fmt.Printf("Count for %s: %d\n", answer.Url, answer.Count)
			total += answer.Count
		}
		fmt.Printf("Total: %d\n", total)
		end <- struct{}{}
	}(answers, end)
	var wg sync.WaitGroup
	k := make(chan struct{}, 5)
	for _, url := range os.Args[1:] {
		k <- struct{}{}
		wg.Add(1)
		go func(url string, answers chan<- Answer, k <-chan struct{}, wg *sync.WaitGroup) {
			answers <- Get(url)
			<-k
			wg.Done()
		}(url, answers, k, &wg)
	}
	wg.Wait()
	close(k)
	close(answers)
	<-end
}

func Get(url string) Answer {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return Answer{}
	}

	scan := bufio.NewScanner(resp.Body)
	var count int

	for bts := scan.Scan(); bts != false; bts = scan.Scan() {
		c := bytes.Count(scan.Bytes(), []byte("Go"))
		count = count + c
	}
	result := Answer{
		Url:   url,
		Count: count,
	}
	return result
}
