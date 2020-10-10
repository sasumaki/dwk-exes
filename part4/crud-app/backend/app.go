package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path"

	"database/sql"

	"github.com/gorilla/mux"
	nats "github.com/nats-io/nats.go"
	"github.com/rs/cors"

	_ "github.com/lib/pq"
)

var (
	host     = "postgres-svc"
	port     = 5432
	user     = "postgres"
	password = os.Getenv("POSTGRES_PASSWORD")
	dbname   = "postgres"
)

var psqlInfo = fmt.Sprintf("host=%s port=%d user=%s "+
	"password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)

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

func openDb() (*sql.DB, error) {
	db, err := sql.Open("postgres", psqlInfo)

	err = db.Ping()
	fmt.Println("Successfully connected to postgres!")

	return db, err
}
func todos(w http.ResponseWriter, r *http.Request) {
	db, err := openDb()

	defer db.Close()
	var data []Todo

	if r.Method == http.MethodPost {
		var p Todo
		err := json.NewDecoder(r.Body).Decode(&p)
		fmt.Println("new todo", p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		insertStatement := `
			INSERT INTO todos (title, done) VALUES ($1, $2)
		`
		_, err = db.Exec(insertStatement, p.Title, p.Done)
		if err != nil {
			fmt.Println("error: ", &err)
			panic(err)
		}
		fmt.Println(os.Getenv("NATS_URL"))
		nc, err := nats.Connect(os.Getenv("NATS_URL"))
		if err != nil {
			panic(err)
		}
		defer nc.Close()
		fmt.Println(nc)
		nc.Publish("todos", []byte(p.Title+" todo added"))
	}
	fmt.Println("getting todos")
	getStatement := `
			SELECT * FROM todos
		`
	rows, err := db.Query(getStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var title string
		var done bool
		err = rows.Scan(&id, &title, &done)
		if err != nil {
			// handle this error
			panic(err)
		}
		data = append(data, Todo{Title: title, Done: done})
	}
	fmt.Println(data)

	json.NewEncoder(w).Encode(data)

}

func meme(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "jajajjajajaja")
}
func health(w http.ResponseWriter, r *http.Request) {
	db, err := openDb()
	defer db.Close()
	if err != nil {
		http.Error(w, "shit's down", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	fmt.Fprintln(w, "Healthy")
}
func toggleTodo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Toggle todo called")
	var data []Todo

	title := path.Base(r.URL.Path)
	db, _ := openDb()
	defer db.Close()

	updateStatement := `
		UPDATE todos
		SET done = true
		WHERE title = $1;
	`
	db.Exec(updateStatement, title)
	nc, _ := nats.Connect(os.Getenv("NATS_URL"))
	defer nc.Close()
	nc.Publish("todos", []byte(title+" todo doned"))
	getStatement := `
			SELECT * FROM todos
		`
	rows, err := db.Query(getStatement)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var title string
		var done bool
		err = rows.Scan(&id, &title, &done)
		if err != nil {
			// handle this error
			panic(err)
		}
		data = append(data, Todo{Title: title, Done: done})
	}
	fmt.Println(data)

	json.NewEncoder(w).Encode(data)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/health", health).Methods("GET", "POST")

	db, _ := openDb()
	defer db.Close()

	statement := `CREATE TABLE IF NOT EXISTS todos (id SERIAL PRIMARY KEY, title TEXT, done BOOL)`
	db.Exec(statement)

	router.HandleFunc("/api/todos/", todos).Methods("GET", "POST")
	router.HandleFunc("/", meme).Methods("GET", "POST")
	router.HandleFunc("/api/todos/{title}", toggleTodo).Methods("GET", "PUT", "POST")
	port := "3000"
	fmt.Println("Server starting in port", port)
	http.ListenAndServe(":"+port, router)
	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "DELETE", "PATCH", "OPTIONS"},
	}).Handler(router)

	http.ListenAndServe(":"+port, handler)

}
