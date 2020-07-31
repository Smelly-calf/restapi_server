all: users
users:
    go build -o bin/users main/main.go

test:
    go test -v test/service_test.go