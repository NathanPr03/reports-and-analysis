package api

import (
	"encoding/json"
	"github.com/NathanPr03/price-control/pkg/db"
	"net/http"
)

type SignUpStats struct {
	Month       string `json:"month"`
	SignUpCount int    `json:"sign_up_count"`
}

func GetSignUpsPerMonth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	dbConnection, err := db.ConnectToDb()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}
	defer dbConnection.Close()

	rows, err := dbConnection.Query(`
        SELECT
            TO_CHAR(sign_up_date, 'YYYY-MM') AS sign_up_month,
            COUNT(*) AS sign_up_count
        FROM
            customer
        GROUP BY
            sign_up_month
        ORDER BY
            sign_up_month
    `)
	if err != nil {
		http.Error(w, "Error executing query", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var signUpStats []SignUpStats

	for rows.Next() {
		var stats SignUpStats
		if err := rows.Scan(&stats.Month, &stats.SignUpCount); err != nil {
			http.Error(w, "Error scanning row", http.StatusInternalServerError)
			return
		}
		signUpStats = append(signUpStats, stats)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(signUpStats)
}

func init() {
	http.HandleFunc("/signups-per-month", GetSignUpsPerMonth)
}
