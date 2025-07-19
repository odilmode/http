package main
import (
	"net/http"
	"fmt"
)


func handleReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	bodyText := "OK"
	_, err := w.Write([]byte(bodyText))
	if err != nil {
		fmt.Printf("Error writing response: %v\n", err)
	}
}

