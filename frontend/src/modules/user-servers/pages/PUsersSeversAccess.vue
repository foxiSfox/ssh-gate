<template>
  <div>
    <div class="header">
      <h2>Управление доступом к серверам</h2>
    </div>

    <!-- Список пользователей и их серверов -->
    <div class="list">
      <div v-for="user in users" :key="user.id" class="list-item">
        <div class="list-item-content">
          <div>
            <h3>{{ user.username }}</h3>
            <div class="user-servers">
              <h4>Доступные серверы:</h4>
              <ul>
                <li v-for="server in userServers[user.id]" :key="server.id" class="server-item">
                  <span>{{ server.ip }}</span>
                  <button class="button button-danger" @click="removeServerFromUser(user.id, server.id)">
                    Удалить доступ
                  </button>
                </li>
              </ul>
            </div>
          </div>
          <button class="button button-primary" @click="showAssignServerModal = true; selectedUser = user">
            Добавить сервер
          </button>
        </div>
      </div>
    </div>

    <!-- Модальное окно назначения сервера -->
    <div v-if="showAssignServerModal" class="modal">
      <div class="modal-content">
        <div class="modal-header">
          <h3 class="modal-title">Добавить сервер для пользователя {{ selectedUser?.username }}</h3>
          <button class="modal-close" @click="showAssignServerModal = false">&times;</button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="assignServer">
            <div class="form-group">
              <label class="form-label" for="server">Выберите сервер</label>
              <select
                id="server"
                v-model="selectedServer"
                class="form-input"
                required
              >
                <option value="">Выберите сервер</option>
                <option v-for="server in availableServers" :key="server.id" :value="server.id">
                  {{ server.ip }}
                </option>
              </select>
            </div>
            <div class="modal-footer">
              <button type="button" class="button" @click="showAssignServerModal = false">
                Отмена
              </button>
              <button type="submit" class="button button-primary">
                Добавить
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { fetchApi } from "@/shared/utils.ts";
interface User {
  id: number
  username: string
  public_key: string
}

interface Server {
  id: number
  ip: string
}

const users = ref<User[]>([])
const servers = ref<Server[]>([])
const userServers = ref<Record<number, Server[]>>({})
const showAssignServerModal = ref(false)
const selectedUser = ref<User | null>(null)
const selectedServer = ref<number | null>(null)

// Загрузка пользователей
const fetchUsers = async () => {
  try {
    const response = await fetchApi('/api/users')
    users.value = await response.json()
    await Promise.all(users.value.map(user => fetchUserServers(user.id)))
  } catch (error) {
    console.error('Ошибка при загрузке пользователей:', error)
  }
}

// Загрузка серверов
const fetchServers = async () => {
  try {
    const response = await fetchApi('/api/servers')
    servers.value = await response.json()
  } catch (error) {
    console.error('Ошибка при загрузке серверов:', error)
  }
}

// Загрузка серверов пользователя
const fetchUserServers = async (userId: number) => {
  try {
    const response = await fetchApi(`/api/users/${userId}/servers`)
    userServers.value[userId] = await response.json()
  } catch (error) {
    console.error('Ошибка при загрузке серверов пользователя:', error)
  }
}

// Вычисляемое свойство для доступных серверов
const availableServers = computed(() => {
  if (!selectedUser.value) return []
  const userServerIds = new Set(userServers.value[selectedUser.value.id]?.map(s => s.id) || [])
  return servers.value.filter(server => !userServerIds.has(server.id))
})

// Назначение сервера пользователю
const assignServer = async () => {
  if (!selectedUser.value || !selectedServer.value) return

  try {
    const response = await fetchApi(
      `/api/users/${selectedUser.value.id}/servers/${selectedServer.value}`,
      {
        method: 'POST'
      }
    )
    if (response.ok) {
      showAssignServerModal.value = false
      selectedServer.value = null
      await fetchUserServers(selectedUser.value.id)
    }
  } catch (error) {
    console.error('Ошибка при назначении сервера:', error)
  }
}

// Удаление сервера у пользователя
const removeServerFromUser = async (userId: number, serverId: number) => {
  if (!confirm('Вы уверены, что хотите удалить доступ к этому серверу?')) return

  try {
    const response = await fetchApi(
      `/api/users/${userId}/servers/${serverId}`,
      {
        method: 'DELETE'
      }
    )
    if (response.ok) {
      await fetchUserServers(userId)
    }
  } catch (error) {
    console.error('Ошибка при удалении сервера:', error)
  }
}

onMounted(() => {
  fetchUsers()
  fetchServers()
})
</script>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.header h2 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
}

.list-item-content {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.list-item-content h3 {
  margin: 0;
  font-size: 1rem;
  font-weight: 500;
}

.user-servers {
  margin-top: 0.5rem;
}

.user-servers h4 {
  margin: 0;
  font-size: 0.875rem;
  font-weight: 500;
  color: #666;
}

.user-servers ul {
  margin: 0.5rem 0 0;
  padding: 0;
  list-style: none;
}

.server-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 0.5rem;
  padding: 0.5rem;
  background-color: #f8f9fa;
  border-radius: 4px;
}

.server-item span {
  font-size: 0.875rem;
}
</style>
