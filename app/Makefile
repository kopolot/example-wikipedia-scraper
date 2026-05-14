.PHONY: build scraper api clean notify

build: scraper api notify

scraper:
	go build -o bin/scraper.bin ./cmd/scraper/main.go

api:
	go build -o bin/api.bin ./cmd/api/main.go

notify:
	go build -o bin/notify.bin ./cmd/notify/main.go

clean:
	rm -f bin/*.bin
