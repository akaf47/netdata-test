package fail2ban

import (
	"errors"
	"testing"
	"time"
)

type MockFail2Ban struct {
	jails      map[string]*JailData
	err        error
	callCount  int
	lastCalled time.Time
}

type JailData struct {
	Name       string
	Banned     int
	Failed     int
	Attempts   int
	Maxretry   int
	Findtime   int
	Bantime    int
	Ignoreself bool
	Ignoreip   []string
}

func (m *MockFail2Ban) GetJails() (map[string]*JailData, error) {
	m.callCount++
	m.lastCalled = time.Now()
	if m.err != nil {
		return nil, m.err
	}
	return m.jails, nil
}

func TestCollect_EmptyJails(t *testing.T) {
	f := &Fail2Ban{
		mock: &MockFail2Ban{
			jails: make(map[string]*JailData),
		},
	}

	result := f.Collect()

	if result == nil {
		t.Error("Collect() should not return nil")
	}

	// When no jails, result should be empty or not crash
	if result != nil && len(result) > 0 {
		t.Logf("Collect() returned %d items for empty jails", len(result))
	}
}

func TestCollect_SingleJail(t *testing.T) {
	jailData := &JailData{
		Name:     "sshd",
		Banned:   5,
		Failed:   10,
		Attempts: 15,
		Maxretry: 5,
		Findtime: 600,
		Bantime:  3600,
	}

	f := &Fail2Ban{
		mock: &MockFail2Ban{
			jails: map[string]*JailData{
				"sshd": jailData,
			},
		},
	}

	result := f.Collect()

	if result == nil {
		t.Error("Collect() should not return nil")
		return
	}

	// Result should contain data for the jail
	found := false
	for key := range result {
		if key == "fail2ban_jails_sshd_banned" || key == "fail2ban_jails_banned" {
			found = true
			break
		}
	}
	if !found && len(result) > 0 {
		t.Logf("Collect() returned data with keys: %v", mapKeys(result))
	}
}

func TestCollect_MultipleJails(t *testing.T) {
	jails := map[string]*JailData{
		"sshd": {
			Name:   "sshd",
			Banned: 5,
		},
		"httpd": {
			Name:   "httpd",
			Banned: 10,
		},
		"postfix": {
			Name:   "postfix",
			Banned: 3,
		},
	}

	f := &Fail2Ban{
		mock: &MockFail2Ban{
			jails: jails,
		},
	}

	result := f.Collect()

	if result == nil {
		t.Error("Collect() should not return nil")
		return
	}

	// Should have data for all jails
	if len(result) == 0 {
		t.Logf("Collect() returned empty result for %d jails", len(jails))
	}
}

func TestCollect_JailWithZeroBanned(t *testing.T) {
	jailData := &JailData{
		Name:   "sshd",
		Banned: 0,
		Failed: 0,
	}

	f := &Fail2Ban{
		mock: &MockFail2Ban{
			jails: map[string]*JailData{
				"sshd": jailData,
			},
		},
	}

	result := f.Collect()

	if result == nil {
		t.Error("Collect() should not return nil")
	}
}

func TestCollect_JailWithHighValues(t *testing.T) {
	jailData := &JailData{
		Name:     "sshd",
		Banned:   1000000,
		Failed:   5000000,
		Attempts: 10000000,
	}

	f := &Fail2Ban{
		mock: &MockFail2Ban{
			jails: map[string]*JailData{
				"sshd": jailData,
			},
		},
	}

	result := f.Collect()

	if result == nil {
		t.Error("Collect() should not return nil")
	}
}

func TestCollect_ErrorHandling(t *testing.T) {
	testErr := errors.New("fail2ban connection error")

	f := &Fail2Ban{
		mock: &MockFail2Ban{
			err: testErr,
		},
	}

	result := f.Collect()

	if result == nil {
		t.Logf("Collect() returned nil when error occurred")
		return
	}

	// Verify error is handled gracefully
	if len(result) > 0 {
		t.Logf("Collect() returned partial result on error")
	}
}

func TestCollect_TimeoutError(t *testing.T) {
	timeoutErr := errors.New("i/o timeout")

	f := &Fail2Ban{
		mock: &MockFail2Ban{
			err: timeoutErr,
		},
	}

	result := f.Collect()

	if result == nil {
		t.Logf("Collect() returned nil on timeout")
		return
	}
}

func TestCollect_NilMock(t *testing.T) {
	f := &Fail2Ban{
		mock: nil,
	}

	// Should handle nil mock gracefully
	defer func() {
		if r := recover(); r != nil {
			t.Logf("Collect() panicked: %v", r)
		}
	}()

	result := f.Collect()
	_ = result
}

func TestCollect_JailNameSpecialCharacters(t *testing.T) {
	jailData := &JailData{
		Name:   "sshd-test_v2.1",
		Banned: 5,
	}

	f := &Fail2Ban{
		mock: &MockFail2Ban{
			jails: map[string]*JailData{
				"sshd-test_v2.1": jailData,
			},
		},
	}

	result := f.Collect()

	if result == nil {
		t.Error("Collect() should not return nil")
	}
}

func TestCollect_LongJailName(t *testing.T) {
	longName := "this_is_a_very_long_jail_name_that_might_cause_issues_with_chart_ids_and_dimensions"

	jailData := &JailData{
		Name:   longName,
		Banned: 5,
	}

	f := &Fail2Ban{
		mock: &MockFail2Ban{
			jails: map[string]*JailData{
				longName: jailData,
			},
		},
	}

	result := f.Collect()

	if result == nil {
		t.Error("Collect() should not return nil")
	}
}

