package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/google/uuid"
)

func readFile(fname string, uuid string) {
	file, err := ioutil.ReadFile("/go/src/app/files/" + fname)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(file), uuid)
	pongFile, err := ioutil.ReadFile("/go/src/app/files/pongs.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(pongFile))
}

func main() {
	uuid := uuid.New().String()

	ping := func() { readFile("hashes.txt", uuid) }

	tick := time.Tick(5000 * time.Millisecond)
	for range tick {
		ping()
	}

}
