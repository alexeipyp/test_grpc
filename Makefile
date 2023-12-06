SERVICE_NAME=test_grpc
INSTALL_PATH=/usr/local/bin
CONFIG_DIR=/etc/$(SERVICE_NAME)/
CONFIG_FILENAME=config.json
SYSTEMD_SETUP_PATH=/etc/systemd/system/
LOGS_DIR=/var/log/$(SERVICE_NAME)/

CONFIG_PATH=./config/config.json
LOGS_PATH=./logs/

generate-proto:
	protoc -I api/proto/v1 api/proto/v1/employee.proto \
	--go_out=pkg/employee/v1  --go_opt=paths=source_relative \
	--go-grpc_out=pkg/employee/v1 --go-grpc_opt=paths=source_relative \
	--doc_out=./docs --doc_opt=markdown,proto_v1.md

lint:
	golangci-lint run ./... --config=./.golangci.yml

lint-fast:
	golangci-lint run ./... --fast --config=./.golangci.yml

run-from-source:
	mkdir -p $(LOGS_PATH)
	CONFIG_PATH=$(CONFIG_PATH) LOGS_PATH=$(LOGS_PATH) go run ./cmd/grpc_server/main.go

run-unittests:
	bash ./scripts/unittests.sh | cat

run-unittests-coverage:
	bash ./scripts/coverage.sh | cat

build-linux:
	GOOS=linux cd ./cmd/grpc_server && go build -o ../../bin/$(SERVICE_NAME)

run-binary:
	mkdir -p $(LOGS_PATH)
	CONFIG_PATH=$(CONFIG_PATH) LOGS_PATH=$(LOGS_PATH) ./bin/$(SERVICE_NAME)

build-and-run: build-linux run-binary

install:
ifneq ($(shell id -u), 0)
	sudo make $@
else
	cp ./bin/$(SERVICE_NAME) $(INSTALL_PATH)
endif

setup-config:
ifneq ($(shell id -u), 0)
	sudo make $@
else
	mkdir -p $(CONFIG_DIR) && cp ./config/$(CONFIG_FILENAME) $(CONFIG_DIR)$(CONFIG_FILENAME)
endif

gen-env-config:
	@echo "[Service]" > ./init/env.conf
	@echo "Environment=\"CONFIG_PATH=$(CONFIG_DIR)$(CONFIG_FILENAME)\"" >> ./init/env.conf
	@echo "Environment=\"LOGS_PATH=$(LOGS_DIR)\"" >> ./init/env.conf
	@echo "ExecStart=\"$(INSTALL_PATH)/$(SERVICE_NAME)\"" >> ./init/env.conf

install-systemd-service:	build-linux	install	setup-config gen-env-config
ifneq ($(shell id -u), 0)
	sudo make $@
else
	mkdir -p $(LOGS_DIR)
	cp ./init/$(SERVICE_NAME).service $(SYSTEMD_SETUP_PATH)
	mkdir -p $(SYSTEMD_SETUP_PATH)$(SERVICE_NAME).service.d
	cp ./init/env.conf $(SYSTEMD_SETUP_PATH)$(SERVICE_NAME).service.d/env.conf
	systemctl daemon-reload
endif