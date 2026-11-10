package config

import "testing"

func TestAILimitsCannotBeDisabledByInvalidEnvironment(t *testing.T) {
	for _, value := range []string{"", "0", "-1", "invalid"} {
		t.Setenv("AI_RATE_PER_MINUTE", value)
		t.Setenv("AI_MAX_CONCURRENT", value)
		t.Setenv("AI_DAILY_LIMIT", value)
		t.Setenv("AI_MAX_TOKENS", value)
		cfg := Load()
		if cfg.AIRatePerMinute != 6 || cfg.AIMaxConcurrent != 3 || cfg.AIDailyLimit != 500 || cfg.AIMaxTokens != 1024 {
			t.Fatalf("invalid defaults for %q: %#v", value, cfg)
		}
	}
	t.Setenv("AI_MAX_TOKENS", "999999")
	if Load().AIMaxTokens != 4096 {
		t.Fatal("output cap exceeded")
	}
}
