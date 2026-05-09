# 列表共享 实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为待办事项应用增加列表共享功能 — 用户生成邀请码分享列表，他人通过邀请码加入，支持 "查看" 和 "编辑" 两种权限

**Architecture:** 新增 list_shares + list_share_members 两张表，4 个 API 端点，前端 ShareDialog + JoinDialog 组件。遵循现有 3 层架构 (handler→service→repository)。权限校验通过在现有 todo/list handler 中增加成员身份检查实现。

**Tech Stack:** Go + Gin + GORM + SQLite (backend), Vue 3 + TypeScript + Pinia + Tailwind CSS (frontend)

---

### Task 1: 后端 — 数据模型与数据库迁移

**Files:**
- Create: `backend/models/list_share.go`
- Modify: `backend/database/database.go:12-18`

- [ ] **Step 1: 创建模型文件**

```go
// backend/models/list_share.go
package models

import "time"

type ListShare struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	ListID     uint      `json:"list_id" gorm:"index;not null"`
	Code       string    `json:"code" gorm:"uniqueIndex;not null;size:8"`
	Permission string    `json:"permission" gorm:"not null;default:'view'"`
	CreatedAt  time.Time `json:"created_at"`
}

type ListShareMember struct {
	ID       uint      `json:"id" gorm:"primaryKey"`
	ShareID  uint      `json:"share_id" gorm:"index;not null"`
	UserID   uint      `json:"user_id" gorm:"index;not null"`
	JoinedAt time.Time `json:"joined_at"`
}
```

- [ ] **Step 2: 注册 AutoMigrate**

修改 `backend/database/database.go` 的 AutoMigrate 调用，在参数列表末尾添加 `&models.ListShare{}, &models.ListShareMember{}`：

```go
return DB.AutoMigrate(&models.User{}, &models.Todo{}, &models.List{}, &models.Attachment{}, &models.Subtask{}, &models.FocusSession{}, &models.ListShare{}, &models.ListShareMember{})
```

- [ ] **Step 3: 编译验证**

Run: `cd backend && go build ./...`
Expected: 编译通过

- [ ] **Step 4: Commit**

```bash
git add backend/models/list_share.go backend/database/database.go
git commit -m "feat: 新增 ListShare 和 ListShareMember 模型，注册 AutoMigrate"
```

---

### Task 2: 后端 — Repository 层

**Files:**
- Create: `backend/repository/list_share.go`

- [ ] **Step 1: 创建 repository 文件**

```go
// backend/repository/list_share.go
package repository

import (
	"crypto/rand"
	"math/big"
	"todo-list/backend/database"
	"todo-list/backend/models"
)

func generateCode() string {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	code := make([]byte, 8)
	for i := range code {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		code[i] = chars[n.Int64()]
	}
	return string(code)
}

func CreateShare(listID uint, permission string) (*models.ListShare, error) {
	// If share already exists, return it
	var existing models.ListShare
	if database.DB.Where("list_id = ?", listID).First(&existing).Error == nil {
		return &existing, nil
	}
	s := &models.ListShare{ListID: listID, Code: generateCode(), Permission: permission}
	if err := database.DB.Create(s).Error; err != nil {
		// Retry once on code collision
		s.Code = generateCode()
		if err2 := database.DB.Create(s).Error; err2 != nil {
			return nil, err2
		}
	}
	return s, nil
}

func FindShareByListID(listID uint) (*models.ListShare, error) {
	var s models.ListShare
	err := database.DB.Where("list_id = ?", listID).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func FindShareByCode(code string) (*models.ListShare, error) {
	var s models.ListShare
	err := database.DB.Where("code = ?", code).First(&s).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func DeleteShare(listID uint) error {
	s, err := FindShareByListID(listID)
	if err != nil {
		return err
	}
	database.DB.Where("share_id = ?", s.ID).Delete(&models.ListShareMember{})
	return database.DB.Delete(s).Error
}

func AddShareMember(shareID uint, userID uint) error {
	// Idempotent: if already a member, do nothing
	var count int64
	database.DB.Model(&models.ListShareMember{}).Where("share_id = ? AND user_id = ?", shareID, userID).Count(&count)
	if count > 0 {
		return nil
	}
	m := &models.ListShareMember{ShareID: shareID, UserID: userID}
	return database.DB.Create(m).Error
}

func FindShareMembers(shareID uint) ([]models.ListShareMember, error) {
	var members []models.ListShareMember
	err := database.DB.Where("share_id = ?", shareID).Find(&members).Error
	return members, err
}

// FindSharedListIDs returns list IDs where the user is a share member
func FindSharedListIDs(userID uint) []uint {
	var ids []uint
	database.DB.Model(&models.ListShareMember{}).
		Joins("JOIN list_shares ON list_shares.id = list_share_members.share_id").
		Where("list_share_members.user_id = ?", userID).
		Pluck("list_shares.list_id", &ids)
	return ids
}

// GetUserPermission returns "" if no access, "view" or "edit" if member
func GetUserPermission(listID uint, userID uint) string {
	var s models.ListShare
	if database.DB.Where("list_id = ?", listID).First(&s).Error != nil {
		return ""
	}
	var count int64
	database.DB.Model(&models.ListShareMember{}).Where("share_id = ? AND user_id = ?", s.ID, userID).Count(&count)
	if count == 0 {
		return ""
	}
	return s.Permission
}
```

