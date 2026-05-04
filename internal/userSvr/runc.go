package main

import (
	userSvr "seckill/internal/userSvr/core/app"
)

func main() {
	userSvr.OnCreate()
	defer userSvr.OnDestory()
	userSvr.RunServer()
}