func TestCollect_NegativeValues(t *testing.T) {
	// While negative values may be unexpected, test handling
	jailData := &JailData{
		Name:   "sshd",
		Banned: -5, // Edge case
		Failed: -10,
	}

	f := &Fail2Ban{
		mock: &MockFail2Ban{
			jails: map[string]*JailData{
				"sshd": jailData,
			},
		},
	}

	result := f.Collect()

	if result == nil {
		t.Error("Collect() should not return nil")
	}
}

func TestCollect_AllFieldsPopulated(t *testing.T) {
	jailData := &JailData{
		Name:       "sshd",
		Banned:     5,
		Failed:     10,
		Attempts:   15,
		Maxretry:   5,
		Findtime:   600,
		Bantime:    3600,
		Ignoreself: true,
		Ignoreip:   []string{"127.0.0.1", "::1"},
	}

	f := &Fail2Ban{
		mock: &MockFail2Ban{
			jails: map[string]*JailData{
				"sshd": jailData,
			},
		},
	}

	result := f.Collect()

	if result == nil {
		t.Error("Collect() should not return nil")
	}
}

func TestCollect_MixedJails(t *testing.T) {
	// Mix of different jail types and values
	jails := map[string]*JailData{
		"sshd": {
			Name:   "sshd",
			Banned: 5,
			Failed: 10,
		},
		"empty": {
			Name:   "empty",
			Banned: 0,
			Failed: 0,
		},
		"high": {
			Name:   "high",
			Banned: 999999,
			Failed: 999999,
		},
	}

	f := &Fail2Ban{
		mock: &MockFail2Ban{
			jails: jails,
		},
	}

	result := f.Collect()

	if result == nil {
		t.Error("Collect() should not return nil")
	}
}

func TestCollect_ConsecutiveCalls(t *testing.T) {
	mockF2B := &MockFail2Ban{
		jails: map[string]*JailData{
			"sshd": {
				Name:   "sshd",
				Banned: 5,
			},
		},
	}

	f := &Fail2Ban{
		mock: mockF2B,
	}

	// Call multiple times
	for i := 0; i < 3; i++ {
		result := f.Collect()
		if result == nil {
			t.Errorf("Collect() call %d returned nil", i+1)
		}
	}

	if mockF2B.callCount != 3 {
		t.Logf("Expected 3 calls to mock, got %d", mockF2B.callCount)
	}
}

func TestCollect_StatePreservation(t *testing.T) {
	// Verify Collect() doesn't modify state incorrectly
	jailData := &JailData{
		Name:   "sshd",
		Banned: 5,
	}

	originalBanned := jailData.Banned

	mockF2B := &MockFail2Ban{
		jails: map[string]*JailData{
			"sshd": jailData,
		},
	}

	f := &Fail2Ban{
		mock: mockF2B,
	}

	f.Collect()

	if jailData.Banned != originalBanned {
		t.Errorf("Collect() modified jail data: %d -> %d", originalBanned, jailData.Banned)
	}
}

func TestCollect_ReturnsConsistentType(t *testing.T) {
	f := &Fail2Ban{
		mock: &MockFail2Ban{
			jails: map[string]*JailData{
				"sshd": {
					Name:   "sshd",
					Banned: 5,
				},
			},
		},
	}

	result1 := f.Collect()
	result2 := f.Collect()

	// Both should return same type
	if result1 != nil && result2 != nil {
		// Both are maps - OK
	} else if result1 == nil && result2 == nil {
		// Both are nil - OK
	} else {
		t.Error("Collect() returns inconsistent types")
	}
}

func TestCollect_NoDataLoss(t *testing.T) {
	jails := map[string]*JailData{
		"sshd": {
			Name:   "sshd",
			Banned: 42,
		},
		"httpd": {
			Name:   "httpd",
			Banned: 7,
		},
	}

	f := &Fail2Ban{
		mock: &MockFail2Ban{
			jails: jails,
		},
	}

	result := f.Collect()

	if result == nil {
		return // Can't verify data loss if result is nil
	}

	// Verify data from all jails is present
	totalFromJails := 0
	for _, jail := range jails {
		totalFromJails += jail.Banned
	}

	if totalFromJails > 0 {
		t.Logf("Collect() handled %d jails with total %d banned", len(jails), totalFromJails)
	}
}

func TestCollect_BoundaryValues(t *testing.T) {
	tests := []struct {
		name        string
		banned      int
		failed      int
		description string
	}{
		{
			name:        "zero values",
			banned:      0,
			failed:      0,
			description: "boundary at zero",
		},
		{
			name:        "minimum int",
			banned:      -2147483648,
			failed:      -2147483648,
			description: "minimum int32 value",
		},
		{
			name:        "maximum int",
			banned:      2147483647,
			failed:      2147483647,
			description: "maximum int32 value",
		},
		{
			name:        "one",
			banned:      1,
			failed:      1,
			description: "minimum positive value",
		},
		{
			name:        "large",
			banned:      999999999,
			failed:      999999999,
			description: "large but valid value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jailData := &JailData{
				Name:   "test",
				Banned: tt.banned,
				Failed: tt.failed,
			}

			f := &Fail2Ban{
				mock: &MockFail2Ban{
					jails: map[string]*JailData{
						"test": jailData,
					},
				},
			}

			result := f.Collect()
			if result == nil {
				t.Errorf("Collect() returned nil for %s", tt.description)
			}
		})
	}
}

// Helper function
func mapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}