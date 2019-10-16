package model

// UserData : for storing user data
type UserData struct {
	Org         string `json:"org"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	OldPassword string `json:"oldPassword"`
	Name        string `json:"name"`
	Mobile      string `json:"mobile"`
	Role        string `json:"role"`
}
