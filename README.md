# echo-rest-api

## Установка зависимостей
dep ensure

## Генерация моков
mockgen -source=store.go -destination=../test/store_mock.go -package=test
mockgen -source=service/category_service.go -destination=test/category_service_mock.go -package=test
mockgen -source=service/product_service.go -destination=test/product_service_mock.go -package=test
