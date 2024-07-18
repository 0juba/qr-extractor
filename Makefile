include .env.local

DOCKER=docker
GOPROXY=https://proxy.golang.org,direct

.PHONY: build-app run-lima run-service run-application run-pgsql fix-permissions-pgsql stop-pgsql

remove-vendor-from-git:
	git rm -rf --cached vendor

ci-local:
	golangci-lint run

run-lima:
	limactl start docker

build-app:
	docker build --tag=devhands:app .

run-service:
	docker run --publish 8080:8080 devhands:app

run-application:
	docker-compose up

run-memcached:
	$(DOCKER) run -it --rm --name devhands_memcached -p 11211:11211 memcached:latest

# https://github.com/docker-library/docs/blob/master/postgres/README.md
# How to configure permissions for PG data directory
run-pgsql:
	docker run -it --rm \
	--name devhands_pgsql \
	--user 1000:1000 \
	-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
	-e POSTGRES_USER=$(DB_USER) \
	-v pgdata:/var/lib/postgresql/data \
	-p 5432:5432 \
	postgres:latest

fix-permissions-pgsql:
	docker run -it --rm --user 1000:1000 -v pgdata:/var/lib/postgresql/data postgres

stop-pgsql:
	 docker stop devhands_pgsql

docker-ps:
	$(DOCKER) ps

start-local-env: rebuild-app
	docker-compose --env-file=.env.local up -d --wait
	docker-compose --env-file=.env.local logs -f app

rm-local-env:
	docker-compose --env-file=.env.local down -v

rebuild-app:
	docker-compose --env-file=.env.local build app

rebuild: rm-local-env rebuild-app start-local-env
	@$(info rebuild app)
