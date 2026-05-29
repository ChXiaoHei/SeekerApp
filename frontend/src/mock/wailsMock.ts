/**
 * Wails API Mock for browser-based preview
 * 仅在浏览器中运行时（非 Wails 应用）注入模拟数据，方便 UI 预览
 */

interface MockApp {
  ClearHistory: () => Promise<void>
  CopyPathToClipboard: (path: string) => Promise<void>
  GetConfig: () => Promise<any>
  GetHistory: () => Promise<string[]>
  GetSearchOptions: () => Promise<any>
  HideWindow: () => Promise<void>
  IsWindowVisible: () => Promise<boolean>
  OpenFile: (path: string) => Promise<void>
  OpenFolder: (path: string) => Promise<void>
  SaveConfig: (config: any) => Promise<void>
  Search: (query: string, opts: any) => Promise<any[]>
  ShowWindow: () => Promise<void>
  ToggleWindow: () => Promise<void>
}

// 判断是否为 App / 可执行程序
const APP_EXTENSIONS = ['exe', 'msi', 'app', 'lnk', 'appx', 'dmg', 'pkg']
export function isAppFile(fileName: string): boolean {
  const ext = fileName.split('.').pop()?.toLowerCase() || ''
  return APP_EXTENSIONS.includes(ext)
}

// Mock 搜索结果数据（涵盖 App / 文件夹 / 各类文件）
const MOCK_FILES = [
  // === Apps（应用） ===
  {
    fileName: 'Chrome.app',
    fullPath: '/Applications/Chrome.app',
    isFolder: false,
    size: 1024 * 1024 * 350,
    dateModified: new Date('2026-05-20T10:30:00').toISOString()
  },
  {
    fileName: 'VSCode.exe',
    fullPath: 'C:\\Program Files\\Microsoft VS Code\\VSCode.exe',
    isFolder: false,
    size: 1024 * 1024 * 120,
    dateModified: new Date('2026-05-22T14:20:00').toISOString()
  },
  {
    fileName: 'Reactotron.app',
    fullPath: '/Applications/Reactotron.app',
    isFolder: false,
    size: 1024 * 1024 * 180,
    dateModified: new Date('2026-05-15T09:00:00').toISOString()
  },
  {
    fileName: 'Everything.exe',
    fullPath: 'C:\\Program Files\\Everything\\Everything.exe',
    isFolder: false,
    size: 1024 * 1024 * 8,
    dateModified: new Date('2026-04-10T11:30:00').toISOString()
  },
  {
    fileName: 'NodeSetup.msi',
    fullPath: 'C:\\Downloads\\NodeSetup.msi',
    isFolder: false,
    size: 1024 * 1024 * 35,
    dateModified: new Date('2026-03-01T16:00:00').toISOString()
  },

  // === Folders（文件夹） ===
  {
    fileName: 'react-app',
    fullPath: '/Users/demo/projects/react-app',
    isFolder: true,
    size: 0,
    dateModified: new Date('2026-05-28T09:00:00').toISOString()
  },
  {
    fileName: 'src',
    fullPath: '/Users/demo/projects/react-app/src',
    isFolder: true,
    size: 0,
    dateModified: new Date('2026-05-27T09:00:00').toISOString()
  },
  {
    fileName: 'Documents',
    fullPath: '/Users/demo/Documents',
    isFolder: true,
    size: 0,
    dateModified: new Date('2026-05-26T18:00:00').toISOString()
  },
  {
    fileName: 'Downloads',
    fullPath: '/Users/demo/Downloads',
    isFolder: true,
    size: 0,
    dateModified: new Date('2026-05-25T18:00:00').toISOString()
  },

  // === Files（普通文件） ===
  {
    fileName: 'README.md',
    fullPath: '/Users/demo/projects/react-app/README.md',
    isFolder: false,
    size: 2048,
    dateModified: new Date('2026-05-20T10:30:00').toISOString()
  },
  {
    fileName: 'package.json',
    fullPath: '/Users/demo/projects/react-app/package.json',
    isFolder: false,
    size: 1024,
    dateModified: new Date('2026-05-25T14:20:00').toISOString()
  },
  {
    fileName: 'App.tsx',
    fullPath: '/Users/demo/projects/react-app/src/App.tsx',
    isFolder: false,
    size: 5338,
    dateModified: new Date('2026-05-28T17:12:00').toISOString()
  },
  {
    fileName: 'main.tsx',
    fullPath: '/Users/demo/projects/react-app/src/main.tsx',
    isFolder: false,
    size: 256,
    dateModified: new Date('2026-05-28T17:19:00').toISOString()
  },
  {
    fileName: 'photo.jpg',
    fullPath: '/Users/demo/Pictures/photo.jpg',
    isFolder: false,
    size: 1024 * 1024 * 2.5,
    dateModified: new Date('2026-05-15T18:00:00').toISOString()
  },
  {
    fileName: 'document.pdf',
    fullPath: '/Users/demo/Documents/document.pdf',
    isFolder: false,
    size: 1024 * 512,
    dateModified: new Date('2026-05-10T11:30:00').toISOString()
  },
  {
    fileName: 'video.mp4',
    fullPath: '/Users/demo/Videos/video.mp4',
    isFolder: false,
    size: 1024 * 1024 * 150,
    dateModified: new Date('2026-05-05T20:45:00').toISOString()
  },
  {
    fileName: 'config.json',
    fullPath: '/Users/demo/.config/app/config.json',
    isFolder: false,
    size: 768,
    dateModified: new Date('2026-05-22T08:15:00').toISOString()
  }
]

