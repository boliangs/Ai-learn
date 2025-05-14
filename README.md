# AI模拟面试系统

基于Go语言开发的AI模拟面试系统，集成了DeepSeek API，支持Word文档解析和智能评分。

## 项目结构

```
.
├── cmd/                    # 主程序入口
│   └── main.go
├── config/                 # 配置文件
│   └── config.yaml
├── internal/              # 内部包
│   ├── controller/        # 控制器层
│   ├── service/          # 服务层
│   ├── repository/       # 数据访问层
│   └── model/            # 数据模型
├── pkg/                   # 公共包
│   ├── utils/            # 工具函数
│   └── middleware/       # 中间件
├── docs/                  # 文档
├── Dockerfile            # Docker构建文件
├── docker-compose.yml    # Docker编排文件
├── go.mod               # Go模块文件
└── README.md            # 项目说明文档
```

## 快速开始

1. 安装依赖
```bash
go mod download
```

2. 配置环境变量
```bash
cp config/config.example.yaml config/config.yaml
# 编辑 config.yaml 文件，填入必要的配置信息
```

3. 启动服务
```bash
go run cmd/main.go
```

## Docker部署

```bash
docker-compose up -d
```

## 主要功能

- 用户管理
- Word文档解析
- AI面试问答
- 智能评分
- 面试记录管理 