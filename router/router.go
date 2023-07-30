package router

import (
	"stockapi/middleware"

	"github.com/gorilla/mux"
)

func Router() *mux.Router{
	router:= mux.NewRouter()

	router.HandleFunc("/api/stock/{id}", middleware.GetStock).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/getAllStock", middleware.GetAllStock).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/newStock", middleware.CreateStock).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/deleteStock/{id}", middleware.DeleteStock).Methods("DELETE", "OPTIONS")
	router.HandleFunc("/api/updateStock/{id}", middleware.UpdateStock).Methods("PUT", "OPTIONS")
	return router
}