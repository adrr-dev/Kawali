APP := mewali-server
BIN := bin/$(APP)

.PHONY: all generate build run test vet tidy css clean

all: generate build

# Wajib sebelum build: hasilkan *_templ.go dari *.templ
generate:
	templ generate

build: generate
	go build -o $(BIN) ./cmd/server

run: build
	./$(BIN)

# Seluruh unit test (SQLite file/memory, tanpa layanan luar)
test:
	go test ./...

test-verbose:
	go test -v ./...

vet:
	go vet ./...

tidy:
	go mod tidy

# Butuh @tailwindcss/cli (npm i -g @tailwindcss/cli) + tailwindcss v4.
# Memakai instalasi global via NODE_PATH agar repo bebas node_modules.
css:
	PATH="$(HOME)/.local/bin:$$PATH" NODE_PATH="$$(npm root -g)" tailwindcss -i web/static/css/input.css -o web/static/css/app.css --minify

clean:
	rm -rf bin
