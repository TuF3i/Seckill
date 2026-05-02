package main

import (
	"log"
	"seckill/internal/userSvr/core/handler"
	usersvr "seckill/internal/userSvr/kitex_gen/usersvr/usersvr"
)

func main() {
	svr := usersvr.NewServer(new(handler.UserSvrImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
