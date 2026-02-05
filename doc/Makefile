# Makefile for Sphinx documentation
#
# Golang 实战开发指南

# You can set these variables from the command line, and also
# from the environment for the first two.
SPHINXOPTS    ?=
SPHINXBUILD   ?= sphinx-build
SOURCEDIR     = source
BUILDDIR      = build

# Put it first so that "make" without argument is like "make help".
help:
	@$(SPHINXBUILD) -M help "$(SOURCEDIR)" "$(BUILDDIR)" $(SPHINXOPTS) $(O)

.PHONY: help Makefile install livehtml clean-all serve

# Catch-all target: route all unknown targets to Sphinx using the new
# "make mode" option.  $(O) is meant as a shortcut for $(SPHINXOPTS).
%: Makefile
	@$(SPHINXBUILD) -M $@ "$(SOURCEDIR)" "$(BUILDDIR)" $(SPHINXOPTS) $(O)

# Install dependencies
install:
	pip install -r requirements.txt

# Live rebuild on changes
livehtml:
	sphinx-autobuild "$(SOURCEDIR)" "$(BUILDDIR)/html" $(SPHINXOPTS) $(O)

# Clean all build files
clean-all:
	rm -rf $(BUILDDIR)/*

# Serve the built documentation
serve:
	cd $(BUILDDIR)/html && python -m http.server 8000

# Build and serve
dev: html serve

# Quick build (faster, less strict)
quick:
	@$(SPHINXBUILD) -M html "$(SOURCEDIR)" "$(BUILDDIR)" -j auto $(SPHINXOPTS) $(O)
