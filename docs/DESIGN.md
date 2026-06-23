# Marmot Ledger Design System

> **Calm Marmot Finance** — 温和可信的土拨鼠财务驾驶舱

Marmot Ledger 的 UI 底色来自 **Wise + Linear + Notion** 的组合，但不直接照搬任何一个产品。

```text
Wise      50%  品牌气质和金融可信感
Linear    35%  组件、表格、后台结构
Notion    15%  留白、低压力、文档感
```

目标：让 Marmot Ledger 看起来像一个长期使用不累、可信、清晰、生活化的个人/家庭财务驾驶舱。

---

## 1. Design Principles

| Principle | Meaning |
| --- | --- |
| Calm | 冷静、克制，不让财务数据制造焦虑 |
| Trustworthy | 清晰、可靠，每个数字都有归属 |
| Human | 生活化的财务，不做冰冷银行后台 |
| Practical | 信息优先，少装饰，多效率 |

---

## 2. Inspiration Blend

### Wise — approachable finance

借鉴：

- 温和可信的金融产品气质
- 多币种、账户、资产、汇率场景的自然心智
- 大面积浅色背景和清晰卡片
- 不像传统银行那么冷，也不像交易软件那么刺激

### Linear — crisp interface structure

借鉴：

- 高信息密度但不混乱
- 表格、筛选器、状态和 hover 的克制表达
- 适合后台系统长期使用的组件节奏
- 清晰的视觉层级和操作路径

### Notion — low-pressure readability

借鉴：

- 留白
- 温和边框
- 文档化的信息组织
- 空状态和说明文字的低压力表达

---

## 3. Brand Assets

品牌头像图片位于根目录：

```text
img/marmot-ledger-1.png
img/marmot-ledger-2.png
```

使用规则：

- README、首页、登录页、空状态、项目介绍图优先使用这两张图片。
- 不要再使用历史项目图片作为 Marmot Ledger 的品牌图。
- 图像展示时保留足够留白，不要裁切到主体边缘。
- 推荐尺寸：README 中 160–220px，登录页/空状态中 96–160px。

---

## 4. Color System

### Background layers

```css
--bg-main:          #f7f5ef;  /* 主背景：暖白/米白 */
--bg-main-soft:     #f8faf7;  /* 备用浅背景 */
--bg-card:          #ffffff;  /* 卡片/面板 */
--bg-sidebar:       linear-gradient(180deg, #1f2933 0%, #111827 100%); /* 侧栏深色 */
--bg-sidebar-solid: #1f2933;  /* tooltip / 单色场景 */
--bg-nav:           #ffffff;  /* 顶部 Header */
--bg-input:         #ffffff;
--bg-hover:         rgba(0, 0, 0, 0.03);
```

主页面不要使用纯白铺满；用暖白/米白作为底色，让财务产品更温和。
侧栏保持深 slate，避免与主内容暖白底色抢层级；顶部 Header 保持白色，并使用 `rgba(47, 125, 92, 0.10)` 作为极轻的品牌底边描线。

### Brand colors

```css
--brand-green:      #2f7d5c;  /* 苔绿 = 品牌 / focus / 主操作 */
--brand-deep:       #1f2933;  /* 深层文字和深色背景 */
--brand-accent:     #2f7d5c;  /* 主操作 / 链接色（默认沿用苔绿） */
--brand-warm:       #f4efe6;  /* 暖色点缀 */
--brand-moss:       #dce9df;  /* 柔和辅助色 / Active 状态点 */
```

主操作色统一使用苔绿 `#2f7d5c`；蓝色仅作为信息/链接辅助色，不再作为主 CTA 颜色。
侧栏 active 项使用 `brand-moss` 的圆点 + 半透明 moss 背景，提示当前位置。

### Financial semantic colors

**项目硬性规则：品牌/focus 使用苔绿，收入红色，支出橙色。**

```css
--brand-green:      #2f7d5c;  /* 苔绿 = 品牌 / focus / 当前状态 */
--income-color:     #ef4444;  /* 红色 = 收入 */
--expense-color:    #f97316;  /* 橙色 = 支出 / 消费 / 流出 */
```

不要使用绿色表达支出，避免和品牌绿色、focus 状态混淆。

应用场景：

- 金额颜色
- 收入/支出标签
- 分类图标
- 图表中的收入/支出系列
- 收支筛选器状态

### Neutral text colors

```css
--text-primary:     #1e293b;
--text-secondary:   #64748b;
--text-tertiary:    #94a3b8;
--text-inverse:     #f1f5f9;
--border:           #e5e7eb;
--border-light:     #f0f0ed;
```

### Status colors

```css
--success:          #f97316;
--warning:          #f59e0b;
--error:            #ef4444;
--info:             #3b82f6;
```

---

## 5. Typography

```css
--font-family:      'PingFang SC', -apple-system, BlinkMacSystemFont,
                    'Segoe UI', Roboto, 'Helvetica Neue', sans-serif;
--font-mono:        'SF Mono', 'Fira Code', 'JetBrains Mono', monospace;
```

| Usage | Size | Weight | Line Height |
| --- | --- | --- | --- |
| Page title | 28–32px | 700 | 1.3 |
| Section title | 20px | 600 | 1.4 |
| Card title | 16px | 600 | 1.4 |
| Body | 14px | 400 | 1.6 |
| Small / meta | 12px | 400 | 1.5 |
| Amount | 14–28px | 500–700 | 1.4 |

