# 后端 API 规范

> 记录 WeClawBotNotify 后端 API 的设计规范，包括错误处理、响应格式、鉴权机制等。

---

## 一、统一响应格式

所有 API 返回统一使用 `pkg/result/response.go` 中定义的格式：

```json
{
  "code": 0,
  "msg": "请求成功",
  "data": { ... }
}
```

### 字段说明

| 字段 | 类型 | 说明 |
|------|------|------|
| `code` | int | 业务错误码，`0` 表示成功，非 `0` 表示失败 |
| `msg` | string | 响应消息，成功时为 "请求成功"，失败时为错误描述 |
| `data` | any | 响应数据，成功时返回业务数据，失败时可选 |

---

## 二、HTTP 状态码规范

### 2.1 基本原则

**除 JWT 鉴权失败外，所有业务错误均返回 HTTP 200，通过 `code` 字段区分错误类型。**

### 2.2 状态码映射

| HTTP 状态码 | 触发条件 | 说明 |
|-------------|----------|------|
| **200** | 业务请求成功 | 包含 `code=0` 的成功响应 |
| **200** | 业务逻辑错误 | 包含非 `0` 的 `code`，如 `110504` 应用不存在 |
| **401** | JWT 鉴权失败 | Access Token 过期、无效或签名错误 |

### 2.3 实现位置

```go
// pkg/result/response.go
func ErrorHandler(ctx context.Context, err error) (int, any) {
    httpCode := http.StatusOK
    response := &Response{}
    switch e := err.(type) {
    case *xerr.CodeError:
        if errors.Is(e, xerr.JwtError) {
            httpCode = http.StatusUnauthorized  // 仅 JWT 错误返回 401
        } else {
            httpCode = http.StatusOK            // 业务错误统一返回 200
        }
        response = &Response{
            Code: e.Code,
            Msg:  e.Message,
        }
    default:
        httpCode = http.StatusOK
        response = &Response{
            Code: -1,
            Msg:  err.Error(),
        }
    }
    return httpCode, response
}
```

### 2.4 设计理由

1. **前端统一处理**：Axios 响应拦截器只需处理 401（触发 Token 刷新），其他错误通过 `code` 判断
2. **简化错误流**：避免 HTTP 状态码与业务错误码双重判断
3. **保持一致性**：所有业务接口行为统一，降低调试成本

---

## 三、错误码规范

### 3.1 错误码定义

错误码定义在 `pkg/xerr/errors_code.go`，按模块划分：

| 错误码范围 | 模块 | 说明 |
|------------|------|------|
| `0` | 通用 | 请求成功 |
| `100000` ~ `100099` | 通用错误 | 参数错误、JWT 错误等 |
| `110100` ~ `110199` | 注册模块 | 注册相关错误 |
| `110200` ~ `110299` | 登录模块 | 登录相关错误 |
| `110300` ~ `110399` | 刷新令牌模块 | Token 刷新相关错误 |
| `110400` ~ `110499` | 密码管理模块 | 密码相关错误 |
| `110500` ~ `110599` | Application 模块 | 应用管理相关错误 |
| `110600` ~ `110699` | Message 模块 | 消息推送相关错误 |
| `200000` ~ `200099` | 第三方服务 | 第三方接口调用错误 |

### 3.2 错误码定义示例

```go
// pkg/xerr/errors_code.go

// Application 管理错误 1105XX
var (
    ApplicationParamNameEmpty = NewCodeError(110500, "应用名称不能为空")
    ApplicationInsertFailed   = NewCodeError(110501, "创建应用失败")
    ApplicationGetIdFailed    = NewCodeError(110502, "获取应用ID失败")
    ApplicationTokenFailed    = NewCodeError(110503, "生成应用Token失败")
    ApplicationNotFound       = NewCodeError(110504, "应用不存在")
    ApplicationNoPermission   = NewCodeError(110505, "无权操作该应用")
    ApplicationQueryFailed    = NewCodeError(110506, "查询应用失败")
    ApplicationDeleteFailed   = NewCodeError(110507, "删除应用失败")
)
```

### 3.3 错误码使用规范

1. **按模块分配范围**：每个模块分配 100 个错误码
2. **语义化命名**：错误码变量名应清晰表达错误含义
3. **中文消息**：错误消息使用中文，便于前端展示和日志排查
4. **不要复用错误码**：不同场景使用不同错误码，即使消息相同

---

## 四、鉴权机制

