package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/space-fold-technologies/aurora-service/app"
)

func main() {
	mode := flag.String("mode", "run", " for the app run mode")
	email := flag.String("email", "", "for admin email address")
	password := flag.String("password", "", "for admin password")
	flag.Parse()
	initializer := app.NewInitializer()
	if *mode == "run" || *mode == "RUN" {
		initializer.Migrate()
		application := app.Application{}
		application.Start()
		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		application.Stop()
	} else if *mode == "setup" || *mode == "SETUP" {
		initializer.Register(*email, *password)
	} else if *mode == "reset" || *mode == "RESET" {
		initializer.Reset()
	}
}
