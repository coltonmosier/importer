aswe-importer: main.go
	go build -o bin/aswe-importer .

clean:
	rm -f bin/aswe-importer
