package main

import (
	"echo-rest-api/api"
	"echo-rest-api/config"
	"echo-rest-api/service"
	"echo-rest-api/store"
	"flag"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

func main() {
	var err error
	log.SetFormatter(&log.JSONFormatter{})
	// Получаем флаги запуска приложения
	configFile := flag.String("c", "config.yaml", "Path to config file")
	flag.Parse()
	// Загружаем конфиг
	var conf *config.Config
	if conf, err = config.NewConfig(*configFile); err != nil {
		log.Fatal(err)
	}
	// Конфигурируем логгер
	log.SetLevel(log.Level(conf.LogLevel))
	log.Info("Starting service with configuration: ", conf.ConfigFile)
	// Создаем сторедж
	store, err := store.NewStore(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	log.Info("Store created successfully")
	// Создаем сервисы
	cs := service.NewCategoryService(store)
	ps := service.NewProductService(store)
	log.Info("Services created successfully")
	// Создаем  Api
	api := api.NewApi(conf, cs, ps)
	log.WithField("address", api.GetApiInfo().Address).
		WithField("mw", api.GetApiInfo().MW).
		WithField("routs", api.GetApiInfo().Routs).
		Info("Starting api")
	log.Fatal(api.Start())
}
