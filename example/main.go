package main

import (
	"log"
	"net/http"
  
	"../middleware"
	"../router"
)

func main() {

	db := pg.New("driver", "source")
	request.SetConnection(db)

	router := router.New()

	//brand
	brand := model.Brand{}
	router.Get("/api/brand/query", request.HandleRequest(brand.HandleQuery))
	router.Get("/api/brand/read/:id", request.HandleRequest(brand.HandleRead))
	router.Post("/api/brand/create", request.HandleRequest(brand.HandleCreate))
	router.Post("/api/brand/update", request.HandleRequest(brand.HandleUpdate))
	router.Post("/api/brand/delete", request.HandleNotImplemented())

	http.Handle("/", router)
	err := http.ListenAndServe(":3003", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
