import { useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { api } from '../api/client'
import ServiceControls from './ServiceControls'
import StatusDisplay from './StatusDisplay'
import LogViewer from './LogViewer'
import './Dashboard.css'

export default function Dashboard() {
  const queryClient = useQueryClient()

  // Query service status
  const { data: statusData, isLoading: statusLoading } = useQuery({
    queryKey: ['service-status'],
    queryFn: api.getStatus,
    refetchInterval: 5000, // Poll every 5 seconds
  })

  // Service control mutations
  const startMutation = useMutation({
    mutationFn: api.startService,
  })

  const stopMutation = useMutation({
    mutationFn: api.stopService,
  })

  const restartMutation = useMutation({
    mutationFn: api.restartService,
  })

  // Invalidate queries on successful mutations
  useEffect(() => {
    if (startMutation.isSuccess || stopMutation.isSuccess || restartMutation.isSuccess) {
      queryClient.invalidateQueries({ queryKey: ['service-status'] })
    }
  }, [startMutation.isSuccess, stopMutation.isSuccess, restartMutation.isSuccess, queryClient])

  const status = statusData?.data

  return (
    <div className="dashboard">
      <div className="dashboard-grid">
        <div className="dashboard-section">
          <h2>Service Status</h2>
          <StatusDisplay status={status} isLoading={statusLoading} />
        </div>

        <div className="dashboard-section">
          <h2>Service Controls</h2>
          <ServiceControls
            onStart={() => startMutation.mutate()}
            onStop={() => stopMutation.mutate()}
            onRestart={() => restartMutation.mutate()}
            isStarting={startMutation.isPending}
            isStopping={stopMutation.isPending}
            isRestarting={restartMutation.isPending}
            serviceState={status?.activeState}
          />
        </div>

        <div className="dashboard-section full-width">
          <h2>Live Logs</h2>
          <LogViewer />
        </div>
      </div>
    </div>
  )
}

