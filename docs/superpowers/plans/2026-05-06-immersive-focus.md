# 沉浸专注 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 构建可配置专注/休息时长的沉浸式计时器，关联待办事项，记录专注历史与统计。

**Architecture:** 后端新增 FocusSession 模型和 REST API（记录/统计），前端新增 FocusView 页面（圆形计时器 + 沉浸模式 + 通知），侧边栏添加入口，统计页集成专注数据。

**Tech Stack:** Go + GORM + Vue 3 + Canvas (圆环计时器) + Web Notification API + Fullscreen API + Web Audio API (完成提示音)

---

## 文件结构

### 后端新建
| 文件 | 职责 |
|------|------|
| `backend/models/focus.go` | FocusSession 模型定义 |
| `backend/handlers/focus.go` | 开始/结束/统计/历史 API |
| `backend/services/focus.go` | 业务逻辑（校验、统计聚合） |
| `backend/repository/focus.go` | 数据库查询 |

### 后端修改
| 文件 | 改动 |
|------|------|
| `backend/router/router.go` | 新增 /api/focus 路由组 |
| `backend/models/todo.go` | 新增 FocusMinutes 累计字段 |

### 前端新建
| 文件 | 职责 |
|------|------|
| `frontend/src/views/FocusView.vue` | 计时器主页面（圆环、控件、设置） |
| `frontend/src/api/focus.ts` | 专注 API 调用 |

### 前端修改
| 文件 | 改动 |
|------|------|
| `frontend/src/router/index.ts` | 新增 /focus 路由 |
| `frontend/src/components/layout/Sidebar.vue` | 新增「沉浸专注」导航按钮 |
| `frontend/src/views/StatsView.vue` | 新增专注统计卡片（今日/累计） |
| `frontend/src/types/todo.ts` | 新增 FocusSession 类型 |

---

### Task 1: 后端 — FocusSession 模型与数据表

**Files:**
- Create: `backend/models/focus.go`
- Modify: `backend/models/todo.go`

- [ ] **Step 1: 创建 FocusSession 模型**

```go
// backend/models/focus.go
package models

import "time"

type FocusSession struct {
	ID             uint      `json:"id" gorm:"primaryKey"`
	UserID         uint      `json:"user_id" gorm:"index;not null"`
	TodoID         *uint     `json:"todo_id" gorm:"index;default:null"`
	TodoTitle      string    `json:"todo_title" gorm:"-"`
	DurationMin    int       `json:"duration_min" gorm:"not null"`
	Completed      bool      `json:"completed" gorm:"default:false"`
	StartedAt      time.Time `json:"started_at"`
	EndedAt        *time.Time `json:"ended_at"`
	CreatedAt      time.Time `json:"created_at"`
}
```

- [ ] **Step 2: 在 Todo 模型添加累计专注字段**

```go
// backend/models/todo.go — 在 SubtaskCompleted 行后添加
FocusMinutes int `json:"focus_minutes" gorm:"default:0"`
```

- [ ] **Step 3: 编译验证**

```bash
cd backend && go build ./...
```
Expected: PASS

- [ ] **Step 4: 提交**

```bash
git add backend/models/focus.go backend/models/todo.go
git commit -m "feat: add FocusSession model and Todo.FocusMinutes"
```

---

### Task 2: 后端 — Repository + Service + Handler

**Files:**
- Create: `backend/repository/focus.go`
- Create: `backend/services/focus.go`
- Create: `backend/handlers/focus.go`

- [ ] **Step 1: Repository 层**

