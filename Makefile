yew:
	go build
	go install github.com/petersalex27/yew

debug:
	go build -tags debug
	go install github.com/petersalex27/yew

clean:
	rm yew