.PHONY: build install clean

build:
	go build -o claude-bell .

install: build
	sudo cp claude-bell /usr/local/bin/claude-bell

clean:
	rm -f claude-bell
