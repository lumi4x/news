mkdir -p ../../../../packages/varsoapi/src/auth
echo src/services/auth/internal/proto: Installing protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
protoc -I=./ --plugin=$(go env GOBIN)/protoc-gen-go --plugin=../../../../common/proto/node_modules/.bin/protoc-gen-ts_proto --go_out=../../ --ts_proto_out=../../../../packages/varsoapi/src/auth authbase.proto authrequests.proto

mkdir -p ../../../www/src/proto/src/auth
cp -R ../../../../packages/varsoapi/src/auth ../../../www/src/proto/src/