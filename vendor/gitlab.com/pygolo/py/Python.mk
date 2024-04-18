PYTHON ?= python3

PYTHON_CONFIG_VAR = $(shell $(PYTHON) -c "import sysconfig; print(sysconfig.get_config_var('$(1)') or '')")
PYTHON_VERSION_SHORT := $(call PYTHON_CONFIG_VAR,py_version_short)
PYTHON_LIBDIR := $(call PYTHON_CONFIG_VAR,LIBDIR)
PYTHON_LIBPC := $(call PYTHON_CONFIG_VAR,LIBPC)

ifneq ($(PYGOLO_NOTAGS),true)
	PYGOLO_TAGS += py$(PYTHON_VERSION_SHORT)
endif

PKG_CONFIG_PATH := $(PYTHON_LIBPC):$(PKG_CONFIG_PATH)
export PKG_CONFIG_PATH

CGO_ENABLED := 1
export CGO_ENABLED

ifeq ($(shell go env GOOS),darwin)
	DYLD_LIBRARY_PATH := $(PYTHON_LIBDIR):$(DYLD_LIBRARY_PATH)
	export DYLD_LIBRARY_PATH

	PYTHON_PYTHONFRAMEWORKPREFIX := $(call PYTHON_CONFIG_VAR,PYTHONFRAMEWORKPREFIX)
	ifneq ($(PYTHON_PYTHONFRAMEWORKPREFIX),)
		CGO_LDFLAGS := -Wl,-rpath,$(PYTHON_PYTHONFRAMEWORKPREFIX)
		export CGO_LDFLAGS
	endif
else
	LD_LIBRARY_PATH := $(PYTHON_LIBDIR):$(LD_LIBRARY_PATH)
	export LD_LIBRARY_PATH
endif

ifeq ($(PYTHON_LIBPC),)
define embed-python
	$(error Embedding Python is not supported on this platform)
endef
define extend-python
	$(error Extending Python is not supported on this platform)
endef
else
define extend-python
	$(eval $@: PYGOLO_TAGS += py_ext)
	$(eval $@: PYGOLO_FLAGS += -buildmode=c-shared)
endef
endif

pygolo-diags:
	@go list -f {{.Version}} -m gitlab.com/pygolo/py 2>/dev/null || true
	@$(SHELL) -c command -v go
	@go version
	@$(SHELL) -c command -v $(PYTHON)
	@$(PYTHON) -V
	@echo PYTHON: $(PYTHON)
	@echo PYTHON_LIBDIR: $(PYTHON_LIBDIR)
ifeq ($(shell go env GOOS),darwin)
	@echo DYLD_LIBRARY_PATH: $(DYLD_LIBRARY_PATH)
	@echo PYTHON_PYTHONFRAMEWORKPREFIX: $(PYTHON_PYTHONFRAMEWORKPREFIX)
else
	@echo LD_LIBRARY_PATH: $(LD_LIBRARY_PATH)
endif
	@echo PYTHON_LIBPC: $(PYTHON_LIBPC)
	@echo PKG_CONFIG_PATH: $(PKG_CONFIG_PATH)
ifneq ($(PYTHON_LIBPC),)
	@cd $(PYTHON_LIBPC) && ls -l python*.pc
endif
ifeq ($(realpath $(PWD)),$(realpath $(CURDIR)))
	pyenv versions 2>/dev/null || true
else
	(cd $(PWD) && pyenv versions 2>/dev/null) || true
endif
	@echo "python-$(PYTHON_VERSION_SHORT)-embed.pc ->" \
		`pkg-config --debug python-$(PYTHON_VERSION_SHORT)-embed 2>&1 | grep -e ^Reading -e found: | awk '{print $$NF}' | xargs echo`
	@echo "python-$(PYTHON_VERSION_SHORT).pc ->" \
		`pkg-config --debug python-$(PYTHON_VERSION_SHORT) 2>&1 | grep -e ^Reading -e found: | awk '{print $$NF}' | xargs echo`
	@echo "python3-embed.pc ->" \
		`pkg-config --debug python3-embed 2>&1 | grep -e ^Reading -e found: | awk '{print $$NF}' | xargs echo`
	@echo "python3.pc ->" \
		`pkg-config --debug python3 2>&1 | grep -e ^Reading -e found: | awk '{print $$NF}' | xargs echo`

# prevent any pygolo-* target from becoming the default target
.DEFAULT_GOAL := $(filter-out pygolo-%,$(.DEFAULT_GOAL))
