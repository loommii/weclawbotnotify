# 04 — 前后端交互策略

> 记录前端与 GoZero 后端的 API 调用方式、数据处理流程、鉴权机制实现方案。

---

## 一、统一响应格式

GoZero 后端所有 API 返回统一格式（来自 `pkg/result/response.go`）：

```json
{
  "code": 0,
  "msg": "请求成功",
  "data": { ... }
}
```

### 前端类型定义

```typescript
// types/api.ts
export interface ApiResponse<T = unknown> {
  code: number;
  msg: string;
  data: T;
}
```

### Axios 响应拦截器 — 剥离外层

```typescript
// lib/axios.ts
api.interceptors.response.use(
  (response) => {
    const { code, msg, data } = response.data as ApiResponse;
    if (code !== 0) {
      // 非 0 错误码：抛出业务异常
      return Promise.reject(new ApiError(code, msg));
    }
    return data; // ← 业务代码直接拿 data，无需 .data.data
  },
  async (error) => { /* 401 处理见下文 */ }
);
```

---

## 二、双令牌鉴权机制

> 对应 UML：`docs/uml/04-sequence-login.puml`、`docs/uml/11-sequence-refresh.puml`

### 2.1 两种 Token

| Token | 有效期 | 作用 | 存储位置 |
|-------|--------|------|----------|
| Access Token | 15 分钟 | 携带在请求头中鉴权 | Pinia + localStorage |
| Refresh Token | 7 天 | Access Token 过期后换取新双令牌（轮换） | localStorage |

### 2.2 请求拦截器 — 附加 Token

```
所有请求 → 拦截器读取 Pinia accessToken → 附加 Authorization: Bearer xxx
```

### 2.3 响应拦截器 — Token 刷新流程

```
请求发送
    │
    ▼
GoZero 返回 code=100001 (Unauthorized)
    │
    ▼
检查是否有 refreshToken
    │
    ├── 无 → 清除 Token → 跳转登录页
    │
    └── 有 → 检查是否正在刷新
              │
              ├── 是 → 加入等待队列
              │         └── 刷新完成后用新 Token 重试
              │
              └── 否 → 设置刷新锁
                        │
                        POST /api/auth/refresh
                        │
                        ├── 成功 → 更新 Pinia + localStorage
                        │          └── 释放队列（所有等待请求重试）
                        │
                        └── 失败 → 清除所有 Token
                                  └── 释放队列（所有等待请求失败）
                                  └── 跳转登录页
```

### 2.4 并发刷新保护实现

```typescript
// lib/axios.ts — 核心逻辑
let isRefreshing = false;
let failedQueue: Array<{
  resolve: (token: string) => void;
  reject: (err: unknown) => void;
}> = [];

function processQueue(error: unknown, token: string | null) {
  failedQueue.forEach(({ resolve, reject }) => {
    if (error) reject(error);
    else resolve(token!);
  });
  failedQueue = [];
}

api.interceptors.response.use(
  (res) => { /* code !== 0 处理 */ },
  async (error) => {
    const config = error.config;

    if (error.response?.data?.code === 100001 && !config._retry) {
      // 非 refresh 接口 → 进入刷新流程
      if (isRefreshing) {
        return new Promise((resolve, reject) => {
          failedQueue.push({ resolve, reject });
        }).then((token) => {
          config.headers.Authorization = `Bearer ${token}`;
          return api(config);
        });
      }

      config._retry = true;
      isRefreshing = true;

      const authStore = useAuthStore();
      const refreshToken = authStore.refreshToken;

      if (!refreshToken) {
        authStore.clearAuth();
        window.location.href = '/login';
        return Promise.reject(error);
      }

      try {
        const res = await authService.refresh(refreshToken);
        authStore.setTokens(res.token, res.refreshToken);
        processQueue(null, res.token);
        config.headers.Authorization = `Bearer ${res.token}`;
        return api(config);
      } catch (err) {
        processQueue(err, null);
        authStore.clearAuth();
        window.location.href = '/login';
        return Promise.reject(err);
      } finally {
        isRefreshing = false;
      }
    }

    return Promise.reject(error);
  }
);
```

---

## 三、Service 层设计

### 3.1 设计原则

- 每个 Service 文件对应 GoZero 一个 `@server group`
- 方法名与 Handler 名对应
- 返回值类型直接使用 GoZero `.api` 文件中定义的类型

### 3.2 Service 文件清单

```typescript
// services/auth.service.ts         ← group: auth（公开路由）
export const authService = {
  login:    (data: LoginReq)    => api.post<LoginResp>('/api/auth/login', data),
  register: (data: RegisterReq) => api.post<RegisterResp>('/api/auth/register', data),
  refresh:  (refreshToken: string) => api.post<RefreshResp>('/api/auth/refresh', { refreshToken }),
};

// services/user.service.ts        ← group: user（ClientAuth 鉴权）
export const userService = {
  getProfile:     ()               => api.get<UserInfo>('/api/user/profile'),
  changePassword: (data: ChangePasswordReq) => api.post<void>('/api/user/password', data),
};

// services/client.service.ts      ← group: client
export const clientService = {
  create:     (data: CreateClientReq) => api.post<CreateClientResp>('/api/client/create', data),
  pollStatus: (clientId: number)      => api.get<PollClientStatusResp>('/api/client/status', { params: { client_id: clientId } }),
  list:       ()                      => api.get<ListClientsResp>('/api/client/list'),
  delete:     (id: number)            => api.delete<void>(`/api/client/${id}`),
};

// services/application.service.ts ← group: application
export const applicationService = {
  create: (data: CreateApplicationReq) => api.post<CreateApplicationResp>('/api/application/create', data),
  list:   ()                           => api.get<ListApplicationsResp>('/api/application/list'),
  delete: (id: number)                 => api.delete<void>(`/api/application/${id}`),
};

// services/message.service.ts     ← group: message
export const messageService = {
  list:   (params: ListMessagesReq) => api.get<ListMessagesResp>('/api/message/list', { params }),
  delete: (id: number)              => api.delete<void>(`/api/message/${id}`),
};
```

