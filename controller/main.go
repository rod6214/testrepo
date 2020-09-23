package main

import (
	"log"

	"github.com/rod6214/testrepo/controller/utils"

	"github.com/rod6214/testrepo/controller/actions"

	"github.com/gorilla/mux"
	"github.com/rod6214/testrepo/controller/items"
	"google.golang.org/grpc"
	// "github.com/southworks/gnalog/demo/controller/items"
)

// -----
// type createItemRequest struct {
// 	ID          string `json:"id"`
// 	Description string `json:"description"`
// }

// type updateItemRequest struct {
// 	ID          string `json:"id"`
// 	Description string `json:"description"`
// }

// func respondJSON(responseWriter http.ResponseWriter, status int, payload interface{}) {
// 	response, err := json.Marshal(payload)
// 	if err != nil {
// 		responseWriter.WriteHeader(http.StatusInternalServerError)
// 		responseWriter.Write([]byte(err.Error()))
// 		return
// 	}
// 	responseWriter.Header().Set("Content-Type", "application/json")
// 	responseWriter.WriteHeader(status)
// 	responseWriter.Write([]byte(response))
// }

// func respondMessage(responseWriter http.ResponseWriter, code int, message string) {
// 	respondJSON(responseWriter, code, map[string]string{"message": message})
// }

// func respondError(responseWriter http.ResponseWriter, code int, message string) {
// 	respondJSON(responseWriter, code, map[string]string{"error": message})
// }

// var itemsClient items.ItemServiceClient

// -----

func main() {
	// actions.TestMethod()
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
	// router.HandleFunc("/", helloWorld)
	// router.HandleFunc("/grpc", testGRPC)
	itemsRouter := router.PathPrefix("/items").Subrouter()
	itemsRouter.HandleFunc("/", controllers.CreateItem).Methods("POST")
	itemsRouter.HandleFunc("/{id}", controllers.ReadItem).Methods("GET")
	itemsRouter.HandleFunc("/", updateItem).Methods("PUT")
	itemsRouter.HandleFunc("/", getAllItems).Methods("GET")
	itemsRouter.HandleFunc("/{id}", deleteItem).Methods("DELETE")
	// itemsRouter.HandleFunc("/", getIds).Methods("GET")
	// log.Println("Listening in port 80")
	// log.Fatal(http.ListenAndServe(":8080", router))
}
