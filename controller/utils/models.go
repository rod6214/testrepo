package utils

type CreateItemRequest struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type UpdateItemRequest struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}
