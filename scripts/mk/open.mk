##
# This makefile define OPEN variable depending on the first available option.
#
# Variables:
#   OPEN
##

ifneq (,$(shell which open 2>/dev/null))
OPEN ?= open
endif
ifneq (,$(shell which xdg-open 2>/dev/null))
OPEN ?= xdg-open
endif
ifeq (,$(OPEN))
OPEN ?= false
endif