// Mock 搜索历史
let mockHistory: string[] = ['react', 'package.json', 'README', 'config']

// Mock Config
let mockConfig = {
  hotkey: 'win+alt+s',
  maxHistory: 50,
  dllPath: '',
  theme: 'dark',
  windowWidth: 600,
  windowHeight: 400
}

const mockApp: MockApp = {
  async ClearHistory() {
    mockHistory = []
    console.log('[Mock] Cleared history')
  },

  async CopyPathToClipboard(path: string) {
    try {
      await navigator.clipboard.writeText(path)
      console.log('[Mock] Copied to clipboard:', path)
    } catch (e) {
      console.warn('[Mock] Failed to copy:', e)
    }
  },

  async GetConfig() {
    return { ...mockConfig }
  },

  async GetHistory() {
    return [...mockHistory]
  },

  async GetSearchOptions() {
    return {
      matchCase: false,
      matchWholeWord: false,
      matchPath: false,
      useRegex: false,
      maxResults: 100
    }
  },

  async HideWindow() {
    console.log('[Mock] HideWindow called')
  },

  async IsWindowVisible() {
    return true
  },

  async OpenFile(path: string) {
    console.log('[Mock] 打开文件:', path)
    alert(`[模拟] 打开文件：${path}`)
  },

  async OpenFolder(path: string) {
    console.log('[Mock] 打开所在文件夹:', path)
    alert(`[模拟] 打开所在文件夹：${path}`)
  },

  async SaveConfig(config: any) {
    mockConfig = { ...config }
    console.log('[Mock] Saved config:', config)
  },

  async Search(query: string, opts: any) {
    console.log('[Mock] Search:', query, opts)

    // 添加到历史
    if (query && !mockHistory.includes(query)) {
      mockHistory.unshift(query)
      if (mockHistory.length > mockConfig.maxHistory) {
        mockHistory.pop()
      }
    }

    // 模拟搜索过滤
    const lowerQuery = query.toLowerCase()
    return MOCK_FILES.filter(file => {
      // 仅 App 和文件夹过滤
      if (opts?.appAndFolderOnly) {
        if (!file.isFolder && !isAppFile(file.fileName)) {
          return false
        }
      }

      const target = opts?.matchPath ? file.fullPath : file.fileName
      const searchIn = opts?.matchCase ? target : target.toLowerCase()
      const searchFor = opts?.matchCase ? query : lowerQuery

      if (opts?.useRegex) {
        try {
          return new RegExp(searchFor).test(searchIn)
        } catch {
          return false
        }
      }

      if (opts?.matchWholeWord) {
        return new RegExp(`\\b${searchFor}\\b`).test(searchIn)
      }

      return searchIn.includes(searchFor)
    })
  },

  async ShowWindow() {
    console.log('[Mock] ShowWindow called')
  },

  async ToggleWindow() {
    console.log('[Mock] ToggleWindow called')
  }
}

/**
 * 初始化 mock，仅在 window.go 不存在时（即浏览器环境）注入
 */
export function setupWailsMock() {
  if (typeof window === 'undefined') return

  if (!(window as any).go) {
    console.log('[Mock] Wails runtime not detected, injecting mock data for preview')
    ;(window as any).go = {
      main: {
        App: mockApp
      }
    }

    // Mock Wails runtime API
    ;(window as any).runtime = {
      ClipboardSetText: async (text: string) => {
        try {
          await navigator.clipboard.writeText(text)
        } catch (e) {
          console.warn('[Mock] Clipboard failed:', e)
        }
      },
      LogInfo: (msg: string) => console.log('[Wails LogInfo]', msg),
      LogError: (msg: string) => console.error('[Wails LogError]', msg),
      LogWarning: (msg: string) => console.warn('[Wails LogWarning]', msg),
      WindowShow: () => console.log('[Mock] WindowShow'),
      WindowHide: () => console.log('[Mock] WindowHide'),
      WindowSetAlwaysOnTop: (b: boolean) => console.log('[Mock] WindowSetAlwaysOnTop:', b)
    }
  }
}
