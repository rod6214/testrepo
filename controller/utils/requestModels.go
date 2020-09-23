package utils

type createItemRequest struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type updateItemRequest struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}
