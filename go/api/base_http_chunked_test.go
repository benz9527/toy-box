package api

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"testing"
	"time"
)

func server(wg *sync.WaitGroup) {
	http.HandleFunc("/xxx", func(w http.ResponseWriter, r *http.Request) {
		flusher := w.(http.Flusher)
		for i := 0; i < 2; i++ {
			fmt.Fprintf(w, "Ben\n")
			flusher.Flush()
			<-time.Tick(1 * time.Second)
		}
	})
	go func() {
		wg.Done()
		_ = http.ListenAndServe(":8421", nil)
	}()
}

func client(wg *sync.WaitGroup) {
	wg.Wait()
	<-time.Tick(1 * time.Second)
	resp, err := http.Get("http://127.0.0.1:8421/xxx")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.TransferEncoding)

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			log.Fatal(err)
		}
		if len(line) > 0 {
			fmt.Print(line)
		}
	}
}

func TestWatch(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	server(&wg)
	client(&wg)
}
