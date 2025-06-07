import { fetchApi } from "@/shared/utils.ts";

export const serversFetch = async () => {
  const response = await fetchApi('/api/servers')
  return await response.json()
}

export const serverCreate = async (server: any) => {
  const response = await fetchApi('/api/servers', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: server
  })

  return await response.json();
}

export const serverDelete = async (id: number) => {
  await fetchApi(`/api/servers/${id}`, {
    method: 'DELETE'
  })
  return true;
}

export const serverUpdate = async (id: number, server: any) => {
  const response = await fetchApi(`/api/servers/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(server)
  })
  return await response.json()
}
