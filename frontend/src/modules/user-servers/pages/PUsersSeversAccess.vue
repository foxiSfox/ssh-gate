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
                  <button class="button button-danger" @click="onRemoveServerFromUser(user.id, server.id)">
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
          <form @submit.prevent="onAssignServer">
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
import { ref, computed } from 'vue'
import { useQuery, useQueries, useMutation, useQueryClient } from '@tanstack/vue-query';
import { usersFetch, serversFetch, userServersFetch, assignServer, removeServerFromUser } from "../api";

interface User {
  id: number
  username: string
  public_key: string
}

interface Server {
  id: number
  ip: string
}

const showAssignServerModal = ref(false)
const selectedUser = ref<User | null>(null)
const selectedServer = defineModel()

const { data: users } = useQuery({
  queryKey: ['users'],
  queryFn: usersFetch,
})

const { data: servers } = useQuery({
  queryKey: ['servers'],
  queryFn: serversFetch,
})

const userServersQueries = useQueries({
  queries: computed(() => {
    return users.value
      ? users?.value?.map((user: User) => ({
          queryKey: ['user-servers', user.id],
          queryFn: () => userServersFetch(user.id),
        }))
      : []
  })
})

const userServers = computed(() => {
  const result: Record<number, any> = {}
  userServersQueries.value.forEach((queryResult, index) => {
    const user = users.value?.[index]
    if (user) {
      result[user.id] = queryResult.data
    }
  })
  return result
})

const queryClient = useQueryClient()
const { mutate: mutateAssignServer } = useMutation({
  mutationFn: assignServer,
  onSuccess: (_, variables) => {
    showAssignServerModal.value = false;
    selectedServer.value = null;
    queryClient.invalidateQueries({ queryKey: ['user-servers', variables.userId] });
  },
})

const onAssignServer = async () => {
  if (!selectedUser.value || !selectedServer.value) {
     return
  }

  mutateAssignServer({
    userId: selectedUser.value.id,
    serverId: selectedServer.value
  })
}

const { mutate: mutateRemoveServerFromUser} = useMutation({
  mutationFn: removeServerFromUser,
  onSuccess: (_, variables) => {
    queryClient.invalidateQueries({ queryKey: ['user-servers', variables.userId] })
  }
})

const onRemoveServerFromUser = async (userId: number, serverId: number) => {
  if (!confirm('Вы уверены, что хотите удалить доступ к этому серверу?')) {
    return
  }
  mutateRemoveServerFromUser({ userId, serverId })
}

/**
 * Получение доступных серверов для пользователя
 */
const availableServers = computed(() => {
  if (!selectedUser.value) {
    return []
  }
  const userServerIds = new Set<number>(
    userServers.value[selectedUser.value.id]?.map((s: Server) => s.id) || []
  )

  return servers.value.filter((server: Server) => !userServerIds.has(server.id))
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
