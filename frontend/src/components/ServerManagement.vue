<template>
  <div>
    <div class="header">
      <h2>Управление серверами</h2>
      <button class="button button-primary" @click="showAddServerModal = true">
        Добавить сервер
      </button>
    </div>

    <!-- Список серверов -->
    <div class="list">
      <div v-for="server in servers" :key="server.id" class="list-item">
        <div class="list-item-content">
          <div>
            <h3>IP: {{ server.ip }}</h3>
          </div>
          <div class="list-item-actions">
            <button class="button" @click="viewServerUsers(server)">
              Просмотр пользователей
            </button>
            <button class="button button-danger" @click="deleteServer(server.id)">
              Удалить
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Модальное окно добавления сервера -->
    <div v-if="showAddServerModal" class="modal">
      <div class="modal-content">
        <div class="modal-header">
          <h3 class="modal-title">Добавить сервер</h3>
          <button class="modal-close" @click="showAddServerModal = false">&times;</button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="addServer">
            <div class="form-group">
              <label class="form-label" for="ip">IP адрес</label>
              <input
                type="text"
                id="ip"
                v-model="newServer.ip"
                class="form-input"
                required
              />
            </div>
            <div class="modal-footer">
              <button type="button" class="button" @click="showAddServerModal = false">
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
import { ref, onMounted } from 'vue'

interface Server {
  id: number
  ip: string
}

const servers = ref<Server[]>([])
const showAddServerModal = ref(false)
const newServer = ref({
  ip: ''
})

// Загрузка серверов
const fetchServers = async () => {
  try {
    const response = await fetch('http://localhost:8080/api/servers')
    servers.value = await response.json()
  } catch (error) {
    console.error('Ошибка при загрузке серверов:', error)
  }
}

// Добавление сервера
const addServer = async () => {
  try {
    const response = await fetch('http://localhost:8080/api/servers', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(newServer.value)
    })
    if (response.ok) {
      showAddServerModal.value = false
      newServer.value = { ip: '' }
      await fetchServers()
    }
  } catch (error) {
    console.error('Ошибка при добавлении сервера:', error)
  }
}

// Удаление сервера
const deleteServer = async (id: number) => {
  if (!confirm('Вы уверены, что хотите удалить этот сервер?')) return
  
  try {
    const response = await fetch(`http://localhost:8080/api/servers/${id}`, {
      method: 'DELETE'
    })
    if (response.ok) {
      await fetchServers()
    }
  } catch (error) {
    console.error('Ошибка при удалении сервера:', error)
  }
}

// Просмотр пользователей сервера
const viewServerUsers = (server: Server) => {
  // TODO: Реализовать просмотр пользователей сервера
}

onMounted(() => {
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
  align-items: center;
}

.list-item-content h3 {
  margin: 0;
  font-size: 1rem;
  font-weight: 500;
}

.list-item-actions {
  display: flex;
  gap: 0.5rem;
}
</style> 