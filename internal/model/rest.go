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

type Discover struct {
	Id             int    `json:"id," db:"id"`
	Name           string `json:"name" db:"name"`
	Gender         string `json:"gender" db:"gender"`
	Age            int    `json:"age" db:"age"`
	DistanceFromMe string `json:"distanceFromMe,omitempty" db:"distance"`
}
