package middleware

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"github.com/joho/godotenv"
)

type response struct {
	ID int64 `json:"id,omitempty"`// omitempty means that if the field is empty, it will not be included in the JSON response
	Message string `json:"message,omitempty"`
}

func CreateConnection() *sql.DB{

	err := godotenv.Load(".env")
	if err != nil {
		panic("failed to load .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}

	err =db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully connected!")
	return db
}
func CloseConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println("Successfully closed connection!")
}

func GetStock(db *sql.DB, id int64) (Stock, error) {

}
func GetAllStock(db *sql.DB) ([]Stock, error) {

}
func CreateStock(w http.ResponseWriter, r *http.Request) (int64, error) {
	
}
func DeleteStock(db *sql.DB, id int64) (int64, error) {

}