import { fetchApi } from '@/shared/utils.ts'
export { usersFetch } from '@/modules/users/api'
export { serversFetch } from '@/modules/servers/api'

import type { ServerDto } from '@/modules/servers/api'

export const userServersFetch = async (userId: number): Promise<ServerDto[]> => {
  const response = await fetchApi(`/api/users/${userId}/servers`)
  return (await response.json()) as ServerDto[]
}

export interface UserServerPayload {
  userId: number
  serverId: number
}

export const assignServer = async ({ userId, serverId }: UserServerPayload): Promise<boolean> => {
  await fetchApi(`/api/users/${userId}/servers/${serverId}`, {
    method: 'POST',
  })
  return true
}

export const removeServerFromUser = async ({ userId, serverId }: UserServerPayload): Promise<boolean> => {
  await fetchApi(`/api/users/${userId}/servers/${serverId}`, {
    method: 'DELETE',
  })
  return true
}
