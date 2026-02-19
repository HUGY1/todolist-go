# Todo List 后端 API

基于 Gin 和 GORM 框架的待办事项 RESTful API 服务，使用 MySQL 存储数据。

## 项目结构

```
.
├── main.go              # 程序入口，初始化数据库、路由并启动服务
├── go.mod               # Go 模块定义
├── config/
│   └── database.go      # 数据库连接配置
├── models/
│   └── todo.go          # Todo 数据模型
├── handlers/
│   └── todo_handler.go  # API 请求处理函数
└── routes/
    └── routes.go        # 路由配置
```

## 数据库配置

- **主机**: test-db-mysql.ns-wzme3ot2.svc
- **端口**: 3306
- **用户**: root
- **密码**: lgzxp6qg
- **数据库**: todolist（需提前创建）
- **表名**: list（程序启动时自动创建）

## 安装与运行

### 1. 创建数据库

在 MySQL 中执行：

```sql
CREATE DATABASE IF NOT EXISTS todolist CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 2. 下载依赖

```bash
go mod tidy
```

### 3. 启动服务

```bash
go run main.go
```

服务默认监听 `http://localhost:8080`。

## API 接口

### 数据模型

**Todo 对象结构**:
```json
{
  "id": 1,
  "value": "待办事项内容",
  "isCompleted": false,
  "createdAt": "2026-02-18T06:30:05.166Z",
  "updatedAt": "2026-02-18T06:30:05.166Z"
}
```

**字段说明**:
- `id` (uint): 待办事项唯一标识，自增主键
- `value` (string): 待办事项内容，最大长度 500 字符
- `isCompleted` (boolean): 是否完成，默认 false
- `createdAt` (string): 创建时间，ISO 8601 格式
- `updatedAt` (string): 更新时间，ISO 8601 格式

---

### 1. 查询所有待办事项

- **接口**: `POST /api/get-todo`
- **请求方法**: POST
- **参数**: 无
- **请求头**: `Content-Type: application/json`（可选）

**成功响应** (HTTP 200):
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "value": "学习 Go",
      "isCompleted": false,
      "createdAt": "2026-02-18T06:30:05.166Z",
      "updatedAt": "2026-02-18T06:30:05.166Z"
    },
    {
      "id": 2,
      "value": "完成项目",
      "isCompleted": true,
      "createdAt": "2026-02-18T06:30:59.031Z",
      "updatedAt": "2026-02-18T06:31:10.123Z"
    }
  ]
}
```

**错误响应** (HTTP 500):
```json
{
  "success": false,
  "message": "查询待办事项失败",
  "error": "数据库错误信息"
}
```

**示例请求**:
```bash
curl -X POST http://localhost:8080/api/get-todo \
  -H "Content-Type: application/json"
```

---

### 2. 添加待办事项

- **接口**: `POST /api/add-todo`
- **请求方法**: POST
- **请求体**:
```json
{
  "value": "待办内容",
  "isCompleted": false
}
```

**请求参数说明**:
- `value` (string, 必填): 待办事项内容
- `isCompleted` (boolean, 可选): 是否完成，默认为 false

**成功响应** (HTTP 201):
```json
{
  "success": true,
  "data": {
    "id": 1,
    "value": "学习 Go",
    "isCompleted": false,
    "createdAt": "2026-02-18T06:30:05.166Z",
    "updatedAt": "2026-02-18T06:30:05.166Z"
  }
}
```

**错误响应** (HTTP 400 - 参数错误):
```json
{
  "success": false,
  "message": "请求参数无效，value 为必填项",
  "error": "Key: 'AddTodoRequest.Value' Error:Field validation for 'Value' failed on the 'required' tag"
}
```

**错误响应** (HTTP 500 - 服务器错误):
```json
{
  "success": false,
  "message": "添加待办事项失败",
  "error": "数据库错误信息"
}
```

**示例请求**:
```bash
curl -X POST http://localhost:8080/api/add-todo \
  -H "Content-Type: application/json" \
  -d '{"value":"学习 Go","isCompleted":false}'
```

---

### 3. 更新待办状态

- **接口**: `POST /api/update-todo/:id`
- **请求方法**: POST
- **路径参数**: 
  - `id` (uint, 必填): 待办事项的唯一标识
- **功能**: 将指定待办事项的 `isCompleted` 状态取反（true → false，false → true）

**成功响应** (HTTP 200):
```json
{
  "success": true,
  "data": {
    "id": 1,
    "value": "学习 Go",
    "isCompleted": true,
    "createdAt": "2026-02-18T06:30:05.166Z",
    "updatedAt": "2026-02-18T06:31:10.123Z"
  }
}
```

**错误响应** (HTTP 400 - 参数格式错误):
```json
{
  "success": false,
  "message": "id 必须是有效的数字",
  "error": "strconv.ParseUint: parsing \"abc\": invalid syntax"
}
```

**错误响应** (HTTP 404 - 资源不存在):
```json
{
  "success": false,
  "message": "待办事项不存在",
  "error": "record not found"
}
```

**错误响应** (HTTP 500 - 服务器错误):
```json
{
  "success": false,
  "message": "更新待办事项失败",
  "error": "数据库错误信息"
}
```

**示例请求**:
```bash
curl -X POST http://localhost:8080/api/update-todo/1 \
  -H "Content-Type: application/json"
```

---

### 4. 删除待办事项

- **接口**: `POST /api/del-todo/:id`
- **请求方法**: POST
- **路径参数**: 
  - `id` (uint, 必填): 待办事项的唯一标识

**成功响应** (HTTP 200):
```json
{
  "success": true,
  "message": "删除成功",
  "data": {
    "id": 1,
    "deleted": true
  }
}
```

**错误响应** (HTTP 400 - 参数格式错误):
```json
{
  "success": false,
  "message": "id 必须是有效的数字",
  "error": "strconv.ParseUint: parsing \"abc\": invalid syntax"
}
```

**错误响应** (HTTP 404 - 资源不存在):
```json
{
  "success": false,
  "message": "待办事项不存在"
}
```

**错误响应** (HTTP 500 - 服务器错误):
```json
{
  "success": false,
  "message": "删除待办事项失败",
  "error": "数据库错误信息"
}
```

**示例请求**:
```bash
curl -X POST http://localhost:8080/api/del-todo/1 \
  -H "Content-Type: application/json"
```

---

## 通用响应格式

所有接口遵循统一的响应格式：

**成功响应**:
```json
{
  "success": true,
  "data": { ... }  // 或 [] 数组
}
```

**错误响应**:
```json
{
  "success": false,
  "message": "错误描述信息",
  "error": "详细错误信息（可选）"
}
```

**HTTP 状态码**:
- `200 OK`: 请求成功（查询、更新、删除）
- `201 Created`: 资源创建成功（添加）
- `400 Bad Request`: 请求参数错误
- `404 Not Found`: 资源不存在
- `500 Internal Server Error`: 服务器内部错误
