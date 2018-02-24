package main

import (
	"io"
	"os"
	"net/http"
	"fmt"
	"time"
	"runtime"
	"flag"
	"io/ioutil"

	"github.com/gorilla/mux"
	"os/exec"
	"strconv"
	"encoding/json"
	"strings"
	"log"
)

var cfg struct {
	httpPort   uint
	bufferSize uint
}

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.UintVar(&cfg.httpPort, "http-port", 8080, "HTTP-demon port")
	flag.UintVar(&cfg.bufferSize, "buffer-size", 1024, "Video buffer size (bytes)")
	flag.Parse()
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/{filepath:[0-9a-zA-Z/-]+}.{extension:mp4|webm|ogg}", mediaHandler).Methods("GET")
	router.HandleFunc("/{filepath:[0-9a-zA-Z/-]+}.{extension:mp4|webm|ogg}.meta", metaHandler).Methods("GET")
	router.HandleFunc("/{filepath:[0-9a-zA-Z/-]+}.{extension:html|css|js|png|jpg}", staticHandler).Methods("GET")
	http.Handle("/", router)
	http.ListenAndServe(fmt.Sprintf(":%d", cfg.httpPort), nil)
}

func mediaHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	fi, err := os.Open("." + req.RequestURI)
	if err != nil {
		log.Print("Video not found: " + req.RequestURI)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	buf := make([]byte, cfg.bufferSize)

	w.Header().Set("Content-Type", fmt.Sprintf("video/%s", params["extension"]))
	w.WriteHeader(http.StatusOK)

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

		//time.Sleep(100 * time.Millisecond) // "very bad Internet connection" emulation
		time.Sleep(5 * time.Millisecond) // "very bad Internet connection" emulation

		// write a chunk
		if _, err := w.Write(buf[:n]); err != nil {
			panic(err)
		}
	}
}

type Meta struct {
	Duration float64 `json:"duration"`
}

func metaHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	filename := fmt.Sprintf("%s.%s", params["filepath"], params["extension"])

	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filename)
	dr, err := cmd.Output()
	if err != nil {
		log.Print("Video not found: " + filename)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	meta := new(Meta)
	drs := string(dr)
	drs = strings.Replace(drs, string('\n'), "", 10)
	d, err := strconv.ParseFloat(drs, 64)
	if err != nil {
		panic(err)
	}
	meta.Duration = d

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response, err := json.Marshal(meta)
	if err != nil {
		panic(err)
	}
	w.Write(response)
}

func staticHandler(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	mimeSufix := "plain"
	switch params["extension"] {
	case "js":
		mimeSufix = "javascript"
	default:
		mimeSufix =params["extension"]
	}

	w.Header().Set("Content-Type", "text/"+mimeSufix)
	w.WriteHeader(http.StatusOK)

	fi, err := os.Open(fmt.Sprintf("cmd/streamdemonsvc/%s.%s",
		params["filepath"],
		params["extension"],
	))
	if err != nil {
		panic(err)
	}

	n, err := ioutil.ReadAll(fi)
	if err != nil && err != io.EOF {
		panic(err)
	}

	w.Write(n)
}
