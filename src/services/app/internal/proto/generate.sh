mkdir -p ../../../../packages/varsoapi/src/app
echo src/services/app/internal/proto: Installing protoc-gen-go
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
protoc -I=./ --plugin=$(go env GOBIN)protoc-gen-go --plugin=../../../../common/proto/node_modules/.bin/protoc-gen-ts_proto --ts_proto_opt=esModuleInterop=true --go_out=../../ --ts_proto_out=../../../../packages/varsoapi/src/app appbase.proto apprequests.proto

mkdir -p ../../../www/src/proto/src/app
cp -R ../../../../packages/varsoapi/src/app ../../../www/src/proto/src/