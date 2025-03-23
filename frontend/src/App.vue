<template>
  <div>
    <nav class="nav">
      <div class="nav-container">
        <h1 class="nav-title">SSH Gate</h1>
        <div class="nav-links">
          <a
            v-for="tab in tabs"
            :key="tab.name"
            :href="tab.href"
            @click.prevent="currentTab = tab.name"
            :class="['nav-link', { active: currentTab === tab.name }]"
          >
            {{ tab.label }}
          </a>
        </div>
      </div>
    </nav>

    <main class="container">
      <UserManagement v-if="currentTab === 'users'" />
      <ServerManagement v-if="currentTab === 'servers'" />
      <UserServerAccess v-if="currentTab === 'access'" />
    </main>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import UserManagement from './components/UserManagement.vue'
import ServerManagement from './components/ServerManagement.vue'
import UserServerAccess from './components/UserServerAccess.vue'

const tabs = [
  { name: 'users', label: 'Пользователи', href: '#' },
  { name: 'servers', label: 'Серверы', href: '#' },
  { name: 'access', label: 'Управление доступом', href: '#' }
]

const currentTab = ref('users')
</script>
