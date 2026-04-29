import axios from 'axios'
import type { Todo, List, Stats, Attachment, CreateTodoRequest, UpdateTodoRequest, CreateListRequest } from '../types/todo'

const api = axios.create({ baseURL: '/api' })

// Todos
export function listTodos(params?: Record<string, string>) {
  return api.get<Todo[]>('/todos', { params })
}

export function getTodo(id: number) {
  return api.get<Todo>(`/todos/${id}`)
}

export function createTodo(data: CreateTodoRequest) {
  return api.post<Todo>('/todos', data)
}

export function updateTodo(id: number, data: UpdateTodoRequest) {
  return api.put<Todo>(`/todos/${id}`, data)
}

export function toggleTodo(id: number) {
  return api.patch<Todo>(`/todos/${id}/toggle`)
}

export function deleteTodo(id: number) {
  return api.delete(`/todos/${id}`)
}

export function reorderTodos(ids: number[]) {
  return api.put('/todos/reorder', { ids })
}

export function exportTodos(format: 'json' | 'csv', listId?: number) {
  return api.get('/todos/export', { params: { format, list_id: listId }, responseType: 'blob' })
}

export function getStats(listId?: number) {
  return api.get<Stats>('/todos/stats', { params: { list_id: listId } })
}

// Attachments
export function listAttachments(todoId: number) {
  return api.get<Attachment[]>(`/todos/${todoId}/attachments`)
}

export function uploadAttachment(todoId: number, file: File) {
  const form = new FormData()
  form.append('file', file)
  return api.post<Attachment>(`/todos/${todoId}/attachments`, form)
}

export function deleteAttachment(id: number) {
  return api.delete(`/attachments/${id}`)
}

export function attachmentUrl(id: number) {
  return `/api/attachments/${id}`
}

// Lists
export function listLists() {
  return api.get<List[]>('/lists')
}

export function createList(data: CreateListRequest) {
  return api.post<List>('/lists', data)
}

export function updateList(id: number, data: Partial<CreateListRequest>) {
  return api.put<List>(`/lists/${id}`, data)
}

export function deleteList(id: number) {
  return api.delete(`/lists/${id}`)
}
