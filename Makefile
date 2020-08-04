all: users
users:
	go build -o bin/users main/main.go

test:
	go test -v test/service_test.go


benchmark:
	cd test && go test -bench=. -run=none