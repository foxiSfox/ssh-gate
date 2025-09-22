<template>
  <div class="global-notifications" aria-live="assertive">
    <transition-group name="global-notifications-fade">
      <div
        v-for="notification in notifications"
        :key="notification.id"
        class="notification"
        :class="`notification--${notification.type}`"
        role="alert"
      >
        <span class="notification__message">{{ notification.message }}</span>
        <button
          type="button"
          class="notification__close"
          aria-label="Dismiss notification"
          @click="dismiss(notification.id)"
        >
          ×
        </button>
      </div>
    </transition-group>
  </div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import type { NotificationEvent } from '@/shared/notifications'
import { subscribe } from '@/shared/notifications'

type NotificationItem = NotificationEvent & { timeoutId: ReturnType<typeof setTimeout> }

const notifications = ref<NotificationItem[]>([])
let unsubscribe: (() => void) | null = null

const removeNotification = (id: number) => {
  const index = notifications.value.findIndex((item) => item.id === id)
  if (index === -1) {
    return
  }

  clearTimeout(notifications.value[index].timeoutId)
  notifications.value.splice(index, 1)
}

const handleNotification = (notification: NotificationEvent) => {
  const duration = notification.duration ?? 5000
  const timeoutId = window.setTimeout(() => removeNotification(notification.id), duration)

  notifications.value.push({ ...notification, timeoutId })
}

onMounted(() => {
  unsubscribe = subscribe(handleNotification)
})

onBeforeUnmount(() => {
  if (unsubscribe) {
    unsubscribe()
    unsubscribe = null
  }

  notifications.value.forEach((item) => {
    clearTimeout(item.timeoutId)
  })
})

const dismiss = (id: number) => {
  removeNotification(id)
}
</script>

<style scoped>
.global-notifications {
  position: fixed;
  top: 1.5rem;
  right: 1.5rem;
  z-index: 1000;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.global-notifications-fade-enter-active,
.global-notifications-fade-leave-active {
  transition: opacity 0.2s ease, transform 0.2s ease;
}

.global-notifications-fade-enter-from,
.global-notifications-fade-leave-to {
  opacity: 0;
  transform: translateY(-8px);
}

.notification {
  min-width: 240px;
  max-width: 360px;
  padding: 0.75rem 1rem;
  border-radius: 6px;
  background-color: #dbeafe;
  color: #1e3a8a;
  box-shadow: 0 8px 16px rgba(15, 23, 42, 0.12);
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
}

.notification--error {
  background-color: #fee2e2;
  color: #7f1d1d;
}

.notification--info {
  background-color: #bfdbfe;
  color: #1d4ed8;
}

.notification--success {
  background-color: #dcfce7;
  color: #166534;
}

.notification__message {
  flex: 1;
  font-size: 0.875rem;
  line-height: 1.4;
}

.notification__close {
  background: transparent;
  border: none;
  color: inherit;
  cursor: pointer;
  font-size: 1.25rem;
  line-height: 1;
  padding: 0;
}

.notification__close:hover {
  opacity: 0.7;
}
</style>
