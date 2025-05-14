#!/bin/bash

# 创建必要的目录
mkdir -p uploads
chmod 755 uploads

# 创建MySQL数据库
mysql -u root -proot <<EOF
CREATE DATABASE IF NOT EXISTS ai_interview CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EOF

# 启动服务
go run cmd/main.go 