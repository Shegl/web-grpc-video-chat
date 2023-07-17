protos:
	protoc --go_out=src/chat --go_opt=paths=source_relative \
        --go-grpc_out=src/chat --go-grpc_opt=paths=source_relative chat.proto

protos-web:
	npx protoc --plugin=./react_spa/node_modules/.bin/protoc-gen-ts_proto \
		--ts_proto_opt=outputClientImpl=grpc-web \
		--ts_proto_out=react_spa/src chat.proto