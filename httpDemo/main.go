package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", handler)
	http.HandleFunc("/headers", HeaderHandle)
	http.HandleFunc("/count", CountHadnler)
	http.HandleFunc("/image", ImageHandler)
	http.HandleFunc("/lissajous", LissajousHandler)
	http.HandleFunc("/mandelbrot", MandelbrotHandler)

	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	Count++
	mu.Unlock()
	fmt.Fprintf(w, "hello")
}
