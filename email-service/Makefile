build:
	# protoc -I. --go_out=plugins=micro:$(GOPATH)/src/github.com/mooncaker816/shipper/consignment-service \
	# 	proto/consignment/consignment.proto
	GOOS=linux GOARCH=amd64 go build -o email-service -a -installsuffix cgo main.go
	docker build -t email-service .

run:
	docker run -p 50054:50051 -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=mdns email-service