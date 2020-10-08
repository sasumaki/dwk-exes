package main

import (
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

func main() {
	ping := func() {
		timestamp := time.Now().UTC().Format(time.UnixDate)
		uuid := uuid.New().String()
		newLine := (timestamp + ": " + uuid)

		f, _ := os.OpenFile("files/hashes.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 7777)
		_, err := fmt.Fprintln(f, newLine)

		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
		err = f.Close()
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	tick := time.Tick(5000 * time.Millisecond)
	for range tick {
		ping()
	}

}
