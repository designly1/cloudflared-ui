const API_BASE_URL = '/api'

export interface ServiceStatus {
  activeState: string
  subState: string
  loadState: string
  description: string
  mainPID: number
  memoryCurrent: number
  cpuUsageNSec: number
}

export interface LogEntry {
  timestamp: string
  message: string
  priority: string
}

export interface Config {
  tunnel?: string
  ingress?: IngressRule[]
  metrics?: string
  loglevel?: string
  [key: string]: any
}

export interface IngressRule {
  hostname?: string
  service: string
  path?: string
}

export interface ApiResponse<T = any> {
  success: boolean
  message?: string
  data?: T
  error?: string
}

async function fetchApi<T>(
  endpoint: string,
  options?: RequestInit
): Promise<ApiResponse<T>> {
  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options?.headers,
    },
  })

  if (!response.ok) {
    const error = await response.json().catch(() => ({
      success: false,
      error: `HTTP ${response.status}: ${response.statusText}`,
    }))
    throw new Error(error.error || 'API request failed')
  }

  return response.json()
}

export const api = {
  // Service control
  startService: () => fetchApi('/service/start', { method: 'POST' }),
  stopService: () => fetchApi('/service/stop', { method: 'POST' }),
  restartService: () => fetchApi('/service/restart', { method: 'POST' }),
  
  // Status
  getStatus: () => fetchApi<ServiceStatus>('/service/status'),
  
  // Logs
  getRecentLogs: () => fetchApi<LogEntry[]>('/service/logs/recent'),
  
  // Config
  getConfig: () => fetchApi<Config>('/config'),
  updateConfig: (config: Config) =>
    fetchApi('/config', {
      method: 'POST',
      body: JSON.stringify(config),
    }),
  
  // Health
  healthCheck: () => fetchApi('/health'),
}

// WebSocket for log streaming
export function createLogWebSocket(
  onMessage: (message: string) => void,
  onError?: (error: Event) => void
): WebSocket {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const wsUrl = `${protocol}//${window.location.hostname}:8080/api/service/logs`
  
  const ws = new WebSocket(wsUrl)
  
  ws.onmessage = (event) => {
    onMessage(event.data)
  }
  
  ws.onerror = (error) => {
    console.error('WebSocket error:', error)
    onError?.(error)
  }
  
  return ws
}

