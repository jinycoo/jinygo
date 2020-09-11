# Jinygo

a golang micro web framework

## 项目简介

> database - mysql elasticsearch 

> cache - redis 

> message queue - rabbitmq 

> log

> web http

> grpc

> tools 工具部分 - 生成项目脚手架

```shell script

go get github.com/jinycoo/jinygo/tools/jiny

jiny new project_name -o auther -m module -p project_path

cd project_path/project_name

vim ./bin/app.toml

go run ./cmd/main.go -conf_path ./bin/app.toml

```
