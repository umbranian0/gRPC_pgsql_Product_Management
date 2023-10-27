# gRPC CRUD Operations in Postgres DB
base for project 
https://www.golinuxcloud.com/go-grpc-crud-api-postgresql-db/


generate gRPC files 
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/product.proto

start or stop DB 
pg_ctl -D ^"C^:^\pgsql^\pgsql^_data^" -l logfile start
pg_ctl -D ^"C^:^\pgsql^\pgsql^_data^" -l logfile stop

# Product Management System 
Backend system with services:
GET /products GetProducts
GET /product GetProduct(ID)
Post /product CreateProduct
Put /product UpdateProduct(id)
Delete /product DeleteProduct(id)