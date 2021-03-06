package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

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

func openDb() (*sql.DB, error) {
	db, err := sql.Open("postgres", psqlInfo)

	err = db.Ping()
	return db, err
}
func getPongCount() int {
	getsSqlStatement := `
		SELECT COUNT(*) FROM pongs
	`
	db, err := openDb()
	defer db.Close()

	pongs := 0
	err = db.QueryRow(getsSqlStatement).Scan(&pongs)
	if err != nil {
		panic(err)
	}
	return pongs
}

func insertPongCount(newValue int) int {
	getsSqlStatement := `
		INSERT INTO pongs (counter) VALUES (1)
	`
	db, err := openDb()
	defer db.Close()
	fmt.Println("value", newValue)
	fmt.Println(getsSqlStatement)
	_, err = db.Exec(getsSqlStatement)
	if err != nil {
		fmt.Println("error: ", &err)
		panic(err)
	}
	pongs := getPongCount()
	return pongs
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Serving pingpong request.")
	pongs := getPongCount()
	fmt.Fprintf(w, "<p>%s</p>", "pong "+strconv.Itoa(pongs))
	insertPongCount(pongs + 1)
}

func pongsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Serving pongs request.")

	_, err := fmt.Fprintln(w, strconv.Itoa(getPongCount()))
	if err != nil {
		panic(err)
	}
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	db, err := openDb()
	defer db.Close()
	if err != nil {
		http.Error(w, "shit's down", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}
	fmt.Fprintln(w, "Healthy")
}

func main() {
	db, _ := openDb()
	defer db.Close()
	statement := `CREATE TABLE IF NOT EXISTS pongs (counter integer DEFAULT 0)`
	db.Exec(statement)
	http.HandleFunc("/health", healthcheck)
	http.HandleFunc("/pingpong", handler)
	http.HandleFunc("/pongs", pongsHandler)
	http.HandleFunc("/", handler)

	port := "8080"
	fmt.Println("Server started in port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
