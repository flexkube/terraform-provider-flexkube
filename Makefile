# Build parameters
CGO_ENABLED=0
LD_FLAGS="-extldflags '-static'"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOBUILD=CGO_ENABLED=$(CGO_ENABLED) $(GOCMD) build -v -buildmode=exe -ldflags $(LD_FLAGS)
GO_PACKAGES=./...
GO_TESTS=^.*$

GOLANGCI_LINT_VERSION=v1.46.0

BIN_PATH=$$HOME/bin

CONTROLLERS=$(shell (grep CONTROLLERS .env 2>/dev/null || echo "1") | cut -d= -f2 2>/dev/null)
WORKERS=$(shell (grep WORKERS .env 2>/dev/null || echo "2") | cut -d= -f2 2>/dev/null)
NODES_CIDR=$(shell (grep NODES_CIDR .env 2>/dev/null || echo "192.168.50.0/24") | cut -d= -f2 2>/dev/null)
FLATCAR_CHANNEL=$(shell (grep FLATCAR_CHANNEL .env 2>/dev/null || echo "stable") | cut -d= -f2 2>/dev/null)
TERRAFORM_ENV=TF_VAR_flatcar_channel=$(FLATCAR_CHANNEL) TF_VAR_controllers_count=$(CONTROLLERS) TF_VAR_workers_count=$(WORKERS) TF_VAR_nodes_cidr=$(NODES_CIDR)
TERRAFORM_BIN=$(TERRAFORM_ENV) /usr/bin/terraform

E2E_IMAGE=flexkube/terraform-provider-flexkube-e2e

E2E_CMD=docker run -it --rm -v /home/core/terraform-provider-flexkube:/root/terraform-provider-flexkube -v /home/core/.ssh:/root/.ssh -v /home/core/.terraform.d:/root/.terraform.d -w /root/terraform-provider-flexkube --net host --entrypoint /bin/bash -e TF_VAR_flatcar_channel=$(FLATCAR_CHANNEL) -e TF_VAR_controllers_count=$(CONTROLLERS) -e TF_VAR_workers_count=$(WORKERS) -e TF_VAR_nodes_cidr=$(NODES_CIDR) $(E2E_IMAGE)

BUILD_CMD=docker run -it --rm -v /home/core/terraform-provider-flexkube:/usr/src/terraform-provider-flexkube -v /home/core/go:/go -v /home/core/.cache:/root/.cache -v /run:/run -w /usr/src/terraform-provider-flexkube $(INTEGRATION_IMAGE)

INTEGRATION_IMAGE=flexkube/terraform-provider-flexkube-integration

INTEGRATION_CMD=docker run -it --rm -v /run:/run -v /home/core/terraform-provider-flexkube:/usr/src/terraform-provider-flexkube -v /home/core/go:/go -v /home/core/.password:/home/core/.password -v /home/core/.ssh:/home/core/.ssh -v /home/core/.cache:/root/.cache -w /usr/src/terraform-provider-flexkube --net host $(INTEGRATION_IMAGE)

VAGRANTCMD=$(TERRAFORM_ENV) vagrant

COVERPROFILE=c.out

CC_TEST_REPORTER_ID=5bc3e58aca2ff47897d533ba92ae8db15ac9fdb83fad3637301ee5d75ccd4143

.PHONY: all
all: build build-test test lint semgrep

.PHONY: download
download:
	$(GOMOD) download

.PHONY: install-golangci-lint
install-golangci-lint:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(BIN_PATH) $(GOLANGCI_LINT_VERSION)

.PHONY: install-cc-test-reporter
install-cc-test-reporter:
	curl -L https://codeclimate.com/downloads/test-reporter/test-reporter-latest-linux-amd64 > $(BIN_PATH)/cc-test-reporter
	chmod +x $(BIN_PATH)/cc-test-reporter

.PHONY: install-ci
install-ci: install-golangci-lint install-cc-test-reporter

.PHONY: build
build:
	$(GOBUILD)

.PHONY: test
test: build-test
	$(GOTEST) -run $(GO_TESTS) $(GO_PACKAGES)

.PHONY: lint
lint:
	golangci-lint run $(GO_PACKAGES)

