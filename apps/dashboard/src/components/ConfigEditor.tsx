import { useState, useEffect } from 'react'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Save, RotateCw } from 'lucide-react'
import { api, type Config, type ApiResponse } from '../api/client'
import './ConfigEditor.css'

export default function ConfigEditor() {
  const queryClient = useQueryClient()
  const [configText, setConfigText] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [successMessage, setSuccessMessage] = useState<string | null>(null)

  // Fetch config
  const { data: configData, isLoading } = useQuery<ApiResponse<Config>>({
    queryKey: ['config'],
    queryFn: api.getConfig,
  })

  // Update config text when data loads
  useEffect(() => {
    if (configData?.data) {
      setConfigText(JSON.stringify(configData.data, null, 2))
    }
  }, [configData])

  // Update config mutation
  const updateMutation = useMutation({
    mutationFn: async (config: Config) => {
      const result = await api.updateConfig(config)
      return result
    },
  })

  // Handle mutation success/error
  useEffect(() => {
    if (updateMutation.isSuccess) {
      queryClient.invalidateQueries({ queryKey: ['config'] })
      setSuccessMessage('Configuration updated successfully')
      setError(null)
      setTimeout(() => setSuccessMessage(null), 3000)
    }
    if (updateMutation.isError) {
      setError(updateMutation.error?.message || 'Failed to update configuration')
      setSuccessMessage(null)
    }
  }, [updateMutation.isSuccess, updateMutation.isError, updateMutation.error, queryClient])

  const handleSave = () => {
    try {
      const config = JSON.parse(configText)
      updateMutation.mutate(config)
    } catch (err) {
      setError('Invalid JSON format')
      setSuccessMessage(null)
    }
  }

  const handleReset = () => {
    if (configData?.data) {
      setConfigText(JSON.stringify(configData.data, null, 2))
      setError(null)
      setSuccessMessage(null)
    }
  }

  if (isLoading) {
    return (
      <div className="config-editor">
        <div className="loading">Loading configuration...</div>
      </div>
    )
  }

  return (
    <div className="config-editor">
      <div className="config-header">
        <h3>Edit Cloudflared Configuration</h3>
        <p className="config-description">
          Edit the cloudflared configuration in JSON format. Changes will be written to the config file.
        </p>
      </div>

      {error && <div className="message error">{error}</div>}
      {successMessage && <div className="message success">{successMessage}</div>}

      <textarea
        className="config-textarea"
        value={configText}
        onChange={(e) => setConfigText(e.target.value)}
        spellCheck={false}
      />

      <div className="config-actions">
        <button
          onClick={handleSave}
          disabled={updateMutation.isPending}
          className="save-button"
        >
          <Save size={18} />
          {updateMutation.isPending ? 'Saving...' : 'Save Configuration'}
        </button>
        <button onClick={handleReset} className="reset-button">
          <RotateCw size={18} />
          Reset
        </button>
      </div>

      <div className="config-warning">
        <strong>Warning:</strong> Incorrect configuration may prevent cloudflared from starting.
        Make sure to restart the service after saving changes.
      </div>
    </div>
  )
}

