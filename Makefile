# Configuration
NAME    := skeleton
RELEASE ?= false

# Computed variables
SHA 	:= $(shell git rev-parse --short HEAD)
VERSION := $(shell cat VERSION)

ifeq ($(RELEASE),true)
	BINDATA_FLAGS :=
else
	BINDATA_FLAGS := -debug
endif

# Lists of files
LAYOUT_FILES := $(shell find handler/frontend/layouts -type f -name '*.tmpl')
TEMPLATE_FILES := $(shell find handler/frontend/templates -type f -name '*.tmpl')

# Targets
all: build

build: static/bindata.go handler/frontend/layouts/bindata.go handler/frontend/templates/bindata.go
	env GO15VENDOREXPERIMENT=1 go build \
		-o $(NAME) \
		-v \
		-ldflags "-X main.revision=$(SHA) -X main.version=$(VERSION)" \
		.

static/bindata.go:
	go-bindata \
		$(BINDATA_FLAGS) \
		-ignore='(\.gitignore$$|\.map$$)' \
		-prefix=$(dir $@) \
		-pkg=static \
		-o $@ \
		$(dir $@)

handler/frontend/layouts/bindata.go: $(LAYOUT_FILES)
	go-bindata \
		$(BINDATA_FLAGS) \
		-ignore='(\.gitignore$$|\.go$$)' \
		-prefix=$(dir $@) \
		-pkg=layouts \
		-o $@ \
		$(dir $@)

handler/frontend/templates/bindata.go: $(TEMPLATE_FILES)
	go-bindata \
		$(BINDATA_FLAGS) \
		-ignore='(\.gitignore$$|\.go$$)' \
		-prefix=$(dir $@) \
		-pkg=templates \
		-o $@ \
		$(dir $@)

clean:
	$(RM) \
		./$(NAME) \
		static/bindata.go \
		handler/frontend/layouts/bindata.go \
		handler/frontend/templates/bindata.go
