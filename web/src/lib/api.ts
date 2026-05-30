// Typed client for the Go backend. Every network call goes through request(),
// which attaches the bearer token and normalizes error handling.

const API_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080'
const TOKEN_KEY = 'lua-academy:token'

export interface User {
  id: string
  username: string
  email: string
  role: 'user' | 'admin'
  createdAt: string
}

export interface AuthResult {
  user: User
  token: string
}

// ---- Course catalogue & content tree (mirrors the Go models) ----
export interface ApiTest {
  id: string
  name: string
  code: string
  position: number
}

export interface ApiExercise {
  id: string
  chapterId: string
  title: string
  difficulty: 'facile' | 'moyen' | 'difficile'
  statement: string
  starter: string
  solution?: string
  hints?: string[]
  position: number
  tests?: ApiTest[]
}

export interface ApiLesson {
  id: string
  chapterId: string
  title: string
  content: string
  position: number
}

export interface ApiChapter {
  id: string
  courseId: string
  title: string
  summary: string
  position: number
  lessons?: ApiLesson[]
  exercises?: ApiExercise[]
}

export interface ApiCourse {
  id: string
  slug: string
  title: string
  summary: string
  priceCents: number
  currency: string
  published: boolean
  position: number
  createdAt: string
  chapters?: ApiChapter[]
}

export interface ApiProgress {
  id: string
  userId: string
  exerciseId: string
  code: string
  completed: boolean
  updatedAt: string
}

export interface ApiNotification {
  id: string
  userId: string
  title: string
  body: string
  link: string
  read: boolean
  createdAt: string
}

/** Error carrying the HTTP status so callers can branch on it. */
export class ApiError extends Error {
  constructor(public status: number, message: string) {
    super(message)
    this.name = 'ApiError'
  }
}

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY)
}

export function setToken(token: string | null): void {
  if (token) localStorage.setItem(TOKEN_KEY, token)
  else localStorage.removeItem(TOKEN_KEY)
}

/** Core fetch wrapper: JSON in, JSON out, bearer auth, uniform errors. */
async function request<T>(method: string, path: string, body?: unknown): Promise<T> {
  const res = await fetch(API_URL + path, {
    method,
    headers: buildHeaders(body !== undefined),
    body: body !== undefined ? JSON.stringify(body) : undefined,
  })
  return parse<T>(res)
}

function buildHeaders(hasBody: boolean): HeadersInit {
  const headers: Record<string, string> = {}
  if (hasBody) headers['Content-Type'] = 'application/json'
  const token = getToken()
  if (token) headers['Authorization'] = `Bearer ${token}`
  return headers
}

/** Parse the response, throwing ApiError with the server message on failure. */
async function parse<T>(res: Response): Promise<T> {
  const text = await res.text()
  const data = text ? JSON.parse(text) : null
  if (!res.ok) {
    throw new ApiError(res.status, data?.error ?? `Erreur ${res.status}`)
  }
  return data as T
}

export const api = {
  register: (username: string, email: string, password: string) =>
    request<AuthResult>('POST', '/api/auth/register', { username, email, password }),
  login: (identifier: string, password: string) =>
    request<AuthResult>('POST', '/api/auth/login', { identifier, password }),
  me: () => request<User>('GET', '/api/auth/me'),

  courses: {
    list: () => request<ApiCourse[]>('GET', '/api/courses'),
    tree: (slug: string) => request<ApiCourse>('GET', `/api/courses/${slug}`),
  },

  progress: {
    list: () => request<ApiProgress[]>('GET', '/api/progress'),
    save: (exerciseId: string, code: string, completed: boolean) =>
      request<ApiProgress>('PUT', `/api/progress/${exerciseId}`, { code, completed }),
  },

  enrollments: {
    mine: () => request<string[]>('GET', '/api/enrollments'),
  },

  account: {
    courses: () => request<ApiCourse[]>('GET', '/api/me/courses'),
    changePassword: (currentPassword: string, newPassword: string) =>
      request<null>('POST', '/api/me/password', { currentPassword, newPassword }),
  },

  notifications: {
    list: () => request<{ notifications: ApiNotification[]; unread: number }>('GET', '/api/notifications'),
    markRead: (id: string) => request<null>('POST', `/api/notifications/${id}/read`),
    markAllRead: () => request<null>('POST', '/api/notifications/read-all'),
  },

  payments: {
    checkout: (courseId: string) =>
      request<{ url: string }>('POST', '/api/payments/checkout', { courseId }),
  },

  admin: {
    users: () => request<User[]>('GET', '/api/admin/users'),
    courses: () => request<ApiCourse[]>('GET', '/api/admin/courses'),
    userProgress: (userId: string) => request<ApiProgress[]>('GET', `/api/admin/users/${userId}/progress`),
    userCourses: (userId: string) => request<ApiCourse[]>('GET', `/api/admin/users/${userId}/courses`),
    grant: (userId: string, courseId: string) =>
      request<null>('POST', '/api/admin/enrollments', { userId, courseId }),
    createCourse: (input: NewCourse) => request<CreatedId>('POST', '/api/admin/courses', input),
    createChapter: (courseId: string, input: NewChapter) =>
      request<CreatedId>('POST', `/api/admin/courses/${courseId}/chapters`, input),
    createLesson: (chapterId: string, input: NewLesson) =>
      request<CreatedId>('POST', `/api/admin/chapters/${chapterId}/lessons`, input),
    createExercise: (chapterId: string, input: NewExercise) =>
      request<CreatedId>('POST', `/api/admin/chapters/${chapterId}/exercises`, input),
    createTest: (exerciseId: string, input: NewTest) =>
      request<CreatedId>('POST', `/api/admin/exercises/${exerciseId}/tests`, input),

    updateCourse: (courseId: string, input: NewCourse) => request<null>('PUT', `/api/admin/courses/${courseId}`, input),
    deleteCourse: (courseId: string) => request<null>('DELETE', `/api/admin/courses/${courseId}`),
    updateChapter: (chapterId: string, input: NewChapter) => request<null>('PUT', `/api/admin/chapters/${chapterId}`, input),
    deleteChapter: (chapterId: string) => request<null>('DELETE', `/api/admin/chapters/${chapterId}`),
    updateLesson: (lessonId: string, input: NewLesson) => request<null>('PUT', `/api/admin/lessons/${lessonId}`, input),
    deleteLesson: (lessonId: string) => request<null>('DELETE', `/api/admin/lessons/${lessonId}`),
    updateExercise: (exerciseId: string, input: NewExercise) => request<null>('PUT', `/api/admin/exercises/${exerciseId}`, input),
    deleteExercise: (exerciseId: string) => request<null>('DELETE', `/api/admin/exercises/${exerciseId}`),
    updateTest: (testId: string, input: NewTest) => request<null>('PUT', `/api/admin/tests/${testId}`, input),
    deleteTest: (testId: string) => request<null>('DELETE', `/api/admin/tests/${testId}`),
  },
}

export interface CreatedId {
  id: string
}

export interface NewCourse {
  slug: string
  title: string
  summary: string
  priceCents: number
  currency: string
  published: boolean
  position: number
}

export interface NewChapter {
  title: string
  summary: string
  position: number
}

export interface NewLesson {
  title: string
  content: string
  position: number
}

export interface NewExercise {
  title: string
  difficulty: 'facile' | 'moyen' | 'difficile'
  statement: string
  starter: string
  solution: string
  hints: string[]
  position: number
}

export interface NewTest {
  name: string
  code: string
  position: number
}
