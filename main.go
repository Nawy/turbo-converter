package main

import (
	"fmt"

	controller "turbo-converter/controller"

	server "github.com/Nawy/go-rest-server"
	model "github.com/Nawy/go-rest-server/model"
)

func main() {
	conf = *loadAppConfig(getPathFromArgs(), &conf)
	initHashGen()
	defer initLogging().Close()

	fmt.Println("image-converter started on port ", conf.Port)
	server.CreateRestServer(conf.Port).SetHandlers(
		model.MakeRequestHandler(model.POST, "/upload").SetHandler(controller.uploadImageHandler),
		model.MakeRequestHandler(model.DELETE, "/delete").SetHandler(controller.deleteImageHandler),
		model.MakeRequestHandler(model.GET, "/status").SetHandler(controller.statusHandler),
	).Start()
}
