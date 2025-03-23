<template>
  <div>
    <div class="header">
      <h2>Управление пользователями</h2>
      <button class="button button-primary" @click="showAddUserModal = true">
        Добавить пользователя
      </button>
    </div>

    <!-- Список пользователей -->
    <div class="list">
      <div v-for="user in users" :key="user.id" class="list-item">
        <div class="list-item-content">
          <div>
            <h3>{{ user.username }}</h3>
            <p class="text-muted">{{ user.public_key }}</p>
          </div>
          <div class="list-item-actions">
            <button class="button" @click="viewUserServers(user)">
              Просмотр серверов
            </button>
            <button class="button button-danger" @click="deleteUser(user.id)">
              Удалить
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Модальное окно добавления пользователя -->
    <div v-if="showAddUserModal" class="modal">
      <div class="modal-content">
        <div class="modal-header">
          <h3 class="modal-title">Добавить пользователя</h3>
          <button class="modal-close" @click="showAddUserModal = false">&times;</button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="addUser">
            <div class="form-group">
              <label class="form-label" for="username">Имя пользователя</label>
              <input
                type="text"
                id="username"
                v-model="newUser.username"
                class="form-input"
                required
              />
            </div>
            <div class="form-group">
              <label class="form-label" for="publicKey">Публичный ключ</label>
              <textarea
                id="publicKey"
                v-model="newUser.public_key"
                class="form-input"
                rows="3"
                required
              ></textarea>
            </div>
            <div class="modal-footer">
              <button type="button" class="button" @click="showAddUserModal = false">
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

interface User {
  id: number
  username: string
  public_key: string
}

const users = ref<User[]>([])
const showAddUserModal = ref(false)
const newUser = ref({
  username: '',
  public_key: ''
})

// Загрузка пользователей
const fetchUsers = async () => {
  try {
    const response = await fetch('http://localhost:8080/api/users')
    users.value = await response.json()
  } catch (error) {
    console.error('Ошибка при загрузке пользователей:', error)
  }
}

// Добавление пользователя
const addUser = async () => {
  try {
    const response = await fetch('http://localhost:8080/api/users', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(newUser.value)
    })
    if (response.ok) {
      showAddUserModal.value = false
      newUser.value = { username: '', public_key: '' }
      await fetchUsers()
    }
  } catch (error) {
    console.error('Ошибка при добавлении пользователя:', error)
  }
}

// Удаление пользователя
const deleteUser = async (id: number) => {
  if (!confirm('Вы уверены, что хотите удалить этого пользователя?')) return
  
  try {
    const response = await fetch(`http://localhost:8080/api/users/${id}`, {
      method: 'DELETE'
    })
    if (response.ok) {
      await fetchUsers()
    }
  } catch (error) {
    console.error('Ошибка при удалении пользователя:', error)
  }
}

// Просмотр серверов пользователя
const viewUserServers = (user: User) => {
  // TODO: Реализовать просмотр серверов пользователя
}

onMounted(() => {
  fetchUsers()
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

.text-muted {
  color: #666;
  font-size: 0.875rem;
  margin: 0.25rem 0 0;
}

.list-item-actions {
  display: flex;
  gap: 0.5rem;
}
</style> 