package ailimit

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var ErrDailyLimit = errors.New("AI daily quota exhausted")

// Limiter is shared by both chat routes. Daily usage counts every upstream
// attempt, including failures/fallbacks, and is persisted before making a call.
// Use one server process per quota file.
type Limiter struct {
	mu                              sync.Mutex
	now                             func() time.Time
	rate, concurrent, daily, active int
	visitors                        map[string][]time.Time
	lastCleanup                     time.Time
	path                            string
}

type usage struct {
	Day      string `json:"day"`
	Attempts int    `json:"attempts"`
}

func New(rate, concurrent, daily int, path string) *Limiter {
	return &Limiter{rate: rate, concurrent: concurrent, daily: daily, path: path, now: time.Now, visitors: make(map[string][]time.Time)}
}

func (l *Limiter) Acquire(ip string) (release func(), retry int, message string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := l.now()
	if now.Sub(l.lastCleanup) >= time.Minute {
		for key, times := range l.visitors {
			if len(times) == 0 || !times[len(times)-1].After(now.Add(-time.Minute)) {
				delete(l.visitors, key)
			}
		}
		l.lastCleanup = now
	}
	times := l.visitors[ip]
	for len(times) > 0 && !times[0].After(now.Add(-time.Minute)) {
		times = times[1:]
	}
	if len(times) >= l.rate {
		return nil, int(times[0].Add(time.Minute).Sub(now).Seconds()) + 1, "提问过于频繁，请稍后重试"
	}
	if l.active >= l.concurrent {
		return nil, 5, "AI 正忙，请稍后重试"
	}
	if len(l.visitors) >= 10000 && len(times) == 0 {
		return nil, 60, "AI 正忙，请稍后重试"
	}
	l.visitors[ip] = append(times, now)
	l.active++
	var once sync.Once
	return func() { once.Do(func() { l.mu.Lock(); l.active--; l.mu.Unlock() }) }, 0, ""
}

func (l *Limiter) ReserveDaily() error {
	l.mu.Lock()
	defer l.mu.Unlock()
	day := l.now().In(time.FixedZone("Asia/Shanghai", 8*3600)).Format("2006-01-02")
	var used usage
	raw, err := os.ReadFile(l.path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		if err = json.Unmarshal(raw, &used); err != nil {
			return err
		}
		_, dateErr := time.Parse("2006-01-02", used.Day)
		if dateErr != nil || used.Attempts < 0 {
			return errors.New("invalid AI usage file")
		}
	}
	if used.Day != day {
		used = usage{Day: day}
	}
	if used.Attempts >= l.daily {
		return ErrDailyLimit
	}
	used.Attempts++
	raw, err = json.Marshal(used)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(l.path), 0700); err != nil {
		return err
	}
	f, err := os.CreateTemp(filepath.Dir(l.path), ".ai-usage-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err = f.Write(raw); err != nil {
		f.Close()
		return err
	}
	if err = f.Sync(); err != nil {
		f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), l.path)
}

func (l *Limiter) DailyRetryAfter() int {
	now := l.now().In(time.FixedZone("Asia/Shanghai", 8*3600))
	next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
	return int(next.Sub(now).Seconds()) + 1
}
