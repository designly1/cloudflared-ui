package systemd

import (
	"context"
	"fmt"

	"github.com/coreos/go-systemd/v22/dbus"
)

const serviceName = "cloudflared.service"

type SystemdService struct {
	conn *dbus.Conn
}

// New creates a new SystemdService with a system D-Bus connection
func New() (*SystemdService, error) {
	conn, err := dbus.NewSystemConnectionContext(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to system D-Bus: %w", err)
	}
	return &SystemdService{conn: conn}, nil
}

// Close closes the D-Bus connection
func (s *SystemdService) Close() {
	if s.conn != nil {
		s.conn.Close()
	}
}

// Start starts the cloudflared service
func (s *SystemdService) Start() error {
	ctx := context.Background()
	responseChan := make(chan string)
	_, err := s.conn.StartUnitContext(ctx, serviceName, "replace", responseChan)
	if err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}
	
	// Wait for the job to complete
	status := <-responseChan
	if status != "done" {
		return fmt.Errorf("start job failed with status: %s", status)
	}
	return nil
}

// Stop stops the cloudflared service
func (s *SystemdService) Stop() error {
	ctx := context.Background()
	responseChan := make(chan string)
	_, err := s.conn.StopUnitContext(ctx, serviceName, "replace", responseChan)
	if err != nil {
		return fmt.Errorf("failed to stop service: %w", err)
	}
	
	// Wait for the job to complete
	status := <-responseChan
	if status != "done" {
		return fmt.Errorf("stop job failed with status: %s", status)
	}
	return nil
}

// Restart restarts the cloudflared service
func (s *SystemdService) Restart() error {
	ctx := context.Background()
	responseChan := make(chan string)
	_, err := s.conn.RestartUnitContext(ctx, serviceName, "replace", responseChan)
	if err != nil {
		return fmt.Errorf("failed to restart service: %w", err)
	}
	
	// Wait for the job to complete
	status := <-responseChan
	if status != "done" {
		return fmt.Errorf("restart job failed with status: %s", status)
	}
	return nil
}

