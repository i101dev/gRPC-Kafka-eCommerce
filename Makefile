proto:
	#export PATH="$PATH:$(go env GOPATH)/bin"
	@protoc --go_out=. --go-grpc_out=. proto/inventory.proto
	@protoc --go_out=. --go-grpc_out=. proto/order.proto
	@protoc --go_out=. --go-grpc_out=. proto/product.proto
	@protoc --go_out=. --go-grpc_out=. proto/user.proto

purge:
	@rm -r proto/*.pb.go

clean:
	docker system prune


.PHONY: proto