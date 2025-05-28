// Package phocus_diagnostics handles system diagnostics and periodic republishing
package phocus_diagnostics

import (
	"log"
	"sync"
	"time"

	mqtt "github.com/wolffshots/phocus/v2/mqtt"
)

// DiagnosticsUpdater implements mqtt.DiagnosticsUpdater interface
type DiagnosticsUpdater struct{}

// UpdateError implements the mqtt.DiagnosticsUpdater interface
func (d *DiagnosticsUpdater) UpdateError(err error) {
	UpdateError(err)
}

// DiagnosticsData holds all the diagnostic information
type DiagnosticsData struct {
	StartTime    time.Time
	LastError    string
	ErrorTime    time.Time
	HealthCounter int64
	Version      string
}

// Global diagnostics data with mutex for thread-safe access
var (
	diagnostics DiagnosticsData
	diagMutex   sync.RWMutex
)

// Initialize sets up the diagnostics system with initial values
func Initialize(version string) {
	diagMutex.Lock()
	defer diagMutex.Unlock()
	
	diagnostics = DiagnosticsData{
		StartTime:     time.Now(),
		LastError:     "",
		ErrorTime:     time.Time{},
		HealthCounter: 0,
		Version:       version,
	}
	
	// Set up diagnostics updater for mqtt package
	updater := &DiagnosticsUpdater{}
	mqtt.SetDiagnosticsUpdater(updater)
}

// UpdateError updates the last error and its timestamp
func UpdateError(err error) {
	diagMutex.Lock()
	defer diagMutex.Unlock()
	
	if err != nil {
		diagnostics.LastError = err.Error()
		diagnostics.ErrorTime = time.Now()
	} else {
		diagnostics.LastError = ""
		diagnostics.ErrorTime = time.Time{}
	}
}

// IncrementHealthCounter increments the health counter (thread-safe)
func IncrementHealthCounter() {
	diagMutex.Lock()
	defer diagMutex.Unlock()
	
	diagnostics.HealthCounter++
}

// GetDiagnostics returns a copy of current diagnostics data
func GetDiagnostics() DiagnosticsData {
	diagMutex.RLock()
	defer diagMutex.RUnlock()
	
	return diagnostics
}

// StartPeriodicRepublisher starts a goroutine that republishes diagnostics every 30 minutes
func StartPeriodicRepublisher(client mqtt.Client) {
	go func() {
		ticker := time.NewTicker(30 * time.Minute)
		defer ticker.Stop()
		
		// Publish immediately on start
		publishDiagnostics(client)
		
		for {
			select {
			case <-ticker.C:
				publishDiagnostics(client)
			}
		}
	}()
	
	log.Println("Started periodic diagnostics republisher (30 minute interval)")
}

// publishDiagnostics publishes all diagnostic data to MQTT topics
func publishDiagnostics(client mqtt.Client) {
	diagMutex.RLock()
	data := diagnostics
	diagMutex.RUnlock()
	
	// Increment health counter each time we publish
	IncrementHealthCounter()
	
	// Get updated health counter value
	diagMutex.RLock()
	healthCounter := diagnostics.HealthCounter
	diagMutex.RUnlock()
	
	// Publish start time
	err := mqtt.Send(client, "phocus/stats/start_time", 0, true, data.StartTime.Format(time.RFC3339), 10*time.Second)
	if err != nil {
		log.Printf("Failed to republish start_time: %v", err)
	}
	
	// Publish version
	err = mqtt.Send(client, "phocus/stats/version", 0, true, data.Version, 10*time.Second)
	if err != nil {
		log.Printf("Failed to republish version: %v", err)
	}
	
	// Publish last error (empty string if no error)
	err = mqtt.Send(client, "phocus/stats/error", 0, true, data.LastError, 10*time.Second)
	if err != nil {
		log.Printf("Failed to republish error: %v", err)
	}
	
	// Publish error timestamp if there was an error
	if !data.ErrorTime.IsZero() {
		err = mqtt.Send(client, "phocus/stats/error_time", 0, true, data.ErrorTime.Format(time.RFC3339), 10*time.Second)
		if err != nil {
			log.Printf("Failed to republish error_time: %v", err)
		}
	} else {
		err = mqtt.Send(client, "phocus/stats/error_time", 0, true, "", 10*time.Second)
		if err != nil {
			log.Printf("Failed to clear error_time: %v", err)
		}
	}
	
	// Publish health counter
	err = mqtt.Send(client, "phocus/stats/health_counter", 0, true, healthCounter, 10*time.Second)
	if err != nil {
		log.Printf("Failed to republish health_counter: %v", err)
	}
	
	// Publish uptime duration
	uptime := time.Since(data.StartTime)
	err = mqtt.Send(client, "phocus/stats/uptime_seconds", 0, true, int64(uptime.Seconds()), 10*time.Second)
	if err != nil {
		log.Printf("Failed to republish uptime: %v", err)
	}
	
	log.Printf("Republished diagnostics - Health Counter: %d, Uptime: %v", healthCounter, uptime.Truncate(time.Second))
}