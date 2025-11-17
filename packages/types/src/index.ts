// Service Status Types
export interface ServiceStatus {
  activeState: string
  subState: string
  loadState: string
  description: string
  mainPID: number
  memoryCurrent: number
  cpuUsageNSec: number
}

// Log Entry Types
export interface LogEntry {
  timestamp: string
  message: string
  priority: string
}

// Config Types
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

// API Response Types
export interface ApiResponse<T = any> {
  success: boolean
  message?: string
  data?: T
  error?: string
}

export interface ErrorResponse {
  success: false
  error: string
}

// Service Control Types
export type ServiceAction = 'start' | 'stop' | 'restart'

export interface ServiceControlRequest {
  action: ServiceAction
}

// WebSocket Message Types
export interface LogStreamMessage {
  type: 'log'
  data: string
  timestamp: number
}

export interface ConnectionMessage {
  type: 'connected' | 'disconnected'
  timestamp: number
}

export type WebSocketMessage = LogStreamMessage | ConnectionMessage

