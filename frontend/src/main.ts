import './assets/main.css'
import {
  MutationCache,
  QueryCache,
  QueryClient,
  VueQueryPlugin,
} from '@tanstack/vue-query'
import router from './router'
import { createApp } from 'vue'
import App from './App.vue'
import { notifyError } from '@/shared/notifications'
import { resolveErrorMessage } from '@/shared/utils'

const handleError = (error: unknown) => {
  notifyError({ message: resolveErrorMessage(error) })
}

const queryClient = new QueryClient({
  queryCache: new QueryCache({
    onError: (error) => handleError(error),
  }),
  mutationCache: new MutationCache({
    onError: (error) => handleError(error),
  }),
})

const app = createApp(App)

app.use(VueQueryPlugin, { queryClient })
app.use(router)

app.config.errorHandler = (err) => {
  handleError(err)
}

app.mount('#app')
