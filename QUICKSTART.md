# 快速启动测试指南

## 1. 检查服务状态

### 前端服务
```bash
# 检查前端是否运行
curl http://localhost:5173
```
应该看到 HTML 页面输出

### 后端服务
```bash
# 检查后端是否运行
curl http://localhost:8080/api/auth/login
```

## 2. 测试 API 接口

### 测试登录
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | jq .
```

成功响应示例：
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "token": "eyJhbGci...",
    "user": {
      "id": 1,
      "username": "admin",
      "realname": "管理员",
      "email": "admin@example.com"
    }
  }
}
```

### 测试获取用户列表（需要 token）
```bash
# 先获取 token
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin123"}' \
  | jq -r '.data.token')

# 使用 token 获取用户列表
curl -s http://localhost:8080/api/users \
  -H "Authorization: Bearer $TOKEN" \
  | jq .
```

### 测试创建用户
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "test123",
    "realname": "测试用户",
    "email": "test@example.com"
  }' \
  | jq .
```

### 测试获取角色列表
```bash
curl -s http://localhost:8080/api/roles \
  -H "Authorization: Bearer $TOKEN" \
  | jq .
```

## 3. 浏览器测试

访问 http://localhost:5173

### 测试流程
1. **登录测试**
   - 用户名：`admin`
   - 密码：`admin123`
   - 点击登录按钮

2. **仪表盘测试**
   - 登录成功后应该看到仪表盘页面
   - 显示统计卡片

3. **用户管理测试**
   - 点击左侧菜单"用户管理"
   - 应该看到用户列表
   - 测试新增用户
   - 测试编辑用户
   - 测试删除用户

4. **角色管理测试**
   - 点击左侧菜单"角色管理"
   - 应该看到角色列表
   - 测试新增角色
   - 测试编辑角色
   - 测试删除角色

5. **权限管理测试**
   - 点击左侧菜单"权限管理"
   - 应该看到开发中的提示页面

6. **退出登录测试**
   - 点击右上角用户头像
   - 点击"退出登录"
   - 应该返回登录页面

## 4. 常见问题排查

### 前端无法访问
```bash
# 检查 Vite 是否运行
ps aux | grep vite

# 重启前端
cd /Users/solariswu/WebstormProjects/react-mng2
npm run dev
```

### 后端无法访问
```bash
# 检查 Go 服务是否运行
ps aux | grep "go run"

# 重启后端
cd /Users/solariswu/WebstormProjects/react-mng2/backend
go run main.go
```

### API 返回 401 错误
- 检查 token 是否过期
- 重新登录获取新的 token

### 数据库问题
```bash
# 删除数据库文件重新初始化
cd /Users/solariswu/WebstormProjects/react-mng2/backend
rm data.db
go run main.go
```

## 5. 开发调试

### 查看前端日志
浏览器 F12 -> Console 标签

### 查看后端日志
后端运行的终端会显示所有请求日志

### 查看网络请求
浏览器 F12 -> Network 标签

## 6. 下一步学习

1. **前端学习路径**
   - React Hooks (useState, useEffect, useContext)
   - React Router (路由管理)
   - Ant Design 组件库
   - Axios 请求库
   - 表单处理和验证

2. **后端学习路径**
   - Gin 框架基础
   - GORM ORM 操作
   - JWT 认证机制
   - RESTful API 设计
   - 中间件使用

3. **推荐实践**
   - 阅读现有代码，理解项目结构
   - 尝试添加新功能
   - 修改现有功能
   - 添加新的路由和页面
