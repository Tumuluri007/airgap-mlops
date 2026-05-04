package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"syscall"
	"time"
)

// Alerter writes structured events to a local log path that Loki picks up,
// and takes the configured action when a mismatch is detected.
//
// All alerting is local. No external webhooks (Slack, PagerDuty, etc.)
// are invoked. This is by design: the cluster has no internet.
type Alerter struct {
	EventLogPath string
	OnMismatchAction string
}

// NewAlerter constructs an Alerter and ensures the event log directory exists.
func NewAlerter(eventLogPath, action string) *Alerter {
	_ = os.MkdirAll(filepath.Dir(eventLogPath), 0o755)
	return &Alerter{
		EventLogPath:     eventLogPath,
		OnMismatchAction: action,
	}
}

type event struct {
	Time     time.Time `json:"time"`
	Reason   string    `json:"reason"`
	Severity string    `json:"severity"`
	Message  string    `json:"message"`
	Expected string    `json:"expected,omitempty"`
	Actual   string    `json:"actual,omitempty"`
	Path     string    `json:"path,omitempty"`
}

// OnMismatch records the binding drift event and takes the configured action.
func (a *Alerter) OnMismatch(res *VerifyResult) {
	e := event{
		Time:     time.Now().UTC(),
		Reason:   "CMBABindingDrift",
		Severity: "critical",
		Message:  "model file hash does not match ModelBinding declaration",
		Expected: res.Expected,
		Actual:   res.Actual,
		Path:     res.Path,
	}
	a.write(e)

	switch a.OnMismatchAction {
	case "terminate":
		// Send SIGTERM to the main container (pid 1 in shared pid namespace).
		log.Printf("onMismatch=terminate: sending SIGTERM to pid 1")
		_ = syscall.Kill(1, syscall.SIGTERM)
	case "crash":
		log.Printf("onMismatch=crash: exiting non-zero to force pod restart")
		os.Exit(2)
	case "alert":
		log.Printf("onMismatch=alert: event written, no termination action")
	default:
		log.Printf("onMismatch=%s: unknown action, defaulting to alert-only", a.OnMismatchAction)
	}
}

// OnError records a verifier error event without taking destructive action.
func (a *Alerter) OnError(err error) {
	e := event{
		Time:     time.Now().UTC(),
		Reason:   "CMBAVerificationError",
		Severity: "warning",
		Message:  fmt.Sprintf("sentinel could not verify model: %v", err),
	}
	a.write(e)
}

func (a *Alerter) write(e event) {
	f, err := os.OpenFile(a.EventLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		log.Printf("cannot open event log %s: %v", a.EventLogPath, err)
		return
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(e); err != nil {
		log.Printf("cannot encode event: %v", err)
	}
}
