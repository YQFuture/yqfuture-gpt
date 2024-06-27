## 根据.api文件生成代码
goctl api go -api desc/training.api -dir ./

## 根据.proto文件生成代码
goctl rpc protoc training.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=. --client=true -m