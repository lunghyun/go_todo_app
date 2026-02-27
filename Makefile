.PHONY: help build build-local up down logs ps test
.DEFAULT_GOAL := help

DOCKER_TAG := latest

build: ## 배포용 도커 이미지 빌드
	docker build -t himello/gotodo:${DOCKER_TAG} \
		--target deploy ./

build-local: ## 로컬 환경용 도커 이미지 빌드
	docker compose build --no-cache

up: ## 자동 새로고침을 사용한 도커 구성기 실행
	docker compose up -d

down: ## 도커 구성기 종료
	docker compose down

logs: ## 도커 구성기 로그 출력
	docker compose logs -f

ps: ## 컨테이너 상태 확인
	docker compose ps -a

test: ## 테스트 실행
	go test -race -shuffle=on ./...

help: ## 옵션 보기
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

cov: ## 테스트 실행으로 커버리지 데이터 저장
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html

include .env

migrate: ## exec: go install github.com/sqldef/sqldef/cmd/mysqldef@latest
	mysqldef -u $(TODO_DB_USER) -p $(TODO_DB_PASSWORD) -h 127.0.0.1 -P 33306 $(TODO_DB_NAME) < ./_tools/mysql/schema.sql

dry-migrate:
	mysqldef -u $(TODO_DB_USER) -p $(TODO_DB_PASSWORD) -h 127.0.0.1 -P 33306 $(TODO_DB_NAME) --dry-run < ./_tools/mysql/schema.sql

get-health:
	curl -i -XGET localhost:18000/health

get-tasks:
	curl -i -XGET localhost:18000/tasks

post-ok-tasks:
	curl -i -XPOST localhost:18000/tasks -d @./handler/testdata/add_task/ok_req.json.golden

post-bad-tasks:
	curl -i -XPOST localhost:18000/tasks -d @./handler/testdata/add_task/bad_req.json.golden

