# React 管理系统脚手架

一个基于 React + Vite + Ant Design + Go + Gin + SQLite 的现代化管理系统脚手架。

## 技术栈

### 前端
- **React 18** - UI 框架
- **Vite** - 构建工具
- **Ant Design 5** - UI 组件库
- **React Router 6** - 路由管理
- **Axios** - HTTP 客户端

### 后端
- **Go 1.21+** - 编程语言
- **Gin** - Web 框架
- **GORM** - ORM 框架
- **SQLite** - 数据库
- **JWT** - 身份认证

## 项目结构

```
react-mng2/
├── src/                    # 前端源码
│   ├── api/               # API 接口
│   ├── components/        # 公共组件
│   ├── context/          # React Context
│   ├── layouts/          # 布局组件
│   ├── pages/            # 页面组件
│   ├── router/           # 路由配置
│   ├── utils/            # 工具函数
│   ├── App.jsx           # 根组件
│   ├── App.css           # 全局样式
│   └── main.jsx          # 入口文件
├── backend/               # 后端源码
│   ├── api/              # API 路由
│   ├── config/           # 配置
│   ├── middleware/       # 中间件
│   ├── models/           # 数据模型
│   ├── services/         # 业务逻辑
│   ├── utils/            # 工具函数
│   └── main.go           # 入口文件
├── public/               # 静态资源
├── package.json          # 前端依赖
└── README.md            # 项目说明
```

## 快速开始

### 前端开发

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build
```

前端服务将运行在 http://localhost:5173

### 后端开发

```bash
# 进入后端目录
cd backend

# 安装依赖
go mod download

# 启动服务器
go run main.go
```

后端服务将运行在 http://localhost:8080

## 功能特性

### 已实现功能
- ✅ 用户登录/登出
- ✅ JWT 身份认证
- ✅ 路由权限控制
- ✅ 用户管理（CRUD）
- ✅ 角色管理（CRUD）
- ✅ 权限管理
- ✅ 响应式布局
- ✅ 统一错误处理
- ✅ API 请求拦截

### 待扩展功能
- ⬜ 文件上传
- ⬜ 数据导出
- ⬜ 操作日志
- ⬜ 系统设置

## 开发指南

### 前端开发规范

1. **组件命名**：使用 PascalCase
2. **文件命名**：组件文件使用 PascalCase，其他文件使用 camelCase
3. **样式管理**：使用 CSS Modules 或 Ant Design 的主题定制
4. **状态管理**：使用 React Context API
5. **API 调用**：统一使用 `src/api` 目录下的接口

### 后端开发规范

1. **目录结构**：按照分层架构组织代码
2. **错误处理**：统一使用自定义错误类型
3. **日志记录**：使用结构化日志
4. **数据验证**：使用 Gin 的验证器
5. **数据库操作**：使用 GORM 进行数据库操作

## API 文档

### 认证相关
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/logout` - 用户登出
- `GET /api/auth/profile` - 获取当前用户信息

### 用户管理
- `GET /api/users` - 获取用户列表
- `GET /api/users/:id` - 获取用户详情
- `POST /api/users` - 创建用户
- `PUT /api/users/:id` - 更新用户
- `DELETE /api/users/:id` - 删除用户

### 角色管理
- `GET /api/roles` - 获取角色列表
- `GET /api/roles/:id` - 获取角色详情
- `POST /api/roles` - 创建角色
- `PUT /api/roles/:id` - 更新角色
- `DELETE /api/roles/:id` - 删除角色

## 默认账号

- 用户名：`admin`
- 密码：`admin123`

## 许可证

MIT