- [ ] **Step 2: 编译验证**

Run: `cd backend && go build ./...`
Expected: 编译通过

- [ ] **Step 3: Commit**

```bash
git add backend/repository/list_share.go
git commit -m "feat: 实现列表共享 Repository 层 — 增删查、成员管理、权限查询"
```

---

### Task 3: 后端 — Service 层

**Files:**
- Create: `backend/services/list_share.go`
- Modify: `backend/services/list.go:9-11` — GetLists 扩展，含共享列表
- Modify: `backend/services/todo.go:11-13,39-48,51-69,71-119,121-129,131-139,142-152,154-158` — 所有读写方法增加权限校验
- Modify: `backend/repository/list.go:8-11` — 新增 FindListsByIDs

- [ ] **Step 1: 创建 list_share service**

```go
// backend/services/list_share.go
package services

import (
	"errors"
	"todo-list/backend/models"
	"todo-list/backend/repository"
)

func CreateShare(userID uint, listID uint, permission string) (*models.ListShare, error) {
	list, err := repository.FindListByID(userID, listID)
	if err != nil {
		return nil, errors.New("list not found")
	}
	if list.UserID != userID {
		return nil, errors.New("only the list owner can share")
	}
	if permission != "view" && permission != "edit" {
		permission = "view"
	}
	return repository.CreateShare(listID, permission)
}

func GetShare(userID uint, listID uint) (*models.ListShare, []models.ListShareMember, error) {
	list, err := repository.FindListByID(userID, listID)
	if err != nil {
		return nil, nil, errors.New("list not found")
	}
	if list.UserID != userID {
		return nil, nil, errors.New("only the list owner can view share info")
	}
	s, err := repository.FindShareByListID(listID)
	if err != nil {
		return nil, nil, errors.New("share not found")
	}
	members, _ := repository.FindShareMembers(s.ID)
	return s, members, nil
}

func DeleteShare(userID uint, listID uint) error {
	list, err := repository.FindListByID(userID, listID)
	if err != nil {
		return errors.New("list not found")
	}
	if list.UserID != userID {
		return errors.New("only the list owner can revoke share")
	}
	return repository.DeleteShare(listID)
}

func JoinShare(userID uint, code string) (*models.List, error) {
	s, err := repository.FindShareByCode(code)
	if err != nil {
		return nil, errors.New("invalid share code")
	}
	// Prevent joining own list
	var list models.List
	if database.DB.First(&list, s.ListID).Error != nil {
		return nil, errors.New("list not found")
	}
	if list.UserID == userID {
		return nil, errors.New("cannot join your own list")
	}
	if err := repository.AddShareMember(s.ID, userID); err != nil {
		return nil, err
	}
	return &list, nil
}
```

需要在文件顶部 import `"todo-list/backend/database"` 和 `"todo-list/backend/models"`（JoinShare 中使用了 `database.DB` 和 `models.List`）。

- [ ] **Step 2: 修改 GetLists — 包含共享列表**

修改 `backend/services/list.go` 的 `GetLists`:

```go
func GetLists(userID uint) ([]models.List, error) {
	own, err := repository.FindAllLists(userID)
	if err != nil {
		return nil, err
	}
	sharedIDs := repository.FindSharedListIDs(userID)
	if len(sharedIDs) == 0 {
		return own, nil
	}
	shared, err := repository.FindListsByIDs(sharedIDs)
	if err != nil {
		return own, nil // best-effort
	}
	return append(own, shared...), nil
}
```

