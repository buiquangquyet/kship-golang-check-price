package main

import (
	"check-price/src/bootstrap"
	"check-price/src/common"
	"check-price/src/common/configs"
	"check-price/src/common/log"
	"check-price/src/core/constant"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	defaultGracefulTimeout = 15 * time.Second
)

func init() {
	var pathConfig string
	flag.StringVar(&pathConfig, "config", "configs/config.yaml", "path to config file")
	var activeInstanceArg string
	flag.StringVar(&activeInstanceArg, "env", "public", "public | private")
	flag.Parse()
	err := configs.LoadConfig(pathConfig)
	if err != nil {
		panic(err)
	}
	if !constant.IsProdEnv() {
		cf, _ := json.Marshal(configs.Get())
		fmt.Printf("configs:[%s]\n", cf)
	}
	log.NewLogger()
}

func main() {
	logger := log.GetLogger().GetZap()
	logger.Debugf("App %s is running", configs.Get().Mode)
	app := fx.New(
		fx.Provide(log.GetLogger().GetZap),
		fx.Invoke(common.InitTracer),

		// storage module
		bootstrap.BuildStorageModules(),

		//build ext service
		bootstrap.BuildExtServicesModules(),

		//build http server
		bootstrap.BuildServiceModule(),

		//http servuc
		bootstrap.BuildControllerModule(),
		bootstrap.BuildValidator(),
		bootstrap.BuildHTTPServerModule(),
	)

	startCtx, cancel := context.WithTimeout(context.Background(), defaultGracefulTimeout)
	defer cancel()
	if err := app.Start(startCtx); err != nil {
		logger.Fatalf(err.Error())
	}

	interruptHandle(app, logger)
}

func interruptHandle(app *fx.App, logger *zap.SugaredLogger) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	logger.Debugf("Listening Signal...")
	s := <-c
	logger.Infof("Received signal: %s. Shutting down Server ...", s)

	stopCtx, cancel := context.WithTimeout(context.Background(), defaultGracefulTimeout)
	defer cancel()

	if err := app.Stop(stopCtx); err != nil {
		logger.Fatalf(err.Error())
	}
}
