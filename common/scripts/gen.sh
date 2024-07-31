## 根据.api文件生成代码
goctl api go -api desc/training.api -dir ./

## 根据.proto文件生成代码
goctl rpc protoc pb/org.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=. --client=true -m

## 通过指定mysql连接地址来生成model代码
goctl model mysql datasource -url="root:mysql_QspKnh@tcp(10.0.0.118:3306)/yqfuture_gpt" -table="ts*,bs*" -dir ./model/orm

## 通过指定mongodb集合名称来生成model代码
goctl model mongo --type dbsavegoodscrawlertitles --dir ./model/mongo
goctl model mongo --type dbpresettingshoptitles --dir ./model/mongo
goctl model mongo --type dborgpermission --dir ./model/mongo

## 生成swagger文档
goctl api plugin -plugin goctl-swagger="swagger -filename training.json" -api desc/training.api -dir .