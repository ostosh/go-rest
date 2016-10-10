package request

import (
	"log"
	"net/http"

	"../util"
)

type DelegateRequest func(w http.ResponseWriter, req *http.Request)


func HandleRequest(h DelegateRequest) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("error: processing request %v", err)
				writer := util.NewJsonWriter(w)
				writer.RootObject(func() {
					writer.KeyValue("status", "error")
				})
			}
		}()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		conn := GetConnection()
		h(w, req)
	}
}

func HandleNotImplemented() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotImplemented)
	}
}
