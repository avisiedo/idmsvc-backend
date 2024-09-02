.PHONY: install-pre-commit
install-pre-commit: $(PRE_COMMIT)

.PHONY: lint
lint: $(PRE_COMMIT) $(GOLANGCI_LINT)
	PATH=$(TOOLS_BIN):$${PATH} $(PYTHON_VENV)/bin/pre-commit run --all-files --verbose
# For troubleshooting
# PATH=$(TOOLS_BIN):$${PATH} $(GOLANGCI_LINT) run --verbose
	PATH=$(TOOLS_BIN):$${PATH} $(GOLANGCI_LINT) config verify
	PATH=$(TOOLS_BIN):$${PATH} $(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: $(PRE_COMMIT) $(GOLANGCI_LINT)
	PATH=$(TOOLS_BIN):$${PATH} $(GOLANGCI_LINT) run --fix


.PHONY: shellcheck
shellcheck:
	shellcheck -x -s bash -P test/scripts/ test/scripts/*.sh test/scripts/*.inc
