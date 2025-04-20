
export const fetchApi = async (url: any, options?: any) => {
  const apiUrl = import.meta.env.VITE_API_URL || '';
  return await fetch(`${apiUrl}${url}`, options);
}
