include $(GOROOT)/src/Make.inc

TARG=github.com/monnand/gotaskqueue
GOFILES=\
	taskqueue.go\
	decorators.go\

include $(GOROOT)/src/Make.pkg
