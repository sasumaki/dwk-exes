package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func readHash(fname string, uuid string) string {
	file, err := ioutil.ReadFile("/go/src/app/files/" + fname)
	if err != nil {
		panic(err)
	}
	fmt.Println()

	return string(file) + ": " + uuid
}

func getPongs() string {
	resp, err := http.Get("http://ponger-svc/pongs")
	if err != nil {
		time.Sleep(60 * time.Second)
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(body)
}

func getResponse(data string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("RRRQST")
		fmt.Fprintln(w, data)
	}
}
func meme(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "fuck you")
}

func main() {
	uuid := uuid.New().String()
	fmt.Println("IM ALIIIIVE")
	router := mux.NewRouter()

	data := readHash("hashes.txt", uuid) + "\n" + getPongs()
	fmt.Println(data)
	router.HandleFunc("/hashes/", getResponse(data)).Methods("GET")
	router.HandleFunc("/", meme).Methods("GET")

	port := "8081"
	fmt.Println("Server starting in port", port)
	http.ListenAndServe(":"+port, router)

	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "DELETE", "PATCH", "OPTIONS"},
	}).Handler(router)

	http.ListenAndServe(":"+port, handler)

}
