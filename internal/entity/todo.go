package entity

type Todo struct {
	ID     uint   `json:"id"`
	UserID uint   `json:"user_id"`
	Title  string `json:"title"`
	Done   bool   `json:"done"`
}
func (Todo) TableName() string {
	return "public.todos"
}