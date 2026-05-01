# 02 — 前端技术栈

> 完整的技术栈清单，包括版本号、用途说明、与后端 GoZero 的对接方式。

---

## 一、总览

```
Vue 3 + Vite + TypeScript     ← 核心三件套
        │
        ├── Vue Router 4                  ← 路由与导航守卫
        ├── Pinia                          ← 客户端状态（Token / User）
        ├── @tanstack/vue-query            ← 服务端状态（API 缓存 / 轮询）
        ├── PrimeVue 4 (Aura)              ← UI 组件库
        ├── Axios                          ← HTTP 客户端（拦截器 / Token 刷新）
        ├── VeeValidate + Zod              ← 表单校验
        ├── TailwindCSS + tailwindcss-primeui ← 页面样式 + 语义色统一
        ├── @vueuse/core                   ← 通用工具（轮询 / 存储 / 剪贴板）
        ├── qrcode                         ← 二维码生成
        ├── Vitest                         ← 单元测试
        └── Playwright                     ← E2E 测试
```

---

## 二、分层明细

### 2.1 核心框架

| 技术 | 版本 | 用途 | 选择理由 |
|------|------|------|----------|
| Vue | 3.5+ | UI 框架 | 国语社区活跃，模板语法简洁 |
| TypeScript | 5.x | 类型安全 | 与 GoZero API 响应类型对齐，消除运行时类型错误 |
| Vite | 5.x | 构建工具 | 秒级 HMR，纯静态产物，与 GoZero embed 完美对接 |

### 2.2 路由

| 技术 | 版本 | 用途 |
|------|------|------|
| Vue Router | 4.x | SPA 页面跳转、路由守卫（未登录自动跳转登录页）、路由懒加载 |

### 2.3 状态管理

| 技术 | 版本 | 用途 | 职责边界 |
|------|------|------|----------|
| Pinia | 2.x | 客户端状态 | accessToken、refreshToken、user、isAuthenticated |
| @tanstack/vue-query | 5.x | 服务端状态 | Application 列表、Client 列表、Message 列表、Client 扫码状态 |

**为什么要分开？**

- Pinia 管 `浏览器里独有的状态` —— 登录后写入，登出后清空
- Vue Query 管 `以服务端为准的状态` —— 增删改后自动失效缓存并重新获取

### 2.4 UI 组件库

| 技术 | 版本 | 主题 | 用途 |
|------|------|------|------|
| PrimeVue | 4.x | Aura | DataTable、Dialog、Form 组件、Layout、Toast、Button 等 80+ 组件 |

### 2.5 HTTP 客户端

| 技术 | 版本 | 用途 |
|------|------|------|
| Axios | 1.x | 请求拦截器附加 Token、响应拦截器处理 401 → 自动 Refresh → 重试 |

### 2.6 表单校验

| 技术 | 版本 | 用途 |
|------|------|------|
| VeeValidate | 4.x | 表单状态管理、错误展示 |
| Zod | 3.x | Schema 校验规则定义 + TypeScript 类型推导 |

### 2.7 样式

| 技术 | 版本 | 用途 |
|------|------|------|
| TailwindCSS | 3.x 或 4.x | 页面级布局（flex/grid/gap/padding）、响应式 |
| tailwindcss-primeui | - | 将 PrimeVue 主题色映射为 Tailwind 语义 class |

### 2.8 工具库

| 技术 | 版本 | 用途 |
|------|------|------|
| @vueuse/core | 10.x | `useIntervalFn`（扫码轮询）、`useStorage`（Token 持久化）、`useClipboard`（复制 Token） |
| qrcode | 1.x | 将二维码 URL 渲染为 Canvas（Client 创建后展示） |

### 2.9 测试

| 技术 | 版本 | 用途 |
|------|------|------|
| Vitest | 1.x | 单元测试（与 Vite 同生态，零配置） |
| @vue/test-utils | 2.x | Vue 组件测试工具 |
| Playwright | 1.x | E2E 测试（关键用户路径） |

