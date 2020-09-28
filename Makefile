SRC := main.go \
	   web.go \
	   db.go \
	   auth.go \
	   worker.go \

all: deps todo

deps:
	go get

todo: $(SRC)
	go build

web: $(SRC)
	go run $(SRC) $@

worker: $(SRC)
	go run $(SRC) $@

test:
	go test
