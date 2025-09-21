export type NotificationType = 'error' | 'info' | 'success'

export interface NotificationPayload {
  message: string
  type?: NotificationType
  duration?: number
  id?: number
}

export interface NotificationEvent extends NotificationPayload {
  id: number
  type: NotificationType
}

type Listener = (notification: NotificationEvent) => void

const listeners = new Set<Listener>()
let counter = 0

export const subscribe = (listener: Listener) => {
  listeners.add(listener)
  return () => listeners.delete(listener)
}

const emit = (notification: NotificationEvent) => {
  listeners.forEach((listener) => listener(notification))
}

export const notifyError = (payload: NotificationPayload) => {
  const id = payload.id ?? ++counter
  emit({
    id,
    type: 'error',
    duration: payload.duration,
    message: payload.message,
  })
}

export const notify = (payload: NotificationPayload) => {
  const id = payload.id ?? ++counter
  emit({
    id,
    type: payload.type ?? 'info',
    duration: payload.duration,
    message: payload.message,
  })
}
