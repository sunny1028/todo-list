<script setup lang="ts">
const emit = defineEmits<{ close: [] }>()

const shortcuts = [
  { key: 'Ctrl+N', desc: '聚焦新建输入框' },
  { key: 'Ctrl+K', desc: '聚焦搜索' },
  { key: 'Esc', desc: '清除筛选并失焦' },
  { key: '?', desc: '打开/关闭快捷键面板' },
]
</script>

<template>
  <Teleport to="body">
    <div
      class="fixed inset-0 z-50 flex items-center justify-center bg-black/30 backdrop-blur-sm"
      @click.self="emit('close')"
    >
      <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-2xl p-6 w-[90%] max-w-sm mx-auto">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-bold text-gray-800 dark:text-gray-100">快捷键</h3>
          <button @click="emit('close')" class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><line x1="18" y1="6" x2="6" y2="18"/><line x1="6" y1="6" x2="18" y2="18"/></svg>
          </button>
        </div>
        <div class="space-y-2">
          <div v-for="s in shortcuts" :key="s.key" class="flex items-center justify-between text-sm">
            <span class="text-gray-500 dark:text-gray-400">{{ s.desc }}</span>
            <kbd class="px-2 py-0.5 text-xs rounded-md bg-gray-100 dark:bg-gray-700 text-gray-600 dark:text-gray-300 font-mono">{{ s.key }}</kbd>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>
