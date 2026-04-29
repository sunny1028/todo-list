<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useTodos } from '../stores/todo'
import type { Stats } from '../types/todo'

const store = useTodos()
const stats = ref<Stats | null>(null)

const priorityLabels: Record<string, string> = { low: '低', medium: '中', high: '高' }
const priorityColors: Record<string, string> = { low: 'bg-blue-500', medium: 'bg-yellow-500', high: 'bg-red-500' }

onMounted(async () => {
  stats.value = await store.fetchStats()
})

function pct(part: number, total: number) {
  if (total === 0) return 0
  return Math.round((part / total) * 100)
}

function tagBarWidth(count: number, max: number) {
  if (max === 0) return '0%'
  return `${Math.round((count / max) * 100)}%`
}
</script>

<template>
  <div v-if="stats" class="max-w-2xl mx-auto">
    <h2 class="text-lg font-bold text-gray-800 dark:text-gray-100 mb-6">统计面板</h2>

    <!-- Overview cards -->
    <div class="grid grid-cols-3 gap-3 mb-6">
      <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-4 text-center">
        <div class="text-2xl font-bold text-gray-800 dark:text-gray-100">{{ stats.total }}</div>
        <div class="text-xs text-gray-400 dark:text-gray-500 mt-1">全部</div>
      </div>
      <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-4 text-center">
        <div class="text-2xl font-bold text-indigo-500">{{ stats.active }}</div>
        <div class="text-xs text-gray-400 dark:text-gray-500 mt-1">未完成</div>
      </div>
      <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-4 text-center">
        <div class="text-2xl font-bold text-emerald-500">{{ stats.completed }}</div>
        <div class="text-xs text-gray-400 dark:text-gray-500 mt-1">已完成</div>
      </div>
    </div>

    <!-- Progress bar -->
    <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-5 mb-6">
      <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-3">完成率</h3>
      <div class="flex items-center gap-3">
        <div class="flex-1 h-3 bg-gray-100 dark:bg-gray-800 rounded-full overflow-hidden">
          <div
            class="h-full bg-emerald-500 rounded-full transition-all duration-500"
            :style="{ width: pct(stats.completed, stats.total) + '%' }"
          />
        </div>
        <span class="text-sm font-bold text-gray-600 dark:text-gray-400">{{ pct(stats.completed, stats.total) }}%</span>
      </div>
    </div>

    <!-- Priority distribution -->
    <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-5 mb-6">
      <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-3">优先级分布</h3>
      <div class="space-y-3">
        <div v-for="(label, key) in priorityLabels" :key="key" class="flex items-center gap-3">
          <span class="w-12 text-xs text-gray-500 dark:text-gray-400">{{ label }}</span>
          <div class="flex-1 h-5 bg-gray-100 dark:bg-gray-800 rounded-full overflow-hidden">
            <div
              :class="priorityColors[key]"
              class="h-full rounded-full transition-all"
              :style="{ width: pct(stats.by_priority?.[key] || 0, stats.total) + '%' }"
            />
          </div>
          <span class="text-xs font-semibold text-gray-500 dark:text-gray-400 w-8 text-right">{{ stats.by_priority?.[key] || 0 }}</span>
        </div>
      </div>
    </div>

    <!-- Tag distribution -->
    <div v-if="Object.keys(stats.by_tag || {}).length > 0" class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-5">
      <h3 class="text-sm font-semibold text-gray-700 dark:text-gray-300 mb-3">标签分布</h3>
      <div class="space-y-2">
        <div v-for="(count, tag) in stats.by_tag" :key="tag" class="flex items-center gap-2">
          <span class="w-20 text-xs text-gray-500 dark:text-gray-400 truncate">{{ tag }}</span>
          <div class="flex-1 h-4 bg-gray-100 dark:bg-gray-800 rounded-full overflow-hidden">
            <div
              class="h-full bg-indigo-400 rounded-full"
              :style="{ width: tagBarWidth(count, Math.max(...Object.values(stats.by_tag || {}))) }"
            />
          </div>
          <span class="text-xs text-gray-400 dark:text-gray-500 w-6 text-right">{{ count }}</span>
        </div>
      </div>
    </div>
  </div>

  <div v-else class="text-center py-20 text-gray-400">加载中…</div>
</template>
