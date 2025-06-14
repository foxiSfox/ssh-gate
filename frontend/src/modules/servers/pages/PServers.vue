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
            <h3>{{ server.ip }}:{{ server.port }}</h3>
          </div>
          <div class="list-item-actions">
            <button class="button" @click="viewServerUsers(server)">
              Просмотр пользователей
            </button>
            <button class="button" @click="editServer(server)">Редактировать</button>
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
            <div class="form-group">
              <label class="form-label" for="port">Порт</label>
              <input
                type="number"
                id="port"
                v-model.number="newServer.port"
                class="form-input"
                min="1"
                max="65535"
              />
            </div>
            <div class="form-group">
              <label class="form-label" for="login">Логин</label>
              <input
                type="text"
                id="login"
                v-model="newServer.login"
                class="form-input"
                required
              />
            </div>
            <div class="form-group">
              <label class="form-label" for="password">Пароль</label>
              <input
                type="password"
                id="password"
                v-model="newServer.password"
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

    <!-- Модальное окно редактирования сервера -->
    <div v-if="showEditServerModal" class="modal">
      <div class="modal-content">
        <div class="modal-header">
          <h3 class="modal-title">Редактировать сервер</h3>
          <button class="modal-close" @click="showEditServerModal = false">&times;</button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="onServerUpdate">
            <div class="form-group">
              <label class="form-label" for="edit-ip">IP адрес</label>
              <input
                type="text"
                id="edit-ip"
                v-model="editedServer.ip"
                class="form-input"
                required
              />
            </div>
            <div class="form-group">
              <label class="form-label" for="edit-port">Порт</label>
              <input
                type="number"
                id="edit-port"
                v-model.number="editedServer.port"
                class="form-input"
                min="1"
                max="65535"
              />
            </div>
            <div class="form-group">
              <label class="form-label" for="edit-login">Логин</label>
              <input
                type="text"
                id="edit-login"
                v-model="editedServer.login"
                class="form-input"
                required
              />
            </div>
            <div class="form-group">
              <label class="form-label" for="edit-password">Пароль</label>
              <input
                type="password"
                id="edit-password"
                v-model="editedServer.password"
                class="form-input"
                required
              />
            </div>
            <div class="modal-footer">
              <button type="button" class="button" @click="showEditServerModal = false">
                Отмена
              </button>
              <button type="submit" class="button button-primary">
                Сохранить
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
import { serversFetch, serverCreate, serverDelete, serverUpdate } from '../api'

interface Server {
  id: number
  ip: string
  port: number
  login: string
  password: string
}

const showAddServerModal = ref(false)
const showEditServerModal = ref(false)
const newServer = ref({
  ip: '',
  port: 22,
  login: '',
  password: ''
})
const editedServer = ref({ id: 0, ip: '', port: 22, login: '', password: '' })

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

const { mutate: mutateServerUpdate } = useMutation({
  mutationFn: ({ id, data }: { id: number; data: any }) => serverUpdate(id, data),
  onSuccess: () => {
    showEditServerModal.value = false
    queryClient.invalidateQueries({ queryKey: ['servers'] })
  },
})

const editServer = (s: Server) => {
  editedServer.value = { ...s }
  showEditServerModal.value = true
}

const onServerUpdate = () => {
  mutateServerUpdate({ id: editedServer.value.id, data: editedServer.value })
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
