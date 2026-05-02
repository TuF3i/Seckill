package user

type RegisterUserReq struct {
	Uid      string `json:"uid"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginReq struct {
	Uid      string `json:"uid"`
	Password string `json:"password"`
}

type LoginResp struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type RefreshAccessTokenResp struct {
	AccessToken string `json:"accessToken"`
}
