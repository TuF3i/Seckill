package lkeygen

func GenAccessTokenKey(uid string) string {
	return "user:" + "token:" + "access:" + uid
}

func GenRefreshTokenKey(uid string) string {
	return "user:" + "token:" + "refresh:" + uid
}
