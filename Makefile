include $(GOROOT)/src/Make.inc

TARG=github.com/monnand/gotaskqueue
GOFILES=\
	taskqueue.go\

include $(GOROOT)/src/Make.pkg
