namespace go usersvr

struct JWTToken {
    1: required string accessToken
    2: required string refreshToken
}

service UserSvr {
    void RegisterUser(1: string uid, 2: string password) // 注册
    JWTToken Login(1: string uid, 2: string password) // 登录
    void Logout(1: string accessToken) // 推出登录
    string RefreshAccessToken(1: string refreshToken) // 刷新accessToken
}