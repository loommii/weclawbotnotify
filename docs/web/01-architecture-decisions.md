# 01 — 架构决策记录 (ADR)

> 记录 WeClawBotNotify 前端架构中每一项关键技术选型的决策背景、候选方案、最终选择及理由。

---

## ADR-001: 选用 Vue 3 作为前端框架

### 决策日期
2026-05

### 状态
已采纳

### 背景
WeClawBotNotify 是一个微信消息通知推送系统，V1 为单用户模式。后端使用 GoZero 框架，前端需要提供一个管理后台（Dashboard），用于管理 Client（微信Bot 绑定）、Application（推送凭证）、查看通知记录、修改密码等操作。

### 候选方案

| 方案 | 描述 |
|------|------|
| **Vue 3 + TypeScript** | 渐进式框架，模板语法直觉，国内社区活跃 |
| React 18 + TypeScript | 函数式 UI 库，生态最大，并发特性先进 |
| Angular 17 | 企业级框架，依赖注入 + RxJS |
| Svelte 4 | 编译时框架，无 Virtual DOM，包体积极小 |

### 决策

**选择 Vue 3 + TypeScript**

### 理由

1. **国内社区优势** —— Vue 在国内中小团队更普及，中文文档、视频教程、技术文章极丰富。对于以国语为主的团队，学习成本和问题排查效率显著优于 React。

2. **模板语法降低上手门槛** —— SFC（单文件组件）的 `<template>` + `<script setup>` + `<style scoped>` 关注点内聚但语法简洁，非专业前端开发者也能快速参与。

3. **生态完整性** —— Vue Router（路由）、Pinia（状态管理）、VueUse（工具库）、PrimeVue（UI 组件库）形成完整闭环，不需要像 React 生态那样在多个竞争方案间反复评估。

4. **本项目复杂度适配** —— V1 仅 6 个页面、CRUD 为主、单用户，Vue 3 的能力完全覆盖且没有过度设计。

5. **Axios 拦截器逻辑等效** —— 双令牌轮换机制在 Axios 层面实现，与框架无关。Vue 3 的 Pinia store 可以在组件外调用（`useAuthStore()`），实现 token 刷新队列的方式与 React 完全等效。

### 影响

- 团队需要熟悉 Vue 3 Composition API（`<script setup>` 语法）
- Vite 作为构建工具（Vue 官方推荐）
- UI 组件库需要在 Vue 生态中选择（见 ADR-003）

---

## ADR-002: 选用 Vite 作为构建工具

### 决策日期
2026-05

### 状态
已采纳

### 背景
前端 SPA 需要构建工具。GoZero 后端通过 `//go:embed` 嵌入静态资源，产物需为纯静态文件（HTML + JS + CSS）。

### 候选方案

| 方案 | 描述 |
|------|------|
| **Vite 5** | 基于 esbuild + Rollup，秒级冷启动，Vue 官方推荐 |
| Webpack 5 | 生态最丰富但配置复杂，构建速度较慢 |
| Turbopack | Next.js 绑定，输出 Node.js 运行时 |

### 决策

**选择 Vite 5**

### 理由

1. **与 GoZero 部署模型天然匹配** —— Vite 构建产物为纯静态文件，`outDir` 直接指向 GoZero 的 `static/` 目录即可。

2. **零配置成本** —— TypeScript 原生编译（esbuild）、CSS 原生支持、`server.proxy` 一行配置解决开发时跨域。

3. **Vue 官方生态** —— `create-vue` 脚手架默认使用 Vite，`@vitejs/plugin-vue` 由 Vue 核心团队维护。

### 影响

- 开发服务器端口默认 5173，需配置 proxy 到 GoZero 后端（:8080）
- 生产构建产物输出路径需与 GoZero 的 `embed` 路径对齐

---

## ADR-003: 选用 PrimeVue 作为 UI 组件库

### 决策日期
2026-05

### 状态
已采纳

### 背景
管理后台需要完整的 UI 组件支持：布局（侧边栏 + 头部 + 内容区）、表格（Client / Application / Message 列表）、表单（登录/注册/改密）、弹窗（创建/删除确认）、二维码展示（Client 创建）。

### 候选方案

| 方案 | 描述 |
|------|------|
| **PrimeVue 4 (Aura)** | 国际最大 Vue UI 库，设计精美，80+ 组件 |
| NaiveUI | 国内优秀 Vue 3 组件库，TypeScript 原生 |
| Element Plus | 饿了么出品，国内 Vue 2 迁移首选 |
| Ant Design Vue | Ant Design 的 Vue 实现 |

### 决策

**选择 PrimeVue 4 (Aura 主题)**

### 理由

1. **设计语言出色** —— Aura 主题风格接近 Stripe / Linear 等新一代 B2B 产品，视觉上比 NaiveUI 和 Element Plus 更精致。

2. **TailwindCSS 深度融合** —— 官方 `tailwindcss-primeui` 插件将主题色映射为 Tailwind 语义 class（`bg-primary`、`text-muted-color`），页面级样式和组件样式视觉完全统一，这是其他 Vue UI 库做不到的。

3. **DataTable 功能最强** —— 支持行过滤、列排序、分页、行选择、行展开、虚拟滚动、骨架屏，覆盖本项目全部列表需求。

4. **官方 Admin 模板（Sakai）** —— 提供开箱即用的布局骨架（侧边栏 + 路由 + 暗黑模式），可快速起手。

5. **国际化社区活跃** —— 400M+ npm 下载量、GitHub 10k+ Star，问题排查资源丰富。

### 唯一的短板

PrimeVue 无内置 QRCode 组件，需引入 `qrcode` 库（6KB）自行封装展示组件——NaiveUI 同样需要。

