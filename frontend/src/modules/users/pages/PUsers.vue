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
            <button class="button button-danger" @click="onUserDelete(user.id)">
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
          <form @submit.prevent="onUserCreate">
            <div class="form-group">
              <label class="form-label" for="username">Имя пользователя</label>
              <input
                type="text"
                id="username"
                v-model="user.username"
                class="form-input"
                required
              />
            </div>
            <div class="form-group">
              <label class="form-label" for="publicKey">Публичный ключ</label>
              <textarea
                id="publicKey"
                v-model="user.public_key"
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
import { ref } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { usersFetch, userCreate, userDelete } from '../api';

interface User {
  id: number
  username: string
  public_key: string
}

const showAddUserModal = ref(false)
const user = ref({
  username: '',
  public_key: ''
})

const { data: users } = useQuery({
  queryKey: ['users'],
  queryFn: usersFetch,
})

const queryClient = useQueryClient()
const { mutate: mutateUserCreate } = useMutation({
  mutationFn: userCreate,
  onSuccess: () => {
    showAddUserModal.value = false
    queryClient.invalidateQueries({ queryKey: ['users'] })
  },
})

const onUserCreate = () => {
  mutateUserCreate(JSON.stringify(user.value))
}

const { mutate: mutateUserDelete } = useMutation({
  mutationFn: userDelete,
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['users'] })
  },
})

const onUserDelete = (id: number) => {
  if (!confirm('Вы уверены, что хотите удалить этого пользователя?')) {
    return
  }
  mutateUserDelete(id)
}


// Просмотр серверов пользователя
const viewUserServers = (user: User) => {
  // TODO: Реализовать просмотр серверов пользователя
}
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
  word-break: break-all;
}

.list-item-actions {
  display: flex;
  gap: 0.5rem;
}
</style>
