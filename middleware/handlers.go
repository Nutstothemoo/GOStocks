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

func createConnection() *sql.DB{

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

func GetStock(w http.ResponseWriter, r *http.Request) {
	params:= mux.Vars(r)
	// fmt.Println(params)
	id, err:= strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int.  %v", err)
	}
	stock, err := getStock(int64(id))
	if err != nil {
		log.Fatalf("Unable to get stock. %v", err)
	}
	json.NewEncoder(w).Encode(stock)
}
func GetAllStock(w http.ResponseWriter, r *http.Request)  {
	stocks, err := getAllStock()
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

	insertID, err := insertStock(stock)
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
func insertStock(stock models.Stock) (int64, error) {
	db:= createConnection()
	defer db.Close()
	sqlStatement := `INSERT INTO stocks (name, price, company) VALUES ($1, $2, $3) RETURNING stockid`
	var id int64
	err:= db.QueryRow(sqlStatement, stock.Name, stock.Price, stock.Company).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	fmt.Println("New record ID is:", id)
	return id, nil	
}
func getAllStock() ([]models.Stock, error) {
	db:= createConnection()
	defer db.Close()
	sqlStatement := `SELECT * FROM stocks`
	rows, err:= db.Query(sqlStatement)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var stocks []models.Stock
	for rows.Next() {
		var stock models.Stock
		err:= rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
		if err != nil {
			return nil, err
		}
		stocks = append(stocks, stock)
	}
	return stocks, err
}
func getStock(id int64) (models.Stock, error) {
	db:= createConnection()
	defer db.Close()
	sqlStatement := `SELECT * FROM stocks WHERE stockid=$1`
	var stock models.Stock
	row:= db.QueryRow(sqlStatement, id)
	err:= row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return stock, nil
	case nil:
		return stock, nil
	default:
		return stock, err
	}	
}
func updateStock(id int64, stock models.Stock) int64 {
	db:= createConnection()
	defer db.Close()
	sqlStatement := `UPDATE stocks SET name=$2, price=$3, company=$4 WHERE stockid=$1`
	res, err:= db.Exec(sqlStatement, id, stock.Name, stock.Price, stock.Company)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	rowsAffected, err:= res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}
	fmt.Printf("Total rows/record affected %v", rowsAffected)
	return rowsAffected
}
func deleteStock(id int64) int64 {
	db:= createConnection()
	defer db.Close()
	sqlStatement := `DELETE FROM stocks WHERE stockid=$1`
	res, err:= db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}
	rowsAffected, err:= res.RowsAffected()
	if err != nil {
		log.Fatalf("Error while checking the affected rows. %v", err)
	}
	fmt.Printf("Total rows/record affected %v", rowsAffected)
	return rowsAffected
}