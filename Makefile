GOCMD        = go
GOPATH       = $(CURDIR)/.gopath~
BIN          = $(GOPATH)/bin
BASE         = $(GOPATH)/src/$(PACKAGE)
GOBUILD      = $(GOCMD) build
GOCLEAN      = $(GOCMD) clean
PACKAGE      = bootstrap
DATE         ?= $(shell date +%F)
PLATFORMS    := linux/amd64 linux/arm windows/amd64
PKG_BIN      =  "dep"

temp = $(subst /, ,$@)
os = $(word 1, $(temp))
arch = $(word 2, $(temp))

V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

$(BASE): ; $(info $(M) setting GOPATH…)
	@mkdir -p $(dir $@)
	@ln -sf $(CURDIR) $@

.PHONY: pre-base
pre-base: $(BASE)

vendor: Gopkg.lock | $(info $(M) retrieving dependencies…)
	$(PKG_BIN) ensure -update
	$(BASE)

all:	clean pre-base release

build: ; $(info $(M) building bootstrap)
	$(GOBUILD) -o bin/$(PACKAGE) -v $(BASE)/src/main.go


.PHONY: clean
clean: ; $(info $(M) cleaning…) @ ## Cleanup everything
	@rm -rf $(GOPATH)
	@rm -rf bin

.PHONY: release $(PLATFORMS)
release: $(PLATFORMS)

$(PLATFORMS):
	GOOS=$(os) GOARCH=$(arch) $(GOBUILD) -o 'bin/$(PACKAGE)_$(os)_$(arch)_$(DATE)'  -v $(BASE)/src/main.go
