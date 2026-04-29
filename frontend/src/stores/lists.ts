import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { List } from '../types/todo'
import * as api from '../api/todos'
import { useToast } from './toast'

export const useLists = defineStore('lists', () => {
  const lists = ref<List[]>([])
  const loading = ref(false)

  async function fetchLists() {
    loading.value = true
    try {
      const res = await api.listLists()
      lists.value = res.data
    } catch {
      // silent
    } finally {
      loading.value = false
    }
  }

  async function addList(name: string, color?: string) {
    const toast = useToast()
    try {
      const res = await api.createList({ name, color })
      lists.value.push(res.data)
      toast.show('列表已创建', 'success')
      return res.data
    } catch {
      toast.show('创建列表失败', 'error')
      return null
    }
  }

  async function updateList(id: number, data: { name?: string; color?: string }) {
    const toast = useToast()
    try {
      const res = await api.updateList(id, data)
      const idx = lists.value.findIndex((l) => l.id === id)
      if (idx !== -1) lists.value[idx] = res.data
      toast.show('列表已更新', 'success')
    } catch {
      toast.show('更新失败', 'error')
    }
  }

  async function removeList(id: number) {
    const toast = useToast()
    try {
      await api.deleteList(id)
      lists.value = lists.value.filter((l) => l.id !== id)
      toast.show('列表已删除', 'success')
    } catch {
      toast.show('删除失败', 'error')
    }
  }

  return { lists, loading, fetchLists, addList, updateList, removeList }
})
