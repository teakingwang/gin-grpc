PBDIR := ./pkg/pb
PROTODIR := ./internal/pkg/proto

.PHONY: proto
proto:
	protoc --go_out=plugins=grpc:$(PBDIR) $(PROTODIR)/*.proto
.PHONY: clean
clean:
	rm -rf $(PBDIR)/*.go
