package handler

import (
	"encoding/json"
	"log"
	"net/http"
	conv "skillbox_final/internal/json"
)

func New() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		structure, err := conv.GetResultData()
		if err != nil {
			answer := conv.ResultTErr{
				Status: false,
				Error:  "Service is unavailable",
			}
			result, err := json.Marshal(answer)
			if err != nil {
				log.Printf("Error marshalling error response: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			_, err = w.Write(result)
			if err != nil {
				log.Printf("Error writing error response: %v", err)
			}
			return
		}

		answer := conv.ResultT{
			Status: true,
			Data:   structure,
		}

		result, err := json.Marshal(answer)
		if err != nil {
			log.Printf("Error marshalling success response: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(result)
		if err != nil {
			log.Printf("Error writing success response: %v", err)
		}
	}
}
