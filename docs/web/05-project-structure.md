# 05 — 前端项目结构设计

> 详细的目录结构设计、文件命名规范、组件拆分原则。

---

## 一、完整目录结构

```
web/                                  # 前端项目根目录
├── index.html                        # SPA 入口 HTML
├── package.json
├── tsconfig.json
├── vite.config.ts                    # Vite 配置（含 proxy 到 GoZero）
├── tailwind.config.ts                # TailwindCSS + tailwindcss-primeui
├── eslint.config.ts
├── .prettierrc
│
├── public/                           # 静态资源（不经过 Vite 处理）
│   └── favicon.ico
│
└── src/
    ├── main.ts                       # 应用入口（注册 PrimeVue / Router / Pinia / VueQuery）
    ├── App.vue                       # 根组件（<RouterView /> + <Toast /> + <ConfirmDialog />）
    ├── vite-env.d.ts
    │
    ├── config/
    │   └── index.ts                  # API_BASE_URL / TOKEN_KEY / APP_NAME
    │
    ├── lib/
    │   ├── axios.ts                  # Axios 实例 + 拦截器（Token 刷新核心逻辑）
    │   ├── query-client.ts           # Vue Query 全局配置（staleTime / retry）
    │   ├── storage.ts                # localStorage 安全封装（try-catch + JSON parse）
    │   └── utils.ts                  # 通用工具函数
    │
    ├── types/
    │   ├── api.ts                    # ApiResponse<T> + ErrorCode 常量
    │   ├── auth.ts                   # UserInfo, TokenPair, LoginReq/Resp
    │   ├── client.ts                 # ClientInfo, CreateClientReq, PollClientStatusResp
    │   ├── application.ts            # ApplicationInfo, CreateApplicationReq
    │   └── message.ts                # MessageInfo, ListMessagesReq
    │
    ├── stores/
    │   ├── auth.ts                   # Pinia auth store（accessToken, refreshToken, user）
    │   └── ui.ts                     # Pinia UI store（sidebar collapsed, theme）
    │
    ├── composables/
    │   ├── useAuth.ts                # 登录/注册/登出/刷新 逻辑封装
    │   ├── useQRPolling.ts           # Client 扫码状态轮询
    │   ├── useApplications.ts        # Application 列表查询
    │   ├── useClients.ts             # Client 列表查询
    │   ├── useMessages.ts            # Message 列表查询
    │   ├── useDeleteApplication.ts   # 删除 Application mutation
    │   ├── useDeleteClient.ts        # 删除 Client mutation
    │   ├── useDeleteMessage.ts       # 删除 Message mutation
    │   ├── useCreateApplication.ts   # 创建 Application mutation
    │   ├── useCreateClient.ts        # 创建 Client mutation
    │   └── usePagination.ts          # 分页参数管理
    │
    ├── services/
    │   ├── auth.service.ts           # login / register / refresh
    │   ├── user.service.ts           # getProfile / changePassword
    │   ├── client.service.ts         # create / pollStatus / list / delete
    │   ├── application.service.ts    # create / list / delete
    │   └── message.service.ts        # list / delete
    │
    ├── views/                        # 页面级组件（对应路由）
    │   ├── auth/
    │   │   ├── LoginView.vue
    │   │   └── RegisterView.vue
    │   ├── dashboard/
    │   │   └── DashboardView.vue
    │   ├── clients/
    │   │   └── ClientsView.vue
    │   ├── applications/
    │   │   └── ApplicationsView.vue
    │   ├── messages/
    │   │   └── MessagesView.vue
    │   └── settings/
    │       └── SettingsView.vue
    │
    ├── components/                   # 可复用组件
    │   ├── layout/
    │   │   ├── AppLayout.vue         # 主布局（Sidebar + Header + Content）
    │   │   ├── AppSidebar.vue        # 侧边栏导航
    │   │   └── AppHeader.vue         # 顶部栏（用户信息、退出）
    │   ├── auth/
    │   │   ├── LoginForm.vue         # 登录表单
    │   │   └── RegisterForm.vue      # 注册表单
    │   ├── clients/
    │   │   ├── ClientTable.vue       # Client 列表表格
    │   │   ├── CreateClientDialog.vue # 创建 Client 弹窗
    │   │   └── QRCodeDisplay.vue     # 二维码展示（含轮询状态）
    │   ├── applications/
    │   │   ├── ApplicationTable.vue  # Application 列表表格
    │   │   ├── CreateAppDialog.vue   # 创建 Application 弹窗
    │   │   └── TokenDisplay.vue      # Token 一次性展示组件
    │   ├── messages/
    │   │   └── MessageTable.vue      # 消息列表表格
    │   └── settings/
    │       └── ChangePasswordForm.vue # 修改密码表单
    │
    ├── schemas/                      # Zod 校验 Schema
    │   ├── login.schema.ts           # 登录表单校验
    │   ├── register.schema.ts        # 注册表单校验
    │   ├── password.schema.ts        # 修改密码校验
    │   └── application.schema.ts     # 创建 Application 校验
    │
    └── router/
        └── index.ts                  # Vue Router 路由表 + 导航守卫
```

