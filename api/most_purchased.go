package api

import (
	"encoding/json"
	"github.com/NathanPr03/price-control/pkg/db"
	"net/http"
)

type PurchasedProduct struct {
	ProductName     string `json:"productName"`
	AmountPurchased int    `json:"amountPurchased"`
}

func MostPurchased(w http.ResponseWriter, r *http.Request) {
	dbConnection, err := db.ConnectToDb()
	if err != nil {
		http.Error(w, "Error connecting to database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	query := "SELECT product_purchased, SUM(amount_purchased) FROM customer_purchases GROUP BY product_id, product_purchased"
	rows, err := dbConnection.Query(query)
	if err != nil {
		http.Error(w, "Error querying database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var purchasedProducts []PurchasedProduct
	for rows.Next() {
		var product PurchasedProduct

		if err := rows.Scan(&product.ProductName, &product.AmountPurchased); err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}

		purchasedProducts = append(purchasedProducts, product)
	}

	response := map[string][]PurchasedProduct{"products": purchasedProducts}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}

func init() {
	http.HandleFunc("/most_purchased", MostPurchased)
}
