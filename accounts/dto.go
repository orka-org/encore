package accounts

type RegisterParams struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegisterResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginParams struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshParams struct {
	RefreshToken string `header:"X-Refresh-Token"`
}
type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ValidateParams struct {
	AccessToken string `json:"access_token"`
}
type ValidateResponse struct {
	Subject  string `json:"subject"`
	Expires  int64  `json:"expires"`
	IssuedAt int64  `json:"issued_at"`
	Issuer   string `json:"issuer"`

	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
