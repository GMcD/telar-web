package models

type Signup2TokenModel struct {
	User         UserSignupTokenModel `json:"user"`
}

type UserSignupTokenModel struct {
	Fullname string `json:"fullName"`
	Email    string `json:"email" `
	Password string `json:"password" `
}
