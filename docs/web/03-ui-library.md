# 03 — UI 组件库选型分析

> 对 Vue 3 生态中主流 UI 组件库的对比分析，以及 PrimeVue 在本项目中的适配评估。

---

## 一、候选库总览

| 维度 | PrimeVue 4 | NaiveUI | Element Plus | Ant Design Vue |
|------|------------|---------|--------------|----------------|
| 组件数 | 80+ | 80+ | 70+ | 60+ |
| 默认主题 | Aura（现代、精致） | 简洁实用 | Element 风格 | Ant Design 风格 |
| 官方维护方 | PrimeTek（土耳其） | 07akioni（个人） | 饿了么团队 | 社区 |
| GitHub Stars | 10k+ | 16k+ | 24k+ | 20k+ |
| npm 总下载 | 400M+ | 中等 | 高 | 中等 |
| TypeScript | 良好 | ⭐⭐⭐⭐⭐ 原生 | ⭐⭐⭐⭐ | ⭐⭐⭐⭐ |
| 中文文档 | ❌ 英文为主 | ✅ 原生中文 | ✅ 原生中文 | ✅ 原生中文 |
| Tailwind 融合 | ⭐⭐⭐⭐⭐ 官方插件 | ❌ | ❌ | ❌ |
| Figma 设计稿 | ✅ PrimeOne | ❌ | ✅ Element UI Kit | ❌ |
| Admin 模板 | Sakai（简洁） | Soybean Admin | Vue Element Admin | Vue Vben Admin |

---

## 二、选择 PrimeVue 的核心优势

### 2.1 设计语言

PrimeVue 的 Aura 主题由专职设计团队维护，风格接近 Stripe / Linear / Vercel 等新一代 B2B 产品设计语言：

- 圆角更大更柔和（rounded-border design token）
- 阴影层次丰富（emphasis / highlight 多级语义色）
- 间距宽松（符合现代「呼吸感」设计趋势）
- 暗黑模式原生支持且视觉统一

### 2.2 TailwindCSS 深度融合

`tailwindcss-primeui` 插件将 PrimeVue 主题的 design token 映射为 TailwindCSS 语义 class：

| Tailwind Class | 含义 |
|---------------|------|
| `bg-primary` | 主色背景（跟随主题） |
| `text-primary-contrast` | 主色上的文字色 |
| `text-muted-color` | 次要文字色 |
| `border-surface` | 边框色 |
| `bg-emphasis` | 悬停态背景 |
| `rounded-border` | 主题圆角 |
| `dark:bg-surface-900` | 暗黑模式表面色 |

这意味着写页面布局时，TailwindCSS 的颜色直接和 PrimeVue 组件的主题色自动同步。NaiveUI 和 Element Plus **做不到这个级别**，需要手动维护 CSS 变量映射。

### 2.3 DataTable — 功能最全的 Vue 表格

本项目的三个列表页（Client、Application、Message）是核心交互，DataTable 直接决定了日常使用体验。PrimeVue DataTable 支持：

| 功能 | PrimeVue | NaiveUI | Element Plus |
|------|----------|---------|--------------|
| 排序（单列/多列/可撤销） | ✅ | ✅ | ✅ |
| 过滤（行内/菜单/全局） | ✅ | ✅ | ✅ |
| 分页（完全自定义模板） | ✅ | ✅ | ✅ |
| 虚拟滚动 | ✅ | ✅ | ✅ |
| 行展开/行编辑 | ✅ | ✅ | ✅ |
| 行选择（单击/复选框） | ✅ | ✅ | ✅ |
| 列调整大小 | ✅ | ✅ | ✅ |
| 加载骨架屏 | ✅ | ✅ | ✅ |
| 导出 CSV | ✅ | ❌ | ❌ |
| 行分组 | ✅ | ❌ | ❌ |
| 冻结列 | ✅ | ❌ | ❌ |
| 列排序（拖拽） | ✅ | ❌ | ❌ |

