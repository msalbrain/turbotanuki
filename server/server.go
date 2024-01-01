package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/get", handleGet)
	http.HandleFunc("/post", handlePost)
	http.HandleFunc("/put", handlePut)
	http.HandleFunc("/delete", handleDelete)
	http.HandleFunc("/patch", handlePatch)

	fmt.Println("Server is running on :13001")
	http.ListenAndServe(":13001", nil)
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	randomStatusResponse(w)
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	randomStatusResponse(w)
}

func handlePut(w http.ResponseWriter, r *http.Request) {
	randomStatusResponse(w)
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	randomStatusResponse(w)
}

func handlePatch(w http.ResponseWriter, r *http.Request) {
	randomStatusResponse(w)
}

func randomStatusResponse(w http.ResponseWriter) {
	// Simulate random behavior
	randomNumber := rand.Intn(3)

	switch randomNumber {
	case 0:
		// Success
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Success"))
	case 1:
		// Redirect
		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte("Redirect"))
	default:
		// Failure
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failure"))
	}
}



