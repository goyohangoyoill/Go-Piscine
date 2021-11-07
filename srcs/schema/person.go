package schema

type Person struct {
	ID        int          `json:"id"`
	Name      string       `json:"name"`
	Course    int          `json:"course"`
	Point     uint         `json:"point"`
	Password  string       `json:"password"`
	Score     SortableEvalRes `json:"score"`
}
