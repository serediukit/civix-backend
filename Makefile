run:
	docker-compose up -d

restart:
	docker-compose down
	docker-compose build --no-cache
	docker-compose up -d