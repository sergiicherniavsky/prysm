package monitor

import (
	"fmt"
	"sync"

	"github.com/sirupsen/logrus"
)

type AttestationStats struct {
	mu              sync.Mutex
	successfulCount int
	failedCount     int
	failureReasons  map[string]int
}

// RecordSuccess increments the successful attestation count.
func (s *AttestationStats) RecordSuccess() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.successfulCount++
}

// RecordFailure increments the failed attestation count and records the reason for failure.
func (s *AttestationStats) RecordFailure(reason string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.failedCount++
	s.failureReasons[reason]++
}

// Summary logs the aggregated attestation statistics.
func (s *AttestationStats) Summary(log *logrus.Entry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	failedReasons := ""
	for reason, count := range s.failureReasons {
		failedReasons += fmt.Sprintf("[Failure reason: %s, Count: %d] ", reason, count)
	}

	log.WithFields(logrus.Fields{
		"successfulCount": s.successfulCount,
		"failedCount":     s.failedCount,
		"failedReasons":   failedReasons,
	}).Info("Aggregated performance of each epoch")
}
