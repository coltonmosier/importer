all: cleaner importer

cleaner: 
	go build -o bin/cleaner ./cmd/cleaner/*

importer: 
	go build -o bin/importer ./cmd/importer/*

clean:
	rm -f bin/cleaner bin/importer
