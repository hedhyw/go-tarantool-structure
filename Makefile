all: test
test:
	vgo test ./... -count=1 -covermode=count -v