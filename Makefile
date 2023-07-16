protos:
	protoc --go_out=src/chat --go_opt=paths=source_relative \
        --go-grpc_out=src/chat --go-grpc_opt=paths=source_relative chat.proto

protos-web:
	npx protoc \
		--js_out=import_style=commonjs,binary:./react_spa/src --ts_out=./react_spa/src \
		--grpc-web_out=import_style=typescript,mode=grpcwebtext:react_spa/src \
		chat.proto

