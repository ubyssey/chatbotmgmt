# Chatbot: For Chatting!

Made with ❤️ (and ☕️) by The Ubyssey

## Running It

Start an instance of MongoDB on your local machine. See `CreateConnection()` in `models/model.go` to see where the Chatbot Management API will dial.

	cd $GOPATH/src/github.com/ubyssey/chatbotmgmt
	go build
	./chatbotmgmt

The Chatbot Management API will be running on your local machine bound to port 3000.

## Caveats

This service _does not_ do any actual chatting; rather, it maintains and exposes backing models (topics and campaigns), which may then be consumed by frontends. Currently, the only frontend is [chatbotfb](https://github.com/ubyssey/chatbotfb)

## Things Not To Do (Yet)

* Run it on the Internet. The project does not have any functional security features, so putting a running instance of it on the web would be poor form.