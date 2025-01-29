package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vladslugin987/marketplace-go/api-service/db"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	StockCount  int     `json:"stock_count"`
}

func main() {
	// initialise db connection
	db.InitDB()
	defer db.DB.Close()

	r := mux.NewRouter()

	// API routes
	r.HandleFunc("/api/health", healthCheck).Methods("GET")
	r.HandleFunc("/api/products", getProducts).Methods("GET")

	// static files for web interface
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	// check db connection
	err := db.DB.Ping()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{
			"status": "database connection error",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": "healthy",
	})
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	// get products from db
	rows, err := db.DB.Query("SELECT id, name, description, price, stock_count FROM products")
	if err != nil {
		log.Printf("Error querying products: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.StockCount)
		if err != nil {
			log.Printf("Error scanning product: %v", err)
			continue
		}
		products = append(products, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
