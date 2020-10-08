package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func readHash(fname string) string {
	file, err := ioutil.ReadFile("/go/src/app/files/" + fname)
	if err != nil {
		panic(err)
	}
	fmt.Println()

	return string(file)
}

func getPongs() string {
	resp, err := http.Get("http://ponger-svc/pongs")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return "Pongs: " + string(body)
}

func getResponse(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RRRQST")
	fmt.Fprintln(w, getData())
}

func meme(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "fuck you")
}
func getData() string {
	return os.Getenv("MESSAGE") + "\n" + readHash("hashes.txt") + "\n" + getPongs()
}
func main() {
	_ = uuid.New().String()
	fmt.Println("IM ALIIIIVE")
	router := mux.NewRouter()

	router.HandleFunc("/hashes/", getResponse).Methods("GET")
	router.HandleFunc("/", meme).Methods("GET")

	port := "8081"
	fmt.Println("Server starting in port", port)
	http.ListenAndServe(":"+port, router)

	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "DELETE", "PATCH", "OPTIONS"},
	}).Handler(router)

	http.ListenAndServe(":"+port, handler)

}
