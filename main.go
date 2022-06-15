package main

import (
	"flag"
	"net/http"

	"github.com/HeRaNO/cdoj-vjudge/config"
	"github.com/HeRaNO/cdoj-vjudge/modules"
	"github.com/HeRaNO/cdoj-vjudge/router"
	"github.com/kataras/iris/v12"
)

func main() {
	initConfigFile := flag.String("c", "./conf/config.yaml", "the path of configure file")

	config.Init(initConfigFile)
	app := router.InitApp()

	go modules.ListenMQ()
	app.Run(iris.Server(&http.Server{Addr: config.SrvAddr}))
}
