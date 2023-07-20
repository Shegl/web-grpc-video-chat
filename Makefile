export CGO_ENABLED=0

protos:
	protoc --go_out=src/chat --go_opt=paths=source_relative \
        --go-grpc_out=src/chat --go-grpc_opt=paths=source_relative chat.proto \
        && protoc --go_out=src/streams --go_opt=paths=source_relative \
		   --go-grpc_out=src/streams --go-grpc_opt=paths=source_relative stream.proto


protos-web:
	protoc \
		--ts_out=react_spa/src \
		--plugin=protoc-gen-ts=./react_spa/node_modules/.bin/protoc-gen-ts \
		-I . chat.proto \
		&& protoc \
			--ts_out=react_spa/src \
			--plugin=protoc-gen-ts=./react_spa/node_modules/.bin/protoc-gen-ts \
			-I . stream.proto


gen-cert:
	cd docker/certs && ./gencert.sh

compose-up:
	docker compose up --attach app
