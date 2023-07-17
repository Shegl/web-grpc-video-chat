protos:
	protoc --go_out=src/chat --go_opt=paths=source_relative \
        --go-grpc_out=src/chat --go-grpc_opt=paths=source_relative chat.proto

protos-web:
	protoc \
		--ts_out=react_spa/src \
		--plugin=protoc-gen-ts=./react_spa/node_modules/.bin/protoc-gen-ts \
		-I . \
		  chat.proto