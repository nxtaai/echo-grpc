.PHONY: lint
lint:
	@echo "Running lint"
	golangci-lint run
	buf lint
	buf format -d --exit-code

.PHONY: lintfix
lintfix:
	@echo "Running lintfix"
	@echo "Automatically fix some lint errors"
	golangci-lint run --fix
	buf format -w

.PHONY: generate
generate:
	@echo "Running generate"
	buf generate

.PHONY: checkgenerate
checkgenerate:
	@echo "Running checkgenerate"
	buf generate
	test -z "$$(git status --porcelain | tee /dev/stderr)"

.PHONY: breaking
breaking:
	@echo "Running breaking"
	buf breaking --against '.git#branch=main'