```go
// backend/repository/focus.go
package repository

import (
	"time"
	"todo-list/backend/database"
	"todo-list/backend/models"
)

func CreateFocusSession(s *models.FocusSession) error {
	return database.DB.Create(s).Error
}

func UpdateFocusSession(s *models.FocusSession) error {
	return database.DB.Save(s).Error
}

func FindFocusSessions(userID uint, limit int) ([]models.FocusSession, error) {
	var sessions []models.FocusSession
	err := database.DB.Where("user_id = ?", userID).
		Order("created_at DESC").Limit(limit).
		Find(&sessions).Error
	if err != nil {
		return nil, err
	}
	// Populate todo titles
	for i := range sessions {
		if sessions[i].TodoID != nil {
			var todo models.Todo
			if database.DB.Select("title").First(&todo, *sessions[i].TodoID).Error == nil {
				sessions[i].TodoTitle = todo.Title
			}
		}
	}
	return sessions, nil
}

type FocusStats struct {
	TodayMinutes   int `json:"today_minutes"`
	TotalMinutes   int `json:"total_minutes"`
	TotalSessions  int64 `json:"total_sessions"`
	StreakDays     int `json:"streak_days"`
}

func GetFocusStats(userID uint) FocusStats {
	var s FocusStats
	today := time.Now().Format("2006-01-02")

	// Today
	database.DB.Model(&models.FocusSession{}).
		Where("user_id = ? AND completed = ? AND date(started_at) = ?", userID, true, today).
		Select("COALESCE(SUM(duration_min), 0)").Scan(&s.TodayMinutes)

	// Total
	database.DB.Model(&models.FocusSession{}).
		Where("user_id = ? AND completed = ?", userID, true).
		Select("COALESCE(SUM(duration_min), 0)").Scan(&s.TotalMinutes)

	// Sessions count
	database.DB.Model(&models.FocusSession{}).
		Where("user_id = ? AND completed = ?", userID, true).Count(&s.TotalSessions)

	// Streak: count consecutive days back from today with at least one completed session
	d := time.Now()
	for {
		dateStr := d.Format("2006-01-02")
		var count int64
		database.DB.Model(&models.FocusSession{}).
			Where("user_id = ? AND completed = ? AND date(started_at) = ?", userID, true, dateStr).
			Count(&count)
		if count == 0 {
			break
		}
		s.StreakDays++
		d = d.AddDate(0, 0, -1)
	}

	return s
}
```

- [ ] **Step 2: Service 层**

```go
// backend/services/focus.go
package services

import (
	"errors"
	"time"
	"todo-list/backend/models"
	"todo-list/backend/repository"
)

func StartFocus(userID uint, todoID *uint, durationMin int) (*models.FocusSession, error) {
	if durationMin <= 0 || durationMin > 180 {
		return nil, errors.New("duration must be 1-180 minutes")
	}
	s := &models.FocusSession{
		UserID:      userID,
		TodoID:      todoID,
		DurationMin: durationMin,
		StartedAt:   time.Now(),
	}
	if err := repository.CreateFocusSession(s); err != nil {
		return nil, err
	}
	return s, nil
}

func CompleteFocus(userID uint, id uint) (*models.FocusSession, error) {
	var s models.FocusSession
	if err := database.DB.Where("user_id = ?", userID).First(&s, id).Error; err != nil {
		return nil, errors.New("session not found")
	}
	if s.Completed {
		return &s, nil
	}
	now := time.Now()
	s.Completed = true
	s.EndedAt = &now
	if err := repository.UpdateFocusSession(&s); err != nil {
		return nil, err
	}
	// Increment todo focus minutes if linked
	if s.TodoID != nil {
		database.DB.Model(&models.Todo{}).Where("id = ?", *s.TodoID).
			Update("focus_minutes", gorm.Expr("focus_minutes + ?", s.DurationMin))
	}
	return &s, nil
}
```

需要在文件顶部添加 import：
```go
import (
	"todo-list/backend/database"
	"gorm.io/gorm"
)
```

- [ ] **Step 3: Handler 层**

```go
// backend/handlers/focus.go
package handlers

import (
	"net/http"
	"strconv"
	"todo-list/backend/services"
	"github.com/gin-gonic/gin"
)

type startFocusRequest struct {
	TodoID      *uint `json:"todo_id"`
	DurationMin int   `json:"duration_min" binding:"required"`
}

func StartFocus(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req startFocusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "duration_min is required"})
		return
	}
	s, err := services.StartFocus(userID, req.TodoID, req.DurationMin)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, s)
}

func CompleteFocus(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	s, err := services.CompleteFocus(userID, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, s)
}

func ListFocusSessions(c *gin.Context) {
	userID := c.GetUint("user_id")
	sessions, err := services.GetFocusSessions(userID, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sessions)
}

func GetFocusStats(c *gin.Context) {
	userID := c.GetUint("user_id")
	stats := services.GetFocusStats(userID)
	c.JSON(http.StatusOK, stats)
}
```

