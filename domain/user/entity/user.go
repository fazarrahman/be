package entity

type User struct {
	Username  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Password  []byte `json:"password"`
	Status    int8   `json:"status"`
}