需要在 `backend/repository/list.go` 添加:

```go
func FindListsByIDs(ids []uint) ([]models.List, error) {
	var lists []models.List
	err := database.DB.Where("id IN ?", ids).Find(&lists).Error
	return lists, err
}
```

- [ ] **Step 3: 修改 GetTodos — 校验共享列表权限**

修改 `backend/services/todo.go` 的 `GetTodos`，在获取 todos 前校验权限。先在文件开头添加辅助函数：

```go
func canAccessList(userID uint, listID uint) bool {
	if listID == 0 {
		return true // all lists (no specific list)
	}
	// Owner check
	list, err := repository.FindListByID(userID, listID)
	if err == nil && list.UserID == userID {
		return true
	}
	// Share member check
	perm := repository.GetUserPermission(listID, userID)
	return perm != "" // view or edit both can read
}

func canWriteList(userID uint, listID uint) bool {
	if listID == 0 {
		return true
	}
	list, err := repository.FindListByID(userID, listID)
	if err == nil && list.UserID == userID {
		return true
	}
	return repository.GetUserPermission(listID, userID) == "edit"
}
```

在 `GetTodos` 函数开头增加校验：

```go
func GetTodos(userID uint, listID uint, status, priority, tag, search string) ([]models.Todo, error) {
	if !canAccessList(userID, listID) {
		return nil, errors.New("access denied")
	}
	return repository.FindAll(userID, listID, status, priority, tag, search)
}
```

- [ ] **Step 4: 修改 CreateTodo/UpdateTodo/ToggleTodo 等写入操作**

在 `CreateTodo`、`UpdateTodo`、`ToggleTodo`、`ArchiveTodo`、`UnarchiveTodo`、`DeleteTodo`、`ReorderTodos` 中增加 `canWriteList` 校验。以 CreateTodo 为例：

```go
func CreateTodo(todo *models.Todo) error {
	if todo.Title == "" {
		return errors.New("title is required")
	}
	if todo.Priority == "" {
		todo.Priority = "medium"
	}
	if !canWriteList(todo.UserID, todo.ListID) {
		return errors.New("access denied")
	}
	maxOrder, _ := repository.GetMaxSortOrder(todo.UserID)
	todo.SortOrder = maxOrder + 1
	return repository.Create(todo)
}
```

对其他写入方法（UpdateTodo, ToggleTodo, ArchiveTodo, UnarchiveTodo, DeleteTodo, ReorderTodos）做同样的处理 — 在函数体开头加入 `canWriteList` 检查。每个方法的 `userID` 和 `listID` 来源：
- `UpdateTodo(id)` → 先从 todo 获取 listID，再校验
- `ToggleTodo(id)` → 从 FindByID 获取 listID
- 其他类似

注意：部分方法需要先从 DB 取出 todo 才能知道 listID。在已有 FindByID 调用的函数中，在 FindByID 之后立即校验。

- [ ] **Step 5: 编译验证**

Run: `cd backend && go build ./...`
Expected: 编译通过

- [ ] **Step 6: Commit**

```bash
git add backend/services/list_share.go backend/services/list.go backend/services/todo.go backend/repository/list.go
git commit -m "feat: 列表共享 Service 层 — 创建/撤销共享、加入列表、权限校验"
```

---

### Task 4: 后端 — Handler 层与路由

**Files:**
- Create: `backend/handlers/list_share.go`
- Modify: `backend/router/router.go:69-73` — 添加 share/join 路由

- [ ] **Step 1: 创建 handler 文件**

```go
// backend/handlers/list_share.go
package handlers

import (
	"net/http"
	"strconv"
	"todo-list/backend/services"

	"github.com/gin-gonic/gin"
)

func CreateShare(c *gin.Context) {
	userID := c.GetUint("user_id")
	listID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var body struct {
		Permission string `json:"permission"`
	}
	c.ShouldBindJSON(&body)
	share, err := services.CreateShare(userID, uint(listID), body.Permission)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, share)
}

func GetShare(c *gin.Context) {
	userID := c.GetUint("user_id")
	listID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	share, members, err := services.GetShare(userID, uint(listID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"share":   share,
		"members": members,
	})
}

func DeleteShare(c *gin.Context) {
	userID := c.GetUint("user_id")
	listID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	if err := services.DeleteShare(userID, uint(listID)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "share revoked"})
}

func JoinList(c *gin.Context) {
	userID := c.GetUint("user_id")
	var body struct {
		Code string `json:"code"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "code is required"})
		return
	}
	list, err := services.JoinShare(userID, body.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, list)
}
```

- [ ] **Step 2: 注册路由**

修改 `backend/router/router.go`，在 `// Lists` 注释块附近（`protected.DELETE("/lists/:id", ...)` 之后）添加：

