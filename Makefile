.PHONY: test test-cover bench test-folder bench-folder tag

MODULES_RAW := $(shell go list -m)
MODULES := $(filter github.com/jarollz/data-structure/%,$(MODULES_RAW))
PKGS := $(addsuffix /...,$(MODULES))

ifeq ($(firstword $(MAKECMDGOALS)),tag)
TAG_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
ifeq ($(origin FOLDER), undefined)
FOLDER := $(word 1,$(TAG_ARGS))
endif
ifeq ($(origin VERSION), undefined)
VERSION := $(word 2,$(TAG_ARGS))
endif
ifneq ($(strip $(TAG_ARGS)),)
.PHONY: $(TAG_ARGS)
$(TAG_ARGS):
	@:
endif
endif

FOLDER_MODULE := github.com/jarollz/data-structure/$(FOLDER)

define ensure_folder
	@test -n "$(FOLDER)" || (echo "FOLDER is required, example: make $(1) FOLDER=list-array" && exit 1)
	@test -f "$(FOLDER)/go.mod" || (echo "Invalid FOLDER: $(FOLDER) (missing $(FOLDER)/go.mod)" && exit 1)
endef

define ensure_version
	@test -n "$(VERSION)" || (echo "VERSION is required, example: make tag FOLDER=list-array VERSION=v1.2.3" && exit 1)
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

tag:
	$(call ensure_folder,tag)
	$(call ensure_version)
	./.scripts/release-tag.sh "$(FOLDER)" "$(VERSION)"
