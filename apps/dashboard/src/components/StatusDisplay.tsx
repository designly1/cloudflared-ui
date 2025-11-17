import { CheckCircle, XCircle, AlertCircle, Loader } from 'lucide-react'
import type { ServiceStatus } from '../api/client'
import './StatusDisplay.css'

interface StatusDisplayProps {
  status?: ServiceStatus
  isLoading: boolean
}

export default function StatusDisplay({ status, isLoading }: StatusDisplayProps) {
  if (isLoading || !status) {
    return (
      <div className="status-display">
        <div className="status-item loading">
          <Loader className="spinning" size={24} />
          <span>Loading status...</span>
        </div>
      </div>
    )
  }

  const getStatusIcon = () => {
    switch (status.activeState) {
      case 'active':
        return <CheckCircle size={24} className="status-icon active" />
      case 'failed':
        return <XCircle size={24} className="status-icon failed" />
      case 'inactive':
        return <AlertCircle size={24} className="status-icon inactive" />
      default:
        return <AlertCircle size={24} className="status-icon unknown" />
    }
  }

  const getStatusText = () => {
    switch (status.activeState) {
      case 'active':
        return 'Running'
      case 'failed':
        return 'Failed'
      case 'inactive':
        return 'Stopped'
      default:
        return status.activeState || 'Unknown'
    }
  }

  const formatMemory = (bytes: number) => {
    if (bytes === 0) return 'N/A'
    const mb = bytes / (1024 * 1024)
    return `${mb.toFixed(2)} MB`
  }

  return (
    <div className="status-display">
      <div className="status-main">
        {getStatusIcon()}
        <div className="status-text">
          <div className="status-label">Status</div>
          <div className="status-value">{getStatusText()}</div>
        </div>
      </div>

      <div className="status-details">
        <div className="status-detail-item">
          <span className="detail-label">State:</span>
          <span className="detail-value">{status.subState}</span>
        </div>
        <div className="status-detail-item">
          <span className="detail-label">Load State:</span>
          <span className="detail-value">{status.loadState}</span>
        </div>
        <div className="status-detail-item">
          <span className="detail-label">PID:</span>
          <span className="detail-value">{status.mainPID || 'N/A'}</span>
        </div>
        <div className="status-detail-item">
          <span className="detail-label">Memory:</span>
          <span className="detail-value">{formatMemory(status.memoryCurrent)}</span>
        </div>
      </div>
    </div>
  )
}