.PHONY: build-test
build-test:
	$(GOTEST) -run=nope $(GO_PACKAGES)

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(OUTPUT_FILE) || true
	rm -f $(OUTPUT_FILE).sig || true

.PHONY: update
update:
	$(GOGET) -u $(GO_PACKAGES)
	$(GOMOD) tidy

.PHONY: update-linters
update-linters:
	# Remove all enabled linters.
	sed -i '/^  enable:/q0' .golangci.yml
	# Then add all possible linters to config.
	golangci-lint linters | grep -E '^\S+:' | cut -d: -f1 | sort | sed 's/^/    - /g' | grep -v -E "($$(grep '^  disable:' -A 100 .golangci.yml  | grep -E '    - \S+$$' | awk '{print $$2}' | tr \\n '|' | sed 's/|$$//g'))" >> .golangci.yml

.PHONY: test-update-linters
test-update-linters:
	@test -z "$$(git status --porcelain)" || (echo "Working directory must be clean to perform this check."; exit 1)
	make update-linters
	@test -z "$$(git status --porcelain)" || (echo "Linter configuration outdated. Run 'make update-linters' and commit generated changes to fix."; exit 1)

.PHONY: all-cover
all-cover: build build-test test-cover lint

.PHONY: test-cover
test-cover: build-test
	$(GOTEST) -run $(GO_TESTS) -coverprofile=$(COVERPROFILE) $(GO_PACKAGES)

.PHONY: cover-browse
cover-browse:
	go tool cover -html=$(COVERPROFILE)

.PHONY: test-cover-browse
test-cover-browse: test-cover cover-browse

