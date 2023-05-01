package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/BillyBones007/pwdm_server/internal/app/servergrpc"
)

func main() {
	app := servergrpc.NewServer()
	closed := make(chan struct{})

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

		<-stop
		app.Shutdown()
		close(closed)
	}()

	app.StartServer()

	<-closed
	log.Println("Server is closed")
}
