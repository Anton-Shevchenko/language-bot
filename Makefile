.PHONY:
.SILENT:

build-image:
	docker build -t languages-telegram-bot go-app

start-container:
	docker run --name language-bot -p 80:80 --env-file go-app/.env languages-telegram-bot