```go
			// List sharing
			protected.POST("/lists/:id/share", handlers.CreateShare)
			protected.GET("/lists/:id/share", handlers.GetShare)
			protected.DELETE("/lists/:id/share", handlers.DeleteShare)
			protected.POST("/lists/join", handlers.JoinList)
```

- [ ] **Step 3: 编译验证**

Run: `cd backend && go build ./...`
Expected: 编译通过

- [ ] **Step 4: Commit**

```bash
git add backend/handlers/list_share.go backend/router/router.go
git commit -m "feat: 列表共享 Handler 层与路由注册"
```

---

### Task 5: 前端 — 类型定义与 API 层

**Files:**
- Modify: `frontend/src/types/todo.ts` — 新增类型
- Create: `frontend/src/api/shares.ts` — API 调用

- [ ] **Step 1: 添加类型定义**

在 `frontend/src/types/todo.ts` 的 List 接口后添加：

```typescript
export interface ListShare {
  id: number
  list_id: number
  code: string
  permission: 'view' | 'edit'
  created_at: string
}

export interface ListShareMember {
  id: number
  share_id: number
  user_id: number
  joined_at: string
}
```

同时在 List 接口中增加字段：

```typescript
export interface List {
  id: number
  name: string
  color: string
  user_id?: number
  shared?: boolean
  permission?: 'view' | 'edit'
  created_at: string
  updated_at: string
}
```

- [ ] **Step 2: 创建 API 文件**

```typescript
// frontend/src/api/shares.ts
import type { List, ListShare, ListShareMember } from '../types/todo'
import api from './index'

export function createShare(listId: number, permission: string) {
  return api.post<ListShare>(`/lists/${listId}/share`, { permission })
}

export function getShare(listId: number) {
  return api.get<{ share: ListShare; members: ListShareMember[] }>(`/lists/${listId}/share`)
}

export function deleteShare(listId: number) {
  return api.delete(`/lists/${listId}/share`)
}

export function joinList(code: string) {
  return api.post<List>('/lists/join', { code })
}
```

- [ ] **Step 3: 类型检查**

Run: `cd frontend && npx vue-tsc --noEmit --project tsconfig.json`
Expected: 编译通过

- [ ] **Step 4: Commit**

```bash
git add frontend/src/types/todo.ts frontend/src/api/shares.ts
git commit -m "feat: 前端共享列表类型定义与 API 层"
```

---

### Task 6: 前端 — ShareDialog 组件

**Files:**
- Create: `frontend/src/components/ui/ShareDialog.vue`

- [ ] **Step 1: 创建 ShareDialog 组件**

