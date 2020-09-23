package actions

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rod6214/testrepo/controller/items"
	"github.com/rod6214/testrepo/controller/utils"
)

type Controller struct {
	Response utils.Response
}

func (controller Controller) CreateItem(responseWriter http.ResponseWriter, request *http.Request) {
	itemsClient := controller.Response.ItemsClient
	relational := request.URL.Query().Get("relational") == "true"
	var createItem utils.CreateItemRequest
	decoder := json.NewDecoder(request.Body)
	if error := decoder.Decode(&createItem); error != nil {
		controller.Response.JSON(responseWriter, http.StatusInternalServerError, error)
		return
	}
	defer request.Body.Close()
	createItemResponse, error := itemsClient.CreateItem(context.Background(), &items.CreateItemRequest{Id: createItem.ID, Description: createItem.Description, Relational: relational})
	if error != nil {
		controller.Response.JSON(responseWriter, http.StatusInternalServerError, error)
		return
	}
	controller.Response.JSON(responseWriter, http.StatusOK, createItemResponse)
}

func (controller Controller) ReadItem(responseWriter http.ResponseWriter, request *http.Request) {
	itemsClient := controller.Response.ItemsClient
	relational := request.URL.Query().Get("relational") == "true"
	id := mux.Vars(request)["id"]
	readItemResponse, error := itemsClient.ReadItem(context.Background(), &items.ReadItemRequest{Id: id, Relational: relational})
	if error != nil {
		controller.Response.JSON(responseWriter, http.StatusInternalServerError, error)
		return
	}
	if readItemResponse.Error {
		controller.Response.Error(responseWriter, http.StatusBadRequest, readItemResponse.Message)
		return
	}
	controller.Response.JSON(responseWriter, http.StatusOK, readItemResponse)
}

func (controller Controller) UpdateItem(responseWriter http.ResponseWriter, request *http.Request) {
	itemsClient := controller.Response.ItemsClient
	relational := request.URL.Query().Get("relational") == "true"
	var updateItem utils.UpdateItemRequest
	decoder := json.NewDecoder(request.Body)
	if error := decoder.Decode(&updateItem); error != nil {
		controller.Response.JSON(responseWriter, http.StatusInternalServerError, error)
		return
	}
	defer request.Body.Close()
	updateItemResponse, error := itemsClient.UpdateItem(context.Background(), &items.UpdateItemRequest{Id: updateItem.ID, Description: updateItem.Description, Relational: relational})
	if error != nil {
		controller.Response.JSON(responseWriter, http.StatusInternalServerError, error)
		return
	}
	if updateItemResponse.Error {
		controller.Response.Error(responseWriter, http.StatusBadRequest, updateItemResponse.Message)
		return
	}
	controller.Response.JSON(responseWriter, http.StatusOK, updateItemResponse)
}

func (controller Controller) DeleteItem(responseWriter http.ResponseWriter, request *http.Request) {
	itemsClient := controller.Response.ItemsClient
	relational := request.URL.Query().Get("relational") == "true"
	id := mux.Vars(request)["id"]
	deleteItemResponse, error := itemsClient.DeleteItem(context.Background(), &items.DeleteItemRequest{Id: id, Relational: relational})
	if error != nil {
		controller.Response.JSON(responseWriter, http.StatusInternalServerError, error)
		return
	}
	if deleteItemResponse.Error {
		controller.Response.Error(responseWriter, http.StatusBadRequest, deleteItemResponse.Message)
		return
	}
	controller.Response.JSON(responseWriter, http.StatusOK, deleteItemResponse)
}

func (controller Controller) GetIds(responseWriter http.ResponseWriter, request *http.Request) {
	itemsClient := controller.Response.ItemsClient
	relational := request.URL.Query().Get("relational") == "true"
	getIdsResponse, error := itemsClient.ListIds(context.Background(), &items.ListIdsRequest{Relational: relational})
	if error != nil {
		controller.Response.JSON(responseWriter, http.StatusInternalServerError, error)
		return
	}
	controller.Response.JSON(responseWriter, http.StatusOK, getIdsResponse)
}
