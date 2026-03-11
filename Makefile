PKL_SRC := $(wildcard pkl/*.pkl)

all: gen run

internal/programs/config/Config.pkl.go: $(PKL_SRC)
	rm -f internal/programs/*/*.pkl.go
	pkl-gen-go pkl/config.pkl

.PHONY: gen
gen: internal/programs/config/Config.pkl.go

.PHONY: sync
sync: gen
	go run main.go sync -d -c test-config.pkl

.PHONY: edit
edit: gen
	go run main.go secret edit -d -c test-config.pkl
