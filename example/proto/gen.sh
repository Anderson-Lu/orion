protoc --go_out=../proto_go --go_opt=paths=source_relative --go-grpc_out=../proto_go --go-grpc_opt=paths=source_relative --grpc-gateway_out=../proto_go --grpc-gateway_opt=paths=source_relative --go-grpc_out=require_unimplemented_servers=true todo/*.proto


