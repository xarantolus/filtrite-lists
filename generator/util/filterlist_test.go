package util

import "testing"

func Test_getName(t *testing.T) {
	tests := []struct {
		line     string
		wantName string
		wantOk   bool
	}{
		{"Title: some list title", "some list title", true},
		{"title: 🦄 list", "🦄 list", true},
		{"  abp:subscribe?location=http%3A%2F%2Fpgl.yoyo.org%2Fadservers%2Fserverlist.php%3Fhostformat%3Dadblockplus%26mimetype%3Dplaintext&title=Peter%20Lowe%27s%20list", "Peter Lowe's list", true},
	}
	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			gotName, gotOk := getName(tt.line)
			if gotName != tt.wantName {
				t.Errorf("getName(%q) gotName = %v, want %v", tt.line, gotName, tt.wantName)
			}
			if gotOk != tt.wantOk {
				t.Errorf("getName(%q) gotOk = %v, want %v", tt.line, gotOk, tt.wantOk)
			}
		})
	}
}
