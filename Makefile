stop:
	docker-compose down

clean:
	docker-compose down -v 

build:
	docker-compose build

up:
	docker-compose up

restart: stop clean build up