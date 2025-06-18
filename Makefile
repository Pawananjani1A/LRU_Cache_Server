include envinronments.env
export

BIN = ../../bin
PACKAGE = lruCache/poc
PROTO_DIR = src/internal/proto/
GRPC_SERVER_DIR = src/cmd/grpc/main.go
GRPC_SERVER_BUILD_DIR = src/build/grpc/main
GRPC_GEN_DIR = .

GIN_SERVER_DIR = src/cmd/gin/main.go
GIN_SERVER_BUILD_DIR = src/build/gin

#--
TAG=dev
#--

DIRECTORIES = $(wildcard poc/*)



grpc_generate_go:
	protoc -I${PROTO_DIR} --go_out=${GRPC_GEN_DIR} --go_opt=module=${PACKAGE} --go-grpc_out=${GRPC_GEN_DIR} --go-grpc_opt=module=${PACKAGE} ${PROTO_DIR}/*.proto

grpc_generate_proto:
	protoc -I${PROTO_DIR} --go_out=${GRPC_GEN_DIR} --go-grpc_out=${GRPC_GEN_DIR} ${PROTO_DIR}/*.proto

grpc_server: grpc_generate_go grpc_build_server

grpc_run_server: grpc_generate_proto grpc_build_server grpc_start_server

grpc_build_server:
	go build -tags ${TAG} -o ${GRPC_SERVER_BUILD_DIR} ./${GRPC_SERVER_DIR}

grpc_start_server:
	chmod +x ${GRPC_SERVER_BUILD_DIR}
	./${GRPC_SERVER_BUILD_DIR}

gin_build_server:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -tags "${TAG}" -o ${GIN_SERVER_BUILD_DIR}/gin_bin ./${GIN_SERVER_DIR}

gin_start_server:
	chmod +x ${GIN_SERVER_BUILD_DIR}
	./${GIN_SERVER_BUILD_DIR}/gin_bin

gin_run_server: gin_build_server gin_start_server