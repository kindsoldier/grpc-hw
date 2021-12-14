
GENDIR=	.

all: go

java:
	mkdir -p ${GENDIR} 
	protoc --proto_path=proto --java_out=${GENDIR}/java proto/*.proto


go:
	mkdir -p ${GENDIR}/services
	protoc --proto_path=proto \
	    --go_out=${GENDIR}/services --go_opt=paths=source_relative \
	    --go-grpc_out=${GENDIR}/services --go-grpc_opt=paths=source_relative \
	    proto/services.proto


cpp:
	mkdir -p ${GENDIR}/cpp
	protoc --proto_path=proto --cpp_out=${GENDIR}/cpp proto/*.proto

serve:
	go run server.go
