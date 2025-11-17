import { useEffect, useRef, useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { Terminal, Trash2 } from 'lucide-react'
import { api, createLogWebSocket } from '../api/client'
import './LogViewer.css'

export default function LogViewer() {
  const [logs, setLogs] = useState<string[]>([])
  const [isConnected, setIsConnected] = useState(false)
  const logEndRef = useRef<HTMLDivElement>(null)
  const wsRef = useRef<WebSocket | null>(null)

  // Fetch recent logs on mount
  const { data: recentLogsData } = useQuery({
    queryKey: ['recent-logs'],
    queryFn: api.getRecentLogs,
  })

  useEffect(() => {
    if (recentLogsData?.data) {
      const logMessages = recentLogsData.data.map(
        (entry) => `[${new Date(entry.timestamp).toLocaleString()}] ${entry.message}`
      )
      setLogs(logMessages)
    }
  }, [recentLogsData])

  // WebSocket connection for live logs
  useEffect(() => {
    const ws = createLogWebSocket(
      (message) => {
        setLogs((prev) => [...prev, message])
      },
      (error) => {
        console.error('WebSocket error:', error)
        setIsConnected(false)
      }
    )

    ws.onopen = () => {
      setIsConnected(true)
    }

    ws.onclose = () => {
      setIsConnected(false)
    }

    wsRef.current = ws

    return () => {
      ws.close()
    }
  }, [])

  // Auto-scroll to bottom
  useEffect(() => {
    logEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [logs])

  const clearLogs = () => {
    setLogs([])
  }

  return (
    <div className="log-viewer">
      <div className="log-viewer-header">
        <div className="log-status">
          <Terminal size={18} />
          <span>
            {isConnected ? (
              <>
                <span className="status-dot connected"></span>
                Connected
              </>
            ) : (
              <>
                <span className="status-dot disconnected"></span>
                Disconnected
              </>
            )}
          </span>
        </div>
        <button onClick={clearLogs} className="clear-button" title="Clear logs">
          <Trash2 size={18} />
          Clear
        </button>
      </div>
      <div className="log-content">
        {logs.length === 0 ? (
          <div className="log-empty">No logs available</div>
        ) : (
          logs.map((log, index) => (
            <div key={index} className="log-entry">
              {log}
            </div>
          ))
        )}
        <div ref={logEndRef} />
      </div>
    </div>
  )
}

