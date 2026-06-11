.PHONY: test test-cover bench

MODULES_RAW := $(shell go list -m)
MODULES := $(filter github.com/jarollz/data-structure/%,$(MODULES_RAW))
PKGS := $(addsuffix /...,$(MODULES))

test:
	go test $(PKGS)

test-cover:
	mkdir -p tmp
	go test $(PKGS) -coverprofile=tmp/coverage.out
	go tool cover -func=tmp/coverage.out | tee tmp/coverage.txt

bench:
	go test $(PKGS) -run ^$$ -bench . -benchmem
