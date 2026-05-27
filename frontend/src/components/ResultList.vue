<script lang="ts" setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { OpenFile, OpenFolder, CopyPathToClipboard } from '../../wailsjs/go/main/App'

interface SearchResult {
  fileName: string
  fullPath: string
  isFolder: boolean
  size: number
  dateModified: string
}

interface ContextMenuState {
  visible: boolean
  x: number
  y: number
  result: SearchResult | null
}

const props = defineProps<{
  results: SearchResult[]
  selectedIndex: number
}>()

const emit = defineEmits<{
  'select': [result: SearchResult]
}>()

// Context menu state
const contextMenu = ref<ContextMenuState>({
  visible: false,
  x: 0,
  y: 0,
  result: null
})

// Format file size
function formatSize(bytes: number): string {
  if (bytes === 0) return '0 B'
  if (props.results[props.selectedIndex]?.isFolder) return ''
  
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return `${(bytes / Math.pow(1024, i)).toFixed(1)} ${units[i]}`
}

// Format date
function formatDate(dateStr: string): string {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
}

// Handle double-click to open file
async function handleOpen(result: SearchResult) {
  try {
    await OpenFile(result.fullPath)
  } catch (e) {
    console.error('Failed to open file:', e)
  }
}

// Handle right-click to show context menu
async function handleContextMenu(e: MouseEvent, result: SearchResult) {
  e.preventDefault()
  contextMenu.value = {
    visible: true,
    x: e.clientX,
    y: e.clientY,
    result: result
  }
}

// Close context menu
function closeContextMenu() {
  contextMenu.value.visible = false
  contextMenu.value.result = null
}

// Context menu actions
async function contextOpen() {
  if (contextMenu.value.result) {
    await handleOpen(contextMenu.value.result)
  }
  closeContextMenu()
}

async function contextOpenFolder() {
  if (contextMenu.value.result) {
    try {
      await OpenFolder(contextMenu.value.result.fullPath)
    } catch (e) {
      console.error('Failed to open folder:', e)
    }
  }
  closeContextMenu()
}

async function contextCopyPath() {
  if (contextMenu.value.result) {
    try {
      await CopyPathToClipboard(contextMenu.value.result.fullPath)
    } catch (e) {
      console.error('Failed to copy path:', e)
    }
  }
  closeContextMenu()
}

// Handle click outside to close context menu
function handleDocumentClick(e: MouseEvent) {
  if (contextMenu.value.visible) {
    const target = e.target as HTMLElement
    if (!target.closest('.context-menu')) {
      closeContextMenu()
    }
  }
}

// Handle escape key to close context menu
function handleKeyDown(e: KeyboardEvent) {
  if (e.key === 'Escape' && contextMenu.value.visible) {
    closeContextMenu()
  }
}

// Lifecycle hooks
onMounted(() => {
  document.addEventListener('click', handleDocumentClick)
  document.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  document.removeEventListener('click', handleDocumentClick)
  document.removeEventListener('keydown', handleKeyDown)
})

// Get icon based on file type
function getIconClass(result: SearchResult): string {
  if (result.isFolder) return 'icon-folder'
  
  const ext = result.fileName.split('.').pop()?.toLowerCase() || ''
  switch (ext) {
    case 'exe':
    case 'msi':
      return 'icon-executable'
    case 'jpg':
    case 'jpeg':
    case 'png':
    case 'gif':
    case 'bmp':
    case 'webp':
      return 'icon-image'
    case 'mp3':
    case 'wav':
    case 'flac':
    case 'ogg':
      return 'icon-audio'
    case 'mp4':
    case 'avi':
    case 'mkv':
    case 'mov':
    case 'webm':
      return 'icon-video'
    case 'pdf':
      return 'icon-pdf'
    case 'doc':
    case 'docx':
    case 'xls':
    case 'xlsx':
    case 'ppt':
    case 'pptx':
      return 'icon-document'
    case 'zip':
    case 'rar':
    case '7z':
    case 'tar':
    case 'gz':
      return 'icon-archive'
    case 'js':
    case 'ts':
    case 'py':
    case 'go':
    case 'java':
    case 'c':
    case 'cpp':
    case 'rs':
    case 'vue':
    case 'jsx':
    case 'tsx':
      return 'icon-code'
    default:
      return 'icon-file'
  }
}
</script>

<template>
  <div class="result-list">
    <div
      v-for="(result, index) in results"
      :key="result.fullPath"
      class="result-item"
      :class="{ selected: index === selectedIndex }"
      @click="emit('select', result)"
      @dblclick="handleOpen(result)"
      @contextmenu="handleContextMenu($event, result)"
    >
      <div class="result-icon" :class="getIconClass(result)">
        <span v-if="result.isFolder">📁</span>
        <span v-else>📄</span>
      </div>
      
      <div class="result-info">
        <div class="result-name">{{ result.fileName }}</div>
        <div class="result-path">{{ result.fullPath }}</div>
      </div>
      
      <div class="result-meta">
        <span class="result-size">{{ formatSize(result.size) }}</span>
        <span class="result-date">{{ formatDate(result.dateModified) }}</span>
      </div>
    </div>
    
    <!-- Context Menu -->
    <div
      v-if="contextMenu.visible"
      class="context-menu"
      :style="{ left: contextMenu.x + 'px', top: contextMenu.y + 'px' }"
    >
      <div class="context-menu-item" @click="contextOpen">
        <span class="context-menu-icon">📂</span>
        <span>Open</span>
      </div>
      <div class="context-menu-item" @click="contextOpenFolder">
        <span class="context-menu-icon">📁</span>
        <span>Open containing folder</span>
      </div>
      <div class="context-menu-separator"></div>
      <div class="context-menu-item" @click="contextCopyPath">
        <span class="context-menu-icon">📋</span>
        <span>Copy full path</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.result-list {
  flex: 1;
  overflow-y: auto;
  padding: 4px;
  position: relative;
}

.result-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  border-radius: 4px;
  cursor: pointer;
  transition: background 0.15s;
}

.result-item:hover {
  background: var(--bg-hover);
}

.result-item.selected {
  background: var(--accent);
  color: white;
}

.result-icon {
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 18px;
}

.result-info {
  flex: 1;
  min-width: 0;
}

.result-name {
  font-weight: 500;
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.result-item.selected .result-name {
  color: white;
}

.result-path {
  font-size: 12px;
  color: var(--text-secondary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.result-item.selected .result-path {
  color: rgba(255, 255, 255, 0.7);
}

.result-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 2px;
  font-size: 12px;
  color: var(--text-secondary);
}

.result-item.selected .result-meta {
  color: rgba(255, 255, 255, 0.7);
}

.result-size {
  font-weight: 500;
}

.result-date {
  opacity: 0.7;
}

/* Context Menu Styles */
.context-menu {
  position: fixed;
  background: var(--bg-secondary);
  border: 1px solid var(--border);
  border-radius: 4px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
  min-width: 180px;
  z-index: 1000;
  padding: 4px 0;
}

.context-menu-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  cursor: pointer;
  color: var(--text-primary);
  transition: background 0.15s;
}

.context-menu-item:hover {
  background: var(--accent);
  color: white;
}

.context-menu-icon {
  font-size: 14px;
  width: 20px;
  text-align: center;
}

.context-menu-separator {
  height: 1px;
  background: var(--border);
  margin: 4px 0;
}
</style>