<script lang="ts" setup>
const props = defineProps<{
  items: string[]
}>()

const emit = defineEmits<{
  'select': [item: string]
}>()

function handleSelect(item: string) {
  emit('select', item)
}
</script>

<template>
  <div class="history-panel">
    <div class="history-header">
      <span>Search History</span>
      <span class="history-count">{{ items.length }} items</span>
    </div>
    
    <div class="history-list">
      <div
        v-for="(item, index) in items"
        :key="index"
        class="history-item"
        @click="handleSelect(item)"
      >
        <span class="history-icon">🔍</span>
        <span class="history-text">{{ item }}</span>
      </div>
    </div>
  </div>
</template>

<style scoped>
.history-panel {
  background: var(--bg-secondary);
  border-bottom: 1px solid var(--border);
  max-height: 200px;
  overflow-y: auto;
}

.history-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 12px;
  font-size: 12px;
  color: var(--text-secondary);
  border-bottom: 1px solid var(--border);
}

.history-count {
  opacity: 0.7;
}

.history-list {
  padding: 4px;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
  transition: background 0.15s;
}

.history-item:hover {
  background: var(--bg-hover);
}

.history-icon {
  font-size: 14px;
}

.history-text {
  flex: 1;
  color: var(--text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>