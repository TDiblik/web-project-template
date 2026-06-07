DB_NAME=project-template-db
DB_PORT=35232
DB_PASSWORD=s0m3C0mpl3xP4ss
DB_IMAGE=postgres:alpine
DB_VOLUME=$(shell pwd)/db-data

.PHONY: api api-install api-build api-update api-clean api-cleanup db db-logs db-stop db-remove db-clean gen-types gen-types-clean gen-types-clean-generated fe fe-install fe-build fe-update fe-clean prod-build prod-publish prod-locally prod-locally-logs prod-locally-stop install update build clean
%:
	@:

# ---------- Backend ----------
api:
	cd ./api && air

api-install:
	cd ./api && go mod tidy

api-build:
	cd ./api && go build -ldflags="-s -w" -trimpath -o ./tmp/main .

api-update:
	cd ./api && go get -u all && go mod tidy && gofmt -w -l .

api-clean:
	@echo "Cleaning backend build artifacts..."
	cd ./api && go clean && rm -rf ./tmp/

api-cleanup:
	@echo "Starting Aggressive Golang codebase cleanup..."
	@echo "Installing required tools..."
	cd ./api && go install github.com/elliot40404/go-clean-unused@latest
	cd ./api && go install golang.org/x/tools/cmd/deadcode@latest
	cd ./api && go install golang.org/x/tools/cmd/goimports@latest
	@echo "Removing unused imports..."
	cd ./api && goimports -w .
	@echo "Removing unused local declarations..."
	cd ./api && go-clean-unused --remove ./...
	@echo "Scanning for empty/useless .go files..."
	@for file in $$(find ./api -type f -name "*.go" -not -path "*/vendor/*"); do \
		ACTUAL_CODE=$$(grep -v '^\s*//' "$$file" | grep -v '^\s*$$' | grep -v '^package '); \
		if [ -z "$$ACTUAL_CODE" ]; then \
			echo "Removing empty file: $$file"; \
			rm "$$file"; \
		fi; \
	done
	@echo "Finding dead code..."
	@cd ./api && deadcode -json ./... > deadcode.json || true
	@echo "Nuking dead functions from orbit (Bottom-Up)..."
	@cd ./api && cat deadcode.json | jq -c '[.[].Funcs[]?] | sort_by(.Position.Line) | reverse | .[]' 2>/dev/null | while read -r item; do \
		FILE=$$(echo $$item | jq -r '.Position.File'); \
		LINE=$$(echo $$item | jq -r '.Position.Line'); \
		NAME=$$(echo $$item | jq -r '.Name'); \
		if [ ! -f "$$FILE" ]; then continue; fi; \
		echo "Deleting $$NAME from $$FILE (Line $$LINE)..."; \
		sed -i.bak "$${LINE},/^}/d" "$$FILE"; \
		rm -f "$$FILE.bak"; \
	done
	@cd ./api && rm -f deadcode.json
	@echo "Formatting code to clean up empty spaces..."
	cd ./api && go fmt ./...
	@echo "Tidying go.mod dependencies..."
	cd ./api && go mod tidy
	@echo "Aggressive cleanup complete. PLEASE CHECK GIT STATUS AND RUN TESTS."

# ---------- Database ----------
db:
	@if [ $$(docker ps -a -q -f name=$(DB_NAME)) ]; then \
		echo "Starting existing database container..."; \
		docker start $(DB_NAME); \
	else \
		echo "Running database container for the first time..."; \
		docker run --name $(DB_NAME) -p $(DB_PORT):5432 -e POSTGRES_PASSWORD=$(DB_PASSWORD) -v $(DB_VOLUME):/var/lib/postgresql/data -d $(DB_IMAGE) -c max_connections=200; \
	fi

db-logs:
	docker logs $(DB_NAME) --follow

db-stop:
	@echo "Stopping database container..."
	docker stop $(DB_NAME)

db-remove:
	@echo "Removing database container..."
	docker stop $(DB_NAME)
	docker rm -f $(DB_NAME)

db-clean: db-remove
	@echo "Removing database volume..."
	rm -rf $(DB_VOLUME)

# ---------- Type Generation ----------
gen-types:
	$(MAKE) gen-types-clean-generated
	openapi-generator generate \
		-i ./api/generated/swagger/swagger.yaml \
		-g typescript-fetch \
		-o ./shared/fe/api-client/src/generated \
		--skip-validate-spec \
		--additional-properties=supportsES6=true \
		--additional-properties=typescriptThreePlus=true \
		--additional-properties=useSingleRequestParameter=true \
		--additional-properties=modelPropertyNaming=camelCase \
		--additional-properties=enumPropertyNaming=PascalCase \
		--global-property apiDocs=false,modelDocs=false
	cd ./shared/fe/api-client && \
	bun install && \
	bun run build

gen-types-clean:
	@echo "Cleaning generated TypeScript types..."
	cd ./shared/fe/api-client && rm -rf ./node_modules ./dist

gen-types-clean-generated:
	@echo "Cleaning generated TypeScript types..."
	cd ./shared/fe/api-client && rm -rf ./src/generated

# ---------- Frontend ----------
fe:
	cd ./fe && bun run dev

fe-install:
	cd ./fe && bun install

fe-build:
	cd ./fe && bun run lint && bun run build

fe-update:
	cd ./fe && bun update --latest
	cd ./shared/fe/api-client && bun update --latest
	cd ./fe && bun run lint

fe-clean:
	@echo "Cleaning frontend build..."
	cd ./fe && rm -rf ./node_modules ./dist

# ---------- Docker build for production ----------
VERSION := $(word 2,$(MAKECMDGOALS))
DOCKER_TAG = v$(shell echo $(VERSION) | sed 's/^v*//')
prod-build:
ifndef VERSION
	$(error VERSION is required. Usage: make prod-build vX.X.X)
endif
	@echo "Generating TypeScript types..."
	$(MAKE) gen-types
	@echo "Building Docker image $(DOCKER_TAG)..."
	cp .gitignore .dockerignore
	docker build -f Dockerfile.app -t docker-registry.dev.tomasdiblik.cz/project-template:$(DOCKER_TAG) .

LOCAL_API_PROD_PORT ?= 35230
prod-locally:
ifndef VERSION
	$(error VERSION is required. Usage: make prod-locally vX.X.X)
endif
	@echo "Ensuring database is running..."
	$(MAKE) db
	$(MAKE) prod-locally-stop
	@echo "Running Docker image $(DOCKER_TAG) locally on port $(LOCAL_API_PROD_PORT)..."
	@docker run -d --env-file api/.env.production -p $(LOCAL_API_PROD_PORT):35230 docker-registry.dev.tomasdiblik.cz/project-template:$(DOCKER_TAG)

prod-locally-logs:
	@CONTAINER_ID=$$(docker ps -q --filter "publish=$(LOCAL_API_PROD_PORT)"); \
	if [ -z "$$CONTAINER_ID" ]; then \
		echo "No running container found on port $(LOCAL_API_PROD_PORT)."; \
	else \
		docker logs -f $$CONTAINER_ID; \
	fi

prod-locally-stop:
	@echo "Stopping local prod container, if it exists..."
	@docker ps -q --filter "publish=$(LOCAL_API_PROD_PORT)" | xargs -r docker stop

# ---------- Combined Targets ----------
install: api-install fe-install gen-types
update:
	$(MAKE) install
	$(MAKE) api-update
	$(MAKE) fe-update
	$(MAKE) build
build: api-build fe-build
clean: api-clean fe-clean gen-types-clean
