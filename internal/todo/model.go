package todo

type TodoItemRequest struct {
	Title string `json:"title" binding:"required"`
	Done  bool   `json:"done"`
}

type TodoItem struct {
	Id    int    `json:"id"`
	Title string `json:"title" binding:"required"`
	Done  bool   `json:"done"`
}
