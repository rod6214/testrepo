package actions

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rod6214/testrepo/controller/items"
	"github.com/rod6214/testrepo/controller/utils"
	// "github.com/rod6214/testrepo/controller/items"
	// "github.com/southworks/gnalog/demo/controller/items"
)

// -----------

type createItemRequest struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type updateItemRequest struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

func TestMethod() {}

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

// ---------------
type Controller struct {
	ItemsClient items.ItemServiceClient
	Response    utils.Response
}

// func helloWorld(responseWriter http.ResponseWriter, request *http.Request) {
// 	var response = utils.Response{items.ItemServiceClient}
// 	respondMessage(responseWriter, http.StatusOK, "Hello world with CircleCI integrated with Pulumi!")
// }

// func testGRPC(responseWriter http.ResponseWriter, request *http.Request) {
// 	testGPRCResponse, error := itemsClient.TestGRPC(context.Background(), &items.TestGRPCRequest{})
// 	if error != nil {
// 		respondJSON(responseWriter, http.StatusInternalServerError, error)
// 		return
// 	}
// 	respondJSON(responseWriter, http.StatusOK, testGPRCResponse)
// }

// func createItem(responseWriter http.ResponseWriter, request *http.Request) {
// 	relational := request.URL.Query().Get("relational") == "true"
// 	var createItem createItemRequest
// 	decoder := json.NewDecoder(request.Body)
// 	if error := decoder.Decode(&createItem); error != nil {
// 		respondJSON(responseWriter, http.StatusInternalServerError, error)
// 		return
// 	}
// 	defer request.Body.Close()
// 	createItemResponse, error := itemsClient.CreateItem(context.Background(), &items.CreateItemRequest{Id: createItem.ID, Description: createItem.Description, Relational: relational})
// 	if error != nil {
// 		respondJSON(responseWriter, http.StatusInternalServerError, error)
// 		return
// 	}
// 	respondJSON(responseWriter, http.StatusOK, createItemResponse)
// }

func (controller Controller) readItem(responseWriter http.ResponseWriter, request *http.Request) {
	// res := utils.Response{ItemsClient: controller.ItemsClient}
	relational := request.URL.Query().Get("relational") == "true"
	id := mux.Vars(request)["id"]
	readItemResponse, error := controller.ItemsClient.ReadItem(context.Background(), &items.ReadItemRequest{Id: id, Relational: relational})
	if error != nil {
		Controller.Response.JSON(responseWriter, http.StatusInternalServerError, error)
		return
	}
	if readItemResponse.Error {
		Controller.Response.Error(responseWriter, http.StatusBadRequest, readItemResponse.Message)
		return
	}
	res.JSON(responseWriter, http.StatusOK, readItemResponse)
}

// func updateItem(responseWriter http.ResponseWriter, request *http.Request) {
// 	relational := request.URL.Query().Get("relational") == "true"
// 	var updateItem updateItemRequest
// 	decoder := json.NewDecoder(request.Body)
// 	if error := decoder.Decode(&updateItem); error != nil {
// 		respondJSON(responseWriter, http.StatusInternalServerError, error)
// 		return
// 	}
// 	defer request.Body.Close()
// 	updateItemResponse, error := itemsClient.UpdateItem(context.Background(), &items.UpdateItemRequest{Id: updateItem.ID, Description: updateItem.Description, Relational: relational})
// 	if error != nil {
// 		respondJSON(responseWriter, http.StatusInternalServerError, error)
// 		return
// 	}
// 	if updateItemResponse.Error {
// 		respondError(responseWriter, http.StatusBadRequest, updateItemResponse.Message)
// 		return
// 	}
// 	respondJSON(responseWriter, http.StatusOK, updateItemResponse)
// }

// func getAllItems(responseWriter http.ResponseWriter, request *http.Request) {
// 	respondJSON(responseWriter, http.StatusOK, items)
// }

// func deleteItem(responseWriter http.ResponseWriter, request *http.Request) {
// 	relational := request.URL.Query().Get("relational") == "true"
// 	id := mux.Vars(request)["id"]
// 	deleteItemResponse, error := itemsClient.DeleteItem(context.Background(), &items.DeleteItemRequest{Id: id, Relational: relational})
// 	if error != nil {
// 		respondJSON(responseWriter, http.StatusInternalServerError, error)
// 		return
// 	}
// 	if deleteItemResponse.Error {
// 		respondError(responseWriter, http.StatusBadRequest, deleteItemResponse.Message)
// 		return
// 	}
// 	respondJSON(responseWriter, http.StatusOK, deleteItemResponse)
// }
// func (controller Controller) New(itemService items.ItemServiceClient) *Controller {
// 	controller.ItemsClient = itemService
// 	return controller
// }

func (controller Controller) getIds(responseWriter http.ResponseWriter, request *http.Request) {
	res := utils.Response{ItemsClient: controller.ItemsClient}

	relational := request.URL.Query().Get("relational") == "true"
	// getIdsResponse, error := itemsClient.ListIds(context.Background(), &items.ListIdsRequest{Relational: relational})
	getIdsResponse, error := controller.ItemsClient.ListIds(context.Background(), &items.ListIdsRequest{Relational: relational})
	if error != nil {
		res.JSON(responseWriter, http.StatusInternalServerError, error)
		return
	}
	res.JSON(responseWriter, http.StatusOK, getIdsResponse)
}
