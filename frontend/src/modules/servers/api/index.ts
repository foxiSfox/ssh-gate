import { fetchApi } from '@/shared/utils.ts'

export interface ServerDto {
  id: number
  ip: string
  port: number
  login: string
  password: string
}

export type ServerPayload = Omit<ServerDto, 'id'>

export const serversFetch = async (): Promise<ServerDto[]> => {
  const response = await fetchApi('/api/servers')
  return (await response.json()) as ServerDto[]
}

export const serverCreate = async (payload: ServerPayload): Promise<ServerDto> => {
  const response = await fetchApi('/api/servers', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload),
  })

  return (await response.json()) as ServerDto
}

export const serverDelete = async (id: number): Promise<boolean> => {
  await fetchApi(`/api/servers/${id}`, {
    method: 'DELETE',
  })
  return true
}

export const serverUpdate = async (id: number, payload: ServerPayload): Promise<ServerDto> => {
  const response = await fetchApi(`/api/servers/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload),
  })
  return (await response.json()) as ServerDto
}
