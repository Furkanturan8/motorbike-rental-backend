package viewmodel

type AuthLoginVM struct {
	Email    string `json:"email" validate:"required_without=Phone,omitempty,max=64,email"`
	Phone    string `json:"phone" validate:"required_without=Email,omitempty,max=11,numeric"`
	Password string `json:"password" validate:"required" label:"Parola"`
}

type AuthTokenVM struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthRefreshVM struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
