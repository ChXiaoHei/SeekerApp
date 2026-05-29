interface HistoryPanelProps {
  items: string[]
  onSelect: (item: string) => void
}

function HistoryPanel({ items, onSelect }: HistoryPanelProps) {
  return (
    <div className="history-panel">
      <div className="history-header">
        <span>搜索历史</span>
        <span className="history-count">{items.length} 条记录</span>
      </div>

      <div className="history-list">
        {items.map((item, index) => (
          <div
            key={index}
            className="history-item"
            onClick={() => onSelect(item)}
          >
            <span className="history-icon">🔍</span>
            <span className="history-text">{item}</span>
          </div>
        ))}
      </div>
    </div>
  )
}

export default HistoryPanel