package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"stockapi/models"
	"strconv"

	"github.com/gorilla/mux"
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

func GetStock(w http.ResponseWriter, r *http.Request) {
	params:= mux.Vars(r)
	// fmt.Println(params)
	id, err:= strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}
	stock, err := GetStock(int64(id))
	if err != nil {
		log.Fatalf("Unable to get stock. %v", err)
	}
	json.NewEncoder(w).Encode(stock)
}
func GetAllStock(w http.ResponseWriter, r *http.Request)  {
	stocks, err := GetAllStock()
	if err != nil {
		log.Fatalf("Unable to get all stock. %v", err)
	}
	json.NewEncoder(w).Encode(stocks)

}
func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock

	err:= json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}

	insertID, err := insertStock(db, stock)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	res:= response{
		ID: insertID,
		Message: "Stock created successfully",
	}
	json.NewEncoder(w).Encode(res)	
}

func DeleteStock(w http.ResponseWriter, r *http.Request) {
	params:= mux.Vars(r)
	id, err:= strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}
	deletedRows:= deleteStock(int64(id))
	msg:= fmt.Sprintf("Stock deleted successfully. Total rows/record affected %v", deletedRows)
	res:= response{
		ID: int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}

func UpdateStock(w http.ResponseWriter, r *http.Request)  {
	params:= mux.Vars(r)
	fmt.Println("params", params)
	id, err:= strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}
	var stock models.Stock
	err = json.NewDecoder(r.Body).Decode(&stock)
	if err != nil {
		log.Fatalf("Unable to decode the request body.  %v", err)
	}
	updatedRows:= updateStock(int64(id), stock)
	msg:= fmt.Sprintf("Stock updated successfully. Total rows/record affected %v", updatedRows)
	res:= response{
		ID: int64(id),
		Message: msg,
	}
	json.NewEncoder(w).Encode(res)
}