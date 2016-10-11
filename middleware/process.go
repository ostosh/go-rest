package request

import (
	"log"
	"net/http"

	"../util"
)

type DelegateRequest func(w http.ResponseWriter, req *http.Request)

//Returns abstract request handler for model delegation
func HandleRequest(h DelegateRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("error: processing request %v", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				 
			}
		}()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		conn := GetConnection()
		h(w, req)
	}
}

//Returns unimplemented request handler
func HandleNotImplemented() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf("error: processing unimplemented request %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotImplemented)
	}
}
