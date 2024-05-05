btags := 

yewc:
	go build -tags "${btags}"
	go install github.com/petersalex27/yew

.PHONY: clean yewc
clean: yew
	rm yew