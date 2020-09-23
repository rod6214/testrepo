package main

import (
	"log"
	"net/http"

	"github.com/rod6214/testrepo/controller/utils"

	"github.com/rod6214/testrepo/controller/actions"

	"github.com/gorilla/mux"
	"github.com/rod6214/testrepo/controller/items"
	"google.golang.org/grpc"
)

func main() {
	log.Println("Starting Controller...")
	var clientConnection *grpc.ClientConn
	clientConnection, error := grpc.Dial(":9000", grpc.WithInsecure())
	if error != nil {
		log.Fatalf("Could not connect: %s", error)
	}
	defer clientConnection.Close()
	itemsClient := items.NewItemServiceClient(clientConnection)
	router := mux.NewRouter().StrictSlash(true)
	// Controllers
	controllers := actions.Controller{Response: utils.Response{ItemsClient: itemsClient}}
	// Bind routes
	itemsRouter := router.PathPrefix("/items").Subrouter()
	itemsRouter.HandleFunc("/", controllers.CreateItem).Methods("POST")
	itemsRouter.HandleFunc("/{id}", controllers.ReadItem).Methods("GET")
	itemsRouter.HandleFunc("/", controllers.UpdateItem).Methods("PUT")
	itemsRouter.HandleFunc("/{id}", controllers.DeleteItem).Methods("DELETE")
	itemsRouter.HandleFunc("/", controllers.GetIds).Methods("GET")
	// Start server
	log.Println("Listening in port 80")
	log.Fatal(http.ListenAndServe(":8080", router))
}
