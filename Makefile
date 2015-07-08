build:
	docker-compose build

test:
	docker-compose run --rm lib sh -c "go run cmd/wait_for_db/main.go && go test ./..."

vet:
	docker-compose run --rm lib go vet ./...

lint:
	docker-compose run --rm lib golint ./...

shell:
	docker-compose run --rm lib bash

clean:
	docker-compose stop; \
	docker-compose rm --force; \
	docker rmi --force godatabase_lib; \
	echo "all clean!" # to silence errors from the previous command

.PHONY: build test vet lint shell clean