### 2.10 代码规范

| 技术 | 用途 |
|------|------|
| ESLint | 代码质量检查 |
| Prettier | 代码格式化 |
| Husky + lint-staged | Git 提交前自动检查和格式化 |

---

## 三、与 GoZero 后端的对接点

| 前端层 | 后端对应 | 对接方式 |
|--------|----------|----------|
| Axios 拦截器 | GoZero ClientAuth / ApplicationAuth 中间件 | Bearer Token 头 |
| Service 层 API | GoZero Handler 路由 | RESTful 一对一映射 |
| Zod Schema 校验规则 | GoZero `xerr` 错误码体系 | 前端预校验 + 后端错误码映射 |
| Vite build outDir | GoZero `//go:embed static/*` | 单二进制部署 |
| Vite dev proxy | GoZero `:port` | 开发时代理到本机后端 |

---

## 四、API 路由映射表

> 前端 Service 层与 GoZero Handler 的对应关系

### 公开路由（无鉴权）

| 前端 Service 方法 | HTTP | GoZero Route | Handler |
|-------------------|------|-------------|---------|
| `authService.register()` | POST | `/api/auth/register` | Register |
| `authService.login()` | POST | `/api/auth/login` | Login |
| `authService.refresh()` | POST | `/api/auth/refresh` | Refresh |

### 管理后台路由（ClientAuth 鉴权）

| 前端 Service 方法 | HTTP | GoZero Route | Handler |
|-------------------|------|-------------|---------|
| `userService.getProfile()` | GET | `/api/user/profile` | GetProfile |
| `userService.changePassword()` | POST | `/api/user/password` | ChangePassword |
| `clientService.create()` | POST | `/api/client/create` | CreateClient |
| `clientService.pollStatus()` | GET | `/api/client/status` | PollClientStatus |
| `clientService.list()` | GET | `/api/client/list` | ListClients |
| `clientService.delete()` | DELETE | `/api/client/:id` | DeleteClient |
| `appService.create()` | POST | `/api/application/create` | CreateApplication |
| `appService.list()` | GET | `/api/application/list` | ListApplications |
| `appService.delete()` | DELETE | `/api/application/:id` | DeleteApplication |
| `msgService.list()` | GET | `/api/message/list` | ListMessages |
| `msgService.delete()` | DELETE | `/api/message/:id` | DeleteMessage |

### 推送 API（ApplicationAuth 鉴权）

| 外部系统调用 | HTTP | GoZero Route | Handler |
|-------------|------|-------------|---------|
| MessageService | POST | `/api/message/push` | CreateMessage |

### 健康检查

| HTTP | GoZero Route | Handler |
|------|-------------|---------|
| GET | `/api/health` | Health |

---

## 五、GoZero 错误码 → 前端映射

> 来源：`pkg/xerr/errors_code.go`

| 错误码 | 常量 | 含义 | 前端处理 |
|--------|------|------|----------|
| `0` | `Success` | 请求成功 | 正常返回 |
| `100000` | `RequestParamError` | 参数错误 | 提示用户检查输入 |
| `100001` | `JwtError` | Unauthorized | 触发 Token 刷新流程 |
| `110100` | `RegisterParamEmpty` | 用户/密码为空 | 表单校验提示 |
| `110101` | `RegisterUsernameExist` | 用户名已存在 | 提示换用户名 |
| `110107` | `RegisterClosed` | 注册已关闭 | 提示 V1 已满 |
| `110202` | `LoginPasswordWrong` | 密码错误 | 提示重新输入 |
| `110300` | `RefreshTokenInvalid` | 刷新令牌无效 | 清除 Token 跳登录 |
| `110301` | `RefreshTokenRevoked` | 令牌已失效 | 清除 Token 跳登录 |
| `110302` | `RefreshTokenExpired` | 令牌已过期 | 清除 Token 跳登录 |
| `110400` | `PasswordTooWeak` | 密码强度不足 | 提示规则（≥8位+大小写+数字） |
