package monitor

import (
	"sync"
	"testing"

	logTest "github.com/sirupsen/logrus/hooks/test"

	"github.com/prysmaticlabs/prysm/v5/testing/require"
)

func TestRecordSuccess(t *testing.T) {
	stats := &AttestationStats{
		failureReasons: make(map[string]int),
	}
	stats.RecordSuccess()
	stats.RecordSuccess()
	stats.RecordSuccess()

	require.Equal(
		t,
		3,
		stats.successfulCount,
		"expected successfulCount to be 1, got %d",
		stats.successfulCount,
	)
}

func TestRecordFailure(t *testing.T) {
	stats := &AttestationStats{
		failureReasons: make(map[string]int),
	}

	reason1 := "network error"
	reason2 := "timeout"
	stats.RecordFailure(reason1)
	stats.RecordFailure(reason2)
	stats.RecordFailure(reason2)

	require.Equal(t, 3, stats.failedCount, "expected failedCount to be 3, got %d", stats.failedCount)

	require.Equal(
		t,
		1,
		stats.failureReasons[reason1],
		"expected failureReasons[reason1] to be 1, got %d", stats.failureReasons[reason1],
	)
	require.Equal(
		t,
		2,
		stats.failureReasons[reason2],
		"expected failureReasons[reason2] to be 2, got %d", stats.failureReasons[reason2],
	)
}

func TestRecordSuccessWithGoroutines(t *testing.T) {
	stats := &AttestationStats{
		failureReasons: make(map[string]int),
	}

	var wg sync.WaitGroup
	numGoroutines := 100

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			stats.RecordSuccess()
		}()
	}

	wg.Wait()

	require.Equal(
		t,
		numGoroutines,
		stats.successfulCount,
		"expected successfulCount to be %d, got %d",
		numGoroutines,
		stats.successfulCount,
	)
}

func TestRecordFailureWithGoroutines(t *testing.T) {
	stats := &AttestationStats{
		failureReasons: make(map[string]int),
	}

	var wg sync.WaitGroup
	numGoroutines := 100
	reason := "network error"

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			stats.RecordFailure(reason)
		}()
	}

	wg.Wait()

	require.Equal(
		t,
		numGoroutines,
		stats.failedCount,
		"expected failedCount to be %d, got %d", numGoroutines, stats.failedCount,
	)

	require.Equal(
		t,
		numGoroutines,
		stats.failureReasons[reason],
		"expected failureReasons[%s] to be %d, got %d", reason, numGoroutines, stats.failureReasons[reason],
	)
}

func TestAttestationStats_Summary(t *testing.T) {
	hook := logTest.NewGlobal()
	stats := &AttestationStats{
		successfulCount: 10,
		failedCount:     5,
		failureReasons: map[string]int{
			"network error": 2,
			"timeout":       3,
		},
	}

	stats.Summary(log)

	wanted1 := "\"Aggregated performance of each epoch\" failedCount=5 failedReasons=\"[Failure reason: network error, Count: 2] [Failure reason: timeout, Count: 3] \" prefix=monitor successfulCount=10"
	require.LogsContain(t, hook, wanted1)

}