.PHONY: test-cover-upload-codecov
test-cover-upload-codecov: SHELL=/bin/bash
test-cover-upload-codecov: test-cover
test-cover-upload-codecov:
	bash <(curl -s https://codecov.io/bash) -f $(COVERPROFILE)

.PHONY: test-cover-upload-codeclimate
test-cover-upload-codeclimate: test-cover
test-cover-upload-codeclimate:
	env CC_TEST_REPORTER_ID=$(CC_TEST_REPORTER_ID) cc-test-reporter after-build -t gocov -p $$(go list -m)

.PHONY: test-cover-upload
test-cover-upload: test-cover-upload-codecov test-cover-upload-codeclimate

.PHONY: libvirt-apply
libvirt-apply: libvirt-download-image
	cd libvirt && $(TERRAFORM_BIN) init && $(TERRAFORM_BIN) apply -auto-approve

.PHONY: libvirt-destroy
libvirt-destroy:
	cd libvirt && $(TERRAFORM_BIN) init && $(TERRAFORM_BIN) destroy -auto-approve

.PHONY: libvirt-download-image
libvirt-download-image:
	((test -f libvirt/flatcar_production_qemu_image.img.bz2 || test -f libvirt/flatcar_production_qemu_image.img) || wget https://$(FLATCAR_CHANNEL).release.flatcar-linux.net/amd64-usr/current/flatcar_production_qemu_image.img.bz2 -O libvirt/flatcar_production_qemu_image.img.bz2) || true
	(test -f libvirt/flatcar_production_qemu_image.img.bz2 && bunzip2 libvirt/flatcar_production_qemu_image.img.bz2 && rm libvirt/flatcar_production_qemu_image.img.bz2) || true
	qemu-img resize libvirt/flatcar_production_qemu_image.img +5G

.PHONY: test-local-apply
test-local-apply:
	mkdir -p local-testing/.terraform/plugins/registry.terraform.io/flexkube-testing/flexkube/0.1.0/linux_amd64/ ~/.local/share/terraform/plugins/registry.terraform.io/flexkube-testing/flexkube/0.1.0/linux_amd64/
	$(GOBUILD) -o local-testing/.terraform/plugins/registry.terraform.io/flexkube-testing/flexkube/0.1.0/linux_amd64/terraform-provider-flexkube
	cp local-testing/.terraform/plugins/registry.terraform.io/flexkube-testing/flexkube/0.1.0/linux_amd64/terraform-provider-flexkube ~/.local/share/terraform/plugins/registry.terraform.io/flexkube-testing/flexkube/0.1.0/linux_amd64/
	cd local-testing && $(TERRAFORM_BIN) init && $(TERRAFORM_BIN) apply -auto-approve

.PHONY: build-e2e
build-e2e:
	docker build -t $(E2E_IMAGE) e2e

.PHONY: vagrant-up
vagrant-up:
	$(VAGRANTCMD) up

.PHONY: vagrant-rsync
vagrant-rsync:
	$(VAGRANTCMD) rsync

.PHONY: vagrant-destroy
vagrant-destroy:
	$(VAGRANTCMD) destroy --force

.PHONY: vagrant
vagrant: SHELL=/bin/bash
vagrant:
	alias vagrant='$(VAGRANTCMD)'

.PHONY: vagrant-e2e-build
vagrant-e2e-build:
	$(VAGRANTCMD) ssh -c "$(BUILD_CMD) make build-e2e"

.PHONY: vagrant-e2e-kubeconfig
vagrant-e2e-kubeconfig:
	scp -P 2222 -o StrictHostKeyChecking=no -i ~/.vagrant.d/insecure_private_key core@127.0.0.1:/home/core/terraform-provider-flexkube/e2e/kubeconfig ./e2e/kubeconfig

.PHONY: vagrant-build-bin
vagrant-build-bin: vagrant-integration-build
	$(VAGRANTCMD) ssh -c "$(BUILD_CMD) make build-bin"

.PHONY: vagrant-e2e-run
vagrant-e2e-run: vagrant-up vagrant-rsync vagrant-build-bin vagrant-e2e-build
	$(VAGRANTCMD) ssh -c "$(E2E_CMD) -c 'make test-e2e-run'"
	make vagrant-e2e-kubeconfig

.PHONY: vagrant-e2e-destroy
vagrant-e2e-destroy:
	$(VAGRANTCMD) ssh -c "$(E2E_CMD) -c 'make test-e2e-destroy'"

.PHONY: vagrant-e2e-shell
vagrant-e2e-shell:
	$(VAGRANTCMD) ssh -c "$(E2E_CMD)"

.PHONY: vagrant-e2e
vagrant-e2e: vagrant-e2e-run vagrant-e2e-destroy vagrant-destroy

.PHONY: vagrant-integration-build
vagrant-integration-build:
	$(VAGRANTCMD) ssh -c "docker build -t $(INTEGRATION_IMAGE) terraform-provider-flexkube/integration"

.PHONY: vagrant-integration-shell
vagrant-integration-shell:
	$(VAGRANTCMD) ssh -c "$(INTEGRATION_CMD) bash"

.PHONY: build-bin
build-bin:
	$(GOBUILD)

.PHONY: test-e2e-run
test-e2e-run: TERRAFORM_BIN=$(TERRAFORM_ENV) /bin/terraform
test-e2e-run:
	helm repo update
	mkdir -p ~/.local/share/terraform/plugins/registry.terraform.io/flexkube-testing/flexkube/0.1.0/linux_amd64/ ~/.terraform.d/plugin-cache/registry.terraform.io/flexkube-testing/flexkube/0.1.0/linux_amd64/
	cp ./terraform-provider-flexkube ~/.terraform.d/plugin-cache/registry.terraform.io/flexkube-testing/flexkube/0.1.0/linux_amd64/
	cp ./terraform-provider-flexkube ~/.local/share/terraform/plugins/registry.terraform.io/flexkube-testing/flexkube/0.1.0/linux_amd64/
	cd e2e && $(TERRAFORM_BIN) init && $(TERRAFORM_BIN) apply -auto-approve

.PHONY: test-e2e-destroy
test-e2e-destroy: TERRAFORM_BIN=$(TERRAFORM_ENV) /bin/terraform
test-e2e-destroy:
	$(TERRAFORM_BIN) -chdir=e2e destroy -auto-approve

.PHONY: semgrep
semgrep: SEMGREP_BIN=semgrep
semgrep:
	@if ! which $(SEMGREP_BIN) >/dev/null 2>&1; then echo "$(SEMGREP_BIN) binary not found, skipping extra linting"; else $(SEMGREP_BIN) --error; fi

.PHONY: test-vagrant
test-vagrant:
	vagrant validate --ignore-provider
