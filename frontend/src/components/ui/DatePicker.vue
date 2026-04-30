<script setup lang="ts">
import { computed, ref } from 'vue'

const model = defineModel<string>()

const open = ref(false)

const quickOptions = [
  { label: '1天后', days: 1 },
  { label: '3天后', days: 3 },
  { label: '1周后', days: 7 },
  { label: '2周后', days: 14 },
  { label: '1月后', days: 30 },
  { label: '3月后', days: 90 },
]

function fmt(d: Date) {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

function selectDays(days: number) {
  const d = new Date()
  d.setDate(d.getDate() + days)
  model.value = fmt(d)
  open.value = false
}

function clear() {
  model.value = ''
  open.value = false
}

const display = computed(() => {
  if (!model.value) return ''
  const today = fmt(new Date())
  if (model.value === today) return '今天'
  const tomorrow = fmt(new Date(Date.now() + 86400000))
  if (model.value === tomorrow) return '明天'
  return model.value
})

// Mini calendar
const calYear = ref(0)
const calMonth = ref(0)

const calDays = computed(() => {
  const first = new Date(calYear.value, calMonth.value, 1)
  const last = new Date(calYear.value, calMonth.value + 1, 0)
  const startDay = first.getDay()
  const days: (number | null)[] = []
  for (let i = 0; i < startDay; i++) days.push(null)
  for (let d = 1; d <= last.getDate(); d++) days.push(d)
  return days
})

const calLabel = computed(() => `${calYear.value}年${calMonth.value + 1}月`)

function prevCal() {
  if (calMonth.value === 0) { calMonth.value = 11; calYear.value-- }
  else calMonth.value--
}

function nextCal() {
  if (calMonth.value === 11) { calMonth.value = 0; calYear.value++ }
  else calMonth.value++
}

function selectDay(d: number) {
  const m = String(calMonth.value + 1).padStart(2, '0')
  const day = String(d).padStart(2, '0')
  model.value = `${calYear.value}-${m}-${day}`
  open.value = false
}

function toggleOpen() {
  open.value = !open.value
  if (open.value) {
    const d = model.value ? new Date(model.value) : new Date()
    calYear.value = d.getFullYear()
    calMonth.value = d.getMonth()
  }
}
</script>

<template>
  <div class="relative">
    <button type="button" @click="toggleOpen"
      class="px-3 py-2.5 border border-gray-200 dark:border-gray-800 rounded-lg text-sm bg-white dark:bg-gray-900 dark:text-gray-100 outline-none hover:border-indigo-400 transition-colors min-w-[110px] text-left whitespace-nowrap">
      <span v-if="display" class="text-gray-800 dark:text-gray-200">{{ display }}</span>
      <span v-else class="text-gray-400 dark:text-gray-500">截止日期</span>
    </button>

    <div v-if="open" class="absolute left-0 top-full mt-1 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-xl shadow-lg z-20 p-3 min-w-[260px]">
      <!-- Quick options -->
      <div class="grid grid-cols-3 gap-1.5 mb-2">
        <button v-for="opt in quickOptions" :key="opt.days" type="button"
          @click="selectDays(opt.days)"
          class="px-2 py-1.5 text-xs rounded-lg bg-gray-50 dark:bg-gray-700/50 text-gray-600 dark:text-gray-300 hover:bg-indigo-50 dark:hover:bg-indigo-900/30 hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors">
          {{ opt.label }}
        </button>
      </div>

      <!-- Mini calendar -->
      <div>
        <div class="flex items-center justify-between mb-1">
          <button type="button" @click="prevCal" class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 px-1 text-sm">&larr;</button>
          <span class="text-xs font-semibold text-gray-600 dark:text-gray-300">{{ calLabel }}</span>
          <button type="button" @click="nextCal" class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 px-1 text-sm">&rarr;</button>
        </div>
        <div class="grid grid-cols-7 mb-1">
          <div v-for="d in ['日','一','二','三','四','五','六']" :key="d" class="text-center text-[10px] text-gray-400 dark:text-gray-600 py-0.5">{{ d }}</div>
        </div>
        <div class="grid grid-cols-7">
          <div v-for="(d, i) in calDays" :key="i"
            @click="d && selectDay(d)"
            class="text-center py-1 text-xs rounded cursor-pointer transition-colors"
            :class="d ? 'hover:bg-indigo-50 dark:hover:bg-indigo-900/30 text-gray-700 dark:text-gray-300' : ''">
            {{ d || '' }}
          </div>
        </div>
      </div>

      <!-- Clear -->
      <div class="border-t border-gray-100 dark:border-gray-700 mt-2 pt-2">
        <button v-if="model" type="button" @click="clear"
          class="w-full px-2 py-1.5 text-xs text-red-400 hover:text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors">
          清除日期
        </button>
        <button type="button" @click="open = false"
          class="w-full px-2 py-1.5 text-xs text-gray-400 hover:text-gray-500 mt-0.5 rounded-lg transition-colors">
          关闭
        </button>
      </div>
    </div>

    <!-- Backdrop -->
    <div v-if="open" class="fixed inset-0 z-10" @click="open = false" />
  </div>
</template>
