package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// Make request
	fmt.Println("running")
	site := os.Args[1]
	fmt.Println(site)
	response, err := http.Get(site)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Create output file
	outFile, err := os.Create("output.html")
	if err != nil {
		log.Fatal(err)
	}
	defer outFile.Close()

	// Copy data from HTTP response to file
	_, err = io.Copy(outFile, response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response.Body)
}