### 4.1 双令牌认证

| Token 类型 | 有效期 | 用途 | 路由分组 |
|------------|--------|------|----------|
| Access Token | 15 分钟 | 管理后台 API 鉴权 | `ClientAuth` 中间件 |
| Refresh Token | 7 天 | 刷新 Access Token | `/api/auth/refresh` |
| Application Token | 永久（直到删除） | 推送 API 鉴权 | `ApplicationAuth` 中间件 |

### 4.2 路由鉴权规则

| 路由前缀 | 鉴权方式 | 说明 |
|----------|----------|------|
| `/api/auth/*` | 无 | 公开路由，注册/登录/刷新 |
| `/api/user/*` | ClientAuth (JWT) | 用户管理，需 Access Token |
| `/api/client/*` | ClientAuth (JWT) | 客户端管理，需 Access Token |
| `/api/application/*` | ClientAuth (JWT) | 应用管理，需 Access Token |
| `/api/message/list` | ClientAuth (JWT) | 查看消息，需 Access Token |
| `/api/message/push` | ApplicationAuth | 推送消息，需 Application Token |
| `/api/health` | 无 | 健康检查，公开 |

### 4.3 Application Token 鉴权规则

Application Token 是应用推送消息的凭证，通过 **JSON Body** 传递：

```
POST /api/message/push
Content-Type: application/json

{
  "app_token": "xxx",
  "title": "CI Build",
  "message": "Build #123 passed"
}
```

**设计理由**：
1. POST 请求放在 body 中，HTTP/2 可压缩
2. token 不会出现在 URL、浏览器历史、代理日志中，安全性更好
3. 鉴权逻辑在 Logic 层实现（需查数据库），不使用 middleware

**实现位置**：`internal/logic/message/createMessageLogic.go`

```go
func (l *CreateMessageLogic) CreateMessage(req *types.CreateMessageReq) (resp *types.CreateMessageResp, err error) {
    // 1. 验证 app_token
    app, err := l.svcCtx.ApplicationsModel.ValidateAppToken(l.ctx, req.AppToken)
    if err != nil {
        return nil, xerr.ApplicationTokenInvalid
    }
    
    // 2. 更新 last_used_at
    // 3. 查询可用 Client
    // 4. 插入消息并推送
}
```

### 4.4 Application Token 软删除规则

Application 支持软删除（`status = -1`），软删除后：

1. **Token 立即失效**：`ApplicationAuth` 中间件验证 Token 时检查 `status`
2. **不删除关联数据**：Messages 表数据保留
3. **前端查询过滤**：列表查询时过滤 `status != -1` 的记录

```go
// ApplicationAuth 中间件验证逻辑（已废弃，改为 Logic 层验证）
token := r.URL.Query().Get("token")
if app.Status == -1 || app.Status == 2 {
    httpx.ErrorCtx(r.Context(), w, xerr.ApplicationTokenInvalid)
    return
}
```

> 注：Application Token 验证已移至 Logic 层实现，不再使用 middleware。

---

## 五、Service 层规范

### 5.1 Logic 文件结构

```go
package application

import (
    "context"
    // ... 其他导入
)

type CreateApplicationLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
}

func NewCreateApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateApplicationLogic {
    return &CreateApplicationLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
    }
}

func (l *CreateApplicationLogic) CreateApplication(req *types.CreateApplicationReq) (*types.CreateApplicationResp, error) {
    // 1. 参数校验
    // 2. 业务逻辑
    // 3. 数据库操作
    // 4. 返回结果
}
```

### 5.2 错误处理规范

```go
// 正确示例：使用预定义错误码
if req.Name == "" {
    return nil, xerr.ApplicationParamNameEmpty
}

// 错误示例：不要直接返回 error 字符串
if req.Name == "" {
    return nil, errors.New("应用名称不能为空")  // 不要这样做
}
```

### 5.3 日志规范

```go
// 请求入口记录 Info 级别日志
l.Infof("创建应用请求: name=%s", req.Name)

// 错误记录 Error 级别日志
l.Errorf("创建应用失败: %v", err)

// 成功记录 Info 级别日志
l.Infof("应用创建成功: appId=%d, name=%s, userId=%d", appId, req.Name, userId)
```

---

## 六、Model 层规范

### 6.1 goctl 生成的代码

- `*_gen.go` 文件由 goctl 自动生成，**不要手动修改**
- 如需扩展方法，在对应的自定义文件中添加

