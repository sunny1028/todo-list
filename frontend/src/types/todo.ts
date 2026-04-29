export interface Todo {
  id: number
  title: string
  description: string
  priority: 'low' | 'medium' | 'high'
  tags: string
  completed: boolean
  archived: boolean
  due_date: string | null
  list_id: number
  sort_order: number
  created_at: string
  updated_at: string
}

export interface Subtask {
  id: number
  todo_id: number
  title: string
  completed: boolean
}

export interface List {
  id: number
  name: string
  color: string
  created_at: string
  updated_at: string
}

export interface Stats {
  total: number
  active: number
  completed: number
  by_priority: Record<string, number>
  by_tag: Record<string, number>
}

export interface CreateTodoRequest {
  title: string
  description?: string
  priority?: string
  tags?: string
  due_date?: string | null
  list_id?: number
}

export interface UpdateTodoRequest {
  title?: string
  description?: string
  priority?: string
  tags?: string
  completed?: boolean
  archived?: boolean
  due_date?: string | null
  list_id?: number
}

export interface CreateListRequest {
  name: string
  color?: string
}

export interface Attachment {
  id: number
  todo_id: number
  filename: string
  mime_type: string
  size: number
  created_at: string
}
