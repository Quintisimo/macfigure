all: gen run

gen: pkl/*.pkl
	pkl-gen-go pkl/config.pkl

run: main.go
	go run main.go s -d -c test-config.pkl
