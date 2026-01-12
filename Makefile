RVGLSM ?= $(LOCALBIN)/rvglsm

.PHONY: rvglsm
rvglsm: $(RVGLSM)
$(RVGLSM): $(LOCALBIN)
	@dagger call binary export --path $(RVGLSM)

.PHONY: .git/hooks .git/hooks/ .git/hooks/pre-commit
.git/hooks .git/hooks/ .git/hooks/pre-commit:
	@cp .githooks/* .git/hooks

.PHONY: release
release:
	@git tag $(SEMVER)
	@git push
	@git push --tags

LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	@mkdir -p $(LOCALBIN)

BIN ?= ~/.local/bin
INSTALL ?= install

.PHONY: install
install: rvglsm
	@$(INSTALL) $(RVGLSM) $(BIN)
