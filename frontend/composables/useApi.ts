import { useAuthStore } from '~/stores/auth'

export const useApi = () => {
  const config = useRuntimeConfig()
  const authStore = useAuthStore()

  const request = async (path: string, options: any = {}) => {
    const url = `${config.public.apiBase}${path}`
    
    // Ensure auth token is rehydrated on client
    authStore.rehydrate()

    // Add auth header if token exists
    const isFormData = options.body instanceof FormData
    const headers: any = {
      ...(isFormData ? {} : { 'Content-Type': 'application/json' }),
      ...options.headers,
    }
    if (authStore.token) {
      headers['Authorization'] = `Bearer ${authStore.token}`
    }

    try {
      const response = await fetch(url, {
        ...options,
        headers,
      })

      if (response.status === 401) {
        authStore.logout()
        throw new Error('Unauthorized')
      }

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.error || 'Request failed')
      }

      // Handle 204 No Content
      if (response.status === 204) {
        return null
      }

      const body = await response.json()
      return body.data
    } catch (error: any) {
      console.error(`API Error on ${path}:`, error)
      throw error
    }
  }

  return {
    get: (path: string, options: any = {}) => request(path, { ...options, method: 'GET' }),
    post: (path: string, body?: any, options: any = {}) => {
      const isFormData = body instanceof FormData
      return request(path, { ...options, method: 'POST', body: isFormData ? body : (body ? JSON.stringify(body) : undefined) })
    },
    put: (path: string, body?: any, options: any = {}) => {
      const isFormData = body instanceof FormData
      return request(path, { ...options, method: 'PUT', body: isFormData ? body : (body ? JSON.stringify(body) : undefined) })
    },
    delete: (path: string, options: any = {}) => request(path, { ...options, method: 'DELETE' }),
  }
}
