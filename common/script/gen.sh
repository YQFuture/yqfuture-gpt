## 根据.api文件生成代码
goctl api go -api desc/training.api -dir ./

## 根据.proto文件生成代码
goctl rpc protoc training.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=. --client=true -m

## 通过指定mysql连接地址来生成model代码
goctl model mysql datasource -url="root:mysql_QspKnh@tcp(192.168.3.118:3306)/qyfuture_gpt" -table="ts*" -dir ./model/orm