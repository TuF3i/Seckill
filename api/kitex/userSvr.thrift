namespace go usersvr

enum UserRole {
    ADMIN = 1,
    SIMPLE_USER = 2
}

enum ClaimType {
    ACCESS = 1,
    REFRESH = 2
}

struct JWTToken {
    1: required string accessToken
    2: required string refreshToken
}

struct JWTClaims {
    1: required string UID
    2: required UserRole Role
    3: required ClaimType Type
}

service UserSvr {
    void RegisterUser(1: string uid, 2: string password) // 注册
    JWTToken Login(1: string uid, 2: string password) // 登录
    void Logout(1: string accessToken) // 推出登录
    string RefreshAccessToken(1: string refreshToken) // 刷新accessToken
    JWTClaims VerifyAccessToken(1: string accessToken) // 验证accessToken
    JWTClaims VerifyRefreshToken(1: string refreshToken) // 验证refreshToken
}