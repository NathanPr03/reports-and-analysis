package api

import "net/http"

func Reports(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Reports"))
}

func init() {
	http.HandleFunc("/reports", Reports)
}
