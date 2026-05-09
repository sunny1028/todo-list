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
