prod:
	docker-compose -f docker-compose.yml up --build -d --scale backend=2 parser=2

dev:
	docker-compose -f docker-compose-dev.yml up --build

.PHONY: prod dev