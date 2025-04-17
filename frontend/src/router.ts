import { createWebHistory, createRouter } from 'vue-router'
import PUsers from './modules/users/pages/PUsers.vue';
import PServers from './modules/servers/pages/PServers.vue';
import PUsersSeversAccess from "@/modules/user-servers/pages/PUsersSeversAccess.vue";

const routes: any[] = [
  { path: '/users', component: PUsers, alias: '/' },
  { path: '/servers', component: PServers },
  { path: '/access', component: PUsersSeversAccess },
]

export default createRouter({
  history: createWebHistory(),
  routes,
})

