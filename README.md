# Chatbot: For Chatting!

Made with ❤️ by The Ubyssey

## Running It

Start an instance of MongoDB on your local machine. See `CreateConnection()` in `models/model.go` to see where the Chatbot will dial.

	cd $GOPATH/src/github.com/ubyssey/chatbot
	go build
	./chatbot

The Chatbot will be running on your local machine bound to port 3000.

## Things Not To Do (Yet)

* Run it. The project does not have any functional security features, so putting a running instance of it on the web 