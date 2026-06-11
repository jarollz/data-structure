.PHONY: test test-cover bench test-folder bench-folder

MODULES_RAW := $(shell go list -m)
MODULES := $(filter github.com/jarollz/data-structure/%,$(MODULES_RAW))
PKGS := $(addsuffix /...,$(MODULES))

ifndef FOLDER
FOLDER :=
endif

FOLDER_MODULE := github.com/jarollz/data-structure/$(FOLDER)

define ensure_folder
	@test -n "$(FOLDER)" || (echo "FOLDER is required, example: make $(1) FOLDER=list-array" && exit 1)
	@test -f "$(FOLDER)/go.mod" || (echo "Invalid FOLDER: $(FOLDER) (missing $(FOLDER)/go.mod)" && exit 1)
endef

test:
	go test $(PKGS)

test-cover:
	mkdir -p tmp
	go test $(PKGS) -coverprofile=tmp/coverage.out
	go tool cover -func=tmp/coverage.out | tee tmp/coverage.txt

bench:
	go test $(PKGS) -run ^$$ -bench . -benchmem

test-folder:
	$(call ensure_folder,test-folder)
	go test $(FOLDER_MODULE)/...

bench-folder:
	$(call ensure_folder,bench-folder)
	go test $(FOLDER_MODULE)/... -run ^$$ -bench . -benchmem
