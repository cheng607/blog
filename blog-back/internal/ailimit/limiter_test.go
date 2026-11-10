package ailimit

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRateConcurrencyAndCleanup(t *testing.T) {
	l := New(2, 1, 10, filepath.Join(t.TempDir(), "quota.json"))
	now := time.Date(2026, 11, 7, 0, 0, 0, 0, time.UTC)
	l.now = func() time.Time { return now }
	release, _, _ := l.Acquire("a")
	if release == nil {
		t.Fatal("first rejected")
	}
	if r, _, _ := l.Acquire("b"); r != nil {
		t.Fatal("concurrency bypassed")
	}
	release()
	release()
	release, _, _ = l.Acquire("a")
	release()
	if r, retry, _ := l.Acquire("a"); r != nil || retry <= 0 {
		t.Fatal("rate bypassed")
	}
	now = now.Add(time.Minute)
	release, _, _ = l.Acquire("a")
	if release == nil {
		t.Fatal("rate did not reset")
	}
	release()
}
func TestDailyQuotaPersistsAndResetsAtShanghaiMidnight(t *testing.T) {
	path := filepath.Join(t.TempDir(), "quota.json")
	now := time.Date(2026, 11, 7, 15, 59, 59, 0, time.UTC)
	l := New(5, 2, 1, path)
	l.now = func() time.Time { return now }
	if err := l.ReserveDaily(); err != nil {
		t.Fatal(err)
	}
	restarted := New(5, 2, 1, path)
	restarted.now = l.now
	if err := restarted.ReserveDaily(); !errors.Is(err, ErrDailyLimit) {
		t.Fatalf("restart lost quota: %v", err)
	}
	now = now.Add(time.Second)
	if err := restarted.ReserveDaily(); err != nil {
		t.Fatal("did not reset at Shanghai midnight", err)
	}
	if err := os.WriteFile(path, []byte("corrupt"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := restarted.ReserveDaily(); err == nil {
		t.Fatal("corrupt quota allowed calls")
	}
}
func TestConcurrentReservationsCannotExceedDailyLimit(t *testing.T) {
	l := New(60, 20, 7, filepath.Join(t.TempDir(), "quota.json"))
	var wg sync.WaitGroup
	var accepted atomic.Int32
	for i := 0; i < 30; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if l.ReserveDaily() == nil {
				accepted.Add(1)
			}
		}()
	}
	wg.Wait()
	if accepted.Load() != 7 {
		t.Fatalf("accepted %d", accepted.Load())
	}
}
