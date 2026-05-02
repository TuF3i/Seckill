package user

type RegisterUserReq struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResp struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshAccessTokenResp struct {
	AccessToken string `json:"accessToken"`
}
