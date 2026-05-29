# GoEverythingToolbar

一个用 Go + Wails + React 实现的 Windows 极速文件搜索工具，参考 [EverythingToolbar](https://github.com/srwi/EverythingToolbar) 的 Launcher 模式，通过调用 [Everything](https://www.voidtools.com/) 引擎实现毫秒级全盘文件搜索。

> [!NOTE]
> Deskband 已在 Windows 11 被弃用，本项目采用 **Launcher 模式**：常驻系统托盘 + 全局快捷键唤起搜索浮窗，兼容 Windows 10/11。

---

## ✨ 功能特性

### 核心功能
- 🚀 **极速搜索**：基于 Everything SDK，毫秒级返回全盘文件
- ⌨️ **全局快捷键**：默认 `Win + Alt + S` 唤起/隐藏搜索窗口
- 🪟 **Launcher 模式**：常驻系统托盘，无边框透明窗口，居中显示
- 📜 **搜索历史**：自动记录最近搜索，可清空、可回溯
- 🎨 **毛玻璃 UI**：半透明蒙层 + backdrop-filter 模糊背景
- 📂 **结果分组**：按"应用 / 文件夹 / 文件"分类展示
- 🔍 **多种过滤选项**：
  - ☑ 仅搜索应用和文件夹（默认开启）
  - Aa 区分大小写
  - W 全词匹配
  - P 匹配完整路径
  - .* 正则表达式

### 交互能力
- ⬆⬇ 上下方向键导航
- ⏎ Enter 打开文件 / 启动应用
- 🖱 右键菜单：打开、打开所在文件夹、复制完整路径
- 🖱 双击直接打开
- `Ctrl + ↑` 唤起搜索历史面板
- `Esc` 隐藏窗口
- 失去焦点自动隐藏

---

## 📸 截图预览

> 以下截图位待补充实际运行效果图，路径建议放在 `docs/images/` 下。

| 主界面（毛玻璃搜索） | 分组结果展示 |
|---|---|
| ![主界面](docs/images/main.png) | ![分组结果](docs/images/grouped.png) |

| 右键菜单 | 系统托盘 |
|---|---|
| ![右键菜单](docs/images/context-menu.png) | ![托盘菜单](docs/images/tray.png) |

> [!TIP]
> 浏览器预览模式无需 Windows 即可体验全部 UI 效果，详见下方"快速开始 > 浏览器预览 UI"。

---

## 🛠 技术架构

| 层面 | 技术选型 |
|------|---------|
| 语言 | Go 1.21+ |
| UI 框架 | [Wails v2](https://wails.io/) (Go + WebView) |
| 前端 | React 18 + TypeScript + Vite 5 |
| 搜索引擎 | [Everything SDK](https://www.voidtools.com/support/everything/sdk/) (Everything64.dll) |
| 全局快捷键 | [golang.design/x/hotkey](https://github.com/golang-design/hotkey) |
| 系统托盘 | [lxn/walk](https://github.com/lxn/walk) (Windows) |
| 数据持久化 | JSON 文件（`%APPDATA%/GoEverythingToolbar/`） |

### 架构图

```
┌─────────────────────────────────────────────────────────┐
│                      Wails App                          │
│  ┌──────────────────────┐    ┌────────────────────────┐ │
│  │   React Frontend     │◄──►│      Go Backend        │ │
│  │  ┌───────────────┐   │    │  ┌──────────────────┐  │ │
│  │  │  SearchBar    │   │    │  │  app.go          │  │ │
│  │  │  ResultList   │   │    │  │  - Search        │  │ │
│  │  │  HistoryPanel │   │    │  │  - OpenFile      │  │ │
│  │  └───────────────┘   │    │  │  - GetHistory    │  │ │
│  │      ↓ Wails Bind    │    │  │  - ShowWindow    │  │ │
│  └──────────────────────┘    │  └────────┬─────────┘  │ │
│                              │           │            │ │
│  ┌───────────────────────────┴───────────▼──────────┐ │
│  │  internal/                                        │ │
│  │  ├── everything/  (Everything SDK 封装)           │ │
│  │  ├── hotkey/      (全局快捷键)                    │ │
│  │  ├── history/     (搜索历史)                      │ │
│  │  ├── config/      (配置管理)                      │ │
│  │  └── tray/        (系统托盘)                      │ │
│  └───────────────────┬──────────────────────────────┘ │
└──────────────────────┼─────────────────────────────────┘
                       ▼
              ┌─────────────────┐    IPC    ┌──────────────┐
              │ Everything64.dll├──────────►│  Everything  │
              └─────────────────┘           │   Process    │
                                            └──────────────┘
```

---

## 📁 项目结构

```
SeekerApp/
├── app.go                          # Wails App 主结构 + 前端绑定方法
├── main.go                         # 入口，配置 Wails 窗口属性
├── go.mod / go.sum                 # Go 模块依赖
├── wails.json                      # Wails 项目配置
├── GoEverythingToolbar.exe         # Windows 可执行文件（交叉编译产物）
│
├── internal/                       # 后端业务模块
│   ├── config/
│   │   └── config.go               # 配置管理（hotkey、maxHistory、theme、dllPath）
│   ├── everything/
│   │   ├── types.go                # SearchOptions / SearchResult 类型定义
│   │   ├── sdk_windows.go          # Everything DLL 调用封装（Windows 实现）
│   │   └── sdk_other.go            # 非 Windows stub
│   ├── history/
│   │   └── history.go              # 搜索历史管理
│   ├── hotkey/
│   │   ├── hotkey_windows.go       # 全局快捷键注册
│   │   └── hotkey_other.go         # 非 Windows stub
│   └── tray/
│       ├── tray_windows.go         # 系统托盘（lxn/walk）
│       └── tray_other.go           # 非 Windows stub
│
├── frontend/                       # React 前端
│   ├── package.json
│   ├── vite.config.ts
│   ├── tsconfig.json
│   ├── index.html
│   ├── src/
│   │   ├── main.tsx                # React 入口
│   │   ├── App.tsx                 # 根组件（状态管理、键盘导航、分组逻辑）
│   │   ├── App.css                 # 全局样式（毛玻璃、分组、菜单等）
│   │   ├── components/
│   │   │   ├── SearchBar.tsx       # 搜索栏 + 选项按钮
│   │   │   ├── ResultList.tsx      # 分组结果列表 + 右键菜单
│   │   │   └── HistoryPanel.tsx    # 搜索历史面板
│   │   └── mock/
│   │       └── wailsMock.ts        # 浏览器预览用 Mock 数据
│   └── wailsjs/                    # Wails 自动生成的绑定代码（勿改）
│
└── build/
    ├── appicon.png
    └── windows/
        ├── icon.ico
        └── wails.exe.manifest      # Windows 清单（DPI、兼容性）
```

---

## 🚀 快速开始

### 前置依赖

| 依赖 | 用途 | 安装方式 |
|------|------|---------|
| Go ≥ 1.21 | 后端编译 | https://go.dev/dl/ |
| Node.js ≥ 18 | 前端构建 | https://nodejs.org/ |
| Wails CLI v2 | 构建工具 | `go install github.com/wailsapp/wails/v2/cmd/wails@latest` |
| Everything | 搜索引擎 | https://www.voidtools.com/ |
| MinGW-w64（仅 macOS 交叉编译需要） | C 编译器 | `brew install mingw-w64` |

### 1. 浏览器预览 UI（macOS / Linux / Windows 通用）

仅预览界面效果，使用 Mock 数据：

```bash
cd frontend
npm install
npm run dev
```

浏览器打开 **http://localhost:5173/**（或 5174）即可看到 UI，所有交互（搜索、分组、右键菜单、历史等）均可体验，搜索结果为模拟数据。

### 2. Windows 完整运行（开发模式）

```powershell
# 1. 安装 Everything 并确保进程运行
# 2. 安装前端依赖
cd frontend
npm install
cd ..

# 3. 启动 Wails dev
wails dev
```

### 3. Windows 构建可执行文件

```powershell
# 在 Windows 系统执行
wails build
# 产物：build/bin/GoEverythingToolbar.exe
```

### 4. macOS 交叉编译为 Windows .exe

```bash
# 需要先安装 mingw-w64
brew install mingw-w64

# 编译
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc \
  go build -o GoEverythingToolbar.exe .
```

> [!WARNING]
> 交叉编译产出的 .exe 仍需在 Windows 上运行，且需要 Everything 进程和 Everything64.dll。

---

## ⚙ 配置说明

配置文件路径：`%APPDATA%/GoEverythingToolbar/config.json`

```json
{
  "hotkey": "win+alt+s",
  "maxHistory": 50,
  "dllPath": "",
  "theme": "dark",
  "windowWidth": 600,
  "windowHeight": 400
}
```

| 字段 | 说明 | 默认值 |
|------|------|--------|
| `hotkey` | 全局快捷键，支持 `ctrl/alt/shift/win + <key>` 组合 | `win+alt+s` |
| `maxHistory` | 搜索历史最大保留条数 | `50` |
| `dllPath` | Everything64.dll 路径（留空则自动搜索 PATH） | `""` |
| `theme` | UI 主题（暂仅深色） | `"dark"` |
| `windowWidth` / `windowHeight` | 窗口尺寸 | `600` / `400` |

历史文件：`%APPDATA%/GoEverythingToolbar/history.json`

---

## 🎮 使用指南

### 快捷键

| 操作 | 快捷键 |
|------|-------|
| 唤起/隐藏窗口 | `Win + Alt + S`（全局，可配置） |
| 关闭窗口 | `Esc` |
| 上下选择结果 | `↑` / `↓` |
| 打开选中项 | `Enter` |
| 切换搜索历史面板 | `Ctrl + ↑` |
| 弹出右键菜单 | 鼠标右键 |
| 双击打开 | 鼠标双击 |

### 搜索语法

支持 Everything 原生语法：

| 示例 | 说明 |
|------|------|
| `*.tsx` | 通配符，搜索所有 .tsx 文件 |
| `path:downloads` | 路径包含 downloads |
| `size:>10mb` | 大于 10MB 的文件 |
| `dm:today` | 今天修改的文件 |
| `ext:exe;msi` | 多扩展名 |

更多详见 [Everything 搜索语法](https://www.voidtools.com/support/everything/searching/)。

### 实战示例

**示例 1：找今天下载的所有 PDF**
```
ext:pdf path:downloads dm:today
```

**示例 2：找 D 盘大于 100MB 的视频**
```
d:\ size:>100mb ext:mp4;mkv;avi
```

**示例 3：用正则找含数字结尾的 .log 文件**
```
.*\d+\.log$
```
（启用 `.*` 正则模式按钮）

**示例 4：仅找应用快速启动**
```
chrome
```
（保持"仅搜索应用和文件夹"开启，结果会聚焦到 .exe / .app / .lnk）

**示例 5：找含中文的项目目录**
```
path:projects 文档
```

---

## 🌐 浏览器 Mock 预览说明

为方便在 macOS/Linux 上预览 UI，项目内置了 `frontend/src/mock/wailsMock.ts`：

- **自动检测环境**：仅当 `window.go` 不存在（即浏览器环境）时注入
- **不影响生产构建**：真实 Wails 应用中 `window.go` 由 Wails 注入，mock 自动跳过
- **模拟数据**：10+ 条覆盖应用、文件夹、各种文件类型的样本
- **模拟过滤**：支持 matchCase / matchPath / useRegex / matchWholeWord / appAndFolderOnly

打开浏览器控制台可以看到所有 mock 调用日志，方便调试前端逻辑。

---

## 📦 跨平台支持

| 模块 | Windows | macOS / Linux |
|------|---------|--------------|
| Everything SDK | ✅ 完整功能 | ❌ stub（返回错误） |
| 全局快捷键 | ✅ | ❌ stub |
| 系统托盘 | ✅ | ❌ stub |
| 浏览器 UI 预览 | ✅ | ✅ |

项目使用 Go build tags（`//go:build windows` / `//go:build !windows`）实现条件编译，非 Windows 平台不会编译失败，但运行时功能受限。

---

## 🛠 常见问题（FAQ）

### Q1: 搜索没有结果？
A: 确认以下三点：
1. Everything 主程序已启动（任务栏右下角能看到放大镜图标）
2. Everything 设置中已勾选 `Tools > Options > Service > Start Everything Service`（必须，否则 DLL 无法连接）
3. `Everything64.dll` 在 PATH 中或通过 `config.json` 的 `dllPath` 指定

### Q2: 全局快捷键不生效？
A: 可能与其他工具冲突（如微信截图、QQ 截图、系统输入法切换）。修改 `%APPDATA%/GoEverythingToolbar/config.json` 中的 `hotkey`，例如：
- `ctrl+shift+space`
- `alt+space`（注意会冲突 macOS Spotlight 风格）
- `win+\``（反引号，少冲突）

### Q3: 窗口失焦后不希望自动隐藏？
A: 编辑 `frontend/src/App.tsx`，移除 `window.addEventListener('blur', handleWindowBlur)` 一行，重新 `wails build`。

### Q4: 浏览器预览中点击搜索结果只弹 alert？
A: 这是 mock 模式的正常行为（无法真的打开 Windows 系统文件）。完整运行需在 Windows 上通过 `wails dev` / `wails build`。

### Q5: macOS 上 `wails dev` 失败？
A: 项目设计为 **Windows 专用**，macOS 上仅能：
- 浏览器预览前端 UI（`cd frontend && npm run dev`）
- 交叉编译为 Windows .exe（需 mingw-w64）

### Q6: 字体看起来有点糊？
A: 检查 `build/windows/wails.exe.manifest` 中是否启用了 DPI 感知。如果显示器是高分屏（≥2K），需要 `<dpiAware>True/PM</dpiAware>` 与 `<dpiAwareness>PerMonitorV2</dpiAwareness>` 同时存在。

### Q7: 程序启动后没有窗口？
A: 这是预期行为。本应用为 **Launcher 模式**，启动后只常驻系统托盘，按全局快捷键 `Win + Alt + S` 才会显示搜索窗口。也可右键托盘图标 → "显示窗口"。

### Q8: 如何彻底退出？
A: 关闭窗口（点 ✕ 或按 `Esc`）只是隐藏，**不会退出进程**。需要通过 **右键系统托盘图标 → 退出**。

### Q9: 历史记录存在哪？想清空怎么办？
A: 路径：`%APPDATA%/GoEverythingToolbar/history.json`。可直接删除该文件，或在历史面板内点 "清空" 按钮。

### Q10: 想改默认搜索结果数 / 字体 / 颜色？
A:
- 搜索结果数：修改 `app.go` 中 `SearchOptions.MaxResults` 默认值
- 字体颜色：修改 `frontend/src/App.css` 中 `:root` 的 CSS 变量
- 透明度：调整 `--bg-overlay` 的 alpha 值

更多开发相关问题请参考 [DEVELOPMENT.md](./DEVELOPMENT.md)。

---

## 📌 后续规划

- [ ] 完整的桌面式原生右键菜单（打开方式、移到废纸篓、属性等）
- [ ] 设置面板（在 UI 内修改快捷键、主题）
- [ ] 浅色主题
- [ ] 搜索结果虚拟滚动（大量结果时优化性能）
- [ ] 文件预览（图片/文本/代码高亮）
- [ ] 收藏夹 / 置顶
- [ ] 多语言切换（中/英）

---

## 📚 进阶文档

- [DEVELOPMENT.md](./DEVELOPMENT.md) — 开发者指南（环境搭建、调试、构建、代码规范）
- [ARCHITECTURE.md](./ARCHITECTURE.md) — 架构说明（模块划分、数据流、技术决策）

---

## 🙏 致谢与参考

本项目站在以下优秀开源项目的肩膀上：

| 项目 | 用途 | 协议 |
|------|------|------|
| [EverythingToolbar](https://github.com/srwi/EverythingToolbar) | 原始 Windows 实现参考（C# WPF） | MIT |
| [Everything](https://www.voidtools.com/) | 文件搜索引擎 | Freeware |
| [Wails](https://wails.io/) | Go + WebView 桌面应用框架 | MIT |
| [React](https://react.dev/) | 前端 UI 框架 | MIT |
| [Vite](https://vitejs.dev/) | 前端构建工具 | MIT |
| [golang.design/x/hotkey](https://github.com/golang-design/hotkey) | 跨平台全局快捷键 | MIT |
| [lxn/walk](https://github.com/lxn/walk) | Windows 原生 UI 库（用于托盘） | BSD-3 |

---

## 📄 License

本项目采用 [MIT License](./LICENSE) 开源协议，欢迎自由使用、修改、分发。

---

## 🙋 反馈与贡献

- 🐛 发现 Bug：欢迎在 Issues 区提交，请附上 OS 版本、Everything 版本、复现步骤
- 💡 功能建议：欢迎在 Issues 区提 Feature Request
- 🔧 代码贡献：请先阅读 [DEVELOPMENT.md](./DEVELOPMENT.md)，按照其中的 **提交规范** 提交 PR
- ⭐ 如果本项目对你有帮助，欢迎 Star 支持
