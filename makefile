jarmuzmessage: *.go
	go build -o jarmuzmessage .

run: jarmuzmessage
	./jarmuzmessage