### 6.2 自定义 Model 文件

```go
// applicationsmodel.go
type ApplicationsModel interface {
    applicationsModel  // 继承 goctl 生成的接口
    withSession(session sqlx.Session) ApplicationsModel
    FindByUserId(ctx context.Context, userId int64) ([]*Applications, error)
    // ... 自定义方法
}
```

### 6.3 软删除查询规范

对于支持软删除的表，所有查询方法必须过滤已删除的记录：

```go
// 正确示例：查询时过滤 status
query := fmt.Sprintf("select %s from %s where `user_id` = ? AND `status` != -1", applicationsRows, m.table)

// 错误示例：忘记过滤
query := fmt.Sprintf("select %s from %s where `user_id` = ?", applicationsRows, m.table)
```

---

## 七、状态字段设计规范

### 7.1 Go 语言状态字段避免 0 值原则

**在 Go 语言中，结构体字段的默认值为 0。因此状态字段应避免使用 0 作为有效状态值，以防止未初始化或默认值导致的歧义。**

### 7.2 设计原则

1. **不使用 0 作为有效状态**：0 是 Go 的零值，应保留为"未设置"或"未知"状态
2. **使用正数表示正常状态**：如 `1` 表示正常/启用
3. **使用负数表示终态**：如 `-1` 表示删除
4. **使用其他正数表示中间状态**：如 `2` 表示禁用/暂停

### 7.3 Application 状态值定义

| 状态值 | 含义 | 说明 |
|--------|------|------|
| `1` | 正常 | 应用可正常查询和推送消息 |
| `2` | 禁用 | 应用被临时停用，Token 失效，不可推送 |
| `-1` | 已删除 | 软删除，Token 失效，不可恢复 |

### 7.4 状态流转规则

```
正常(1) ←→ 禁用(2)  : 可互相切换
  ↓
已删除(-1)          : 终态，不可恢复
```

- **正常 → 禁用**：管理员手动禁用应用
- **禁用 → 正常**：管理员重新启用应用
- **正常 → 已删除**：软删除应用
- **禁用 → 已删除**：软删除应用

### 7.5 各状态下的行为

| 行为 | 正常(1) | 禁用(2) | 已删除(-1) |
|------|---------|---------|------------|
| Token 推送 | ✅ 允许 | ❌ 拒绝 | ❌ 拒绝 |
| 前端列表展示 | ✅ 展示 | ✅ 展示（标记禁用） | ❌ 不展示 |
| 可恢复 | - | ✅ 可启用 | ❌ 不可恢复 |

### 7.6 Message 状态码定义

| 状态码 | 含义 | 说明 |
|--------|------|------|
| `100` | 待推送（pending） | 消息已创建，等待推送 |
| `200` | 推送成功（sent） | 所有 Client 推送成功 |
| `300` | 推送失败（failed） | 所有 Client 推送失败 |
| `301` | 失败-微信API错误 | 微信 iLink 接口返回错误 |
| `302` | 失败-超时 | 推送请求超时 |
| `303` | 失败-Client已过期 | Client session 已过期 |
| `304` | 失败-部分成功 | 部分 Client 成功，部分失败 |
| `-1` | 已删除（软删除） | 软删除 |

**设计规则**：
- `1XX` - 初始/中间状态
- `2XX` - 成功状态
- `3XX` - 失败状态（细分失败原因）
- `-1` - 软删除（终态）

**无可用 Client 时不插入消息，直接返回错误码 110601**

### 7.7 代码实现示例

```go
// ApplicationAuth 中间件验证逻辑
if app.Status == -1 || app.Status == 2 {
    httpx.ErrorCtx(r.Context(), w, xerr.JwtError)  // 返回 401
    return
}

// Model 查询逻辑 - 过滤已删除的记录
query := fmt.Sprintf("select %s from %s where `user_id` = ? AND `status` != -1", applicationsRows, m.table)
```

---

## 版本历史

| 日期 | 版本 | 变更 |
|------|------|------|
| 2026-05-01 | 1.0 | 初始版本，记录 API 响应格式、错误码、鉴权等规范 |
| 2026-05-01 | 1.1 | 添加 Go 状态字段避免 0 值原则，Application 增加禁用状态 |
| 2026-05-01 | 1.2 | Application Token 改为 Query 参数传递，添加 Message 状态码定义 |
| 2026-05-01 | 1.3 | Application Token 改为 JSON Body 传递，鉴权移至 Logic 层 |
