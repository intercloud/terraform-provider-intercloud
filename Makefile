TEST?=./...
PKG_NAME?=intercloud
PROVIDER_VERSION?=1.0.0-beta
WEBSITE_REPO=github.com/hashicorp/terraform-website

export GO111MODULE := on

default: run-dev

build: fmtcheck
	@go install

fmtcheck:
	@echo "==> Checking source code against gofmt..."
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

fmt:
	@echo "==> Fixing source code with gofmt..."
	@gofmt -w -s ./$(PKG_NAME)

test: fmtcheck
	@go test -i $(TEST) || exit 1
	@go test $(TEST) $(TESTARGS) -timeout=120s -parallel=4

testacc: fmtcheck
	TF_ACC=1 TF_SCHEMA_PANIC_ON_ERROR=1 go test $(TEST) -v $(TESTARGS) -timeout 240m

run-dev: clean
	@echo "==> Building and applying..."
	@go build -o terraform-provider-intercloud
	@terraform init
	@terraform plan

lint:
	@echo "==> Checking source code..."
	@GOGC=30 golangci-lint run --new-from-rev=HEAD~ ./$(PKG_NAME)/...
	@$(MAKE) provider-lint
	
provider-lint-all:
	@tfproviderlint -c 1 ./$(PKG_NAME)

provider-lint:
	@tfproviderlint -c 1 \
	  -AT001 \
		-AT002 \
		-AT005 \
		-AT006 \
		-AT007 \
		-AT008 \
		-R002 \
		-R004 \
		-R006 \
		-R012 \
		-R013 \
		-R014 \
		-S001 \
		-S002 \
		-S003 \
		-S004 \
		-S005 \
		-S007 \
		-S008 \
		-S009 \
		-S010 \
		-S011 \
		-S012 \
		-S013 \
		-S014 \
		-S015 \
		-S016 \
		-S017 \
		-S019 \
		-S021 \
		-S025 \
		-S026 \
		-S027 \
		-S028 \
		-S029 \
		-S030 \
		-S031 \
		-S032 \
		-S033 \
		-S034 \
		-S035 \
		-S036 \
		-S037 \
		-V002 \
		-V003 \
		-V004 \
		-V006 \
		-V007 \
		-V008 \
	 ./$(PKG_NAME)

pre-commit-install:
	@echo "==> Installing pre-commit hooks..."
	@pre-commit install
	@pre-commit install --hook-type commit-msg

pre-commit-autoupdate:
	@pre-commit autoupdate

tools:
	@go install github.com/bflad/tfproviderlint/cmd/tfproviderlint
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint
	@go install github.com/goreleaser/goreleaser
	@go mod tidy

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	@echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	@git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	@echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	@git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

release-snapshot: fmtcheck tools clean
	@goreleaser build --rm-dist --snapshot

clean:
	@rm -rf ./terraform-provider-intercloud ./dist

.PHONY: build test run-dev pre-commit-install pre-commit-autoupdate lint tools \
	fmtcheck provider-lint provider-lint-all fmt website website-test release-snapshot clean