### 2.4 双模式架构

| 模式 | 用途 | 本项目选择 |
|------|------|-----------|
| **Styled Mode** | 直接用 Aura 等预置主题，开箱即用 | ✅ 使用此模式 |
| **Unstyled Mode** | 组件只有逻辑无样式，用 TailwindCSS 全局控制外观 | 暂不使用 |

### 2.5 官方 Admin 模板 (Sakai)

```bash
git clone https://github.com/primefaces/sakai-vue.git
```

Sakai 自带：侧边栏 + 头部 + 路由配置 + 暗黑模式切换 + Dashboard 示例。代码量比 NaiveUI 的 Soybean Admin 更少，更适合本项目 V1 规模。

---

## 三、与项目需求的适配矩阵

| 页面/功能 | 需要的组件 | PrimeVue 是否内置 | 备注 |
|-----------|-----------|------------------|------|
| 登录页 | InputText, Password, Button | ✅ | `Password` 组件含显隐切换 |
| 注册页 | InputText, Password, Button | ✅ | 同上 |
| 主布局 | PanelMenu, Toolbar, Avatar | ✅ | 或自定义 Slot |
| Client 列表 | DataTable, Button, Dialog, Tag | ✅ | 删除用 `ConfirmDialog` |
| 创建 Client | Dialog, InputText, QRCode 展示 | ⚠️ QRCode 需 `qrcode` 库 | 封装 `QRCodeDisplay.vue` |
| Application 列表 | DataTable, Button, Tag | ✅ | 同上 |
| 创建 Application | Dialog, InputText, Textarea, Message | ✅ | Token 展示用 `Dialog` + `Text` copyable |
| 消息列表 | DataTable, Tag, Button | ✅ | Tag 按优先级着色 |
| 修改密码 | Password (3个), Button | ✅ | 旧密码 + 新密码 + 确认 |
| 消息提示 | Toast | ✅ | 操作成功/失败反馈 |

**结论：覆盖率 100%，仅 QRCode 需要额外引入 `qrcode` 库（6KB）—— NaiveUI 同样需要。**

---

## 四、唯一的短板：无内置 QRCode

解决方案 —— 封装 `QRCodeDisplay.vue`：

```vue
<template>
  <div class="flex flex-col items-center gap-4">
    <canvas ref="canvasRef" />
    <Tag :severity="statusSeverity">{{ statusText }}</Tag>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import QRCode from 'qrcode';
import Tag from 'primevue/tag';

const props = defineProps<{
  url: string;
  status: 'pending' | 'bound' | 'error';
}>();

const canvasRef = ref<HTMLCanvasElement>();

const statusSeverity = computed(() => {
  switch (props.status) {
    case 'bound': return 'success';
    case 'error': return 'danger';
    default: return 'warn';
  }
});

onMounted(async () => {
  if (canvasRef.value && props.url) {
    await QRCode.toCanvas(canvasRef.value, props.url, { width: 200 });
  }
});
</script>
```

封装成本约 30 行，且比 Ant Design 内置的 `<QRCode>` 更灵活（支持颜色/Logo 等高级定制）。

---

## 五、不选其他库的理由

### 为什么不选 NaiveUI？

- TailwindCSS 融合不如 PrimeVue（需手动维护 CSS 变量映射）
- 视觉效果不如 PrimeVue 精致
- 无官方 Figma 设计稿

### 为什么不选 Element Plus？

- 设计风格偏传统「后台管理系统」，不如 PrimeVue 现代化
- TypeScript 类型推导不如 NaiveUI 和 PrimeVue
- 无 TailwindCSS 融合

### 为什么不选 Ant Design Vue？

- 版本迭代慢（v4 迟迟未稳定）
- 部分组件 API 与 React 版本不一致
- 本质是 Ant Design 的社区移植，非官方维护
