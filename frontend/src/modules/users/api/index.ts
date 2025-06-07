import { fetchApi } from "@/shared/utils.ts";

export const usersFetch = async () => {
  const response = await fetchApi('/api/users')
  return response.json();
}

export const userCreate = async (user: any) => {
  const response = await fetchApi('/api/users', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: user
  })
  return response.json();
}

export const userDelete = async (id: number) => {
  await fetchApi(`/api/users/${id}`, {
    method: 'DELETE'
  })
  return true;
}

export const userUpdate = async (id: number, user: any) => {
  const response = await fetchApi(`/api/users/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(user)
  })
  return await response.json()
}

