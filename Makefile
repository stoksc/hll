
VERSION?="0.0.1"
GOFMT_FILES?=$$(find . -not -path "./vendor/*" -type f -name '*.go')


default: bin


fmt:
	@gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/fmtcheck.sh'"\

test: fmtcheck
	sh -c "'$(CURDIR)/scripts/test.sh'"

cover: fmtcheck
	sh -c "'$(CURDIR)/scripts/cover.sh'"


.NOTPARALLEL:

.PHONY: fmt fmtcheck test cover