需要在 service 层补充 `GetFocusSessions` 和 `GetFocusStats` 代理函数：
```go
func GetFocusSessions(userID uint, limit int) ([]models.FocusSession, error) {
	return repository.FindFocusSessions(userID, limit)
}

func GetFocusStats(userID uint) repository.FocusStats {
	return repository.GetFocusStats(userID)
}
```

- [ ] **Step 4: 编译验证**

```bash
cd backend && go build ./...
```
Expected: PASS

- [ ] **Step 5: 提交**

```bash
git add backend/repository/focus.go backend/services/focus.go backend/handlers/focus.go
git commit -m "feat: add focus session CRUD, stats, and handlers"
```

---

### Task 3: 后端 — 路由注册

**Files:**
- Modify: `backend/router/router.go`

- [ ] **Step 1: 在 protected 路由组添加 focus 路由**

在 `backend/router/router.go` 的 Lists 路由之前插入：

```go
// Focus
protected.POST("/focus/start", handlers.StartFocus)
protected.PATCH("/focus/:id/complete", handlers.CompleteFocus)
protected.GET("/focus/sessions", handlers.ListFocusSessions)
protected.GET("/focus/stats", handlers.GetFocusStats)
```

实际插入位置：在 Subtasks 路由组闭合后、Attachments 路由之前。用以下命令：

```bash
cd backend && sed -i '/protected.GET("\/todos\/:id\/attachments"/i\
		// Focus\
		protected.POST("/focus/start", handlers.StartFocus)\
		protected.PATCH("/focus/:id/complete", handlers.CompleteFocus)\
		protected.GET("/focus/sessions", handlers.ListFocusSessions)\
		protected.GET("/focus/stats", handlers.GetFocusStats)\
' router/router.go
```

- [ ] **Step 2: 编译验证**

```bash
cd backend && go build ./...
```
Expected: PASS

- [ ] **Step 3: 提交**

```bash
git add backend/router/router.go
git commit -m "feat: add focus API routes"
```

---

### Task 4: 前端 — 类型定义与 API 层

**Files:**
- Create: `frontend/src/api/focus.ts`
- Modify: `frontend/src/types/todo.ts`

- [ ] **Step 1: 添加 FocusSession 类型**

在 `frontend/src/types/todo.ts` 末尾添加：

```typescript
export interface FocusSession {
  id: number
  user_id: number
  todo_id: number | null
  todo_title: string
  duration_min: number
  completed: boolean
  started_at: string
  ended_at: string | null
  created_at: string
}

export interface FocusStats {
  today_minutes: number
  total_minutes: number
  total_sessions: number
  streak_days: number
}
```

- [ ] **Step 2: 创建 API 层**

```typescript
// frontend/src/api/focus.ts
import api from './index'
import type { FocusSession, FocusStats } from '../types/todo'

export function startFocus(data: { todo_id?: number | null; duration_min: number }) {
  return api.post<FocusSession>('/focus/start', data)
}
export function completeFocus(id: number) {
  return api.patch<FocusSession>(`/focus/${id}/complete`)
}
export function listSessions() {
  return api.get<FocusSession[]>('/focus/sessions')
}
export function getFocusStats() {
  return api.get<FocusStats>('/focus/stats')
}
```

- [ ] **Step 3: 编译验证**

```bash
cd frontend && npx vue-tsc --noEmit 2>&1 | head -5
```
Expected: no errors from focus files

- [ ] **Step 4: 提交**

```bash
git add frontend/src/types/todo.ts frontend/src/api/focus.ts
git commit -m "feat: add focus types and API layer"
```

---

### Task 5: 前端 — FocusView 计时器页面（核心）

**Files:**
- Create: `frontend/src/views/FocusView.vue`

- [ ] **Step 1: 创建 FocusView.vue**

