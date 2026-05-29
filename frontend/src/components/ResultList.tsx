import { useState, useEffect, useCallback } from 'react'
import { OpenFile, OpenFolder, CopyPathToClipboard } from '../../wailsjs/go/main/App'

interface SearchResult {
  fileName: string
  fullPath: string
  isFolder: boolean
  size: number
  dateModified: string
}

export interface ResultGroup {
  title: string
  icon: string
  items: SearchResult[]
}

interface ContextMenuState {
  visible: boolean
  x: number
  y: number
  result: SearchResult | null
}

interface ResultListProps {
  groups: ResultGroup[]
  selectedIndex: number
  onSelect: (result: SearchResult) => void
}

// 根据扩展名返回 emoji 图标
function getFileIcon(result: SearchResult): string {
  if (result.isFolder) return '📁'
  const ext = result.fileName.split('.').pop()?.toLowerCase() || ''
  switch (ext) {
    case 'app':
    case 'dmg':
    case 'pkg':
      return '🍎'
    case 'exe':
    case 'msi':
    case 'lnk':
    case 'appx':
      return '🚀'
    case 'jpg':
    case 'jpeg':
    case 'png':
    case 'gif':
    case 'bmp':
    case 'webp':
      return '🖼️'
    case 'mp3':
    case 'wav':
    case 'flac':
      return '🎵'
    case 'mp4':
    case 'avi':
    case 'mkv':
    case 'mov':
      return '🎬'
    case 'pdf':
      return '📕'
    case 'doc':
    case 'docx':
      return '📘'
    case 'xls':
    case 'xlsx':
      return '📗'
    case 'ppt':
    case 'pptx':
      return '📙'
    case 'zip':
    case 'rar':
    case '7z':
    case 'tar':
    case 'gz':
      return '🗜️'
    case 'js':
    case 'ts':
    case 'tsx':
    case 'jsx':
    case 'py':
    case 'go':
    case 'java':
    case 'c':
    case 'cpp':
    case 'rs':
      return '💻'
    case 'md':
      return '📝'
    case 'json':
      return '⚙️'
    default:
      return '📄'
  }
}

function formatSize(bytes: number, isFolder: boolean): string {
  if (isFolder) return ''
  if (bytes === 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${units[i]}`
}

function formatDate(dateStr: string): string {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

function ResultList({ groups, selectedIndex, onSelect }: ResultListProps) {
  const [contextMenu, setContextMenu] = useState<ContextMenuState>({
    visible: false,
    x: 0,
    y: 0,
    result: null
  })

  const handleOpen = useCallback(async (result: SearchResult) => {
    try {
      await OpenFile(result.fullPath)
    } catch (e) {
      console.error('Failed to open file:', e)
    }
  }, [])

  const handleContextMenu = useCallback((e: React.MouseEvent, result: SearchResult) => {
    e.preventDefault()
    setContextMenu({
      visible: true,
      x: e.clientX,
      y: e.clientY,
      result
    })
  }, [])

  const closeContextMenu = useCallback(() => {
    setContextMenu(prev => ({ ...prev, visible: false, result: null }))
  }, [])

  const contextOpen = useCallback(async () => {
    if (contextMenu.result) {
      await handleOpen(contextMenu.result)
    }
    closeContextMenu()
  }, [contextMenu.result, handleOpen, closeContextMenu])

  const contextOpenFolder = useCallback(async () => {
    if (contextMenu.result) {
      try {
        await OpenFolder(contextMenu.result.fullPath)
      } catch (e) {
        console.error('Failed to open folder:', e)
      }
    }
    closeContextMenu()
  }, [contextMenu.result, closeContextMenu])

  const contextCopyPath = useCallback(async () => {
    if (contextMenu.result) {
      try {
        await CopyPathToClipboard(contextMenu.result.fullPath)
      } catch (e) {
        console.error('Failed to copy path:', e)
      }
    }
    closeContextMenu()
  }, [contextMenu.result, closeContextMenu])

  // Handle click outside and escape key
  useEffect(() => {
    const handleDocumentClick = (e: MouseEvent) => {
      if (contextMenu.visible) {
        const target = e.target as HTMLElement
        if (!target.closest('.context-menu')) {
          closeContextMenu()
        }
      }
    }

    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === 'Escape' && contextMenu.visible) {
        closeContextMenu()
      }
    }

    document.addEventListener('click', handleDocumentClick)
    document.addEventListener('keydown', handleKeyDown)

    return () => {
      document.removeEventListener('click', handleDocumentClick)
      document.removeEventListener('keydown', handleKeyDown)
    }
  }, [contextMenu.visible, closeContextMenu])

  // 按扁平索引计算每个 item 的全局序号
  let flatCursor = 0

  return (
    <div className="result-list">
      {groups.map((group) => (
        <div key={group.title} className="result-group">
          <div className="result-group-title">
            <span>{group.icon}</span>
            <span>{group.title}</span>
            <span className="result-group-count">({group.items.length})</span>
          </div>
          {group.items.map((result) => {
            const flatIndex = flatCursor++
            const isSelected = flatIndex === selectedIndex
            return (
              <div
                key={result.fullPath}
                className={`result-item ${isSelected ? 'selected' : ''}`}
                onClick={() => onSelect(result)}
                onDoubleClick={() => handleOpen(result)}
                onContextMenu={(e) => handleContextMenu(e, result)}
              >
                <div className="result-icon">
                  <span>{getFileIcon(result)}</span>
                </div>

                <div className="result-info">
                  <div className="result-name">{result.fileName}</div>
                  <div className="result-path">{result.fullPath}</div>
                </div>

                <div className="result-meta">
                  <span className="result-size">{formatSize(result.size, result.isFolder)}</span>
                  <span className="result-date">{formatDate(result.dateModified)}</span>
                </div>
              </div>
            )
          })}
        </div>
      ))}

      {/* Context Menu */}
      {contextMenu.visible && (
        <div
          className="context-menu"
          style={{ left: contextMenu.x, top: contextMenu.y }}
        >
          <div className="context-menu-item" onClick={contextOpen}>
            <span className="context-menu-icon">📂</span>
            <span>打开</span>
          </div>
          <div className="context-menu-item" onClick={contextOpenFolder}>
            <span className="context-menu-icon">📁</span>
            <span>打开所在文件夹</span>
          </div>
          <div className="context-menu-separator" />
          <div className="context-menu-item" onClick={contextCopyPath}>
            <span className="context-menu-icon">📋</span>
            <span>复制完整路径</span>
          </div>
        </div>
      )}
    </div>
  )
}

export default ResultList
