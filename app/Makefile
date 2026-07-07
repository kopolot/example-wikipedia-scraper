.PHONY: build scraper api clean notify migrate

build: scraper api notify

scraper:
	go build -o bin/scraper.bin ./cmd/scraper/main.go

api:
	go build -o bin/api.bin ./cmd/api/main.go

notify:
	go build -o bin/notify.bin ./cmd/notify/main.go

migrate:
	go build -o bin/migrate.bin ./cmd/migrate/main.go

clean:
	rm -f bin/*.bin
