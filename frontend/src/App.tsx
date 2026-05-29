import { useState, useEffect, useCallback, useRef, useMemo } from 'react'
import SearchBar from './components/SearchBar'
import ResultList, { ResultGroup } from './components/ResultList'
import HistoryPanel from './components/HistoryPanel'
import { Search, GetHistory, ShowWindow, HideWindow, OpenFile } from '../wailsjs/go/main/App'
import './App.css'

interface SearchResult {
  fileName: string
  fullPath: string
  isFolder: boolean
  size: number
  dateModified: string
}

// 判断是否为 App / 可执行程序
const APP_EXTENSIONS = ['exe', 'msi', 'app', 'lnk', 'appx', 'dmg', 'pkg']
function isAppFile(fileName: string): boolean {
  const ext = fileName.split('.').pop()?.toLowerCase() || ''
  return APP_EXTENSIONS.includes(ext)
}

function App() {
  const [query, setQuery] = useState('')
  const [results, setResults] = useState<SearchResult[]>([])
  const [history, setHistory] = useState<string[]>([])
  const [showHistory, setShowHistory] = useState(false)
  const [selectedIndex, setSelectedIndex] = useState(-1)
  const [isLoading, setIsLoading] = useState(false)
  const [errorMessage, setErrorMessage] = useState('')

  // Search options
  const [matchCase, setMatchCase] = useState(false)
  const [matchWholeWord, setMatchWholeWord] = useState(false)
  const [matchPath, setMatchPath] = useState(false)
  const [useRegex, setUseRegex] = useState(false)
  const [appAndFolderOnly, setAppAndFolderOnly] = useState(true)

  const searchTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  // Perform search
  const performSearch = useCallback(async () => {
    if (!query.trim()) {
      setResults([])
      return
    }

    setIsLoading(true)
    setErrorMessage('')

    try {
      const opts = {
        matchCase,
        matchWholeWord,
        matchPath,
        useRegex,
        appAndFolderOnly,
        maxResults: 100
      }

      const searchResults = await Search(query, opts)
      setResults(searchResults || [])
      setSelectedIndex(searchResults && searchResults.length > 0 ? 0 : -1)
      setShowHistory(false)
    } catch (err: any) {
      setErrorMessage(err.message || 'Search failed')
      setResults([])
    } finally {
      setIsLoading(false)
    }
  }, [query, matchCase, matchWholeWord, matchPath, useRegex, appAndFolderOnly])

  // 按 App / 文件夹 / 文件 分组
  const groupedResults = useMemo<ResultGroup[]>(() => {
    const apps: SearchResult[] = []
    const folders: SearchResult[] = []
    const files: SearchResult[] = []

    for (const r of results) {
      if (r.isFolder) {
        folders.push(r)
      } else if (isAppFile(r.fileName)) {
        apps.push(r)
      } else {
        files.push(r)
      }
    }

    const groups: ResultGroup[] = []
    if (apps.length) groups.push({ title: '应用', icon: '🚀', items: apps })
    if (folders.length) groups.push({ title: '文件夹', icon: '📁', items: folders })
    if (files.length) groups.push({ title: '文件', icon: '📄', items: files })
    return groups
  }, [results])

  // 扁平化的结果列表，用于键盘导航
  const flatResults = useMemo(
    () => groupedResults.flatMap(g => g.items),
    [groupedResults]
  )

  // Debounced search
  useEffect(() => {
    if (searchTimerRef.current) clearTimeout(searchTimerRef.current)
    searchTimerRef.current = setTimeout(performSearch, 300)
    return () => {
      if (searchTimerRef.current) clearTimeout(searchTimerRef.current)
    }
  }, [query, performSearch])

  // Load history
  const loadHistory = useCallback(async () => {
    try {
      const hist = await GetHistory()
      setHistory(hist || [])
    } catch (e) {
      console.error('Failed to load history:', e)
    }
  }, [])

  // Handle history item click
  const selectHistoryItem = useCallback((item: string) => {
    setQuery(item)
    setShowHistory(false)
    setTimeout(() => performSearch(), 0)
  }, [performSearch])

  // Handle result selection
  const selectResult = useCallback(async (result: SearchResult) => {
    try {
      await OpenFile(result.fullPath)
    } catch (e) {
      console.error('Failed to open file:', e)
    }
  }, [])

  // Handle keyboard navigation
  const handleKeyDown = useCallback((e: KeyboardEvent) => {
    if (e.key === 'Escape') {
      HideWindow()
      return
    }

    // Show history on Ctrl+Up（需先判断，避免被 ArrowUp 拦截）
    if (e.key === 'ArrowUp' && e.ctrlKey) {
      e.preventDefault()
      loadHistory()
      setShowHistory(prev => !prev)
      return
    }

    if (e.key === 'ArrowDown') {
      e.preventDefault()
      setSelectedIndex(prev => Math.min(prev + 1, flatResults.length - 1))
      return
    }

    if (e.key === 'ArrowUp') {
      e.preventDefault()
      setSelectedIndex(prev => Math.max(prev - 1, 0))
      return
    }

    if (e.key === 'Enter' && selectedIndex >= 0) {
      e.preventDefault()
      const result = flatResults[selectedIndex]
      if (result) {
        selectResult(result)
      }
      return
    }
  }, [flatResults, selectedIndex, selectResult, loadHistory])

  // Window blur handler
  const handleWindowBlur = useCallback(() => {
    HideWindow()
  }, [])

  // Lifecycle
  useEffect(() => {
    loadHistory()
    window.addEventListener('keydown', handleKeyDown)
    window.addEventListener('blur', handleWindowBlur)
    ShowWindow()

    return () => {
      window.removeEventListener('keydown', handleKeyDown)
      window.removeEventListener('blur', handleWindowBlur)
    }
  }, [loadHistory, handleKeyDown, handleWindowBlur])

  return (
    <div className={`app-container ${isLoading ? 'is-loading' : ''}`}>
      <div className="main-panel">
        <SearchBar
          value={query}
          onChange={setQuery}
          matchCase={matchCase}
          matchWholeWord={matchWholeWord}
          matchPath={matchPath}
          useRegex={useRegex}
          appAndFolderOnly={appAndFolderOnly}
          showHistory={showHistory}
          onToggleHistory={() => setShowHistory(prev => !prev)}
          onMatchCaseChange={setMatchCase}
          onMatchWholeWordChange={setMatchWholeWord}
          onMatchPathChange={setMatchPath}
          onUseRegexChange={setUseRegex}
          onAppAndFolderOnlyChange={setAppAndFolderOnly}
        />

        {showHistory && history.length > 0 && (
          <HistoryPanel items={history} onSelect={selectHistoryItem} />
        )}

        {errorMessage && (
          <div className="error-message">{errorMessage}</div>
        )}

        {!errorMessage && flatResults.length > 0 && (
          <ResultList
            groups={groupedResults}
            selectedIndex={selectedIndex}
            onSelect={selectResult}
          />
        )}

        {!errorMessage && query && !isLoading && flatResults.length === 0 && (
          <div className="empty-state">未找到相关结果</div>
        )}
      </div>
    </div>
  )
}

export default App
