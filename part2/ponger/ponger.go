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

func openDb() *sql.DB {
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected to postgres!")

	return db
}
func getPongCount() int {
	getsSqlStatement := `
		SELECT COUNT(*) FROM pongs
	`
	db := openDb()
	defer db.Close()

	pongs := 0
	err := db.QueryRow(getsSqlStatement).Scan(&pongs)
	if err != nil {
		panic(err)
	}
	return pongs
}

func insertPongCount(newValue int) int {
	getsSqlStatement := `
		INSERT INTO pongs (counter) VALUES (1)
	`
	db := openDb()
	defer db.Close()
	fmt.Println("value", newValue)
	fmt.Println(getsSqlStatement)
	_, err := db.Exec(getsSqlStatement)
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

func main() {
	db := openDb()
	defer db.Close()
	statement := `CREATE TABLE IF NOT EXISTS pongs (counter integer DEFAULT 0)`
	db.Exec(statement)

	http.HandleFunc("/pingpong", handler)
	http.HandleFunc("/pongs", pongsHandler)

	port := "3000"
	fmt.Println("Server started in port", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
