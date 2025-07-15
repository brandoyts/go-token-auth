start:
	make stop && docker compose up --build

stop:
	docker compose down --volumes