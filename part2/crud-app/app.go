package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"html/template"
	"image/jpeg"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type TodoPageData struct {
	Todos []Todo
}

type TodoForm struct {
	Todo string
}
type Todo struct {
	Title string
	Done  bool
}

func todos(w http.ResponseWriter, r *http.Request) {
	data := TodoPageData{
		Todos: []Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: true},
			{Title: "Task 3", Done: true},
		},
	}
	tmpl := template.Must(template.ParseFiles("forms.html"))

	if r.Method != http.MethodPost {
		fmt.Println("this weird shit")
		tmpl.Execute(w, data)
		return
	}
	fmt.Println("???")
	details := TodoForm{
		Todo: r.FormValue("todo"),
	}

	// do something with details
	_ = details

	tmpl.Execute(w, details)

}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Serving request.")
	existingImageFile, err := os.Open("files/image.jpg")
	if err != nil {
		panic(err)
	}
	defer existingImageFile.Close()

	loadedImage, err := jpeg.Decode(existingImageFile)
	var buff bytes.Buffer

	jpeg.Encode(&buff, loadedImage, nil)

	// Encode the bytes in the buffer to a base64 string
	encodedString := base64.StdEncoding.EncodeToString(buff.Bytes())
	// You can embed it in an html doc with this string
	htmlImage := "<img src=\"data:image/jpeg;base64," + encodedString + "\" />"
	fmt.Fprintln(w, htmlImage)
}

func downloadFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("Received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	existingImageFile, err := os.Open("files/image.jpg")
	if err != nil {
		downloadFile("https://picsum.photos/1200", "files/image.jpg")
	}
	defer existingImageFile.Close()

	router := mux.NewRouter()
	router.HandleFunc("/view/", viewHandler).Methods("GET")
	router.HandleFunc("/todos/", todos).Methods("GET", "POST")
	port := "8080"
	fmt.Println("Server starting in port", port)
	http.ListenAndServe(":"+port, router)
	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "DELETE", "PATCH", "OPTIONS"},
	}).Handler(router)

	http.ListenAndServe(":"+port, handler)

}
