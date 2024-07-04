package model

type Login struct {
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Long     float64 `json:"long"`
	Lat      float64 `json:"lat"`
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

type Swipe struct {
	UserID     int  `json:"userID"`
	SwipeRight bool `json:"swipeRight"`
}

type Match struct {
	Matched bool `json:"matched"`
	MatchID int  `json:"matchID,omitempty"`
}

type UserPassword struct {
	ID       int    `db:"id"`
	Password string `db:"password"`
}
