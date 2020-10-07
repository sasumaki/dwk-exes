package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	pongs = 0
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Serving request.")

	fmt.Fprintf(w, "<p>%s</p>", "pong "+strconv.Itoa(pongs))
	pongs = pongs + 1
	f, _ := os.OpenFile("files/pongs.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 7777)
	_, err := fmt.Fprintln(f, "Ping Pongs: "+strconv.Itoa(pongs))
	if err != nil {
		panic(err)
	}
}

func main() {
	http.HandleFunc("/", handler)

	port := "3000"
	fmt.Println("Server started in port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
