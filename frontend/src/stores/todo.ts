import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { Todo, CreateTodoRequest, UpdateTodoRequest, Stats } from '../types/todo'
import * as todoApi from '../api/todos'
import { useToast } from './toast'

export const useTodos = defineStore('todos', () => {
  const todos = ref<Todo[]>([])
  const loading = ref(false)
  const error = ref('')
  const currentListId = ref<number>(0)

  function setList(id: number) {
    currentListId.value = id
  }

  async function fetchTodos(params?: Record<string, string>) {
    loading.value = true
    error.value = ''
    try {
      const merged = { ...params }
      if (currentListId.value > 0 && !merged.list_id) {
        merged.list_id = String(currentListId.value)
      }
      const res = await todoApi.listTodos(merged)
      todos.value = res.data
    } catch (e) {
      error.value = '加载待办失败'
    } finally {
      loading.value = false
    }
  }

  async function addTodo(data: CreateTodoRequest) {
    const toast = useToast()
    try {
      if (currentListId.value > 0 && !data.list_id) {
        data.list_id = currentListId.value
      }
      const res = await todoApi.createTodo(data)
      if (currentListId.value === 0 || res.data.list_id === currentListId.value) {
        todos.value.push(res.data)
      }
      toast.show('已添加', 'success')
    } catch {
      toast.show('添加失败', 'error')
    }
  }

  async function editTodo(id: number, data: UpdateTodoRequest) {
    const toast = useToast()
    try {
      const res = await todoApi.updateTodo(id, data)
      const idx = todos.value.findIndex((t) => t.id === id)
      if (idx !== -1) todos.value[idx] = res.data
      toast.show('已更新', 'success')
    } catch {
      toast.show('更新失败', 'error')
    }
  }

  async function toggle(id: number) {
    const toast = useToast()
    try {
      const res = await todoApi.toggleTodo(id)
      const idx = todos.value.findIndex((t) => t.id === id)
      if (idx !== -1) todos.value[idx] = res.data
    } catch {
      toast.show('操作失败', 'error')
    }
  }

  async function removeTodo(id: number) {
    const toast = useToast()
    try {
      await todoApi.deleteTodo(id)
      todos.value = todos.value.filter((t) => t.id !== id)
      toast.show('已删除', 'success')
    } catch {
      toast.show('删除失败', 'error')
    }
  }

  async function clearCompleted() {
    const toast = useToast()
    const completed = todos.value.filter((t) => t.completed)
    try {
      await Promise.all(completed.map((t) => todoApi.deleteTodo(t.id)))
      todos.value = todos.value.filter((t) => !t.completed)
      toast.show(`已清除 ${completed.length} 个已完成项`, 'success')
    } catch {
      toast.show('清除失败', 'error')
    }
  }

  async function reorder(ids: number[]) {
    try {
      await todoApi.reorderTodos(ids)
    } catch {
      useToast().show('排序失败', 'error')
    }
  }

  async function downloadExport(format: 'json' | 'csv') {
    const toast = useToast()
    try {
      const res = await todoApi.exportTodos(format, currentListId.value || undefined)
      const blob = new Blob([res.data])
      const url = URL.createObjectURL(blob)
      const a = document.createElement('a')
      a.href = url
      a.download = `todos.${format}`
      a.click()
      URL.revokeObjectURL(url)
      toast.show('导出成功', 'success')
    } catch {
      toast.show('导出失败', 'error')
    }
  }

  async function fetchStats(listId?: number): Promise<Stats> {
    const res = await todoApi.getStats(listId || currentListId.value || undefined)
    return res.data
  }

  return { todos, loading, error, currentListId, setList, fetchTodos, addTodo, editTodo, toggle, removeTodo, clearCompleted, reorder, downloadExport, fetchStats }
})
