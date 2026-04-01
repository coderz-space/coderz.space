SHELL := /bin/sh

SERVER_DIR := server
WEB_DIR := web
SERVER_COMPOSE := docker compose --env-file $(SERVER_DIR)/.env -f $(SERVER_DIR)/docker-compose.yml
WEB_IMAGE := coderz-space-web
WEB_CONTAINER := coderz-space-web
SERVER_NETWORK := coderz-space-bootcamp_default

.PHONY: help server-dev server-prod web-dev web-prod

help:
	@echo "Available targets:"
	@echo "  make server-dev   Run the Go API locally with Dockerized Postgres"
	@echo "  make server-prod  Run the API as a Docker Compose service"
	@echo "  make web-dev      Run the Next.js app locally in dev mode"
	@echo "  make web-prod     Run the Next.js app as a Docker container"

server-dev:
	@$(MAKE) -C $(SERVER_DIR) docker-up
	@$(MAKE) -C $(SERVER_DIR) migrate-up
	@$(MAKE) -C $(SERVER_DIR) run

server-prod:
	@$(SERVER_COMPOSE) up -d postgres
	@$(MAKE) -C $(SERVER_DIR) migrate-up
	@$(SERVER_COMPOSE) up --build api

web-dev:
	@npm --prefix $(WEB_DIR) run dev

web-prod:
	@$(SERVER_COMPOSE) up -d postgres
	@$(MAKE) -C $(SERVER_DIR) migrate-up
	@$(SERVER_COMPOSE) up -d --build api
	@docker build -t $(WEB_IMAGE) -f $(WEB_DIR)/Dockerfile $(WEB_DIR)
	@docker rm -f $(WEB_CONTAINER) >/dev/null 2>&1 || true
	@docker run --rm --name $(WEB_CONTAINER) \
		--network $(SERVER_NETWORK) \
		--env-file $(WEB_DIR)/.env.production \
		-p 3000:3000 \
		$(WEB_IMAGE)
