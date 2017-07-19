GOOS ?= linux
GOARCH ?= amd64

all: motd/motd lambda/expirer/handler.zip lambda/initiator/handler.zip

motd/motd:
	env GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ./motd/motd ./motd/...

lambda/expirer/handler.zip:
	cd lambda/expirer/ && make

lambda/initiator/handler.zip:
	cd lambda/initiator/ && make

clean:
	rm -f lambda/expirer/handler.*
	rm -f lambda/initiator/handler.*
	rm -f motd/motd