金额和资产数字使用 monospace 字体，提高对齐和可信感。

---

## 6. Layout & Spacing

使用 4px 网格：

```text
4 / 8 / 12 / 16 / 24 / 32 / 48
```

布局规则：

- 页面最大宽度：`1200px`
- 页面左右 padding：桌面 `24px`，移动端 `16px`
- 卡片 padding：`24px`
- 侧栏宽度：展开 `260px`，折叠 `70px`
- 顶部栏高度：`70px`
- 页面区块间距：`24px` 或 `32px`

---

## 7. Components

### Cards

```css
--card-radius:      16px;
--card-padding:     24px;
--card-shadow:      0 2px 8px rgba(0,0,0,0.03),
                    0 8px 24px rgba(0,0,0,0.02);
```

规则：

- 白底
- 16px 圆角
- 轻阴影
- 可使用顶部 2px 品牌色线条
- 不要厚边框、重阴影、强渐变

### Buttons

- 主按钮：苔绿 `#2f7d5c`，用于创建、保存、确认等主操作。
- 次按钮：白底 + 轻描边，用于取消、刷新、普通辅助动作。
- 编辑按钮：暖中性色 `#f4efe6` 背景 + `#6b5b49` 文本，避免与主操作或危险操作抢层级。
- 危险按钮：仅删除/危险操作使用红色 `#ef4444` 文本 + 浅红背景。
- 高度：卡片内操作默认 `32–36px`，Header / Dialog 主操作 `40px`，重要 CTA 可用 `44px`。
- 圆角：`10px`

不要用收入红色/支出橙色作为普通按钮色；不要用蓝色作为管理模块主操作色。

### Forms

- 输入框高度：`36px`
- 边框：`1px solid var(--border)`
- 聚焦：浅蓝 ring
- 表单标签：14px / 600
- 错误提示短句，不写长段落

### Tables

采用 Linear 式表格：

- 表头透明或极浅背景
- 表头文字 12px / 600 / `text-tertiary`
- 行高 44px
- 金额列右对齐
- 金额使用 monospace
- 轻分割线
- hover 使用 `bg-hover`
- 不使用厚重边框

### Tags

- 圆角：6px
- 背景浅色，文字使用语义色
- 收入标签红色
- 支出标签绿色

### Empty states

- 使用品牌头像或轻插图
- 一句话说明当前为空
- 一个明确行动按钮
- 不要只显示“暂无数据”

### Charts

- 图表低噪音
- 色彩不超过 5 种
- 收入/支出必须使用项目语义色
- 饼图最多展示前 5 项，其余归入“其他”
- 趋势图用线条为主，不要大面积强填充

---

## 8. Dark Mode

```css
--bg-main:          #0f172a;
--bg-card:          #1e293b;
--bg-sidebar:       #020617;
--text-primary:     #f1f5f9;
--text-secondary:   #94a3b8;
--text-tertiary:    #64748b;
--border:           #334155;
```

财务语义色不随深色主题反转：

```css
--income-color:     #ef4444;
--expense-color:    #f97316;
```

---

## 9. Responsive Rules

| Breakpoint | Behavior |
| --- | --- |
| > 1024px | 展开侧栏，内容区多列卡片 |
| 768–1024px | 折叠侧栏，卡片 2 列 |
| < 768px | 单列布局，表格转卡片列表 |

移动端规则：

- 避免横向滚动
- 金额、账户、时间优先展示
- 操作按钮收进更多菜单
- 统计图高度控制在 240px 内

---

## 10. Do / Don't

| Do | Don't |
| --- | --- |
| 使用暖白背景和白色卡片 | 纯白全屏或黑灰后台感 |
| 收入红色、支出橙色 | 红支出、绿收入 |
| 金额使用 monospace | 金额使用花哨字体 |
| 表格轻分割、轻 hover | 厚边框表格 |
| 一个页面一个主操作 | 多个同级强按钮抢焦点 |
| 图表低噪音 | 图表颜色过多 |
| 空状态使用品牌图和行动按钮 | 只写“暂无数据” |
| 品牌图使用 `img/marmot-ledger-*` | 使用历史项目图片 |

---

## 11. For AI Agents

生成 Marmot Ledger UI 时必须遵守：

1. 主背景使用暖白/米白：`#f7f5ef` 或接近色。
2. 卡片白底、16px 圆角、轻阴影。
3. 收入金额/标签使用红色 `#ef4444`。
4. 支出金额/标签使用绿色 `#f97316`。
5. 主操作按钮使用蓝色 `#3b82f6`。
6. 金额数字使用 monospace 并右对齐。
7. 表格采用轻分割线和 hover，不做厚边框。
8. 图表减少颜色，只强调趋势和占比。
9. 品牌头像优先使用 `img/marmot-ledger-1.png` 和 `img/marmot-ledger-2.png`。
10. 页面气质应接近 Wise 的温和可信、Linear 的清晰克制、Notion 的低压力可读。

---

## 12. Product Tone

Marmot Ledger should feel like:

```text
A calm personal finance cockpit:
warm, trustworthy, clear, and practical.
```

中文表达：

```text
一个温和可信、清晰实用、适合长期使用的个人与家庭财务驾驶舱。
```
