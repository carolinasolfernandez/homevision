NAME=homevision

.PHONY: build
## build: Compile the application.
build:
	@go build -o $(NAME)

.PHONY: run
## run: Build and Run.
run: build
	@./$(NAME)

.PHONY: clean
## clean: Clean project and previous builds.
clean:
	@rm -f $(NAME)

.PHONY: deps
## deps: Download modules
deps:
	@go mod download

.PHONY: test
## test: Run tests
test:
	@go test ./...
