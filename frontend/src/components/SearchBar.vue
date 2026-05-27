<script lang="ts" setup>
import { ref, computed, onMounted } from 'vue'

const props = defineProps<{
  modelValue: string
  matchCase: boolean
  matchWholeWord: boolean
  useRegex: boolean
  showHistory: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'update:matchCase': [value: boolean]
  'update:matchWholeWord': [value: boolean]
  'update:useRegex': [value: boolean]
  'toggleHistory': []
}>()

const inputRef = ref<HTMLInputElement | null>(null)

function updateValue(e: Event) {
  const target = e.target as HTMLInputElement
  emit('update:modelValue', target.value)
}

function toggleMatchCase() {
  emit('update:matchCase', !props.matchCase)
}

function toggleMatchWholeWord() {
  emit('update:matchWholeWord', !props.matchWholeWord)
}

function toggleUseRegex() {
  emit('update:useRegex', !props.useRegex)
}

function focusInput() {
  inputRef.value?.focus()
}

// Auto-focus on mount
onMounted(() => {
  focusInput()
})
</script>

<template>
  <div class="search-bar">
    <div class="search-input-wrapper">
      <input
        ref="inputRef"
        type="text"
        class="search-input"
        :value="modelValue"
        @input="updateValue"
        placeholder="Search files..."
        autofocus
      />
      <button
        class="history-toggle"
        :class="{ active: showHistory }"
        @click="emit('toggleHistory')"
        title="Toggle search history (Ctrl+Up)"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <path d="M12 8v4l3 3"/>
          <circle cx="12" cy="12" r="10"/>
        </svg>
      </button>
    </div>
    
    <div class="search-options">
      <button
        class="option-btn"
        :class="{ active: matchCase }"
        @click="toggleMatchCase"
        title="Match case"
      >
        Aa
      </button>
      <button
        class="option-btn"
        :class="{ active: matchWholeWord }"
        @click="toggleMatchWholeWord"
        title="Match whole word"
      >
        W
      </button>
      <button
        class="option-btn"
        :class="{ active: useRegex }"
        @click="toggleUseRegex"
        title="Use regular expression"
      >
        .*
      </button>
    </div>
  </div>
</template>

<style scoped>
.search-bar {
  display: flex;
  flex-direction: column;
  padding: 8px;
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border);
}

.search-input-wrapper {
  display: flex;
  align-items: center;
  gap: 8px;
}

.search-input {
  flex: 1;
  padding: 10px 12px;
  font-size: 16px;
  border: none;
  border-radius: 4px;
  background: var(--bg-primary);
  color: var(--text-primary);
  outline: none;
}

.search-input::placeholder {
  color: var(--text-secondary);
}

.search-input:focus {
  box-shadow: 0 0 0 2px var(--accent);
}

.history-toggle {
  padding: 8px;
  border: none;
  border-radius: 4px;
  background: transparent;
  color: var(--text-secondary);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.history-toggle:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.history-toggle.active {
  background: var(--accent);
  color: white;
}

.search-options {
  display: flex;
  gap: 4px;
  margin-top: 8px;
}

.option-btn {
  padding: 4px 8px;
  font-size: 12px;
  border: 1px solid var(--border);
  border-radius: 4px;
  background: transparent;
  color: var(--text-secondary);
  cursor: pointer;
  transition: all 0.2s;
}

.option-btn:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.option-btn.active {
  background: var(--accent);
  color: white;
  border-color: var(--accent);
}
</style>