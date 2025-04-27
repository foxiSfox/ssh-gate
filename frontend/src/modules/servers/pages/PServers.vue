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
            <button class="button button-danger" @click="onServerDelete(server.id)">
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
          <form @submit.prevent="onServerCreate">
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
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { serversFetch, serverCreate, serverDelete } from '../api'

interface Server {
  id: number
  ip: string
}

const showAddServerModal = ref(false)
const newServer = ref({
  ip: ''
})

const { data: servers } = useQuery({
  queryKey: ['servers'],
  queryFn: serversFetch,
})

const queryClient = useQueryClient()
const { mutate: mutateServerCreate } = useMutation({
  mutationFn: serverCreate,
  onSuccess: () => {
    showAddServerModal.value = false
    queryClient.invalidateQueries({ queryKey: ['servers'] })
  },
});

const onServerCreate = () => {
  mutateServerCreate(JSON.stringify(newServer.value))
}

const { mutate: mutateServerDelete } = useMutation({
  mutationFn: serverDelete,
  onSuccess: () => {
    queryClient.invalidateQueries({ queryKey: ['servers'] })
  }
})

const onServerDelete = (id: number) => {
  mutateServerDelete(id);
}

// Просмотр пользователей сервера
const viewServerUsers = (server: Server) => {
  // TODO: Реализовать просмотр пользователей сервера
}

onMounted(() => {
  serversFetch()
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
