import { Play, Square, RotateCw } from 'lucide-react'
import './ServiceControls.css'

interface ServiceControlsProps {
  onStart: () => void
  onStop: () => void
  onRestart: () => void
  isStarting: boolean
  isStopping: boolean
  isRestarting: boolean
  serviceState?: string
}

export default function ServiceControls({
  onStart,
  onStop,
  onRestart,
  isStarting,
  isStopping,
  isRestarting,
  serviceState,
}: ServiceControlsProps) {
  const isActive = serviceState === 'active'
  const isInactive = serviceState === 'inactive' || serviceState === 'failed'

  return (
    <div className="service-controls">
      <button
        onClick={onStart}
        disabled={isStarting || isActive || isRestarting}
        className="control-button start"
      >
        <Play size={20} />
        {isStarting ? 'Starting...' : 'Start'}
      </button>

      <button
        onClick={onStop}
        disabled={isStopping || isInactive || isRestarting}
        className="control-button stop"
      >
        <Square size={20} />
        {isStopping ? 'Stopping...' : 'Stop'}
      </button>

      <button
        onClick={onRestart}
        disabled={isRestarting || isInactive}
        className="control-button restart"
      >
        <RotateCw size={20} className={isRestarting ? 'spinning' : ''} />
        {isRestarting ? 'Restarting...' : 'Restart'}
      </button>
    </div>
  )
}

