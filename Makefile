CONFIG_PATH = $(shell pwd)/.generation/
# CONFIG_PATH=${HOME}/.generation/

.PHONY: init
init:
	mkdir -p ${CONFIG_PATH}

.PHONY: gencert
gencert:
	cfssl gencert \
		-initca test/ca-csr.json | cfssljson -bare ca
	cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=test/ca-config.json \
		-profile=client \
		test/client-csr.json | cfssljson -bare client
	cfssl gencert \
		-ca=ca.pem \
		-ca-key=ca-key.pem \
		-config=test/ca-config.json \
		-profile=server \
		test/server-csr.json | cfssljson -bare server
	mv *.pem *.csr ${CONFIG_PATH}
	# cfssl gencert \
	# 	-initca test/ca-csr.json | cfssljson -bare ca
	# cfssl gencert \
	# 	-ca=ca.pem \
	# 	-ca-key=ca-key.pem \
	# 	-config=test/ca-config.json \
	# 	-profile=server \
	# 	test/server-csr.json | cfssljson -bare server
	# mv *.pem *.csr ${CONFIG_PATH}

.PHONY: test
test:
	go test -race ./...
	
.PHONY: compile
compile:
	protoc ./api/v1/keyvalue/*.proto \
		--go_out=. \
		--go-grpc_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_opt=paths=source_relative \
		--proto_path=.
#	protoc ./v1/*.proto --go_out=. --go_opt=paths=source_relative --proto_path=.


