package systemd

import (
	"context"
	"fmt"
)

type ServiceStatus struct {
	ActiveState   string `json:"activeState"`
	SubState      string `json:"subState"`
	LoadState     string `json:"loadState"`
	Description   string `json:"description"`
	MainPID       uint32 `json:"mainPID"`
	MemoryCurrent uint64 `json:"memoryCurrent"`
	CPUUsageNSec  uint64 `json:"cpuUsageNSec"`
}

// Status retrieves the current status of the cloudflared service
func (s *SystemdService) Status() (*ServiceStatus, error) {
	ctx := context.Background()
	props, err := s.conn.GetUnitPropertiesContext(ctx, s.serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to get unit properties: %w", err)
	}

	status := &ServiceStatus{
		ActiveState: getStringProp(props, "ActiveState"),
		SubState:    getStringProp(props, "SubState"),
		LoadState:   getStringProp(props, "LoadState"),
		Description: getStringProp(props, "Description"),
		MainPID:     getUint32Prop(props, "MainPID"),
	}

	// Get memory usage if available
	if mem, ok := props["MemoryCurrent"].(uint64); ok {
		status.MemoryCurrent = mem
	}

	// Get CPU usage if available
	if cpu, ok := props["CPUUsageNSec"].(uint64); ok {
		status.CPUUsageNSec = cpu
	}

	return status, nil
}

// Helper functions to safely extract properties
func getStringProp(props map[string]interface{}, key string) string {
	if val, ok := props[key].(string); ok {
		return val
	}
	return ""
}

func getUint32Prop(props map[string]interface{}, key string) uint32 {
	if val, ok := props[key].(uint32); ok {
		return val
	}
	return 0
}

