package main

import (
	"io"
	"os"
	"net/http"
	"log"
	"fmt"
	"time"
)

func main() {
	http.HandleFunc("/", rootHandler)
	log.Fatal(http.ListenAndServe(":9901", nil))
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "video/mp4")
	w.WriteHeader(http.StatusOK)

	fi, err := os.Open("videoplayback.mp4")
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 1024)

	counter := 0
	for {
		// read a chunk
		n, err := fi.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}

		fmt.Println("Write a chunk", counter)
		counter++

		time.Sleep(100 * time.Millisecond) // "very bad Internet connection" emulation

		// write a chunk
		if _, err := w.Write(buf[:n]); err != nil {
			panic(err)
		}
	}
}
