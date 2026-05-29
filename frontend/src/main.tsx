import React from 'react'
import ReactDOM from 'react-dom/client'
import App from './App'
import { setupWailsMock } from './mock/wailsMock'

// 浏览器环境下注入 mock，方便 UI 预览（Wails 应用中此函数会自动跳过）
setupWailsMock()

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>
)
