package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/DenisGubenko/ideasymbols/db"
	"github.com/DenisGubenko/ideasymbols/generator"
	"github.com/DenisGubenko/ideasymbols/http"
	"github.com/DenisGubenko/ideasymbols/utils"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func main() {
	storage, err := setupStorage()
	if err != nil {
		logrus.Errorf(`%+v`, errors.WithStack(err))
		os.Exit(2)
	}

	gen := generator.NewBasicGenerator(storage)

	go func() {
		err = gen.Start()
		if err != nil {
			logrus.Errorf(`%+v`, errors.WithStack(err))
			os.Exit(2)
		}
	}()

	router := http.NewRouterServer(storage)

	err = startServer(router)
	if err != nil {
		logrus.Errorf(`%+v`, errors.WithStack(err))
		os.Exit(2)
	}

	setupShutdown(gen, router)
}

func setupStorage() (db.Storage, error) {
	storage, err := db.NewPostgresStorage(
		utils.GetEnvVariable("POSTGRES_USER"),
		utils.GetEnvVariable("POSTGRES_PASSWORD"),
		utils.GetEnvVariable("POSTGRES_DB"),
		utils.GetEnvVariable("POSTGRES_HOST"),
		utils.GetIntEnvVariable("POSTGRES_PORT"),
		utils.GetEnvVariable("POSTGRES_SSLMODE"))
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return storage, nil
}

func startServer(router http.Server) error {
	logrus.Info("Start http server")
	if err := router.Start(); err != nil {
		logrus.Info("Http server started error occurred")
		return errors.WithStack(err)
	}

	logrus.Info("Http server ended successfully")
	return nil
}

func setupShutdown(gen generator.Generator, server http.Server) {
	gracefulStop := make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		if sig != nil {
			gen.Stop()
			err := server.Stop()
			if err != nil {
				logrus.Errorf(`%+v`, errors.WithStack(err))
			}
		}
	}()
}
