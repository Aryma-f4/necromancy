# Build the Go application
build:
	go mod tidy
	go build -o necromancy

run: build
	./necromancy
