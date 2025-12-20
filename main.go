package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

type User struct{
	Name string `json: "name"` 
}

var cacheMutex sync.RWMutex
var userCache = make(map[int] User)

func main(){
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)
	
	mux.HandleFunc("POST /users", createUser)
	mux.HandleFunc("GET /users/{id}", getUser)
	mux.HandleFunc("DELETE /users/{id}", deleteUser)
	fmt.Println("server listening to 8080")

	http.ListenAndServe(":8080", mux)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func deleteUser(w http.ResponseWriter, r *http.Request){
	index, err := strconv.Atoi(r.PathValue("id"))
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	cacheMutex.Lock()
	_, ok := userCache[index]
	cacheMutex.Unlock()
	
	if !ok {
		http.Error(w, "user not found", http.StatusNotFound)
	}
	cacheMutex.Lock()
	delete(userCache, index)
	cacheMutex.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func getUser(w http.ResponseWriter, r *http.Request){
	index, err := strconv.Atoi(r.PathValue("id"))
	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return 
	}
	user, ok := userCache[index]
		if !ok {
		http.Error(w, "user Not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(user)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return 
	}
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func createUser(w http.ResponseWriter, r *http.Request){
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if user.Name == ""{
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	cacheMutex.Lock()
	userCache[len(userCache) + 1] = user
	cacheMutex.Unlock()

	w.WriteHeader(http.StatusNoContent)
}