```vue
<script setup lang="ts">
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import { useTodos } from '../stores/todo'
import * as focusApi from '../api/focus'
import { useToast } from '../stores/toast'
import type { FocusSession } from '../types/todo'

const store = useTodos()
const toast = useToast()

// Timer state
const STATE_IDLE = 'idle'
const STATE_FOCUS = 'focus'
const STATE_BREAK = 'break'

const state = ref(STATE_IDLE)
const secondsLeft = ref(0)
const totalSeconds = ref(0)
const session = ref<FocusSession | null>(null)
let intervalId: ReturnType<typeof setInterval> | null = null

// Settings
const focusMin = ref(25)
const breakMin = ref(5)
const linkedTodoId = ref<number | null>(null)
const showSettings = ref(false)
const isFullscreen = ref(false)

// Stats
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

const circleDash = computed(() => {
  const r = 140
  const c = 2 * Math.PI * r
  return `${(progress.value * c).toFixed(0)} ${c}`
})

// Canvas circle
const canvas = ref<HTMLCanvasElement | null>(null)
function draw() {
  const c = canvas.value
  if (!c) return
  const ctx = c.getContext('2d')
  if (!ctx) return
  const w = c.width
  const h = c.height
  ctx.clearRect(0, 0, w, h)

  const cx = w / 2; const cy = h / 2; const r = 140
  // Background ring
  ctx.beginPath()
  ctx.arc(cx, cy, r, 0, Math.PI * 2)
  ctx.strokeStyle = 'rgba(148,163,184,0.15)'
  ctx.lineWidth = 12
  ctx.stroke()

  // Progress ring
  const endAngle = -Math.PI / 2 + progress.value * Math.PI * 2
  ctx.beginPath()
  ctx.arc(cx, cy, r, -Math.PI / 2, endAngle)
  const color = state.value === STATE_FOCUS ? '#6366f1' : '#10b981'
  ctx.strokeStyle = color
  ctx.lineWidth = 12
  ctx.lineCap = 'round'
  ctx.stroke()
}

watch(progress, () => draw())
watch(state, () => draw())

function startTimer(type: string) {
  state.value = type
  totalSeconds.value = (type === STATE_FOCUS ? focusMin.value : breakMin.value) * 60
  secondsLeft.value = totalSeconds.value
  // Start backend session for focus
  if (type === STATE_FOCUS) {
    focusApi.startFocus({
      todo_id: linkedTodoId.value,
      duration_min: focusMin.value,
    }).then(res => { session.value = res.data }).catch(() => {})
  }
  intervalId = setInterval(tick, 1000)
}

function tick() {
  secondsLeft.value--
  if (secondsLeft.value <= 0) {
    clearInterval(intervalId!)
    intervalId = null
    onTimerEnd()
  }
}

function onTimerEnd() {
  const audio = new Audio('data:audio/wav;base64,UklGRnoGAABXQVZFZm10IBAAAAABAAEAQB8AAEAfAAABAAgAZGF0YQoGAACAf39/f4B/f3+Af39/gH9/f4B/f3+Af39/gH9/f4B/f3+Af39/gH9/f4B/f3+Af39/gH9/f4B/f3+Af39/gH9/f4B/f3+Af39/gH9/f4B/f3+Af39/gH9/f4B/f3+Af39/gH9/f4B/f3+Af39/gH9/f4B/f3+Af39/gA==')
  audio.play().catch(() => {})
  // Notification
  if ('Notification' in window && Notification.permission === 'granted') {
    new Notification(state.value === STATE_FOCUS ? '专注完成！' : '休息结束！', {
      body: state.value === STATE_FOCUS ? '干得漂亮，休息一下吧' : '准备好开始新的专注了吗？',
    })
  }
  // Complete backend session
  if (session.value) {
    focusApi.completeFocus(session.value.id).catch(() => {})
  }
  loadStats()
  // Auto-switch
  if (state.value === STATE_FOCUS) {
    startTimer(STATE_BREAK)
  } else {
    state.value = STATE_IDLE
    totalSeconds.value = 0
    session.value = null
  }
}

function pause() {
  if (intervalId) { clearInterval(intervalId); intervalId = null }
}
function resume() {
  if (secondsLeft.value > 0 && !intervalId) {
    intervalId = setInterval(tick, 1000)
  }
}
function reset() {
  if (intervalId) { clearInterval(intervalId); intervalId = null }
  state.value = STATE_IDLE
  secondsLeft.value = 0
  totalSeconds.value = 0
  session.value = null
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

async function loadStats() {
  try {
    const res = await focusApi.getFocusStats()
    stats.value = res.data
  } catch {}
}
async function loadSessions() {
  try {
    const res = await focusApi.listSessions()
    sessions.value = res.data
  } catch {}
}

// Request notification permission
function requestNotification() {
  if ('Notification' in window && Notification.permission === 'default') {
    Notification.requestPermission()
  }
}

// Init
store.fetchTodos()
loadStats()
loadSessions()
requestNotification()

// Keep canvas in sync on resize
function onResize() { draw() }
window.addEventListener('resize', onResize)
onBeforeUnmount(() => {
  window.removeEventListener('resize', onResize)
  if (intervalId) clearInterval(intervalId)
})
</script>

<template>
  <div class="max-w-md mx-auto py-4" :class="{ 'max-w-full': isFullscreen }">
    <!-- Stats strip -->
    <div class="flex justify-between text-center mb-6 px-4">
      <div>
        <div class="text-lg font-bold text-indigo-600">{{ stats.today_minutes }}</div>
        <div class="text-[11px] text-gray-400">今日分钟</div>
      </div>
      <div>
        <div class="text-lg font-bold text-gray-700 dark:text-gray-200">{{ stats.total_sessions }}</div>
        <div class="text-[11px] text-gray-400">总次数</div>
      </div>
      <div>
        <div class="text-lg font-bold text-emerald-600">{{ stats.streak_days }}</div>
        <div class="text-[11px] text-gray-400">连续天数</div>
      </div>
      <div>
        <div class="text-lg font-bold text-gray-700 dark:text-gray-200">{{ stats.total_minutes }}</div>
        <div class="text-[11px] text-gray-400">累计分钟</div>
      </div>
    </div>

    <!-- Circle Timer -->
    <div class="flex justify-center mb-6 relative">
      <canvas ref="canvas" width="320" height="320" class="block" />
      <div class="absolute inset-0 flex flex-col items-center justify-center pointer-events-none">
        <span class="text-5xl font-mono font-bold text-gray-800 dark:text-gray-100 tracking-wider">{{ displayTime }}</span>
        <span v-if="state === STATE_FOCUS" class="text-sm text-indigo-500 mt-1">专注中</span>
        <span v-else-if="state === STATE_BREAK" class="text-sm text-emerald-500 mt-1">休息中</span>
        <span v-else class="text-sm text-gray-400 mt-1">准备开始</span>
      </div>
    </div>

    <!-- Controls -->
    <div class="flex items-center justify-center gap-4 mb-6">
      <template v-if="state === STATE_IDLE">
        <button @click="startTimer(STATE_FOCUS)" class="px-8 py-3 bg-indigo-500 text-white rounded-full text-lg font-semibold hover:bg-indigo-600 transition-colors">
          开始专注
        </button>
      </template>
      <template v-else>
        <button @click="reset" class="w-10 h-10 rounded-full bg-gray-200 dark:bg-gray-700 flex items-center justify-center text-gray-500 hover:bg-gray-300 transition-colors">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="2" width="20" height="20" rx="2"/></svg>
        </button>
        <template v-if="intervalId">
          <button @click="pause" class="w-14 h-14 rounded-full bg-indigo-500 text-white flex items-center justify-center hover:bg-indigo-600 transition-colors">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><rect x="6" y="4" width="4" height="16"/><rect x="14" y="4" width="4" height="16"/></svg>
          </button>
        </template>
        <template v-else>
          <button @click="resume" class="w-14 h-14 rounded-full bg-indigo-500 text-white flex items-center justify-center hover:bg-indigo-600 transition-colors">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor"><polygon points="5,3 19,12 5,21"/></svg>
          </button>
        </template>
        <button @click="onTimerEnd" class="w-10 h-10 rounded-full bg-gray-200 dark:bg-gray-700 flex items-center justify-center text-gray-500 hover:bg-gray-300 transition-colors">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><polygon points="11,5 6,9 2,9 2,15 6,15 11,19 11,5"/><line x1="23" y1="9" x2="17" y2="15"/><line x1="17" y1="9" x2="23" y2="15"/></svg>
        </button>
      </template>
    </div>

    <!-- Settings row -->
    <div class="px-4">
      <div v-if="state === STATE_IDLE" class="flex items-center gap-3 flex-wrap">
        <div class="flex items-center gap-1.5">
          <label class="text-xs text-gray-400">专注</label>
          <select v-model.number="focusMin" class="px-2 py-1.5 border border-gray-200 dark:border-gray-700 rounded-lg text-sm bg-white dark:bg-gray-800 dark:text-gray-100">
            <option :value="15">15分钟</option>
            <option :value="25">25分钟</option>
            <option :value="30">30分钟</option>
            <option :value="45">45分钟</option>
            <option :value="60">60分钟</option>
          </select>
        </div>
        <div class="flex items-center gap-1.5">
          <label class="text-xs text-gray-400">休息</label>
          <select v-model.number="breakMin" class="px-2 py-1.5 border border-gray-200 dark:border-gray-700 rounded-lg text-sm bg-white dark:bg-gray-800 dark:text-gray-100">
            <option :value="5">5分钟</option>
            <option :value="10">10分钟</option>
            <option :value="15">15分钟</option>
            <option :value="30">30分钟</option>
          </select>
        </div>
        <div class="flex items-center gap-1.5">
          <label class="text-xs text-gray-400">关联</label>
          <select v-model="linkedTodoId" class="px-2 py-1.5 border border-gray-200 dark:border-gray-700 rounded-lg text-sm bg-white dark:bg-gray-800 dark:text-gray-100 max-w-[160px]">
            <option :value="null">不关联</option>
            <option v-for="t in store.activeTodos" :key="t.id" :value="t.id">{{ t.title }}</option>
          </select>
        </div>
        <button @click="toggleFullscreen" class="ml-auto px-3 py-1.5 text-xs border border-gray-200 dark:border-gray-700 rounded-lg text-gray-500 hover:text-indigo-500 transition-colors">
          沉浸
        </button>
      </div>
    </div>

    <!-- Sessions list -->
    <div class="mt-8 px-4">
      <h3 class="text-sm font-semibold text-gray-500 dark:text-gray-400 mb-2">最近专注</h3>
      <div v-if="sessions.length === 0" class="text-xs text-gray-300 dark:text-gray-700">暂无记录</div>
      <div v-for="s in sessions.slice(0, 10)" :key="s.id" class="flex items-center justify-between py-2 border-b border-gray-50 dark:border-gray-800/50 last:border-0">
        <div class="flex-1 min-w-0">
          <span class="text-sm text-gray-700 dark:text-gray-300">{{ s.todo_title || '未关联任务' }}</span>
          <span class="text-xs text-gray-400 ml-2">{{ new Date(s.started_at).toLocaleTimeString('zh-CN', { hour: '2-digit', minute: '2-digit' }) }}</span>
        </div>
        <span class="text-sm font-semibold" :class="s.completed ? 'text-emerald-500' : 'text-gray-400'">{{ s.duration_min }}分钟</span>
      </div>
    </div>
  </div>
</template>
```

