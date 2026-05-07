<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'
import { useTodos } from '../stores/todo'
import type { Stats, ReviewStats, FocusStats } from '../types/todo'
import * as api from '../api/todos'
import * as focusApi from '../api/focus'

const store = useTodos()
const stats = ref<Stats | null>(null)
const review = ref<ReviewStats | null>(null)
const focusStats = ref<FocusStats | null>(null)

const priorityLabels: Record<string, string> = { low: '低', medium: '中', high: '高' }
const priorityColors: Record<string, string> = { low: 'bg-blue-500', medium: 'bg-yellow-500', high: 'bg-red-500' }

async function load() {
  stats.value = await store.fetchStats()
  focusApi.getFocusStats().then(r => { focusStats.value = r.data }).catch(() => {})
  api.getReviewStats(store.currentListId || undefined).then(r => { review.value = r.data }).catch(() => {})
}

onMounted(load)
watch(() => store.currentListId, load)

function pct(part: number, total: number) {
  if (total === 0) return 0
  return Math.round((part / total) * 100)
}

function tagBarWidth(count: number, max: number) {
  if (max === 0) return '0%'
  return `${Math.round((count / max) * 100)}%`
}

function focusBarHeight(minutes: number, max: number) {
  if (max === 0) return '4px'
  return `${Math.max(4, Math.round((minutes / max) * 100))}%`
}

const maxFocusMins = () => {
  if (!review.value?.daily_focus?.length) return 1
  return Math.max(...review.value.daily_focus.map(d => d.minutes), 1)
}

// Score ring
const scoreColor = () => {
  if (!review.value) return 'text-gray-300'
  const s = review.value.mastery_score
  if (s >= 80) return 'text-emerald-500'
  if (s >= 50) return 'text-amber-500'
  return 'text-red-400'
}
const scoreLabel = () => {
  if (!review.value) return ''
  const s = review.value.mastery_score
  if (s >= 80) return '优秀'
  if (s >= 50) return '良好'
  return '待提升'
}
</script>

