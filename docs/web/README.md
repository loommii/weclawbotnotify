# WeClawBotNotify 前端技术文档

> Vue 3 管理后台的架构设计、技术选型、开发规范。

---

## 文档索引

| 编号 | 文档 | 内容 | 建议读者 | 状态 |
|------|------|------|----------|------|
| 01 | [架构决策记录 (ADR)](./01-architecture-decisions.md) | 每一项技术选型的背景、候选方案、选择理由 | 全员 | ✅ |
| 02 | [技术栈清单](./02-tech-stack.md) | 完整技术栈 + 版本号 + 与 GoZero 对接方式 | 开发人员 | ✅ |
| 03 | [UI 组件库选型](./03-ui-library.md) | PrimeVue vs NaiveUI vs Element Plus 对比 | 全员 | ✅ |
| 04 | [前后端交互策略](./04-api-interaction.md) | 双令牌刷新、Service 层、Vue Query 数据流 | 开发人员 | ✅ |
| 05 | [项目结构设计](./05-project-structure.md) | 目录结构、组件层次、命名规范、路由设计 | 开发人员 | ✅ |
| 06 | [视觉风格指南](./06-style-guide.md) | 基于 Linear 风格的配色/字体/间距/页面模式 | 设计+开发 | ✅ |

---

## 文档记录形式说明

### 为什么选择 Markdown 而非 UML？

本项目 `docs/` 下的文档采用 **分层记录** 策略：

```
docs/
├── uml/                              ← 架构图/流程图（PlantUML）
│   ├── 01-use-case.puml              ← 用例图
│   ├── 02-component.puml             ← 组件图
│   ├── 03-class.puml                 ← 类图
│   ├── 04-sequence-login.puml        ← 序列图
│   ├── ...
│   └── 12-sequence-changepassword.puml
│
└── web/                              ← 前端技术文档（Markdown）
    ├── README.md                     ← 本文档（索引）
    ├── 01-architecture-decisions.md  ← 架构决策记录
    ├── 02-tech-stack.md              ← 技术栈清单
    ├── 03-ui-library.md              ← UI 库选型分析
    ├── 04-api-interaction.md         ← 前后端交互策略
    └── 05-project-structure.md       ← 项目结构设计
```

### UML vs Markdown 的分工

| 文档类型 | 用什么 | 原因 |
|----------|--------|------|
| **系统架构图**（组件关系、部署拓扑） | PlantUML (`.puml`) | 可视化表达组件间关系一目了然 |
| **交互流程**（登录流程、Token 刷新流程） | PlantUML 序列图 | 时序逻辑用文字描述容易遗漏分支 |
| **状态机**（Client 生命周期、认证状态流转） | PlantUML 状态图 | 状态变迁可视化优于文字 |
| **技术选型理由**（为什么选 A 不选 B） | Markdown | 需要结构化论证、对比表格 |
| **技术栈清单**（版本号、用途、对接方式） | Markdown | 列表 + 表格形式最直观 |
| **代码示例**（目录结构、命名规范、实现方案） | Markdown | 代码块嵌入方便，可直接复制 |
| **API 映射表**（前端 Service ↔ GoZero Handler） | Markdown | 表格是最佳载体 |

### 前端文档是否需要 UML？

**不需要为前端单独画 UML。原因：**

1. **后端 UML 已经覆盖了交互流程**。例如 `04-sequence-login.puml` 中 `participant "Web 管理后台 (浏览器)"` 已经表达了前端的参与角色和交互时序。前端开发人员看这些序列图就能理解全链路。

2. **前端文档的重心是「选型」和「实现」**。Markdown 更适合承载对比表格、代码示例、命名规范这类文本密集型内容。

3. **如果未来需要前端架构图**（如组件树、数据流），可以直接在 Markdown 中嵌入 Mermaid，不需要单独 `.puml` 文件。示例：

   ````markdown
   ```mermaid
   graph TD
     App.vue --> AuthLayout
     App.vue --> AppLayout
     AppLayout --> Sidebar
     AppLayout --> RouterView
   ```
   ````

---

## 文档编写规范

### 命名规则

- 文件名格式：`{序号}-{英文短名}.md`
- 序号递增，表示建议的阅读顺序
- 英文短名使用 kebab-case

### Markdown 规范

- 使用 `##` 作为一级标题（`#` 留给文件顶部大标题）
- 表格对齐使用 `|------|------|`
- 代码块必须标注语言：```` ```typescript ````
- 引用后端代码时，使用相对路径如 `pkg/xerr/errors_code.go`
- 引用 UML 文档时，使用相对路径如 `../uml/04-sequence-login.puml`

### 文档更新时机

| 触发条件 | 需更新的文档 |
|----------|-------------|
| 新增技术选型 | `01-architecture-decisions.md`（新增 ADR） |
| 技术栈版本升级 | `02-tech-stack.md` |
| UI 组件库变更 | `03-ui-library.md` |
| API 路由变更 | `02-tech-stack.md`（API 映射表部分）、`04-api-interaction.md` |
| 目录结构调整 | `05-project-structure.md` |
| 新增错误码 | `02-tech-stack.md`（错误码映射部分） |
| 视觉风格迭代 | `06-style-guide.md` |

### 文档维护原则

1. **单一真相来源** —— 同一个信息只在一处定义，通过链接引用而非复制粘贴
2. **与代码同步** —— 目录结构、API 路由、错误码必须与代码仓库一致
3. **决策可追溯** —— ADR 记录为什么做这个选择，即使后来改了，也能理解当时背景

---

## 快速入门（新成员）

1. 先看 [架构决策记录](./01-architecture-decisions.md) —— 理解「为什么这么选」
2. 再看 [技术栈清单](./02-tech-stack.md) —— 知道「用了哪些东西」
3. 接着看 [前后端交互策略](./04-api-interaction.md) —— 理解「前后端怎么通」
4. 然后看 [项目结构设计](./05-project-structure.md) —— 知道「代码放哪里」
5. 接着看 [视觉风格指南](./06-style-guide.md) —— 知道「做成什么样」
6. [UI 组件库选型](./03-ui-library.md) 作为参考资料，用到组件时查 PrimeVue 官方文档

---

## 相关资源

- [PrimeVue 官方文档](https://primevue.org)
- [Tailwind CSS + PrimeVue 集成指南](https://primevue.org/tailwind)
- [Sakai Admin 模板](https://github.com/primefaces/sakai-vue)
- [Vue Router 文档](https://router.vuejs.org)
- [Pinia 文档](https://pinia.vuejs.org)
- [TanStack Vue Query 文档](https://tanstack.com/query/latest)
- [后端 UML 文档](../uml/)
