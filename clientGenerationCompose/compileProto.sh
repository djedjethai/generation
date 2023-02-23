PROTO_DIR=/home/jerome/Documents/projects/generationProject/clientGeneration/pkg/proto/keyvalue 

# that generate only the protoBuf file(but does not compile the grpc)
# PACKAGE=github.com/djedjethai/generation0/pkg/handlers/grcp/proto/keyvalue
# protoc -I${PROTO_DIR} --go_opt=module=${PACKAGE} --go_out=${PROTO_DIR} ${PROTO_DIR}/*.proto
 
protoc --proto_path=$PROTO_DIR \
 --go_out=$PROTO_DIR --go_opt=paths=source_relative \
 --go-grpc_out=$PROTO_DIR --go-grpc_opt=paths=source_relative \
 $PROTO_DIR/*.proto