### 影响

- 文档以英文为主，团队需要阅读英文 API 文档
- 引入 `qrcode` 作为二维码生成依赖
- 如需中文日期/时间组件，需配置国际化

---

## ADR-004: 选用 Pinia 作为客户端状态管理

### 决策日期
2026-05

### 状态
已采纳

### 背景
前端需要管理极少的客户端状态：Access Token、Refresh Token、当前用户信息、认证状态标记。不需要 Redux 级别的状态管理方案。

### 候选方案

| 方案 | 描述 |
|------|------|
| **Pinia** | Vue 官方推荐，Setup Store 语法与 Composition API 一致 |
| Vuex 4 | Vue 3 兼容版本，但已进入维护模式 |
| 纯 Composition API + provide/inject | 无额外依赖，但 token 持久化需自行处理 |

### 决策

**选择 Pinia**

### 理由

1. **Vue 官方推荐** —— Vuex 已进入维护模式，Pinia 是 Vue 3 的状态管理标准答案。

2. **Setup Store 语法** —— `defineStore('auth', () => { ... })` 与 Vue 3 Composition API 心智模型完全一致，学习成本为零。

3. **DevTools 集成** —— Pinia 完整集成 Vue DevTools，支持时间旅行调试。

4. **持久化简单** —— Token 通过 Pinia 的 `$subscribe` 或手动调用 `localStorage.setItem` 即可持久化。

### 影响

- 仅需一个 `auth` store，不超过 30 行代码
- 项目不引入 Vuex

---

## ADR-005: 选用 @tanstack/vue-query 管理服务端状态

### 决策日期
2026-05

### 状态
已采纳

### 背景
前端需要管理来自后端 API 的数据：Application 列表、Client 列表、Client 扫码状态（轮询）、Message 通知记录。这些数据的共同特征是「以服务端为准」，需要缓存策略、后台刷新、乐观更新。

### 候选方案

| 方案 | 描述 |
|------|------|
| **@tanstack/vue-query** | 服务端状态管理专用库 |
| Pinia 自行管理所有状态 | 把 API 数据也放入 Pinia，手动维护缓存逻辑 |
| 组件内直接请求 | 无统一缓存管理 |

### 决策

**选择 @tanstack/vue-query v5**

### 理由

1. **服务端状态 vs 客户端状态分离** —— Pinia 管理 token/user 等客户端状态，Vue Query 管理 API 数据的缓存/失效/重新获取。各司其职，不混淆。

2. **轮询支持** —— Client 扫码状态需要每 2 秒轮询直到绑定完成，Vue Query 的 `refetchInterval` + 条件停止机制刚好匹配。

3. **缓存失效语义清晰** —— 创建/删除操作后用 `queryClient.invalidateQueries()` 标记相关缓存失效，自动触发重新获取，不需要手动管理多处状态。

4. **同一团队出品** —— 与 React 生态的 TanStack Query 共享设计理念，社区资源可互通参考。

### 影响

- 所有 API 请求通过 `useQuery` / `useMutation` 封装
- 缓存 key 命名需规范化（如 `['applications']`、`['client-status', clientId]`）

---

## ADR-006: 选用 Axios 作为 HTTP 客户端

### 决策日期
2026-05

### 状态
已采纳

### 背景
双令牌轮换机制需要：请求拦截器自动附加 Access Token、响应拦截器检测 401 并自动刷新 Token、刷新期间的并发请求排队机制。

### 候选方案

| 方案 | 描述 |
|------|------|
| **Axios** | 拦截器链 + Promise 队列模型 |
| 原生 fetch | 无拦截器，需全部自行封装 |
| ky | 基于 fetch 的轻量封装 |

### 决策

**选择 Axios**

### 理由

1. **拦截器架构** —— 响应拦截器统一剥离 GoZero 的 `{ code, msg, data }` 外层包装、统一拦截业务错误码（`code !== 0`）、统一处理 401 → Token 刷新。

2. **并发刷新队列** —— `isRefreshing` 锁 + `failedQueue` 队列在 Axios 的 Promise 链中实现最自然。

3. **请求/响应转换** —— `transformRequest` / `transformResponse` 可以统一处理数据格式。

### 影响

- Service 层所有 API 函数通过 Axios 实例调用
- 拦截器中错误码映射需与 GoZero 的 `xerr` 定义对齐

---

## ADR-007: 选用 TailwindCSS + tailwindcss-primeui 作为样式方案

### 决策日期
2026-05

### 状态
已采纳

### 背景
PrimeVue 负责组件级 UI，但页面布局、间距、响应式需要页面级样式方案。

### 候选方案

| 方案 | 描述 |
|------|------|
| **TailwindCSS + tailwindcss-primeui** | 原子化 CSS，官方插件统一色值 |
| 纯 PrimeVue | 仅用组件样式，布局靠内联 style |
| SCSS Modules | 组件级 scoped 样式 |

### 决策

**选择 TailwindCSS + tailwindcss-primeui 插件**

### 理由

1. **语义色统一** —— `tailwindcss-primeui` 插件将 PrimeVue 的 design token 映射为 Tailwind class（`bg-primary`、`text-muted-color`、`border-surface`），写布局时颜色自动和组件同步。

2. **不冲突** —— TailwindCSS 管布局（flex、gap、padding），PrimeVue 管组件内部。两者职责分明。

### 影响

- 需安装 `tailwindcss-primeui` 插件
- Tailwind v4 用 `@import "tailwindcss-primeui"`；v3 用 `plugins: [PrimeUI]`

---

## 版本历史

| 日期 | 版本 | 变更 |
|------|------|------|
| 2026-05-01 | 1.0 | 初始版本，记录全部 7 项 ADR |
