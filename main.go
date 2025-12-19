package main

import (
	"fmt"
	"net/http"
)

func main(){
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	
	mux.HandleFunc("POST /users", createUser)
	fmt.Println("server listening to 8080")

	http.ListenAndServe(":8080", mux)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func createUser(w http.ResponseWriter, r *http.Request){
	
}