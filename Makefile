build:
	docker-compose build

test:
	docker-compose run --rm lib sh -c 'waiter $$TEST_DATABASE_URL && go test -v ./...'

vet:
	docker-compose run --rm lib go vet ./...

lint:
	docker-compose run --rm lib golint ./...

shell:
	docker-compose run --rm lib bash

godoc:
	docker-compose run --rm --service-ports lib godoc -http :9000

clean:
	docker-compose stop; \
	docker-compose rm --force; \
	docker rmi --force godatabase_lib; \
	echo "all clean!" # to silence errors from the previous command

.PHONY: build test vet lint shell godoc clean
