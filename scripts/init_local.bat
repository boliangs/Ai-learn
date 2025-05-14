@echo off

REM 创建必要的目录
if not exist uploads mkdir uploads

REM 创建MySQL数据库
mysql -u root -proot -e "CREATE DATABASE IF NOT EXISTS ai_interview CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"

REM 启动服务
go run cmd/main.go 