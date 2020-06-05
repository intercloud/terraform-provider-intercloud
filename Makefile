TEST?=./...
PKG_NAME?=intercloud
PROVIDER_VERSION?=dev
WEBSITE_REPO=github.com/hashicorp/terraform-website

default: run-dev

build: fmtcheck
	go install

fmtcheck:
	@echo "==> Checking source code against gofmt..."
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -w -s ./$(PKG_NAME)

test: fmtcheck
	go test -i $(TEST) || exit 1
	go test $(TEST) $(TESTARGS) -timeout=120s -parallel=4

testacc: fmtcheck
	TF_ACC=1 TF_SCHEMA_PANIC_ON_ERROR=1 go test $(TEST) -v $(TESTARGS) -timeout 240m

run-dev:
	@echo "==> Building and applying..."
	go build -o terraform-provider-intercloud
	terraform init
	terraform plan

install: release
	mkdir -p $(HOME)/.terraform.d/plugins/
	cp -R ./bin/* $(HOME)/.terraform.d/plugins/

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
	pre-commit install
	pre-commit install --hook-type commit-msg

pre-commit-autoupdate:
	pre-commit autoupdate

tools:
	GO111MODULE=on go install github.com/bflad/tfproviderlint/cmd/tfproviderlint
	GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint
	GO111MODULE=on go install github.com/mitchellh/gox

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)


# make release version=1.0.0 
# make release version=1.0.0 prerelease=beta
release: clean
	sed -i "-e" -e 's/Version = ".*"/Version = "$(PROVIDER_VERSION)"/g' -e 's/Prerelease = ".*"/Prerelease = "$(PRERELEASE)"/g' ./version/version.go
	gox -output ./bin/{{.OS}}_{{.Arch}}/terraform-provider-intercloud_v$(PROVIDER_VERSION)$(if $(PRERELEASE),-$(PRERELEASE),)
	gox -output ./bin/solaris_amd64/terraform-provider-intercloud_v$(PROVIDER_VERSION)$(if $(PRERELEASE),-$(PRERELEASE),) -osarch="solaris/amd64"
	rm -f ./version/version.go && mv ./version/version.go-e ./version/version.go

clean: ;@rm -rf ./terraform-provider-intercloud  rm -rf ./bin/*  rm ./bin

zip:
	@cd bin; \
	zip -r windows_386.zip windows_386; \
	zip -r windows_amd64.zip windows_amd64; \
	tar -czvf darwin_386.tar.gz darwin_386; \
	tar -czvf darwin_amd64.tar.gz darwin_amd64; \
	tar -czvf freebsd_386.tar.gz freebsd_386; \
	tar -czvf freebsd_amd64.tar.gz freebsd_amd64; \
	tar -czvf freebsd_arm.tar.gz freebsd_arm; \
	tar -czvf linux_386.tar.gz linux_386; \
	tar -czvf linux_amd64.tar.gz linux_amd64; \
	tar -czvf linux_arm.tar.gz linux_arm; \
	tar -czvf openbsd_386.tar.gz openbsd_386; \
	tar -czvf openbsd_amd64.tar.gz openbsd_amd64; \
	tar -czvf solaris_amd64.tar.gz solaris_amd64

.PHONY: build test run-dev pre-commit-install pre-commit-autoupdate lint tools \
	fmtcheck provider-lint provider-lint-all fmt website website-test release clean \
	zip
