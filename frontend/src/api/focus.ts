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
