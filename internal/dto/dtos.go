package dto

type TodoItem struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}
