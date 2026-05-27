<script lang="ts" setup>
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue'
import SearchBar from './components/SearchBar.vue'
import ResultList from './components/ResultList.vue'
import HistoryPanel from './components/HistoryPanel.vue'

import { Search, GetHistory, ShowWindow, HideWindow, OpenFile } from '../wailsjs/go/main/App'

// Types
interface SearchResult {
  fileName: string
  fullPath: string
  isFolder: boolean
  size: number
  dateModified: string
}

// State
const query = ref('')
const results = ref<SearchResult[]>([])
const history = ref<string[]>([])
const showHistory = ref(false)
const selectedIndex = ref(-1)
const isLoading = ref(false)
const errorMessage = ref('')

// Search options
const matchCase = ref(false)
const matchWholeWord = ref(false)
const useRegex = ref(false)

// Debounce timer
let searchTimer: ReturnType<typeof setTimeout> | null = null

// Perform search
async function performSearch() {
  if (!query.value.trim()) {
    results.value = []
    return
  }

  isLoading.value = true
  errorMessage.value = ''

  try {
    const opts = {
      matchCase: matchCase.value,
      matchWholeWord: matchWholeWord.value,
      useRegex: useRegex.value,
      maxResults: 100
    }
    
    const searchResults = await Search(query.value, opts)
    results.value = searchResults || []
    selectedIndex.value = results.value.length > 0 ? 0 : -1
    showHistory.value = false
  } catch (err: any) {
    errorMessage.value = err.message || 'Search failed'
    results.value = []
  } finally {
    isLoading.value = false
  }
}

// Watch query with debounce
watch(query, () => {
  if (searchTimer) clearTimeout(searchTimer)
  searchTimer = setTimeout(performSearch, 300)
})

// Load history
async function loadHistory() {
  try {
    history.value = await GetHistory() || []
  } catch (e) {
    console.error('Failed to load history:', e)
  }
}

// Handle history item click
function selectHistoryItem(item: string) {
  query.value = item
  showHistory.value = false
  nextTick(() => performSearch())
}

// Handle result selection - opens the file
async function selectResult(result: SearchResult) {
  try {
    await OpenFile(result.fullPath)
  } catch (e) {
    console.error('Failed to open file:', e)
  }
}

// Handle keyboard navigation
function handleKeyDown(e: KeyboardEvent) {
  if (e.key === 'Escape') {
    HideWindow()
    return
  }

  if (e.key === 'ArrowDown') {
    e.preventDefault()
    if (selectedIndex.value < results.value.length - 1) {
      selectedIndex.value++
    }
    return
  }

  if (e.key === 'ArrowUp') {
    e.preventDefault()
    if (selectedIndex.value > 0) {
      selectedIndex.value--
    }
    return
  }

  if (e.key === 'Enter' && selectedIndex.value >= 0) {
    e.preventDefault()
    const result = results.value[selectedIndex.value]
    if (result) {
      selectResult(result)
    }
    return
  }

  // Show history on Ctrl+Up
  if (e.key === 'ArrowUp' && e.ctrlKey) {
    e.preventDefault()
    loadHistory()
    showHistory.value = !showHistory.value
    return
  }
}

// Window focus handling
function handleWindowBlur() {
  HideWindow()
}

// Lifecycle
onMounted(() => {
  loadHistory()
  window.addEventListener('keydown', handleKeyDown)
  window.addEventListener('blur', handleWindowBlur)
  ShowWindow()
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
  window.removeEventListener('blur', handleWindowBlur)
})
</script>

<template>
  <div class="app-container" :class="{ 'is-loading': isLoading }">
    <SearchBar
      v-model="query"
      :match-case="matchCase"
      :match-whole-word="matchWholeWord"
      :use-regex="useRegex"
      :show-history="showHistory"
      @toggle-history="showHistory = !showHistory"
      @update:match-case="matchCase = $event"
      @update:match-whole-word="matchWholeWord = $event"
      @update:use-regex="useRegex = $event"
    />
    
    <HistoryPanel
      v-if="showHistory && history.length > 0"
      :items="history"
      @select="selectHistoryItem"
    />
    
    <div v-if="errorMessage" class="error-message">
      {{ errorMessage }}
    </div>
    
    <ResultList
      v-else-if="results.length > 0"
      :results="results"
      :selected-index="selectedIndex"
      @select="selectResult"
    />
    
    <div v-else-if="query && !isLoading" class="empty-state">
      No results found
    </div>
  </div>
</template>

<style>
:root {
  --bg-primary: #1e1e1e;
  --bg-secondary: #2d2d2d;
  --bg-hover: #3d3d3d;
  --text-primary: #e0e0e0;
  --text-secondary: #a0a0a0;
  --accent: #0078d4;
  --accent-hover: #1a86d9;
  --border: #404040;
  --error: #f44336;
}

* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

html, body {
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  font-size: 14px;
  background: transparent;
  overflow: hidden;
}

.app-container {
  width: 100%;
  height: 100vh;
  background: var(--bg-primary);
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.error-message {
  padding: 16px;
  color: var(--error);
  text-align: center;
}

.empty-state {
  padding: 32px;
  color: var(--text-secondary);
  text-align: center;
}

.is-loading .app-container::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 2px;
  background: var(--accent);
  animation: loading 1s ease-in-out infinite;
}

@keyframes loading {
  0% { transform: translateX(-100%); }
  100% { transform: translateX(100%); }
}
</style>

