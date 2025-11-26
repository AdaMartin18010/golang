import axios from 'axios'

// API客户端配置
export const apiClient = axios.create({
  baseURL: '/api/v1',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// 请求拦截器
apiClient.interceptors.request.use(
  (config) => {
    // 可以在这里添加认证token等
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
apiClient.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    // 统一错误处理
    console.error('API Error:', error)
    return Promise.reject(error)
  }
)

// API方法
export const analysisAPI = {
  analyzeCFG: (code: string) =>
    apiClient.post('/analysis/cfg', { code }),
  analyzeConcurrency: (code: string) =>
    apiClient.post('/analysis/concurrency', { code }),
  analyzeTypes: (code: string) =>
    apiClient.post('/analysis/types', { code }),
  getHistory: () =>
    apiClient.get('/analysis/history'),
}

export const patternsAPI = {
  list: (category?: string) =>
    apiClient.get('/patterns', { params: { category } }),
  get: (name: string) =>
    apiClient.get(`/patterns/${name}`),
  generate: (pattern: string, parameters: Record<string, string>) =>
    apiClient.post('/patterns/generate', { pattern, parameters }),
}

export const projectsAPI = {
  list: () =>
    apiClient.get('/projects'),
  create: (name: string, description: string, path: string) =>
    apiClient.post('/projects', { name, description, path }),
  get: (id: string) =>
    apiClient.get(`/projects/${id}`),
  delete: (id: string) =>
    apiClient.delete(`/projects/${id}`),
}

