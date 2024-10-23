build:
	@cd cmd && go build -o ../bin/summary .

run: build
	@./bin/summary
