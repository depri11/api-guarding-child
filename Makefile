ENV?=dev

init:
	cp -f env/${ENV} .env

test:
	go test -v