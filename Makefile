PKL_SRC := $(wildcard pkl/*.pkl)

all: gen run

programs/config/Config.pkl.go: $(PKL_SRC)
	rm -f programs/*/*.pkl.go
	pkl-gen-go pkl/config.pkl

.PHONY: gen
gen: programs/config/Config.pkl.go

.PHONY: run
run: gen
	go run main.go s -d -c test-config.pkl
