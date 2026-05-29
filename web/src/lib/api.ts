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
}
