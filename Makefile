all: pkl run

pkl: pkl/*.pkl
	pkl-gen-go pkl/config.pkl

run: main.go
	go run main.go --dry-run
