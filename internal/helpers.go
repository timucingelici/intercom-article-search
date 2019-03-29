package internal

import (
	"encoding/json"
	"log"
	"net/http"
)

func sendResponse(data interface{}, w http.ResponseWriter) {
	o, err := json.Marshal(data)

	if err != nil {
		log.Println("Failed to marshall the object : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = w.Write(o)

	if err != nil {
		log.Println("Failed to send a response with response writer : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
