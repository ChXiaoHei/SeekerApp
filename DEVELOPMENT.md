# 开发者指南（DEVELOPMENT）

本文档面向希望参与 **GoEverythingToolbar** 开发、调试或二次定制的开发者，提供从环境搭建到发布构建的全流程指引。

> 如果只想快速使用，请阅读 [README.md](./README.md)；如需了解架构设计，请阅读 [ARCHITECTURE.md](./ARCHITECTURE.md)。

---

## 📋 目录

- [1. 环境准备](#1-环境准备)
- [2. 首次拉取与初始化](#2-首次拉取与初始化)
- [3. 开发模式详解](#3-开发模式详解)
- [4. 目录结构详解](#4-目录结构详解)
- [5. 跨平台编译](#5-跨平台编译)
- [6. 调试技巧](#6-调试技巧)
- [7. 构建与发布](#7-构建与发布)
- [8. 代码规范](#8-代码规范)
- [9. 测试](#9-测试)
- [10. 常见开发问题 FAQ](#10-常见开发问题-faq)

---

## 1. 环境准备

### 1.1 必需软件清单

| 工具 | 版本要求 | 用途 | 备注 |
|------|---------|------|------|
| Go | ≥ 1.21 | 后端编译 | 必装 |
| Node.js | ≥ 18 LTS | 前端构建 | 必装 |
| npm / pnpm | npm ≥ 9 / pnpm ≥ 8 | 前端包管理 | 二选一 |
| Wails CLI | v2.x | 构建工具 | 必装 |
| Everything | ≥ 1.4.1 | 搜索引擎（运行时依赖） | 仅 Windows 运行需要 |
| MinGW-w64 | 任意 | C 编译器 | 仅 macOS 交叉编译 Windows 时需要 |
| Git | ≥ 2.20 | 版本控制 | 必装 |

### 1.2 安装命令

**Windows（PowerShell）**

```powershell
# Go
winget install GoLang.Go

# Node.js
winget install OpenJS.NodeJS.LTS

# Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 验证
wails doctor
```

**macOS（Homebrew）**

```bash
brew install go node mingw-w64

go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 把 $GOPATH/bin 加入 PATH
echo 'export PATH="$HOME/go/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc

wails doctor
```

**Linux（Ubuntu/Debian）**

```bash
sudo apt update
sudo apt install -y golang nodejs npm gcc-mingw-w64
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

### 1.3 Everything 客户端配置

> [!IMPORTANT]
> Everything **客户端进程必须运行**，且必须启用 Service 模式，否则 SDK 调用会失败。

1. 从 https://www.voidtools.com/ 下载安装
2. 启动 Everything
3. 菜单 → `Tools > Options > Service`
4. 勾选 **Start Everything Service**
5. 重启 Everything

**验证：** 任务管理器中应同时看到 `Everything.exe` 和 `Everything Service`。

### 1.4 Everything64.dll 准备

应用启动时按以下顺序查找 DLL：

1. `config.json` 中 `dllPath` 指定的绝对路径
2. 应用同目录
3. 系统 `PATH` 中的目录

**推荐做法：** 将 `Everything64.dll`（在 Everything SDK 压缩包内）放到应用同目录，避免 PATH 污染。

---

## 2. 首次拉取与初始化

```bash
# 1. 拉取代码
git clone <your-repo-url>
cd SeekerApp

# 2. 安装前端依赖
cd frontend
npm install
cd ..

# 3. 验证 Go 模块
go mod tidy

# 4. 启动 Wails dev（Windows 推荐）
wails dev

# 或仅启动前端（macOS/Linux 浏览器预览）
cd frontend && npm run dev
```

### 2.1 首次拉取检查清单

- [ ] `wails doctor` 全部 ✅
- [ ] `go version` 返回 ≥ 1.21
- [ ] `node -v` 返回 ≥ v18
- [ ] `frontend/node_modules` 已生成
- [ ] `frontend/wailsjs/` 已生成（首次 `wails dev` 自动）

---

## 3. 开发模式详解

项目支持两种开发模式，建议根据 OS 选用：

| 模式 | 命令 | 平台 | 完整度 | 适用场景 |
|------|------|------|--------|----------|
| Wails Dev | `wails dev` | Windows | ⭐⭐⭐⭐⭐ 完整 | 调试后端逻辑、全功能联调 |
| 浏览器 Mock | `cd frontend && npm run dev` | 全平台 | ⭐⭐⭐ 仅 UI | 调试 UI、组件、样式 |

### 3.1 Wails Dev 模式

```powershell
wails dev
```

特性：
- 启动 WebView2 窗口
- 前端代码热重载（HMR）
- Go 代码改动需重启
- 真实调用 Everything SDK、全局快捷键、系统托盘

**调试参数：**

```powershell
wails dev -loglevel debug    # 详细日志
wails dev -devserver localhost:34115   # 自定义 devserver 端口
wails dev -browser           # 同时开浏览器（前后端分离调试）
```

### 3.2 浏览器 Mock 模式

```bash
cd frontend
npm run dev
```

特性：
- Vite 极速启动（< 1s）
- 自动注入 `wailsMock.ts`，模拟 `window.go.main.App.XXX()` 调用
- 所有 UI 交互完整可用
- 控制台可看到 mock 调用日志

**Mock 数据修改位置：** `frontend/src/mock/wailsMock.ts`

### 3.3 两种模式差异表

| 能力 | Wails Dev | 浏览器 Mock |
|------|-----------|-------------|
| 真实文件搜索 | ✅ | ❌（模拟数据） |
| 全局快捷键 | ✅ | ❌ |
| 系统托盘 | ✅ | ❌ |
| 打开文件/文件夹 | ✅ | ⚠️ alert 提示 |
| 搜索历史持久化 | ✅（落盘） | ⚠️ 内存（刷新丢失） |
| UI 样式 | ✅ | ✅ |
| 键盘导航、右键菜单 | ✅ | ✅ |
| 启动速度 | ~3s | < 1s |
| 跨平台 | 仅 Windows | 全平台 |

---

## 4. 目录结构详解

### 4.1 后端目录（`internal/`）

| 目录 | 文件 | 职责 |
|------|------|------|
| `config/` | `config.go` | 配置加载/保存（`%APPDATA%/.../config.json`） |
| `everything/` | `types.go` | `SearchOptions` / `SearchResult` 数据结构 |
| `everything/` | `sdk_windows.go` | Everything DLL 的 syscall 封装 |
| `everything/` | `sdk_other.go` | 非 Windows stub（编译占位） |
| `history/` | `history.go` | 搜索历史 LRU + JSON 落盘 |
| `hotkey/` | `hotkey.go` | 公共接口定义 |
| `hotkey/` | `hotkey_windows.go` | golang.design/x/hotkey 注册 |
| `hotkey/` | `hotkey_other.go` | 非 Windows stub |
| `tray/` | `tray_windows.go` | lxn/walk 系统托盘菜单 |
| `tray/` | `tray_other.go` | 非 Windows stub |

### 4.2 根目录文件

| 文件 | 职责 |
|------|------|
| `main.go` | Wails 应用入口，配置窗口属性（Frameless、AlwaysOnTop、StartHidden、Transparent） |
| `app.go` | Wails App 主结构，定义前端可调用的方法（Search/OpenFile/GetHistory/...） |
| `go.mod` / `go.sum` | Go 模块依赖 |
| `wails.json` | Wails 项目配置（产物名、icon、frontend 路径） |

### 4.3 前端目录（`frontend/`）

| 路径 | 职责 |
|------|------|
| `package.json` | npm 依赖 + 脚本 |
| `vite.config.ts` | Vite + React 插件配置 |
| `tsconfig.json` | TS 编译配置（jsx: react-jsx） |
| `index.html` | HTML 入口，挂载 `#root` |
| `src/main.tsx` | React 入口，注入 wailsMock |
| `src/App.tsx` | 根组件：状态、键盘导航、分组逻辑 |
| `src/App.css` | 全局样式：毛玻璃、分组、菜单 |
| `src/components/SearchBar.tsx` | 搜索栏 + 过滤选项按钮 |
| `src/components/ResultList.tsx` | 分组结果列表 + 右键菜单 |
| `src/components/HistoryPanel.tsx` | 历史面板 |
| `src/mock/wailsMock.ts` | 浏览器环境 mock |
| `wailsjs/` | **Wails 自动生成**，请勿手动修改 |

### 4.4 构建目录

| 路径 | 用途 |
|------|------|
| `build/appicon.png` | 应用图标源 |
| `build/windows/icon.ico` | Windows 图标 |
| `build/windows/wails.exe.manifest` | DPI、UAC、兼容性清单 |
| `build/bin/` | `wails build` 输出目录 |

---

## 5. 跨平台编译

### 5.1 build tag 约定

```go
//go:build windows
// +build windows

package hotkey
// ... Windows 实现
```

```go
//go:build !windows
// +build !windows

package hotkey
// ... stub 实现，避免非 Windows 编译失败
```

### 5.2 文件命名约定

| 后缀 | 编译条件 |
|------|---------|
| `xxx_windows.go` | 仅 Windows 编译 |
| `xxx_darwin.go` | 仅 macOS 编译 |
| `xxx_linux.go` | 仅 Linux 编译 |
| `xxx_other.go`（自定义） | 配合 `//go:build !windows` 使用 |

### 5.3 为什么 macOS 上 `go build` 也能通过

- 所有 Windows-only 包都有对应的 `_other.go` stub
- syscall 调用全部隔离在 `sdk_windows.go` 中
- 这样 IDE 索引、单元测试、CI 都可以在非 Windows 跑

### 5.4 交叉编译为 Windows .exe（macOS/Linux）

```bash
# 安装 mingw-w64
brew install mingw-w64    # macOS
sudo apt install mingw-w64    # Linux

# 编译
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 \
  CC=x86_64-w64-mingw32-gcc \
  go build -o GoEverythingToolbar.exe .
```

> [!WARNING]
> 交叉编译仅产出 .exe，但运行仍需 Windows 系统 + Everything 客户端 + Everything64.dll。

---

## 6. 调试技巧

### 6.1 后端调试

**方式 1：详细日志**

```powershell
wails dev -loglevel debug
```

**方式 2：在 Go 代码中加日志**

```go
import "log"

func (a *App) Search(query string, opts everything.SearchOptions) ([]everything.SearchResult, error) {
    log.Printf("[Search] query=%q opts=%+v", query, opts)
    // ...
}
```

日志输出位置：`wails dev` 启动的终端。

**方式 3：Delve 调试器**

```bash
# 安装 Delve
go install github.com/go-delve/delve/cmd/dlv@latest

# 启动调试
dlv debug --headless --listen=:2345 --api-version=2
```

IDE（GoLand / VSCode）接入 `:2345` 端口即可断点调试。

### 6.2 前端调试

**方式 1：右键检查（推荐）**

`wails dev` 模式下，在窗口内按 `F12` 或右键 → 检查（需在 main.go 开启 `Debug.OpenInspectorOnStartup`）。

**方式 2：React DevTools**

浏览器 Mock 模式下，安装 [React Developer Tools](https://chrome.google.com/webstore/detail/react-developer-tools/fmkadmapgofadopljbjfkapdkoienihi) 扩展。

**方式 3：Mock 数据修改**

修改 `frontend/src/mock/wailsMock.ts` 中的 `MOCK_FILES` 数组，刷新浏览器立即生效。

### 6.3 常见构建错误与解决方案

| 错误 | 原因 | 解决 |
|------|------|------|
| `Everything64.dll not found` | DLL 未在 PATH/同目录 | 复制 DLL 到 exe 同目录，或配置 `config.json:dllPath` |
| `hotkey: registration failed` | 快捷键已被占用 | 改 `config.json:hotkey` |
| `walk: Failed to create main window` | 缺少 manifest 或 32 位运行 | 确保 build/windows/wails.exe.manifest 存在 |
| `npm ERR! ERESOLVE` | 依赖冲突 | `rm -rf node_modules package-lock.json && npm install` |
| `wails: command not found` | `$GOPATH/bin` 未在 PATH | 见 1.2 节配置 |
| `cgo: C compiler "gcc" not found` | macOS 缺 mingw | `brew install mingw-w64` |

---

## 7. 构建与发布

### 7.1 开发构建（带调试符号）

```powershell
wails build
```

产物：`build/bin/GoEverythingToolbar.exe`（~15-20 MB）

### 7.2 生产构建（精简体积）

```powershell
wails build -clean -platform windows/amd64 -ldflags "-s -w"
```

`-ldflags "-s -w"` 去掉调试符号，可减小 ~30% 体积。

### 7.3 NSIS 安装包

```powershell
wails build -nsis
```

产物：`build/bin/GoEverythingToolbar-amd64-installer.exe`

### 7.4 产物体积优化建议

| 措施 | 节省 |
|------|------|
| `-ldflags "-s -w"` | ~30% |
| `upx --best app.exe` | 额外 ~50% |
| 移除未使用的前端依赖 | 视情况 |
| 关闭 Vite sourcemap | ~5-10% |

> [!WARNING]
> UPX 可能被部分杀软误报为病毒，发布到公共渠道前需评估风险。

### 7.5 资源嵌入

`Everything64.dll` **不建议嵌入** exe：
- DLL 来自第三方，需保留独立可替换
- 嵌入后增加约 200KB，且需运行时解压到临时目录
- 推荐做法：发布包中独立放置 DLL，或要求用户自行安装 Everything

---

## 8. 代码规范

### 8.1 Go 代码规范

- 使用 `gofmt` / `goimports` 自动格式化
- 通过 `go vet ./...` 静态检查
- 推荐工具：`golangci-lint run ./...`
- 公开 API 必须有注释，遵循 Godoc 规范
- 错误处理：`if err != nil { return nil, fmt.Errorf("xxx: %w", err) }`

### 8.2 TypeScript 代码规范

- 函数与组件使用 **箭头函数 + 类型注解**
- 优先 `const` / `let`，禁用 `var`
- 严格 null 检查（`tsconfig.strict: true`）
- 组件文件名 PascalCase，普通文件 camelCase
- 不使用 `any`，必要时用 `unknown` 再窄化

### 8.3 提交信息约定（Conventional Commits）

```
<type>(<scope>): <subject>

<body>

<footer>
```

**type 取值：**

| type | 含义 |
|------|------|
| `feat` | 新功能 |
| `fix` | Bug 修复 |
| `docs` | 文档变更 |
| `style` | 仅格式调整（无逻辑变化） |
| `refactor` | 重构（无新功能、无 Bug 修复） |
| `perf` | 性能优化 |
| `test` | 测试相关 |
| `chore` | 构建/工具/依赖变更 |

**示例：**

```
feat(search): 支持仅搜索应用和文件夹过滤项

新增 appAndFolderOnly 字段到 SearchOptions
默认开启，可在 SearchBar 切换

Closes #12
```

---

## 9. 测试

### 9.1 当前测试覆盖

- Go 后端：暂无单元测试（后续补齐）
- 前端：暂无单元测试（后续补齐）
- 手动测试：见下方 checklist

### 9.2 添加 Go 单元测试

```go
// internal/history/history_test.go
package history

import "testing"

func TestAddDuplicate(t *testing.T) {
    h := New(10)
    h.Add("foo")
    h.Add("foo")
    if len(h.List()) != 1 {
        t.Fatalf("expect dedup, got %d", len(h.List()))
    }
}
```

运行：

```bash
go test ./...
go test -v -cover ./internal/...
```

### 9.3 手动测试 checklist

启动后逐项验证：

- [ ] `Win+Alt+S` 唤起窗口
- [ ] 输入查询，结果按"应用/文件夹/文件"分组显示
- [ ] ↑/↓ 切换选中项
- [ ] Enter 打开文件
- [ ] 右键菜单出现，"打开/打开所在文件夹/复制路径"功能正常
- [ ] `Esc` 隐藏窗口
- [ ] 失焦自动隐藏
- [ ] `Ctrl+↑` 打开历史面板，点击历史项回填搜索词
- [ ] 取消勾选"仅搜索应用和文件夹"，能搜到普通文件
- [ ] 右键托盘 → 退出，进程完全结束（任务管理器确认）

---

## 10. 常见开发问题 FAQ

### Q1: 改完前端代码不生效？
A: 检查 `wails dev` 终端是否有 HMR 提示。若无，重启 `wails dev`。Vite 偶发缓存问题可执行：
```bash
cd frontend
rm -rf node_modules/.vite
```

### Q2: 改完 Go 代码不生效？
A: Go 代码不支持热重载，需 `Ctrl+C` 停止 `wails dev` 后重新运行。

### Q3: 全局快捷键被占用怎么办？
A:
1. 查找占用者：Windows 自带工具无法直接查，可用 [HotkeysList](https://www.uwe-sieber.de/files/hotkeyslist.zip)
2. 修改 `config.json:hotkey` 为其它组合
3. 重启应用生效

### Q4: 如何切换为其他快捷键？
A: 编辑 `%APPDATA%/GoEverythingToolbar/config.json`：
```json
{ "hotkey": "ctrl+shift+space" }
```
支持的修饰符：`ctrl` / `alt` / `shift` / `win`，按键名见 `golang.design/x/hotkey` 文档。

### Q5: 前端 `window.go is undefined` 报错？
A:
- Wails 模式：检查 `wailsjs/` 是否生成；删后重新 `wails dev`
- 浏览器模式：检查 `main.tsx` 是否调用了 `setupWailsMock()`

### Q6: 如何添加新的前端 → Go 调用？
A:
1. 在 `app.go` 中新增方法（首字母大写）：
   ```go
   func (a *App) MyNewMethod(arg string) (string, error) { ... }
   ```
2. 重新 `wails dev`，`wailsjs/go/main/App.js` 会自动生成绑定
3. 前端调用：
   ```ts
   import { MyNewMethod } from '../wailsjs/go/main/App'
   await MyNewMethod('hello')
   ```
4. 若使用浏览器 mock，需在 `wailsMock.ts` 中补充对应 mock 实现

### Q7: 想添加新的搜索过滤项怎么办？
A:
1. `internal/everything/types.go` 添加字段到 `SearchOptions`
2. `internal/everything/sdk_windows.go` 处理新字段（调用对应 DLL 函数）
3. `frontend/src/components/SearchBar.tsx` 添加切换按钮
4. `frontend/src/App.tsx` 把新字段传入 Search 调用
5. `frontend/src/mock/wailsMock.ts` 同步 mock 行为

### Q8: 打包后启动闪退？
A: 大概率原因：
- 缺 `Everything64.dll`（最常见）
- 缺 WebView2 Runtime（Win10 老版本）→ 安装 [Evergreen Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/)
- manifest 引发的 UAC 异常 → 检查 `build/windows/wails.exe.manifest`

查看崩溃信息：`eventvwr.msc` → Windows 日志 → 应用程序

### Q9: 想加入 macOS / Linux 支持？
A: 工作量较大，主要难点：
- Everything 仅 Windows，需替换为 `mdfind`（macOS）/ `locate`（Linux）
- 全局快捷键和托盘 lib 已经跨平台，无需大改
- 详见 [ARCHITECTURE.md > 跨平台策略](./ARCHITECTURE.md#7-跨平台策略)

### Q10: 如何贡献代码？
A:
1. Fork → 新建 feature 分支
2. 修改并提交（遵循 [Conventional Commits](#83-提交信息约定conventional-commits)）
3. 推到自己 fork
4. 提 Pull Request，描述清楚 **变更动机** 和 **测试方法**

---

## 📚 相关文档

- [README.md](./README.md) — 项目介绍与快速开始
- [ARCHITECTURE.md](./ARCHITECTURE.md) — 架构设计深度解析
- [Wails 官方文档](https://wails.io/docs/introduction)
- [Everything SDK 文档](https://www.voidtools.com/support/everything/sdk/)
