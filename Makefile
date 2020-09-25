SRC := main.go \
	   web.go \
	   db.go \
	   auth.go

all: deps todo

deps:
	go get

todo: main.go web.go db.go
	go build

run: $(SRC)
	go run $(SRC)
