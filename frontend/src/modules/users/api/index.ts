import { fetchApi } from '@/shared/utils.ts'

export interface UserDto {
  id: number
  username: string
  public_key: string
}

export type UserPayload = Pick<UserDto, 'username' | 'public_key'>

export const usersFetch = async (): Promise<UserDto[]> => {
  const response = await fetchApi('/api/users')
  return (await response.json()) as UserDto[]
}

export const userCreate = async (payload: UserPayload): Promise<UserDto> => {
  const response = await fetchApi('/api/users', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload),
  })
  return (await response.json()) as UserDto
}

export const userDelete = async (id: number): Promise<boolean> => {
  await fetchApi(`/api/users/${id}`, {
    method: 'DELETE',
  })
  return true
}

export const userUpdate = async (id: number, payload: UserPayload): Promise<UserDto> => {
  const response = await fetchApi(`/api/users/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(payload),
  })
  return (await response.json()) as UserDto
}
