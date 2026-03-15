import axios from 'axios'

const api = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://127.0.0.1:8080',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export const authApi = {
  login: (data: { username: string; password: string }) =>
    api.post('/auth/login', data),
  register: (data: { username: string; email: string; password: string }) =>
    api.post('/auth/register', data)
}

export const customerApi = {
  getList: (params: any) => api.get('/api/customers', { params }),
  getDetail: (id: number) => api.get(`/api/customers/${id}`),
  create: (data: any) => api.post('/api/customers', data),
  update: (id: number, data: any) => api.put(`/api/customers/${id}`, data),
  delete: (id: number) => api.delete(`/api/customers/${id}`),
  getHealth: (id: number) => api.get(`/api/customers/${id}/health`),
  calculateHealth: (id: number) => api.post(`/api/customers/${id}/calculate-health`)
}

export const contractApi = {
  getList: (params: any) => api.get('/api/contracts', { params }),
  getDetail: (id: number) => api.get(`/api/contracts/${id}`),
  create: (data: any) => api.post('/api/contracts', data),
  update: (id: number, data: any) => api.put(`/api/contracts/${id}`, data),
  delete: (id: number) => api.delete(`/api/contracts/${id}`),
  submit: (id: number) => api.post(`/api/contracts/${id}/submit`),
  approve: (id: number, data: any) => api.post(`/api/contracts/${id}/approve`, data),
  getExpiring: (params: any) => api.get('/api/contracts/expiring', { params }),
  getPerformance: (id: number) => api.get(`/api/contracts/${id}/performance`),
  createPerformance: (id: number, data: any) => api.post(`/api/contracts/${id}/performance`, data)
}

export const interactionApi = {
  getList: (params: any) => api.get('/api/interactions', { params }),
  getDetail: (id: number) => api.get(`/api/interactions/${id}`),
  create: (data: any) => api.post('/api/interactions', data),
  update: (id: number, data: any) => api.put(`/api/interactions/${id}`, data),
  delete: (id: number) => api.delete(`/api/interactions/${id}`),
  getTimeline: (params: any) => api.get('/api/interactions/timeline', { params }),
  getStats: (params: any) => api.get('/api/interactions/stats', { params })
}

export const dashboardApi = {
  getOverview: () => api.get('/api/dashboard/overview'),
  getCustomerGrowth: (params: any) => api.get('/api/dashboard/customer-growth', { params }),
  getContractDistribution: () => api.get('/api/dashboard/contract-distribution'),
  getInteractionHeatmap: (params: any) => api.get('/api/dashboard/interaction-heatmap', { params }),
  getRenewalAlerts: () => api.get('/api/dashboard/renewal-alerts'),
  getTopCustomers: () => api.get('/api/dashboard/top-customers'),
  getRecentActivities: () => api.get('/api/dashboard/recent-activities')
}

export default api
