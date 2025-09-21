export interface ApiErrorPayload {
  message: string
  [key: string]: unknown
}

export class ApiError extends Error {
  status: number
  payload?: unknown

  constructor(message: string, status: number, payload?: unknown) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.payload = payload
  }
}

const extractMessage = (body: unknown, fallback: string): { message: string; payload?: unknown } => {
  if (typeof body === 'string') {
    const trimmed = body.trim()
    return { message: trimmed.length > 0 ? trimmed : fallback, payload: trimmed }
  }

  if (body && typeof body === 'object') {
    const possibleMessage = (body as Record<string, unknown>).message
    if (typeof possibleMessage === 'string' && possibleMessage.trim().length > 0) {
      return { message: possibleMessage.trim(), payload: body }
    }
    return { message: fallback, payload: body }
  }

  return { message: fallback, payload: body }
}

const createApiError = async (response: Response): Promise<ApiError> => {
  const fallbackMessage = `Request failed with status ${response.status}`
  try {
    const cloned = response.clone()
    const contentType = cloned.headers.get('content-type') || ''

    if (contentType.includes('application/json')) {
      const json = await cloned.json()
      const { message, payload } = extractMessage(json, fallbackMessage)
      return new ApiError(message, response.status, payload)
    }

    const text = await cloned.text()
    const { message, payload } = extractMessage(text, fallbackMessage)
    return new ApiError(message, response.status, payload)
  } catch (error) {
    return new ApiError(fallbackMessage, response.status)
  }
}

export const fetchApi = async (url: string, options?: RequestInit) => {
  const apiUrl = import.meta.env.VITE_API_URL || ''
  const response = await fetch(`${apiUrl}${url}`, options)

  if (!response.ok) {
    throw await createApiError(response)
  }

  return response
}

export const resolveErrorMessage = (error: unknown, fallback = 'Unexpected error'): string => {
  if (error instanceof ApiError) {
    return error.message
  }

  if (error instanceof Error) {
    return error.message || fallback
  }

  if (typeof error === 'string') {
    return error || fallback
  }

  return fallback
}
