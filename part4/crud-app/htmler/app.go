package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"image/jpeg"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"time"

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

func todosRequest(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("forms.html"))
	var todos []Todo

	if r.Method != http.MethodPost {
		fmt.Println("GET")
		res, err := http.Get("http://todo-api-svc/api/todos/")
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		fmt.Println(body)
		fmt.Println("json parsed: ")
		err = json.Unmarshal([]byte(body), &todos)
		fmt.Println(todos)
		if err != nil {
			fmt.Println("error:", err)
		}
	} else {
		fmt.Println("POST")

		fmt.Println("new todo", Todo{
			Title: r.FormValue("todo"),
			Done:  false,
		})
		postJson, err := json.Marshal(Todo{
			Title: r.FormValue("todo"),
		})
		if err != nil {
			fmt.Println("error:", err)
		}
		res, err := http.Post("http://todo-api-svc/api/todos/", "application/json", bytes.NewBuffer(postJson))
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		fmt.Println(body)
		err = json.Unmarshal([]byte(body), &todos)
		fmt.Println("todoos")
		fmt.Println(todos)

	}
	data := TodoPageData{
		Todos: todos,
	}
	tmpl.Execute(w, data)

}
func todoUpdate(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("forms.html"))

	var todos []Todo
	fmt.Println("Update requested.")
	title := path.Base(r.URL.Path)
	fmt.Println(title)
	res, err := http.Get("http://todo-api-svc/api/todos/" + title)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(body)
	err = json.Unmarshal([]byte(body), &todos)
	fmt.Println("todoos")
	fmt.Println(todos)

	data := TodoPageData{
		Todos: todos,
	}
	tmpl.Execute(w, data)
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
func health(w http.ResponseWriter, r *http.Request) {
	fmt.Println("READINESS CHECK")
	_, err := http.Get("http://todo-api-svc/api/health")
	if err != nil {
		fmt.Println("NOT RDY")
		fmt.Println(err)
		time.Sleep(60 * time.Second)
		http.Error(w, "cant connet api", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "gucci")
}

func main() {
	existingImageFile, err := os.Open("files/image.jpg")
	if err != nil {
		downloadFile("https://picsum.photos/1200", "files/image.jpg")
	}
	defer existingImageFile.Close()

	router := mux.NewRouter()
	router.HandleFunc("/view/", viewHandler).Methods("GET")
	router.HandleFunc("/todos/", todosRequest).Methods("GET", "POST")
	router.HandleFunc("/todos/{title}", todoUpdate).Methods("GET", "PUT")
	router.HandleFunc("/health", health).Methods("GET")

	port := "8080"
	fmt.Println("Server starting in port", port)
	http.ListenAndServe(":"+port, router)
	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "DELETE", "PATCH", "OPTIONS"},
	}).Handler(router)

	http.ListenAndServe(":"+port, handler)

}
