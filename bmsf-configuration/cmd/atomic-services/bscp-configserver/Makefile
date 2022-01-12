SERVER = bk-bscp-configserver
ARGS   = run

# debug build
PWD         = $(shell pwd)
LOCALBUILD  = $(PWD)/build
OUTPUT_DIR ?= $(LOCALBUILD)

# debug version
LDVersionFLAG ?= "-X bk-bscp/pkg/version.VERSION=DEBUG \
	-X bk-bscp/pkg/version.BUILDTIME=DEBUG \
	-X bk-bscp/pkg/version.GITHASH=DEBUG"

BINDIR = ${OUTPUT_DIR}/$(SERVER)
BIN    = $(BINDIR)/$(SERVER)

PRO_DIR = ${BSCP_DIR}
# if BSCP_DIR doesn't exist, the PRO_DIR will be based on GOPATH
ifeq ($(PRO_DIR),)
	PRO_DIR = $(firstword $(subst :, ,${GOPATH}))/src/bk-bscp
endif

SCRIPTS = $(PRO_DIR)/scripts

export GO111MODULE=on

default:
	@echo -e "\e[34;1mBuilding $(SERVER)...\033[0m"
	go build -ldflags ${LDVersionFLAG} -o $(BIN) main.go
	@cp -r etc $(BINDIR)
	@sh $(SCRIPTS)/daemon-control/generator.sh '$(SERVER)' '$(ARGS)' '$(BINDIR)'
	@echo -e "\e[34;1mBuild $(SERVER) success!\n\033[0m"

clean:
	@rm -rf $(BINDIR) $(LOCALBUILD)
