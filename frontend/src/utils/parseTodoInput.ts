export interface ParsedTodoInput {
  title: string
  dueDate: string | null      // YYYY-MM-DD
  priority: string | null     // 'low' | 'medium' | 'high'
  tags: string | null         // comma-separated
  effort: string | null       // 'easy' | 'medium' | 'hard'
  recurrence: string | null   // 'daily' | 'weekly' | 'monthly'
}

const PRIORITY_MAP: Record<string, string> = {
  '高': 'high', 'h': 'high', '1': 'high',
  '中': 'medium', 'm': 'medium', '2': 'medium',
  '低': 'low', 'l': 'low', '3': 'low',
}

const EFFORT_MAP: Record<string, string> = {
  '简单': 'easy', 'easy': 'easy',
  '中等': 'medium', 'medium': 'medium',
  '困难': 'hard', 'hard': 'hard',
}

const RECURRENCE_MAP: Record<string, string> = {
  '每天': 'daily', 'daily': 'daily',
  '每周': 'weekly', 'weekly': 'weekly',
  '每月': 'monthly', 'monthly': 'monthly',
}

// (?:^|\s+)MATCH means the marker must be at start or preceded by whitespace
const TAG_RE = /(?:^|\s+)(#[\w\u4e00-\u9fff-]+)/g
const PRIORITY_RE = /(?:^|\s+)(!(?:高|中|低|[hml]|[123]))/g
const EFFORT_RE = /(?:^|\s+)(@(?:简单|中等|困难|easy|medium|hard))/g
const RECURRENCE_RE = /(?:^|\s+)(\*(?:每天|每周|每月|daily|weekly|monthly))/g
const DATE_WORD_RE = /(?:^|\s+)(今天|明天|后天|[下]?周[一二三四五六日]|\d{4}-\d{2}-\d{2}|\d{1,2}\/\d{1,2}|\d{1,2}月\d{1,2}日)/g

function fmt(d: Date): string {
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

const DAY_MAP: Record<string, number> = {
  '一': 1, '二': 2, '三': 3, '四': 4, '五': 5, '六': 6, '日': 0,
}

function resolveDateWord(token: string): string | null {
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  if (token === '今天') return fmt(today)
  if (token === '明天') return fmt(new Date(today.getTime() + 86400000))
  if (token === '后天') return fmt(new Date(today.getTime() + 172800000))

  // 下周一 ~ 下周日
  if (token.startsWith('下') && token.length === 3) {
    const dayChar = token[2]
    const targetDay = DAY_MAP[dayChar]
    const currentDay = today.getDay()
    let diff = targetDay - currentDay
    if (diff <= 0) diff += 7
    diff += 7 // always next week
    return fmt(new Date(today.getTime() + diff * 86400000))
  }

  // 周一 ~ 周日
  if (token.startsWith('周') && token.length === 2) {
    const dayChar = token[1]
    const targetDay = DAY_MAP[dayChar]
    const currentDay = today.getDay()
    let diff = targetDay - currentDay
    if (diff <= 0) diff += 7
    return fmt(new Date(today.getTime() + diff * 86400000))
  }

  // YYYY-MM-DD
  if (/^\d{4}-\d{2}-\d{2}$/.test(token)) return token

  // M/D
  if (/^\d{1,2}\/\d{1,2}$/.test(token)) {
    const [m, d] = token.split('/').map(Number)
    return `${today.getFullYear()}-${String(m).padStart(2, '0')}-${String(d).padStart(2, '0')}`
  }

  // M月D日
  if (/^\d{1,2}月\d{1,2}日$/.test(token)) {
    const m = parseInt(token.split('月')[0])
    const d = parseInt(token.split('月')[1].replace('日', ''))
    return `${today.getFullYear()}-${String(m).padStart(2, '0')}-${String(d).padStart(2, '0')}`
  }

  return null
}

/**
 * Parse natural language from the todo title input.
 * Extracts tags (#tag), priority (!高), effort (@简单),
 * recurrence (*每天), and date keywords.
 */
export function parseTodoInput(input: string): ParsedTodoInput {
  let text = input.trim()
  const tags: string[] = []
  let priority: string | null = null
  let effort: string | null = null
  let recurrence: string | null = null
  let dueDate: string | null = null

  // 1. Extract tags (multiple allowed)
  text = text.replace(TAG_RE, (_, m: string) => {
    tags.push(m.slice(1)) // remove # prefix
    return ''
  })

  // 2. Extract priority (last wins)
  text = text.replace(PRIORITY_RE, (_, m: string) => {
    const key = m.slice(1) // remove ! prefix
    priority = PRIORITY_MAP[key] || null
    return ''
  })

  // 3. Extract effort (last wins)
  text = text.replace(EFFORT_RE, (_, m: string) => {
    const key = m.slice(1) // remove @ prefix
    effort = EFFORT_MAP[key] || null
    return ''
  })

  // 4. Extract recurrence (last wins)
  text = text.replace(RECURRENCE_RE, (_, m: string) => {
    const key = m.slice(1) // remove * prefix
    recurrence = RECURRENCE_MAP[key] || null
    return ''
  })

  // 5. Extract date keywords (last wins)
  text = text.replace(DATE_WORD_RE, (_, token: string) => {
    const d = resolveDateWord(token)
    if (d) dueDate = d
    return ''
  })

  return {
    title: text.replace(/\s+/g, ' ').trim(),
    dueDate,
    priority,
    tags: tags.length > 0 ? tags.join(',') : null,
    effort,
    recurrence,
  }
}
