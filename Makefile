-include .env
export

GIN_MODE=test
BUILD_TIME?=$(shell TZ=${TZ} date '+%Y-%m-%d %H:%M:%S')

LDFLAGS=$(shell echo \
	"-X 'osoc-dialog/pkg/application.buildVersionTime=${BUILD_TIME}'" \
)

DEFAULT_GOAL := help
.PHONY: help
help:
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-27s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: init-utils
init-utils: ## Init utils for work space
	go install golang.org/x/tools/cmd/goimports@latest

.PHONY: test
test: ## Test project.
	go test -v ./...

.PHONY: test-short
test-short: ## Test project.
	go test -short -v ./...

.PHONY: test-race
test-race: ## Test project with race detection.
	go test -v -race ./...

.PHONY: test-cover
test-cover: ## Test with coverage and open result in a browser
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

.PHONY: fmt
fmt: ## Format golang files with goimports
	find . -name \*.go -not -path \*/wire_gen.go -exec goimports -w {} \;

.PHONY: wire
wire: ## Actualize dependency-injection
	wire ./...

.PHONY: mod
mod: ## Remove unused modules
	go mod tidy -v

.PHONY: finalcheck
finalcheck: wire fmt mod lint test-short swagger-gen ## Make a final complex check before the commit

.PHONY: run
run: ## Run project for local
	go run -race -ldflags="${LDFLAGS}" ./cmd/${APP_NAME}/.
	#go run -ldflags="${LDFLAGS}" ./cmd/${APP_NAME}/. 2> trace.out

.PHONY: debug
debug: ## Run all container without app
	docker-compose up -d --scale app=0
	make run

.PHONY: init
init: ## Init project
	docker-compose -f docker-compose.prod.yml -f docker-compose.yml up --build -d
	sleep 5
	make migrate

.PHONY: run-with-shards
run-with-shards: ## Run project with shards
	docker-compose -f docker-compose.shard.yml -f docker-compose.yml up --build --scale mysql=0

.PHONY: report
report: ## Generate report by pprof. usage: make report type=heap|profile|block|mutex|trace
	curl -s http://${APP_HOST}:${APP_PORT}/debug/pprof/$(type) > ./$(type).out
ifeq ($(type),trace)
	go tool trace -http=:8080 ./$(type).out
else
	go tool pprof -http=:8080 ./$(type).out
endif

.PHONY: watch
watch: ## Run in live-reload mode
	make stop
	docker-compose up

.PHONY: stop
stop: ## Remove containers but keep volumes
	docker-compose down --remove-orphans

.PHONY: clear
clear: ### Remove containers and volumes
	docker-compose down --remove-orphans --volumes

.PHONY: rebuild
rebuild: ## Rebuild by Docker Compose
	make stop
	docker-compose build --no-cache

# make goose cmd="create database sql"
# make goos cmd="up"
# make goos cmd="down"
# Pay attention, --network parameter must be the same as network in docker-compose.yml file.
.PHONY: goose
goose: ## Work with migration
	docker run -ti --rm -u $(shell id -u) --workdir=/home --network=${APP_NAME}_default -v $(shell pwd):/home jerray/goose goose -dir=migrations mysql "$(MY_USER):$(MY_PASSWORD)@($(MY_HOST):$(MY_PORT))/$(MY_DB_NAME)?parseTime=$(MY_PARSE_TIME)" $(cmd)

.PHONY: goose-shard
goose-shard: ## Work with migration
	docker run -ti --rm -u $(shell id -u) --workdir=/home --network=${APP_NAME}_default -v $(shell pwd):/home jerray/goose goose -dir=migrations mysql "test:pzjqUkMnc7vfNHET@($(shard):3306)/test?parseTime=true" $(cmd)

.PHONY: migrate-shard
migrate-shard: ## Migration up
	make goose-shard shard="db-node-1" cmd="up"
	make goose-shard shard="db-node-2" cmd="up"

.PHONY: migrate
migrate: ## Migration up
	make goose cmd="up"

.PHONY: compile
compile: ## Make binary and docs
	go build -ldflags="${LDFLAGS}" -o bin/${APP_NAME} cmd/${APP_NAME}/main.go cmd/${APP_NAME}/wire_gen.go

## https://goswagger.io/use/spec.html
.PHONY: swagger-gen
swagger-gen: ## Validate swagger spec and generate its JSON version.
ifneq ($(CI),)
	swagger generate spec -i ./swagger.yml -o ./swagger/api.swagger.json
	swagger validate ./swagger/api.swagger.json
else
	docker run --rm -it -e GOPATH=${HOME}/go:/go -v ${HOME}:${HOME} -w $(shell pwd) quay.io/goswagger/swagger validate ./swagger.yml
	docker run --rm -it -e GOPATH=${HOME}/go:/go -v ${HOME}:${HOME} -w $(shell pwd) quay.io/goswagger/swagger generate spec -i ./swagger.yml -o ./swagger/api.swagger.json
endif

.PHONY: run-cluster
run-cluster: ## run mysql cluster
	docker run --rm -v $(shell pwd)/my.cnf:/etc/mysql/conf.d/my.cnf -d -e MYSQL_ROOT_PASSWORD=root -e CLUSTER_NAME=cluster_osoc-dialog -p 3307:3306 --name=osoc-dialog_node1 --net=osoc-dialog_default percona/percona-xtradb-cluster:5.7
	sleep 30
	docker run --rm -d -e MYSQL_ROOT_PASSWORD=root -e CLUSTER_NAME=cluster_osoc-dialog -e CLUSTER_JOIN=osoc-dialog_node1 -p 3308:3306 --name=osoc-dialog_node2 --net=osoc-dialog_default percona/percona-xtradb-cluster:5.7
	sleep 5
	docker run --rm -v $(shell pwd)/my.cnf:/etc/mysql/conf.d/my.cnf -d -e MYSQL_ROOT_PASSWORD=root -e CLUSTER_NAME=cluster_osoc-dialog -e CLUSTER_JOIN=osoc-dialog_node1 -p 3309:3306 --name=osoc-dialog_node3 --net=osoc-dialog_default percona/percona-xtradb-cluster:5.7
	docker exec -i osoc-dialog_node1 mysql -uroot -proot -e "create database osoc-dialog;"
	docker run -ti --rm -u $(shell id -u) --workdir=/home --network=${APP_NAME}_default -v $(shell pwd):/home jerray/goose goose -dir=migrations mysql "root:root@(osoc-dialog_node1:3306)/osoc-dialog" up

#docker run --rm -d -e MYSQL_ROOT_PASSWORD=root -e CLUSTER_NAME=cluster_osoc-dialog -p 3307:3306 --name=osoc-dialog_node1 --net=osoc-dialog_default percona/percona-xtradb-cluster:5.7
#docker run --rm -d -e MYSQL_ROOT_PASSWORD=root -e CLUSTER_NAME=cluster_osoc-dialog -e CLUSTER_JOIN=osoc-dialog_node1 -p 3308:3306 --name=osoc-dialog_node2 --net=osoc-dialog_default percona/percona-xtradb-cluster:5.7
#test cluster
#for node1: create database osoc-dialog;
#for migrations
#MY_HOST=osoc-dialog_node1
#MY_PORT=3306
