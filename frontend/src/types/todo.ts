export interface Todo {
  id: number
  title: string
  description: string
  priority: 'low' | 'medium' | 'high'
  effort: '' | 'easy' | 'medium' | 'hard'
  recurrence: '' | 'daily' | 'weekly' | 'monthly'
  tags: string
  completed: boolean
  archived: boolean
  due_date: string | null
  list_id: number
  sort_order: number
  subtask_count: number
  subtask_completed: number
  created_at: string
  updated_at: string
}

export interface Subtask {
  id: number
  todo_id: number
  title: string
  completed: boolean
  sort_order: number
}

export interface List {
  id: number
  name: string
  color: string
  created_at: string
  updated_at: string
}

export interface DailyTrend {
  date: string
  created: number
  completed: number
}

export interface Stats {
  total: number
  active: number
  completed: number
  by_priority: Record<string, number>
  by_tag: Record<string, number>
  daily_trends?: DailyTrend[]
}

export interface CreateTodoRequest {
  title: string
  description?: string
  priority?: string
  effort?: string
  recurrence?: string
  tags?: string
  due_date?: string | null
  list_id?: number
}

export interface UpdateTodoRequest {
  title?: string
  description?: string
  priority?: string
  effort?: string
  recurrence?: string
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

export interface YesterdaySummary {
  tasks_created: number
  tasks_completed: number
  focus_minutes: number
}

export interface WeeklyReport {
  tasks_created: number
  tasks_completed: number
  focus_minutes: number
  best_day: string
}

export interface DailyFocusItem {
  date: string
  minutes: number
}

export interface ReviewStats {
  yesterday_summary: YesterdaySummary
  weekly_report: WeeklyReport
  daily_focus: DailyFocusItem[]
  mastery_score: number
}

export interface FocusStats {
  today_minutes: number
  total_minutes: number
  total_sessions: number
  streak_days: number
}
