init:
	if ! hash dep 2>/dev/null; then go get -u github.com/golang/dep/cmd/dep; fi
	if ! hash sql-migrate 2>/dev/null; then go get -v github.com/rubenv/sql-migrate/...; fi
	if ! hash mockgen 2>/dev/null; then go get github.com/golang/mock/mockgen; fi
	if ! hash swagger 2>/dev/null; then go get -u github.com/go-swagger/go-swagger/cmd/swagger; ficlear
	dep ensure

run: db-run db-migrate
	dep ensure
	go run main.go -c=config/config.yaml

test: db-run db-migrate
	dep ensure
	go test -v ./test/
	$(MAKE) db-stop

bench: db-run db-migrate
	go test -bench=. -benchmem ./test/
	$(MAKE) db-stop

spec:
	swagger generate spec -m -o ./spec/api.json

spec-ui: spec
	swagger serve -F=swagger ./spec/api.json

db-run: db-stop
	 docker run --rm -d -p 5433:5432 --name echo-rest-api-db postgres
	 sleep 2

db-stop:
	docker container stop echo-rest-api-db >/dev/null 2>&1 || exit 0

db-migrate:
	sql-migrate up -config=store/store.yaml

.PHONY: run test db-run db-stop db-migrate spec spec-ui