- [ ] **Step 2: 编译验证**

```bash
cd frontend && npx vite build 2>&1 | tail -5
```
Expected: built successfully

- [ ] **Step 3: 提交**

```bash
git add frontend/src/views/FocusView.vue
git commit -m "feat: add FocusView with circle timer, controls, and session history"
```

---

### Task 6: 前端 — 路由与导航入口

**Files:**
- Modify: `frontend/src/router/index.ts`
- Modify: `frontend/src/components/layout/Sidebar.vue`

- [ ] **Step 1: 添加 /focus 路由**

```typescript
// frontend/src/router/index.ts — 在 routes 数组中添加
{ path: '/focus', component: () => import('../views/FocusView.vue') },
```

实际插入位置：`/stats` 路由之后。

- [ ] **Step 2: 侧边栏添加导航按钮**

在 Sidebar.vue 的「统计」按钮之后添加：

```html
<button @click="router.push('/focus')"
  class="w-full text-left px-3 py-2 rounded-xl text-sm font-medium transition-colors flex items-center gap-2"
  :class="route.path === '/focus' ? 'bg-indigo-50 dark:bg-indigo-900/30 text-indigo-700 dark:text-indigo-300' : 'text-gray-600 dark:text-gray-400 hover:bg-gray-50 dark:hover:bg-gray-800'">
  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><polyline points="12 6 12 12 16 14"/></svg>
  沉浸专注
</button>
```

