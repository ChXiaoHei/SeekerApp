import { useEffect, useRef } from 'react'

interface SearchBarProps {
  value: string
  onChange: (value: string) => void
  matchCase: boolean
  matchWholeWord: boolean
  matchPath: boolean
  useRegex: boolean
  appAndFolderOnly: boolean
  showHistory: boolean
  onToggleHistory: () => void
  onMatchCaseChange: (value: boolean) => void
  onMatchWholeWordChange: (value: boolean) => void
  onMatchPathChange: (value: boolean) => void
  onUseRegexChange: (value: boolean) => void
  onAppAndFolderOnlyChange: (value: boolean) => void
}

function SearchBar({
  value,
  onChange,
  matchCase,
  matchWholeWord,
  matchPath,
  useRegex,
  appAndFolderOnly,
  showHistory,
  onToggleHistory,
  onMatchCaseChange,
  onMatchWholeWordChange,
  onMatchPathChange,
  onUseRegexChange,
  onAppAndFolderOnlyChange
}: SearchBarProps) {
  const inputRef = useRef<HTMLInputElement>(null)

  useEffect(() => {
    inputRef.current?.focus()
  }, [])

  return (
    <div className="search-bar">
      <div className="search-input-wrapper">
        <input
          ref={inputRef}
          type="text"
          className="search-input"
          value={value}
          onChange={(e) => onChange(e.target.value)}
          placeholder="搜索应用、文件夹和文件..."
          autoFocus
        />
        <button
          className={`history-toggle ${showHistory ? 'active' : ''}`}
          onClick={onToggleHistory}
          title="切换搜索历史 (Ctrl+↑)"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
            <path d="M12 8v4l3 3" />
            <circle cx="12" cy="12" r="10" />
          </svg>
        </button>
      </div>

      <div className="search-options">
        <button
          className={`option-btn-primary ${appAndFolderOnly ? 'active' : ''}`}
          onClick={() => onAppAndFolderOnlyChange(!appAndFolderOnly)}
          title="仅搜索应用和文件夹"
        >
          <span>{appAndFolderOnly ? '☑' : '☐'}</span>
          <span>仅搜索应用和文件夹</span>
        </button>

        <div className="option-divider" />

        <button
          className={`option-btn ${matchCase ? 'active' : ''}`}
          onClick={() => onMatchCaseChange(!matchCase)}
          title="区分大小写"
        >
          Aa
        </button>
        <button
          className={`option-btn ${matchWholeWord ? 'active' : ''}`}
          onClick={() => onMatchWholeWordChange(!matchWholeWord)}
          title="全词匹配"
        >
          W
        </button>
        <button
          className={`option-btn ${matchPath ? 'active' : ''}`}
          onClick={() => onMatchPathChange(!matchPath)}
          title="匹配路径"
        >
          P
        </button>
        <button
          className={`option-btn ${useRegex ? 'active' : ''}`}
          onClick={() => onUseRegexChange(!useRegex)}
          title="使用正则表达式"
        >
          .*
        </button>
      </div>
    </div>
  )
}

export default SearchBar
