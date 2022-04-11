package util

import "testing"

func Test_getName(t *testing.T) {
	tests := []struct {
		line     string
		wantName string
		wantOk   bool
	}{
		{"! Title: some list title", "some list title", true},
		{"! title: 🦄 list", "🦄 list", true},
		{"!  abp:subscribe?location=http%3A%2F%2Fpgl.yoyo.org%2Fadservers%2Fserverlist.php%3Fhostformat%3Dadblockplus%26mimetype%3Dplaintext&title=Peter%20Lowe%27s%20list", "Peter Lowe's list", true},
		{"# Title: EasyList Germany", "EasyList Germany", true},
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

func Test_makeListTitle(t *testing.T) {
	tests := []struct {
		name    string
		wantOut string
	}{
		{"bromite-extended", "Bromite Extended"},
		{"german", "German"},
		{"here are spaces", "Here Are Spaces"},
		{"username_lists", "Username Lists"},
		{"Username-Bromite_List", "Username Bromite List"},
		{"bromite-lite", "Bromite Lite"},
		{"ublock", "uBlock"},
		{"adblock", "AdBlock"},
		{"Bromite 4pda", "Bromite 4pda"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := MakeListTitle(tt.name); gotOut != tt.wantOut {
				t.Errorf("MakeListTitle(%q) = %q, want %q", tt.name, gotOut, tt.wantOut)
			}
		})
	}
}
