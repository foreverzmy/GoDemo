package main

import (
	"fmt"
	"net/http"
	"sync"
)

var mu sync.Mutex

// Count 计数器
var Count int

// CountHadnler echoes the number of calls so far.
func CountHadnler(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "Count %d\n", Count)
	mu.Unlock()
}
