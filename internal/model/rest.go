package model

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreatedUser struct {
	Id       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Age      int    `json:"age"`
}

type Result struct {
	Result any    `json:"result,omitempty"`
	Error  string `json:"error,omitempty"`
}

type Token struct {
	Token string `json:"token"`
}