```vue
<script setup lang="ts">
import { onMounted, ref } from 'vue'
import * as shareApi from '../../api/shares'
import type { ListShare, ListShareMember } from '../../types/todo'
import { useToast } from '../../stores/toast'

const props = defineProps<{ listId: number; listName: string }>()
const emit = defineEmits<{ close: [] }>()

const share = ref<ListShare | null>(null)
const members = ref<ListShareMember[]>([])
const permission = ref<'view' | 'edit'>('edit')
const loading = ref(false)

onMounted(async () => {
  try {
    const res = await shareApi.getShare(props.listId)
    share.value = res.data.share
    members.value = res.data.members
  } catch {
    // not shared yet
  }
})

async function handleCreate() {
  loading.value = true
  try {
    const res = await shareApi.createShare(props.listId, permission.value)
    share.value = res.data
    useToast().show('共享链接已生成', 'success')
  } catch {
    useToast().show('创建共享失败', 'error')
  } finally {
    loading.value = false
  }
}

async function handleRevoke() {
  try {
    await shareApi.deleteShare(props.listId)
    share.value = null
    members.value = []
    useToast().show('共享已撤销', 'success')
  } catch {
    useToast().show('撤销失败', 'error')
  }
}

function copyLink() {
  if (!share.value) return
  const link = `${window.location.origin}/join?code=${share.value.code}`
  navigator.clipboard.writeText(link)
  useToast().show('链接已复制', 'success')
}
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/30" @click.self="emit('close')">
    <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl p-6 w-[380px] max-w-[90vw]">
      <h3 class="text-base font-bold text-gray-800 dark:text-gray-100 mb-4">共享列表「{{ listName }}」</h3>

      <!-- Not shared yet: create -->
      <div v-if="!share" class="space-y-3">
        <div class="flex items-center gap-2">
          <label class="text-xs text-gray-500">权限</label>
          <select v-model="permission" class="px-2.5 py-1.5 border border-gray-200 dark:border-gray-700 rounded-lg text-sm bg-white dark:bg-gray-800 dark:text-gray-100 outline-none">
            <option value="view">仅查看</option>
            <option value="edit">可编辑</option>
          </select>
        </div>
        <button @click="handleCreate" :disabled="loading"
          class="w-full py-2 bg-indigo-500 text-white rounded-lg text-sm font-semibold hover:bg-indigo-600 disabled:opacity-40 transition-colors">
          {{ loading ? '创建中…' : '生成共享链接' }}
        </button>
      </div>

      <!-- Shared: show code & members -->
      <div v-else class="space-y-4">
        <div class="flex items-center gap-2">
          <span class="text-xs text-gray-400">邀请码</span>
          <code class="px-3 py-1.5 bg-gray-100 dark:bg-gray-700 rounded-lg text-lg font-mono font-bold text-indigo-600 dark:text-indigo-400 tracking-wider select-all">{{ share.code }}</code>
        </div>
        <button @click="copyLink"
          class="w-full py-2 border border-indigo-200 dark:border-indigo-800 text-indigo-600 dark:text-indigo-400 rounded-lg text-sm font-medium hover:bg-indigo-50 dark:hover:bg-indigo-900/30 transition-colors">
          复制分享链接
        </button>
        <div>
          <div class="text-xs text-gray-400 mb-2">成员 ({{ members.length }})</div>
          <div v-if="members.length === 0" class="text-xs text-gray-300 dark:text-gray-600">暂无成员</div>
          <div v-for="m in members" :key="m.id" class="text-sm text-gray-600 dark:text-gray-400 py-1">
            ID: {{ m.user_id }}
          </div>
        </div>
        <button @click="handleRevoke"
          class="w-full py-2 text-red-500 rounded-lg text-sm font-medium hover:bg-red-50 dark:hover:bg-red-900/20 transition-colors">
          撤销共享
        </button>
      </div>

      <button @click="emit('close')"
        class="w-full py-2 mt-3 text-gray-400 rounded-lg text-sm hover:text-gray-600 dark:hover:text-gray-300 transition-colors">
        关闭
      </button>
    </div>
  </div>
</template>
```

- [ ] **Step 2: 类型检查**

Run: `cd frontend && npx vue-tsc --noEmit --project tsconfig.json`
Expected: 编译通过

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/ui/ShareDialog.vue
git commit -m "feat: ShareDialog 组件 — 创建共享、复制链接、查看成员、撤销共享"
```

---

### Task 7: 前端 — JoinDialog 组件

**Files:**
- Create: `frontend/src/components/ui/JoinDialog.vue`

- [ ] **Step 1: 创建 JoinDialog 组件**

```vue
<script setup lang="ts">
import { ref } from 'vue'
import * as shareApi from '../../api/shares'
import { useLists } from '../../stores/lists'
import { useToast } from '../../stores/toast'

const emit = defineEmits<{ close: [] }>()

const code = ref('')
const loading = ref(false)
const listStore = useLists()