- [ ] **Step 3: 编译验证**

```bash
cd frontend && npx vite build 2>&1 | tail -5
```
Expected: built successfully

- [ ] **Step 4: 提交**

```bash
git add frontend/src/router/index.ts frontend/src/components/layout/Sidebar.vue
git commit -m "feat: add /focus route and sidebar navigation"
```

---

### Task 7: 前端 — 统计页集成专注数据

**Files:**
- Modify: `frontend/src/views/StatsView.vue`

- [ ] **Step 1: 在统计页添加专注统计卡片**

在 StatsView.vue 的 `<script>` 中添加：

```typescript
import { onMounted, ref } from 'vue'
import * as focusApi from '../api/focus'
import type { FocusStats } from '../types/todo'

const focusStats = ref<FocusStats | null>(null)
async function loadFocusStats() {
  try { const r = await focusApi.getFocusStats(); focusStats.value = r.data } catch {}
}
// 在 onMounted 和 watch 中调用 loadFocusStats()
```

在概览卡片区域（3 个卡片之后）添加专注卡片：

```html
<div v-if="focusStats" class="grid grid-cols-4 gap-3 mb-6">
  <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-4 text-center">
    <div class="text-2xl font-bold text-indigo-500">{{ focusStats.today_minutes }}</div>
    <div class="text-xs text-gray-400 dark:text-gray-500 mt-1">今日专注(分)</div>
  </div>
  <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-4 text-center">
    <div class="text-2xl font-bold text-gray-800 dark:text-gray-100">{{ focusStats.total_sessions }}</div>
    <div class="text-xs text-gray-400 dark:text-gray-500 mt-1">总次数</div>
  </div>
  <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-4 text-center">
    <div class="text-2xl font-bold text-emerald-500">{{ focusStats.streak_days }}</div>
    <div class="text-xs text-gray-400 dark:text-gray-500 mt-1">连续天数</div>
  </div>
  <div class="bg-white dark:bg-gray-900 rounded-xl border border-gray-200 dark:border-gray-800 p-4 text-center">
    <div class="text-2xl font-bold text-gray-800 dark:text-gray-100">{{ focusStats.total_minutes }}</div>
    <div class="text-xs text-gray-400 dark:text-gray-500 mt-1">累计(分)</div>
  </div>
</div>
```

