package views

type AuthRegisterWithEmail struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthLoginWithEmail struct {
	Username string `json:"username"`
	Password string `json:"password"`
}