async function handleJoin() {
  if (!code.value.trim()) return
  loading.value = true
  try {
    await shareApi.joinList(code.value.trim().toUpperCase())
    await listStore.fetchLists()
    useToast().show('已加入共享列表', 'success')
    emit('close')
  } catch {
    useToast().show('加入失败，请检查邀请码', 'error')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black/30" @click.self="emit('close')">
    <div class="bg-white dark:bg-gray-800 rounded-2xl shadow-xl p-6 w-[360px] max-w-[90vw]">
      <h3 class="text-base font-bold text-gray-800 dark:text-gray-100 mb-4">加入共享列表</h3>
      <input v-model="code" @keydown.enter="handleJoin"
        placeholder="输入邀请码"
        class="w-full px-3 py-2.5 border border-gray-200 dark:border-gray-700 rounded-lg text-sm bg-white dark:bg-gray-800 dark:text-gray-100 outline-none focus:border-indigo-400 mb-3 uppercase"
        maxlength="8" />
      <button @click="handleJoin" :disabled="loading || !code.trim()"
        class="w-full py-2 bg-indigo-500 text-white rounded-lg text-sm font-semibold hover:bg-indigo-600 disabled:opacity-40 transition-colors mb-2">
        {{ loading ? '加入中…' : '加入' }}
      </button>
      <button @click="emit('close')"
        class="w-full py-2 text-gray-400 rounded-lg text-sm hover:text-gray-600 dark:hover:text-gray-300 transition-colors">
        取消
      </button>
    </div>
  </div>
</template>
```

- [ ] **Step 2: 类型检查**

Run: `cd frontend && npx vue-tsc --noEmit --project tsconfig.json`
Expected: 编译通过

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/ui/JoinDialog.vue
git commit -m "feat: JoinDialog 组件 — 通过邀请码加入共享列表"
```

---

### Task 8: 前端 — 侧边栏与权限展示

**Files:**
- Modify: `frontend/src/components/layout/Sidebar.vue` — 共享按钮 + 加入按钮 + 共享列表图标
- Modify: `frontend/src/types/todo.ts` — 确认 List 类型含 shared/permission 字段

- [ ] **Step 1: 更新侧边栏 — 添加共享按钮和加入按钮**

修改 `Sidebar.vue` 的 `<script>` 部分，添加：

```typescript
import ShareDialog from '../ui/ShareDialog.vue'
import JoinDialog from '../ui/JoinDialog.vue'

const shareListId = ref<number | null>(null)
const showJoinDialog = ref(false)
```

在自定义列表按钮内，增加共享图标按钮（在删除按钮旁边）：

```html
<!-- 在 deleteList 按钮之后，</button> 模板结束前 -->
<button
  @click.stop="shareListId = list.id"
  class="opacity-0 group-hover:opacity-100 text-gray-400 hover:text-indigo-500 transition-all ml-0.5"
  title="共享列表">
  <svg width="12" height="12" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="18" cy="5" r="3"/><circle cx="6" cy="12" r="3"/><circle cx="18" cy="19" r="3"/><line x1="8.59" y1="13.51" x2="15.42" y2="17.49"/><line x1="15.41" y1="6.51" x2="8.59" y2="10.49"/></svg>
</button>
```

在侧边栏底部的「加入列表」处，在视图链接之前添加：

```html
<button @click="showJoinDialog = true"
  class="w-full text-left px-3 py-2 rounded-xl text-sm font-medium transition-colors flex items-center gap-2 text-gray-400 hover:text-indigo-500 hover:bg-indigo-50 dark:hover:bg-indigo-900/20">
  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="16"/><line x1="8" y1="12" x2="16" y2="12"/></svg>
  加入列表
</button>
```

在模板末尾（`</aside>` 之前）添加：

```html
<ShareDialog v-if="shareListId" :list-id="shareListId" :list-name="listStore.lists.find(l => l.id === shareListId)?.name || ''" @close="shareListId = null" />
<JoinDialog v-if="showJoinDialog" @close="showJoinDialog = false" />
```

- [ ] **Step 2: 共享列表标识**

在自定义列表按钮上，如果 list 有 `shared` 标记，显示图标：

```html
<!-- 在列表名前面 -->
<span v-if="list.shared" class="text-xs mr-1" title="共享列表">👥</span>
```

当前后端 `GET /api/lists` 返回的列表中不包含 `shared` 字段。需要在 `GetLists` service 中标记。修改 `backend/services/list.go`:

```go
func GetLists(userID uint) ([]models.List, error) {
	own, err := repository.FindAllLists(userID)
	if err != nil {
		return nil, err
	}
	sharedIDs := repository.FindSharedListIDs(userID)
	sharedMap := make(map[uint]bool)
	for _, id := range sharedIDs {
		sharedMap[id] = true
	}
	if len(sharedIDs) > 0 {
		shared, err := repository.FindListsByIDs(sharedIDs)
		if err == nil {
			own = append(own, shared...)
		}
	}
	// Mark shared lists in response — handled by frontend checking list.user_id !== current user
	return own, nil
}
```

前端侧边栏中判断共享列表：比较 `list.user_id` 与当前用户 ID。需要把 auth store 暴露出来。在 `Sidebar.vue` 的 `<script>` 中，`auth` store 已导入。在列表按钮中：

```html
<span v-if="list.user_id && list.user_id !== auth.user?.id" class="text-[10px] ml-1" title="共享列表">👥</span>
```

- [ ] **Step 3: 权限限制 UI**

只读列表下隐藏新增表单和编辑/删除按钮。需要在 `TodoList.vue` 中判断当前列表权限。

修改 `Sidebar.vue` 的 `selectList` 函数，存储当前列表的权限：

```typescript
const currentPermission = ref<'view' | 'edit' | ''>('')

function selectList(id: number) {
  todoStore.setList(id)
  todoStore.fetchTodos()
  router.replace({ query: {} })
  // Check if this is a shared list with view-only permission
  const list = listStore.lists.find(l => l.id === id)
  if (list && list.user_id !== undefined && list.user_id !== auth.user?.id) {
    // TODO: get permission from API
    currentPermission.value = list.permission || 'view'
  } else {
    currentPermission.value = ''
  }
}
```

等等，List 接口需要包含 `user_id` 和 `permission` 字段。让我更新类型：

```typescript
export interface List {
  id: number
  name: string
  color: string
  user_id: number
  permission?: 'view' | 'edit'
  created_at: string
  updated_at: string
}
```

后端 `GetLists` 需要为共享列表填充 `permission`。修改 `backend/services/list.go`:

```go
func GetLists(userID uint) ([]models.List, error) {
	own, err := repository.FindAllLists(userID)
	if err != nil {
		return nil, err
	}
	sharedIDs := repository.FindSharedListIDs(userID)
	if len(sharedIDs) > 0 {
		shared, err := repository.FindListsByIDs(sharedIDs)
		if err == nil {
			own = append(own, shared...)
		}
	}
	return own, nil
}
```

List 模型需要加 `Permission` 虚拟字段。修改 `backend/models/list.go`:

```go
type List struct {
	UserID     uint      `json:"user_id" gorm:"index;not null;default:0"`
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"not null"`
	Color     string    `json:"color" gorm:"default:'#6366f1'"`
	Permission string   `json:"permission" gorm:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
```

在 `GetLists` 中填充 permission：

```go
func GetLists(userID uint) ([]models.List, error) {
	own, err := repository.FindAllLists(userID)
	if err != nil {
		return nil, err
	}
	sharedIDs := repository.FindSharedListIDs(userID)
	if len(sharedIDs) > 0 {
		shared, err := repository.FindListsByIDs(sharedIDs)
		if err == nil {
			for i := range shared {
				shared[i].Permission = repository.GetUserPermission(shared[i].ID, userID)
			}
			own = append(own, shared...)
		}
	}
	return own, nil
}
```

前端在 `Sidebar.vue` 中通过 computed 暴露 `currentPermission`，然后在 `TodoList.vue` 中根据权限隐藏 TodoForm 和编辑/删除按钮。

修改 `Sidebar.vue` 提供权限：

```typescript
const currentPermission = computed(() => {
  const list = listStore.lists.find(l => l.id === todoStore.currentListId)
  if (!list) return ''
  if (list.user_id === auth.user?.id) return ''  // owner, full access
  return list.permission || 'view'
})
// Expose to child components via provide/inject or just use in template
```

最简单的方法：在 `Sidebar.vue` 中通过 `provide` 暴露 permission，`TodoList.vue` 和 `TodoItem.vue` 通过 `inject` 读取。或者直接在 `TodoList.vue` 中检查：

```typescript
// In TodoList.vue
import { useLists } from '../stores/lists'
import { useAuth } from '../stores/auth'

const listStore = useLists()
const auth = useAuth()
const isReadonly = computed(() => {
  const list = listStore.lists.find(l => l.id === store.currentListId)
  if (!list || list.user_id === undefined) return false
  return list.user_id !== auth.user?.id && list.permission === 'view'
})
```

- [ ] **Step 4: 类型检查**

Run: `cd frontend && npx vue-tsc --noEmit --project tsconfig.json`
Expected: 编译通过

- [ ] **Step 5: Commit**

```bash
git add frontend/src/components/layout/Sidebar.vue backend/services/list.go backend/models/list.go
git commit -m "feat: 侧边栏共享按钮、加入列表入口、权限标识"
```

---

### Task 9: 前端 — TodoList/TodoItem 权限禁用

**Files:**
- Modify: `frontend/src/components/TodoList.vue:225-226` — 隐藏 TodoForm 当只读
- Modify: `frontend/src/components/TodoItem.vue:265-278` — 禁用编辑/删除按钮当只读

- [ ] **Step 1: TodoList 中隐藏新增表单**

在 `TodoList.vue` 的 `<script>` 中添加：

```typescript
import { useLists } from '../stores/lists'
import { useAuth } from '../stores/auth'

const listStore = useLists()
const authStore = useAuth()
const isReadonly = computed(() => {
  const list = listStore.lists.find(l => l.id === store.currentListId)
  if (!list) return false
  return list.user_id !== authStore.user?.id && list.permission === 'view'
})
```

在模板中，将 `<TodoForm>` 包裹在条件中：

```html
<TodoForm v-if="!isReadonly" @created="() => {}" />
```

- [ ] **Step 2: TodoItem 中禁用编辑/删除**

在 `TodoItem.vue` 中添加 readonly check。通过 inject 或直接检查 list store。

最简单的：由于 TodoItem 是纯展示组件，不直接访问 list store，可以通过 prop 传递 readonly。

在 `TodoList.vue` 中传递 prop：

```html
<TodoItem :todo="asTodo(item)" :readonly="isReadonly" />
```

在 `TodoItem.vue` 中接收 prop：

```typescript
const props = defineProps<{ todo: Todo; readonly?: boolean }>()
```

在模板中条件禁用编辑和删除按钮：

```html
<button v-if="!props.readonly" @click.stop="startEdit" ...>
<button v-if="!props.readonly" @click.stop="showConfirm = true" ...>
```

- [ ] **Step 3: 类型检查和构建**

Run: `cd frontend && npx vue-tsc --noEmit --project tsconfig.json`
Expected: 编译通过
Run: `cd frontend && npx vite build`
Expected: 构建成功

- [ ] **Step 4: Commit**

```bash
git add frontend/src/components/TodoList.vue frontend/src/components/TodoItem.vue
git commit -m "feat: 只读共享列表下隐藏编辑/删除/新增功能"
```

---

### Task 10: 端到端验证与收尾

- [ ] **Step 1: 重建后端并启动**

```bash
cd backend && go build -o todo-back . && kill $(pgrep todo-back) 2>/dev/null; ./todo-back &
```

- [ ] **Step 2: API 烟雾测试**

```bash
# 1. 获取 token
TOKEN=$(curl -s http://localhost:8080/api/auth/uuid -H "Content-Type: application/json" -d '{"uuid":"test-share-1"}' | python3 -c "import sys,json; print(json.load(sys.stdin)['token'])")

# 2. 创建列表
LIST_ID=$(curl -s http://localhost:8080/api/lists -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{"name":"测试共享"}' | python3 -c "import sys,json; print(json.load(sys.stdin)['id'])")

# 3. 创建共享
curl -s http://localhost:8080/api/lists/$LIST_ID/share -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" -d '{"permission":"edit"}' | python3 -m json.tool

# 4. 另一个用户加入
TOKEN2=$(curl -s http://localhost:8080/api/auth/uuid -H "Content-Type: application/json" -d '{"uuid":"test-share-2"}' | python3 -c "import sys,json; print(json.load(sys.stdin)['token'])")
# (用上面返回的 code)
curl -s http://localhost:8080/api/lists/join -H "Authorization: Bearer $TOKEN2" -H "Content-Type: application/json" -d '{"code":"XXXX"}' | python3 -m json.tool
```

- [ ] **Step 3: 浏览器测试**

打开 `http://localhost:5173/`，验证：
- 自定义列表悬停显示共享按钮
- 点击共享按钮打开弹窗 → 生成邀请码 → 可复制链接
- 另一个浏览器窗口/隐身模式 → 加入列表 → 侧边栏显示共享列表
- 只读列表不显示添加待办和编辑/删除按钮

- [ ] **Step 4: Commit（如有剩余改动）**

```bash
git add -A && git commit -m "chore: 端到端测试通过，修复最后的问题"
```
