SRC := main.go \
	   web.go \
	   db.go

todo: main.go web.go db.go
	go build

run: $(SRC)
	go run $(SRC)

all: todo
