package main

import "testing"

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
			if gotOut := makeListTitle(tt.name); gotOut != tt.wantOut {
				t.Errorf("makeListTitle(%q) = %q, want %q", tt.name, gotOut, tt.wantOut)
			}
		})
	}
}
