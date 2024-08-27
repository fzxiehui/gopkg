
.PHONY: default help

default: test

help:
	@echo 'Management commands for Project:'
	@echo
	@echo 'Usage:'
	@echo '    make test            Run tests on a compiled project.'
	@echo

test:
	go test -v -count=1 ./...