<template>
  <div v-if="stats" class="max-w-2xl mx-auto">
    <router-link to="/" class="inline-flex items-center gap-1 text-xs text-gray-400 hover:text-indigo-500 transition-colors mb-3">
      <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polyline points="15 18 9 12 15 6"/></svg>
      返回列表
    </router-link>
    <h2 class="text-lg font-bold text-gray-800 dark:text-gray-100 mb-6">统计面板</h2>

    <!-- Row 1: Overview + Mastery Score -->
    <div class="grid grid-cols-4 gap-3 mb-4">
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
      <div v-if="review" class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-3 text-center flex flex-col items-center justify-center">
        <div class="text-xl font-bold" :class="scoreColor()">{{ review.mastery_score }}</div>
        <div class="text-[10px] text-gray-400">掌控评分</div>
        <div class="text-[10px] mt-0.5" :class="scoreColor()">{{ scoreLabel() }}</div>
      </div>
    </div>

    <!-- Row 2: Focus stats -->
    <div v-if="focusStats" class="grid grid-cols-4 gap-3 mb-4">
      <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-3 text-center">
        <div class="text-lg font-bold text-indigo-500">{{ focusStats.today_minutes }}</div>
        <div class="text-[10px] text-gray-400 mt-0.5">今日专注(分)</div>
      </div>
      <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-3 text-center">
        <div class="text-lg font-bold text-gray-800 dark:text-gray-100">{{ focusStats.total_sessions }}</div>
        <div class="text-[10px] text-gray-400 mt-0.5">总次数</div>
      </div>
      <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-3 text-center">
        <div class="text-lg font-bold text-emerald-500">{{ focusStats.streak_days }}</div>
        <div class="text-[10px] text-gray-400 mt-0.5">连续天数</div>
      </div>
      <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-3 text-center">
        <div class="text-lg font-bold text-gray-800 dark:text-gray-100">{{ focusStats.total_minutes }}</div>
        <div class="text-[10px] text-gray-400 mt-0.5">累计(分)</div>
      </div>
    </div>

    <!-- Row 3: Yesterday + Weekly -->
    <div v-if="review" class="grid grid-cols-2 gap-3 mb-4">
      <!-- Yesterday Summary -->
      <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-4">
        <h3 class="text-xs font-semibold text-gray-500 dark:text-gray-400 mb-3">昨日小结</h3>
        <div class="space-y-2">
          <div class="flex justify-between items-center">
            <span class="text-xs text-gray-400">新增任务</span>
            <span class="text-sm font-semibold text-gray-700 dark:text-gray-300">{{ review.yesterday_summary.tasks_created }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-xs text-gray-400">完成任务</span>
            <span class="text-sm font-semibold text-emerald-500">{{ review.yesterday_summary.tasks_completed }}</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-xs text-gray-400">专注时长</span>
            <span class="text-sm font-semibold text-indigo-500">{{ review.yesterday_summary.focus_minutes }}分钟</span>
          </div>
        </div>
      </div>

      <!-- Weekly Report -->
      <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-4">
        <h3 class="text-xs font-semibold text-gray-500 dark:text-gray-400 mb-3">本周报告</h3>
        <div class="space-y-2">
          <div class="flex justify-between items-center">
            <span class="text-xs text-gray-400">新增 / 完成</span>
            <span class="text-sm font-semibold text-gray-700 dark:text-gray-300">{{ review.weekly_report.tasks_created }} / <span class="text-emerald-500">{{ review.weekly_report.tasks_completed }}</span></span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-xs text-gray-400">专注时长</span>
            <span class="text-sm font-semibold text-indigo-500">{{ review.weekly_report.focus_minutes }}分钟</span>
          </div>
          <div class="flex justify-between items-center">
            <span class="text-xs text-gray-400">最佳工作日</span>
            <span class="text-sm font-semibold text-amber-500">{{ review.weekly_report.best_day }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- Daily trends (created/completed) -->
    <div v-if="stats.daily_trends?.length" class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-5 mb-4">
      <h3 class="text-xs font-semibold text-gray-500 dark:text-gray-400 mb-3">近7天任务趋势</h3>
      <div class="flex items-end justify-between gap-1" style="height: 80px">
        <div v-for="t in stats.daily_trends" :key="t.date" class="flex-1 flex flex-col items-center gap-1">
          <div class="w-full flex flex-col items-center gap-0.5" style="height: 55px; justify-content: flex-end">
            <div v-if="t.completed > 0" class="w-full max-w-[20px] rounded-t bg-emerald-400 dark:bg-emerald-600 transition-all" :style="{ height: Math.min(t.completed * 15, 40) + 'px' }" />
            <div v-if="t.created > 0" class="w-full max-w-[20px] rounded-t bg-indigo-300 dark:bg-indigo-700 transition-all" :style="{ height: Math.min((t.created - t.completed) * 10, 30) + 'px' }" />
          </div>
          <span class="text-[9px] text-gray-400 mt-1">{{ t.date.slice(5) }}</span>
        </div>
      </div>
      <div class="flex items-center gap-4 mt-3 text-[10px] text-gray-400">
        <span class="flex items-center gap-1"><span class="w-2.5 h-2.5 rounded-sm bg-indigo-300" /> 新增</span>
        <span class="flex items-center gap-1"><span class="w-2.5 h-2.5 rounded-sm bg-emerald-400" /> 完成</span>
      </div>
    </div>

    <!-- Completion rate -->
    <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-5 mb-4">
      <h3 class="text-xs font-semibold text-gray-500 dark:text-gray-400 mb-3">完成率</h3>
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
    <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-5 mb-4">
      <h3 class="text-xs font-semibold text-gray-500 dark:text-gray-400 mb-3">优先级分布</h3>
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
      <h3 class="text-xs font-semibold text-gray-500 dark:text-gray-400 mb-3">标签分布</h3>
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

    <!-- Focus trend chart -->
    <div v-if="review?.daily_focus?.length" class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-5 mt-4">
      <h3 class="text-xs font-semibold text-gray-500 dark:text-gray-400 mb-3">近7天专注趋势（分钟）</h3>
      <div class="flex items-end justify-between gap-1" style="height: 80px">
        <div v-for="d in review.daily_focus" :key="d.date" class="flex-1 flex flex-col items-center gap-1">
          <div class="w-full flex flex-col justify-end" style="height: 60px">
            <div
              class="w-full rounded-t bg-indigo-400 dark:bg-indigo-600 transition-all mx-auto"
              :style="{ height: d.minutes > 0 ? focusBarHeight(d.minutes, maxFocusMins()) : '4px', maxWidth: '28px' }"
            />
          </div>
          <span class="text-[9px] text-gray-400 mt-1">{{ d.date.slice(5) }}</span>
          <span v-if="d.minutes > 0" class="text-[9px] font-semibold text-indigo-500">{{ d.minutes }}</span>
        </div>
      </div>
    </div>
  </div>

  <div v-else class="text-center py-20 text-gray-400">加载中…</div>
</template>
