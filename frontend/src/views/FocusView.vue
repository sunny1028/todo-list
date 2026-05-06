<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { useTodos } from '../stores/todo'
import * as focusApi from '../api/focus'
import type { FocusSession } from '../types/todo'

const store = useTodos()

// Timer state
type State = 'idle' | 'focus' | 'break'
const state = ref<State>('idle')
const secondsLeft = ref(0)
const totalSeconds = ref(0)
const session = ref<FocusSession | null>(null)
let intervalId: ReturnType<typeof setInterval> | null = null

// Settings
const focusMin = ref(25)
const breakMin = ref(5)
const linkedTodoId = ref<number | null>(null)
const isFullscreen = ref(false)

// Stats & history
const stats = ref({ today_minutes: 0, total_minutes: 0, total_sessions: 0, streak_days: 0 })
const sessions = ref<FocusSession[]>([])

const progress = computed(() => {
  if (totalSeconds.value === 0) return 0
  return (totalSeconds.value - secondsLeft.value) / totalSeconds.value
})

const displayTime = computed(() => {
  const m = Math.floor(secondsLeft.value / 60)
  const s = secondsLeft.value % 60
  return `${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
})

// Canvas circle
const canvas = ref<HTMLCanvasElement | null>(null)
function draw() {
  const c = canvas.value
  if (!c) return
  const ctx = c.getContext('2d')
  if (!ctx) return
  const w = c.width; const h = c.height
  ctx.clearRect(0, 0, w, h)

  const cx = w / 2; const cy = h / 2; const r = 140

  // Background ring
  ctx.beginPath()
  ctx.arc(cx, cy, r, 0, Math.PI * 2)
  ctx.strokeStyle = 'rgba(148,163,184,0.12)'
  ctx.lineWidth = 10
  ctx.stroke()

  // Progress ring
  if (progress.value > 0) {
    const endAngle = -Math.PI / 2 + progress.value * Math.PI * 2
    ctx.beginPath()
    ctx.arc(cx, cy, r, -Math.PI / 2, endAngle)
    ctx.strokeStyle = state.value === 'focus' ? '#6366f1' : '#10b981'
    ctx.lineWidth = 10
    ctx.lineCap = 'round'
    ctx.stroke()
  }
}

watch(progress, () => draw())
watch(state, () => draw())

function tick() {
  secondsLeft.value--
  if (secondsLeft.value <= 0) {
    clearInterval(intervalId!)
    intervalId = null
    onTimerEnd()
  }
}

function startTimer(type: 'focus' | 'break') {
  state.value = type
  totalSeconds.value = (type === 'focus' ? focusMin.value : breakMin.value) * 60
  secondsLeft.value = totalSeconds.value
  if (type === 'focus') {
    focusApi.startFocus({
      todo_id: linkedTodoId.value,
      duration_min: focusMin.value,
    }).then(res => { session.value = res.data }).catch(() => {})
  }
  intervalId = setInterval(tick, 1000)
  draw()
}

function onTimerEnd() {
  // Sound - web audio beep
  try {
    const ctx = new AudioContext()
    const osc = ctx.createOscillator()
    const gain = ctx.createGain()
    osc.connect(gain); gain.connect(ctx.destination)
    osc.frequency.value = 880; osc.type = 'sine'
    gain.gain.setValueAtTime(0.3, ctx.currentTime)
    gain.gain.exponentialRampToValueAtTime(0.01, ctx.currentTime + 0.3)
    osc.start(); osc.stop(ctx.currentTime + 0.3)
  } catch { /* no audio */ }

  // Notification
  if ('Notification' in window && Notification.permission === 'granted') {
    new Notification(state.value === 'focus' ? '专注完成！' : '休息结束！', {
      body: state.value === 'focus' ? '干得漂亮，休息一下吧。' : '准备好开始新专注了吗？',
    })
  }

  // Complete backend session
  if (session.value) {
    focusApi.completeFocus(session.value.id).catch(() => {})
  }

  loadStats()
  loadSessions()

  if (state.value === 'focus') {
    startTimer('break')
  } else {
    state.value = 'idle'
    totalSeconds.value = 0
    session.value = null
    draw()
  }
}

function pause() { if (intervalId) { clearInterval(intervalId); intervalId = null } }
function resume() { if (secondsLeft.value > 0 && !intervalId) { intervalId = setInterval(tick, 1000) } }
function reset() {
  if (intervalId) { clearInterval(intervalId); intervalId = null }
  state.value = 'idle'; secondsLeft.value = 0; totalSeconds.value = 0; session.value = null
  draw()
}
function toggleFullscreen() {
  if (!document.fullscreenElement) {
    document.documentElement.requestFullscreen()
    isFullscreen.value = true
  } else {
    document.exitFullscreen()
    isFullscreen.value = false
  }
}
function skip() { clearInterval(intervalId!); intervalId = null; onTimerEnd() }

async function loadStats() {
  try { const r = await focusApi.getFocusStats(); stats.value = r.data } catch {}
}
async function loadSessions() {
  try { const r = await focusApi.listSessions(); sessions.value = r.data } catch {}
}

store.fetchTodos()
loadStats()
loadSessions()
if ('Notification' in window && Notification.permission === 'default') {
  Notification.requestPermission()
}

function onResize() { draw() }
window.addEventListener('resize', onResize)
onBeforeUnmount(() => {
  window.removeEventListener('resize', onResize)
  if (intervalId) clearInterval(intervalId)
})
</script>

<template>
  <div class="max-w-md mx-auto py-4">
    <!-- Stats strip -->
    <div class="flex justify-center gap-8 text-center mb-6">
      <div>
        <div class="text-xl font-bold text-indigo-600 dark:text-indigo-400">{{ stats.today_minutes }}</div>
        <div class="text-[11px] text-gray-400">今日分钟</div>
      </div>
      <div>
        <div class="text-xl font-bold text-gray-700 dark:text-gray-200">{{ stats.total_sessions }}</div>
        <div class="text-[11px] text-gray-400">总次数</div>
      </div>
      <div>
        <div class="text-xl font-bold text-emerald-600 dark:text-emerald-400">{{ stats.streak_days }}</div>
        <div class="text-[11px] text-gray-400">连续天数</div>
      </div>
    </div>

    <!-- Circle Timer -->
    <div class="flex justify-center mb-6 relative">
      <canvas ref="canvas" width="320" height="320" class="block" />
      <div class="absolute inset-0 flex flex-col items-center justify-center pointer-events-none">
        <span class="text-5xl font-mono font-bold text-gray-800 dark:text-gray-100 tracking-wider tabular-nums" style="font-variant-numeric: tabular-nums">{{ displayTime }}</span>
        <span v-if="state === 'focus'" class="text-sm text-indigo-500 dark:text-indigo-400 mt-1 font-medium">专注中</span>
        <span v-else-if="state === 'break'" class="text-sm text-emerald-500 dark:text-emerald-400 mt-1 font-medium">休息</span>
        <span v-else class="text-sm text-gray-400 mt-1">准备开始</span>
      </div>
    </div>

    <!-- Controls -->
    <div class="flex items-center justify-center gap-4 mb-8">
      <template v-if="state === 'idle'">
        <button @click="startTimer('focus')" class="px-10 py-3 bg-indigo-500 text-white rounded-full text-lg font-semibold hover:bg-indigo-600 transition-colors shadow-sm">
          开始专注
        </button>
      </template>
      <template v-else>
        <button @click="reset" class="w-10 h-10 rounded-full bg-gray-200 dark:bg-gray-700 flex items-center justify-center text-gray-500 dark:text-gray-400 hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors" title="结束">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="2" width="20" height="20" rx="2"/></svg>
        </button>
        <template v-if="intervalId">
          <button @click="pause" class="w-14 h-14 rounded-full bg-indigo-500 text-white flex items-center justify-center hover:bg-indigo-600 transition-colors shadow-sm" title="暂停">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="currentColor"><rect x="6" y="4" width="4" height="16" rx="1"/><rect x="14" y="4" width="4" height="16" rx="1"/></svg>
          </button>
        </template>
        <template v-else>
          <button @click="resume" class="w-14 h-14 rounded-full bg-indigo-500 text-white flex items-center justify-center hover:bg-indigo-600 transition-colors shadow-sm" title="继续">
            <svg width="22" height="22" viewBox="0 0 24 24" fill="currentColor"><polygon points="6,3 20,12 6,21"/></svg>
          </button>
        </template>
        <button @click="skip" class="w-10 h-10 rounded-full bg-gray-200 dark:bg-gray-700 flex items-center justify-center text-gray-500 dark:text-gray-400 hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors" title="跳过">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="13 19 22 12 13 5 13 19"/><polygon points="2 19 11 12 2 5 2 19"/></svg>
        </button>
      </template>
    </div>

    <!-- Settings -->
    <div class="px-4">
      <div v-if="state === 'idle'" class="flex items-center gap-3 flex-wrap justify-center">
        <div class="flex items-center gap-1.5">
          <label class="text-xs text-gray-400">专注</label>
          <select v-model.number="focusMin" class="px-2 py-1.5 border border-gray-200 dark:border-gray-700 rounded-lg text-sm bg-white dark:bg-gray-800 dark:text-gray-100 outline-none cursor-pointer">
            <option :value="15">15分钟</option>
            <option :value="25">25分钟</option>
            <option :value="30">30分钟</option>
            <option :value="45">45分钟</option>
            <option :value="60">60分钟</option>
          </select>
        </div>
        <div class="flex items-center gap-1.5">
          <label class="text-xs text-gray-400">休息</label>
          <select v-model.number="breakMin" class="px-2 py-1.5 border border-gray-200 dark:border-gray-700 rounded-lg text-sm bg-white dark:bg-gray-800 dark:text-gray-100 outline-none cursor-pointer">
            <option :value="5">5分钟</option>
            <option :value="10">10分钟</option>
            <option :value="15">15分钟</option>
          </select>
        </div>
        <div class="flex items-center gap-1.5">
          <label class="text-xs text-gray-400">关联</label>
          <select v-model="linkedTodoId" class="px-2 py-1.5 border border-gray-200 dark:border-gray-700 rounded-lg text-sm bg-white dark:bg-gray-800 dark:text-gray-100 outline-none cursor-pointer max-w-[150px]">
            <option :value="null">不关联</option>
            <option v-for="t in store.todos.filter(t => !t.completed)" :key="t.id" :value="t.id">{{ t.title }}</option>
          </select>
        </div>
        <button @click="toggleFullscreen" class="px-3 py-1.5 text-xs border border-gray-200 dark:border-gray-700 rounded-lg text-gray-500 dark:text-gray-400 hover:text-indigo-500 dark:hover:text-indigo-400 transition-colors">
          沉浸
        </button>
      </div>
    </div>

    <!-- Session history -->
    <div class="mt-10 px-4">
      <h3 class="text-sm font-semibold text-gray-500 dark:text-gray-400 mb-3">最近专注</h3>
      <div v-if="sessions.length === 0" class="text-xs text-gray-300 dark:text-gray-700">暂无记录</div>
      <div v-for="s in sessions.slice(0, 10)" :key="s.id" class="flex items-center justify-between py-2.5 border-b border-gray-50 dark:border-gray-800/50 last:border-0">
        <div class="flex-1 min-w-0">
          <div class="text-sm text-gray-700 dark:text-gray-300 truncate">{{ s.todo_title || '未关联任务' }}</div>
          <div class="text-[11px] text-gray-400">{{ new Date(s.started_at).toLocaleString('zh-CN', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' }) }}</div>
        </div>
        <span class="text-sm font-semibold ml-3" :class="s.completed ? 'text-emerald-500' : 'text-gray-400'">{{ s.duration_min }}分钟</span>
      </div>
    </div>
  </div>
</template>