---

## 二、组件层次关系

```
App.vue
└── <RouterView />
    │
    ├── 未登录 → AuthLayout
    │   ├── LoginView.vue
    │   │   └── LoginForm.vue          # VeeValidate + Zod
    │   └── RegisterView.vue
    │       └── RegisterForm.vue       # VeeValidate + Zod
    │
    └── 已登录 → AppLayout
        ├── AppSidebar.vue             # PanelMenu
        ├── AppHeader.vue              # Toolbar + Avatar + Dropdown
        └── <RouterView />
            ├── DashboardView.vue
            ├── ClientsView.vue
            │   ├── CreateClientDialog.vue
            │   │   └── QRCodeDisplay.vue
            │   └── ClientTable.vue
            ├── ApplicationsView.vue
            │   ├── CreateAppDialog.vue
            │   │   └── TokenDisplay.vue
            │   └── ApplicationTable.vue
            ├── MessagesView.vue
            │   └── MessageTable.vue
            └── SettingsView.vue
                └── ChangePasswordForm.vue
```

---

## 三、命名规范

### 3.1 文件命名

| 类型 | 规范 | 示例 |
|------|------|------|
| 页面组件 | `PascalCase + View` 后缀 | `LoginView.vue`, `DashboardView.vue` |
| 业务组件 | `PascalCase` | `ClientTable.vue`, `QRCodeDisplay.vue` |
| 布局组件 | `PascalCase + Layout` 后缀 | `AppLayout.vue` |
| Service | `kebab-case.service.ts` | `auth.service.ts` |
| Composable | `use + PascalCase.ts` | `useQRPolling.ts` |
| Store | `kebab-case.ts` | `auth.ts` |
| Schema | `kebab-case.schema.ts` | `login.schema.ts` |
| Type | `kebab-case.ts` | `api.ts` |

### 3.2 组件 Script 结构

```vue
<script setup lang="ts">
// 1. 类型导入
import type { ClientInfo } from '@/types/client';

// 2. Props / Emits
const props = defineProps<{ ... }>();
const emit = defineEmits<{ ... }>();

// 3. Composables
const { data, isLoading } = useQuery(...);

// 4. 本地状态
const visible = ref(false);

// 5. 计算属性
const filtered = computed(() => ...);

// 6. 方法
function handleCreate() { ... }
</script>
```

---

## 四、路由设计

```typescript
// router/index.ts
const routes = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/LoginView.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('@/views/auth/RegisterView.vue'),
    meta: { requiresAuth: false },
  },
  {
    path: '/',
    component: () => import('@/components/layout/AppLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      { path: '', name: 'Dashboard',