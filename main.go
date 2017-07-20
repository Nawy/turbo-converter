package main

import (
	"fmt"

	config "turbo-converter/config"
	controller "turbo-converter/controller"
	service "turbo-converter/service"

	server "github.com/Nawy/go-rest-server"
	model "github.com/Nawy/go-rest-server/model"
)

func main() {
	config.InitConfig()
	service.InitHashGen()
	defer service.InitLogging().Close()

	fmt.Println("image-converter started on port ", config.GlobalConfig().Port)
	server.CreateRestServer(config.GlobalConfig().Port).SetHandlers(
		model.MakeRequestHandler(model.POST, "/upload").SetHandler(controller.UploadImageHandler),
		model.MakeRequestHandler(model.DELETE, "/delete").SetHandler(controller.DeleteImageHandler),
		model.MakeRequestHandler(model.GET, "/status").SetHandler(controller.StatusHandler),
	).Start()
}
