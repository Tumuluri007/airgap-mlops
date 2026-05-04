// Package main implements the cmba-sentinel runtime sidecar.
//
// The sentinel runs alongside the model server in every ml-serving pod.
// Every recheckIntervalSeconds (default 300) it re-hashes the mounted
// model file and compares to the ModelBinding's expected hash. On
// mismatch the sentinel:
//   - increments ModelBinding.status.mismatchEvents
//   - writes a structured event to /var/log/cmba/events.log
//   - sets a Kubernetes Event with reason CMBABindingDrift
//   - takes the configured action (terminate, alert, or crash)
//
// All alerting is local. No webhooks to Slack, no PagerDuty calls.
// Memory footprint target: <5 MB at steady state.
package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	bindingFile     = flag.String("binding-file", "/etc/cmba/binding.yaml", "Path to mounted ModelBinding")
	modelPath       = flag.String("model-path", "/models", "Mount path of the model file")
	recheckInterval = flag.Int("recheck-interval", 300, "Recheck interval in seconds (minimum 60)")
	onMismatch      = flag.String("on-mismatch", "terminate", "Action on mismatch: terminate|alert|crash")
	eventLogPath    = flag.String("event-log", "/var/log/cmba/events.log", "Local event log path")
)

func main() {
	flag.Parse()

	if *recheckInterval < 60 {
		log.Printf("recheck-interval %d below minimum 60s; clamping", *recheckInterval)
		*recheckInterval = 60
	}

	v := NewVerifier(*bindingFile, *modelPath)
	a := NewAlerter(*eventLogPath, *onMismatch)

	log.Printf("cmba-sentinel starting: model=%s binding=%s interval=%ds onMismatch=%s",
		*modelPath, *bindingFile, *recheckInterval, *onMismatch)

	ticker := time.NewTicker(time.Duration(*recheckInterval) * time.Second)
	defer ticker.Stop()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Initial verification on startup.
	doCheck(v, a)

	for {
		select {
		case <-ticker.C:
			doCheck(v, a)
		case sig := <-stop:
			log.Printf("received signal %s; sentinel exiting", sig)
			return
		}
	}
}

func doCheck(v *Verifier, a *Alerter) {
	res, err := v.Verify()
	if err != nil {
		log.Printf("verification error: %v", err)
		a.OnError(err)
		return
	}
	if !res.Match {
		log.Printf("CMBA MISMATCH: expected=%s actual=%s", res.Expected, res.Actual)
		a.OnMismatch(res)
		return
	}
	log.Printf("CMBA OK: model verified (%s)", res.Actual)
}
