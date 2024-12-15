1. загружаем необходимые пакеты и исполняемые файлы
go get google.golang.org/protobuf/cmd/protoc-gen-go 
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
go install google.golang.org/protobuf/cmd/protoc-gen-go 
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc

2. серрелиализуем структуры 
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --experimental_allow_proto3_optional exchange/*.proto