- [ ] **Step 2: 编译验证**

```bash
cd frontend && npx vite build 2>&1 | tail -5
```
Expected: built successfully

- [ ] **Step 3: 提交**

```bash
git add frontend/src/views/StatsView.vue
git commit -m "feat: integrate focus stats into statistics page"
```

---

### Task 8: 后端 — 启动验证与 E2E 测试

- [ ] **Step 1: 重启后端**

```bash
cd backend && go build -o todo-back . && pkill -f todo-back; sleep 1 && ./todo-back &
```

- [ ] **Step 2: API 烟雾测试**

```bash
# 获取 token 并测试创建 session
TOKEN=$(curl -s -X POST http://localhost:8080/api/auth/uuid -H 'Content-Type: application/json' -d '{"uuid":"focus-test"}' | python3 -c "import sys,json; print(json.load(sys.stdin).get('token',''))")
# 创建
curl -s -X POST http://localhost:8080/api/focus/start -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' -d '{"duration_min":25}' | python3 -c "import sys,json; d=json.load(sys.stdin); print(f'Created session {d[\"id\"]}')"
# 完成
curl -s -X PATCH "http://localhost:8080/api/focus/1/complete" -H "Authorization: Bearer $TOKEN" | python3 -c "import sys,json; d=json.load(sys.stdin); print(f'Completed: {d[\"completed\"]}')"
# 统计
curl -s http://localhost:8080/api/focus/stats -H "Authorization: Bearer $TOKEN" | python3 -c "import sys,json; d=json.load(sys.stdin); print(f'Today: {d[\"today_minutes\"]}min, Total: {d[\"total_minutes\"]}min')"
```

Expected: Created session 1, Completed: true, Today: 25min, Total: 25min

- [ ] **Step 3: MCP 浏览器验证**

1. 打开 `http://localhost:5173/focus`
2. 点击「开始专注」→ 圆环开始倒计时
3. 等待计时结束 → 提示音 + 通知
4. 自动切换到休息倒计时
5. 在统计页看到今日专注数据

- [ ] **Step 4: 提交**

```bash
git add . && git commit -m "chore: final verification, E2E smoke test passed"
```

---

## Review Checklist

1. **Spec coverage**: ✅ 计时器、配置、沉浸模式、关联待办、统计历史全覆盖
2. **Placeholder scan**: ✅ 无 TBD/TODO
3. **Type consistency**: ✅ FocusSession 前后端字段一致，API 路径一致
