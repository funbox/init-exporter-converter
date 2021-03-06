################################################################################

# This Makefile generated by GoMakeGen 1.3.1 using next command:
# gomakegen .
#
# More info: https://kaos.sh/gomakegen

################################################################################

.DEFAULT_GOAL := help
.PHONY = fmt vet all clean git-config deps help

################################################################################

all: init-exporter-converter ## Build all binaries

init-exporter-converter: ## Build init-exporter-converter binary
	go build init-exporter-converter.go

install: ## Install all binaries
	cp init-exporter-converter /usr/bin/init-exporter-converter

uninstall: ## Uninstall all binaries
	rm -f /usr/bin/init-exporter-converter

git-config: ## Configure git redirects for stable import path services
	git config --global http.https://pkg.re.followRedirects true

deps: git-config ## Download dependencies
	go get -d -v github.com/funbox/init-exporter
	go get -d -v pkg.re/essentialkaos/ek.v12
	go get -d -v pkg.re/essentialkaos/go-simpleyaml.v2

fmt: ## Format source code with gofmt
	find . -name "*.go" -exec gofmt -s -w {} \;

vet: ## Runs go vet over sources
	go vet -composites=false -printfuncs=LPrintf,TLPrintf,TPrintf,log.Debug,log.Info,log.Warn,log.Error,log.Critical,log.Print ./...

clean: ## Remove generated files
	rm -f init-exporter-converter

help: ## Show this info
	@echo -e '\n\033[1mSupported targets:\033[0m\n'
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[33m%-25s\033[0m %s\n", $$1, $$2}'
	@echo -e ''
	@echo -e '\033[90mGenerated by GoMakeGen 1.3.1\033[0m\n'

################################################################################
