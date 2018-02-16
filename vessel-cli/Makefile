build:  
	GOOS=linux GOARCH=amd64 go build
	docker build -t vessel-cli .

run:  
	docker run -e MICRO_REGISTRY=mdns vessel-cli