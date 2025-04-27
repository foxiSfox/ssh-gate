import { fetchApi } from "@/shared/utils.ts";

export const usersFetch = async () => {
  const response = await fetchApi('/api/users')
  return response.json();
}

export const serversFetch = async () => {
  const response = await fetchApi('/api/servers')
  return response.json();
}

export const userServersFetch = async (userId: number) => {
  const response = await fetchApi(`/api/users/${userId}/servers`)
  return response.json();
}

export const assignServer = async (payload: any) => {
  const { userId, serverId } = payload;
  await fetchApi(
    `/api/users/${userId}/servers/${serverId}`,
    {
      method: 'POST'
    }
  )
  return true;
}

export const removeServerFromUser = async (payload: any) => {
  const { userId, serverId } = payload;
  await fetchApi(
    `/api/users/${userId}/servers/${serverId}`,
    {
      method: 'DELETE'
    }
  )
  return true;
}
