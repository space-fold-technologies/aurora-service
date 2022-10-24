package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/space-fold-technologies/aurora-service/app"
	"github.com/space-fold-technologies/aurora-service/app/core/database"
	"github.com/space-fold-technologies/aurora-service/app/core/logging"
)

func main() {
	mode := flag.String("mode", "run", " for the app run mode")
	email := flag.String("email", "", "for admin email address")
	password := flag.String("password", "", "for admin password")
	flag.Parse()

	if *mode == "run" || *mode == "RUN" {
		if err := database.NewMigrationHandler().Migrate(""); err != nil {
			logging.GetInstance().Error(err)
		}
		application := app.Application{}
		application.Start()
		c := make(chan os.Signal, 2)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		application.Stop()
	} else if *mode == "setup" || *mode == "SETUP" {
		initializer := app.Initializer{}
		initializer.Register(*email, *password)
	} else if *mode == "reset" || *mode == "RESET" {
		initializer := app.Initializer{}
		initializer.Reset()
	}
}
