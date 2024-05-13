include .env

DOCKER=docker --context=lima-docker

.PHONY: build-app run-lima run-service run-application run-pgsql fix-permissions-pgsql stop-pgsql

ci-local:
	golangci-lint run

run-lima:
	limactl start docker

build-app:
	docker --context=lima-docker build --tag=devhands:app .

run-service:
	docker --context=lima-docker run --publish 8080:8080 devhands:app

run-application:
	docker-compose --context=lima-docker up

run-memcached:
	$(DOCKER) run -it --rm --name devhands_memcached -p 11211:11211 memcached:latest

# https://github.com/docker-library/docs/blob/master/postgres/README.md
# How to configure permissions for PG data directory
run-pgsql:
	docker --context=lima-docker run -it --rm \
	--name devhands_pgsql \
	--user 1000:1000 \
	-e POSTGRES_PASSWORD=$(DB_PASSWORD) \
	-e POSTGRES_USER=$(DB_USER) \
	-v pgdata:/var/lib/postgresql/data \
	-p 5432:5432 \
	postgres:latest

fix-permissions-pgsql:
	docker --context=lima-docker run -it --rm --user 1000:1000 -v pgdata:/var/lib/postgresql/data postgres

stop-pgsql:
	 docker --context=lima-docker stop devhands_pgsql

docker-ps:
	$(DOCKER) ps