---

## 四、Vue Query 数据流

### 4.1 查询 Hook 封装

```typescript
// composables/useApplications.ts
export function useApplications() {
  return useQuery({
    queryKey: ['applications'],
    queryFn: applicationService.list,
    staleTime: 60_000, // 60 秒内不重新请求
  });
}

// composables/useQRPolling.ts
export function useQRPolling(clientId: Ref<number | null>) {
  return useQuery({
    queryKey: ['client-status', clientId],
    queryFn: () => clientService.pollStatus(clientId.value!),
    refetchInterval: (query) => {
      if (query.state.data?.status === 'bound') return false;
      if (query.state.data?.status === 'error') return false;
      return 2000; // 每 2 秒轮询
    },
    enabled: computed(() => !!clientId.value),
  });
}
```

### 4.2 变更 Hook 封装

```typescript
// composables/useDeleteApplication.ts
export function useDeleteApplication() {
  const queryClient = useQueryClient();
  const toast = useToast();

  return useMutation({
    mutationFn: applicationService.delete,
    onSuccess: () => {
      // 失效缓存 → 自动重新获取列表
      queryClient.invalidateQueries({ queryKey: ['applications'] });
      toast.add({ severity: 'success', summary: '已删除', life: 3000 });
    },
  });
}
```

### 4.3 缓存策略

| 数据类型 | staleTime | 策略理由 |
|----------|-----------|----------|
| Application 列表 | 60s | 变更频率极低 |
| Client 列表 | 30s | 绑定状态可能变化 |
| Client 扫码状态 | 0（实时） | 需要轮询，不可缓存 |
| Message 列表 | 15s | 外部推送可能随时产生新消息 |
| 用户信息 | Infinity | 登录期间不变 |

---

## 五、数据流全景图

```
┌─────────────────────────────────────────────────────────────────┐
│                      Vue Component                               │
│  使用 useQuery / useMutation 获取数据和操作                     │
└──────────┬──────────────────────────────────┬───────────────────┘
           │ data / mutate                    │
           ▼                                  ▼
┌──────────────────────┐          ┌───────────────────────┐
│  @tanstack/vue-query │          │     Pinia Store        │
│  (服务端状态缓存)     │          │  (客户端状态：auth)    │
│                      │          │                       │
│  - queryKey 管理     │          │  - accessToken        │
│  - staleTime 控制    │          │  - refreshToken       │
│  - 后台重新验证       │          │  - user               │
│  - 乐观更新          │          │  - isAuthenticated     │
└──────────┬───────────┘          └───────────────────────┘
           │ queryFn / mutationFn
           ▼
┌──────────────────────────────────────────────────────────┐
│                    Service Layer                          │
│  auth.service / user.service / client.service / ...      │
│  每个方法 → 一个 GoZero Handler                          │
└──────────┬───────────────────────────────────────────────┘
           │ Axios 实例
           ▼
┌──────────────────────────────────────────────────────────┐
│                 Axios Interceptors                        │
│                                                           │
│  Request:  attach Bearer accessToken                      │
│  Response: code !== 0 → throw ApiError                    │
│           code === 100001 → refresh token → retry        │
└──────────┬───────────────────────────────────────────────┘
           │ HTTPS
           ▼
┌──────────────────────────────────────────────────────────┐
│              GoZero HTTPServer (:8080)                    │
│                                                           │
│  ClientAuth (JWT RS256)      ApplicationAuth (appToken)  │
└──────────────────────────────────────────────────────────┘
```

---

## 六、错误码处理策略

> 错误码定义来源：`pkg/xerr/errors_code.go`

### 6.1 前端错误码常量

```typescript
// types/api.ts
export const ErrorCode = {
  UNAUTHORIZED:         100001, // Access Token 过期/无效 → 触发刷新
  REGISTER_CLOSED:      110107, // V1 仅单用户
  LOGIN_PASSWORD_WRONG: 110202, // 密码错误
  REFRESH_INVALID:      110300, // Refresh Token 无效
  REFRESH_REVOKED:      110301, // Refresh Token 已轮换
  REFRESH_EXPIRED:      110302, // Refresh Token 已过期
  PASSWORD_TOO_WEAK:    110400, // 密码强度不足
} as const;
```

### 6.2 错误处理分层

| 层级 | 处理内容 |
|------|----------|
| Axios 响应拦截器 | `code === 100001` 触发刷新；`code !== 0` 抛出 `ApiError` |
| Vue Query `onError` | 显示 Toast 错误提示 |
| Vue Query `onSuccess` | 失效相关缓存、显示成功 Toast |
| 表单组件 | Zod 预校验防止无效请求；提交失败显示后端错误信息 |

### 6.3 修改密码后的安全联动

> 对应 UML：`docs/uml/12-sequence-changepassword.puml`

修改密码成功后，GoZero 后端会撤销所有 Refresh Token（`revoked = 1`）。前端必须：

```typescript
// 改密成功后
authStore.clearAuth();                              // 清除 Pinia
queryClient.clear();                                // 清除所有缓存
router.push('/login');                              // 跳转登录页
toast.add({ severity: 'info', summary: '密码修改成功，请重新登录', life: 5000 });
```
