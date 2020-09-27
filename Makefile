SRC := main.go \
	   web.go \
	   db.go \
	   auth.go \
	   worker.go \

all: deps todo

deps:
	go get

todo: main.go web.go db.go
	go build

web: $(SRC)
	go run $(SRC) $@

worker: $(SRC)
	go run $(SRC) $@
