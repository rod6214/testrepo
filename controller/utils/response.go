package utils

import (
	"encoding/json"
	"net/http"

	"github.com/rod6214/testrepo/controller/items"
)

type Response struct {
	itemsClient *items.ItemServiceClient
}

// type createItemRequest struct {
// 	ID          string `json:"id"`
// 	Description string `json:"description"`
// }

// type updateItemRequest struct {
// 	ID          string `json:"id"`
// 	Description string `json:"description"`
// }

func (res *Response) New(itemsClient *items.ItemServiceClient) *Response {
	res.itemsClient = itemsClient
	return res
}

func (res *Response) JSON(responseWriter http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		responseWriter.Write([]byte(err.Error()))
		return
	}
	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.WriteHeader(status)
	responseWriter.Write([]byte(response))
}

func (res *Response) Message(responseWriter http.ResponseWriter, code int, message string) {
	res.JSON(responseWriter, code, map[string]string{"message": message})
}

func (res *Response) Error(responseWriter http.ResponseWriter, code int, message string) {
	res.JSON(responseWriter, code, map[string]string{"error": message})
}

// var itemsClient items.ItemServiceClient
