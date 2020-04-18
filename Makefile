all: linters tests

linters:
	golangci-lint run --disable-all \
		--enable gosimple \
		--enable errcheck \
		--enable golint \
		--enable govet \
		--enable staticcheck \
		--enable gocritic \
		--enable gosec \
		--enable scopelint \
		--enable prealloc \
		--enable maligned \
		--enable ineffassign \
		--enable unparam \
		--enable deadcode \
		--enable unused \
		--enable varcheck \
		--enable unconvert \
		--enable misspell \
		--enable goconst \
		--enable gochecknoinits \
		--enable gochecknoglobals \
		--enable nakedret \
		--enable gocyclo \
		--enable goimports \
		--enable gofmt

tests:
	go test tic_tac_toe/game
