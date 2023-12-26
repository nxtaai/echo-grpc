.PHONY: check
check: lint checkgenerate checkbreaking checktidy

.PHONY: lint
lint:
	@echo "Running lint"
	golangci-lint run
	buf lint
	buf format -d --exit-code
	npx prettier --check .

.PHONY: lintfix
lintfix:
	@echo "Running lintfix"
	@echo "Automatically fix some lint errors"
	golangci-lint run --fix
	buf format -w
	npx prettier --write .

.PHONY: generate
generate:
	@echo "Running generate"
	buf generate

.PHONY: checkgenerate
checkgenerate:
	@echo "Running checkgenerate"
	buf generate
	test -z "$$(git status --porcelain | tee /dev/stderr)"

.PHONY: checkbreaking
checkbreaking:
	@echo "Running checkbreaking"
	buf breaking --against '.git#branch=main'

.PHONY: checktidy
checktidy:
	@echo "Running checktidy"
	go mod tidy
	test -z "$$(git status --porcelain | tee /dev/stderr)"
