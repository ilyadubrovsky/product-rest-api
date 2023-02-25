package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/ilyadubrovsky/product-rest-api/internal/config"
	"github.com/ilyadubrovsky/product-rest-api/internal/product"
	storage "github.com/ilyadubrovsky/product-rest-api/internal/storage/mongodb"
	"github.com/ilyadubrovsky/product-rest-api/pkg/client/mongodb"
	"github.com/ilyadubrovsky/product-rest-api/pkg/logging"
)

func main() {
	logger := logging.GetLogger()

	logger.Info("read application config")
	cfg, err := config.GetConfig()
	if err != nil {
		logger.Panic(err)
	}

	logger.Info("initializing mongodb client")
	database, err := mongodb.NewClient(context.TODO(), cfg.Host, cfg.Port, cfg.Database,
		cfg.Username, cfg.Password, cfg.AuthDB)

	if err != nil {
		logger.Panic(err)
	}

	logger.Info("initializing mongodb products collection")
	productsStorage := storage.NewProductsMongo(database, logger)

	logger.Info("initializing product service")
	service := product.NewService(logger, cfg, productsStorage)

	logger.Info("initializing router")
	router := gin.Default()

	logger.Info("initializing handler")
	handler := product.NewHandler(logger, cfg, service)

	logger.Info("register handler")
	handler.Register(router)

	logger.Info("router run")

	router.Run(":8080")
}
