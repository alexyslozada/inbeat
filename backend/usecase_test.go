package main

import "testing"

func TestUseCase_cleanUserName(t *testing.T) {
	table := []struct {
		name     string
		userName string
		want     string
	}{
		{name: "With @", userName: "@alexys_lozada", want: "alexys_lozada"},
		{name: "Without @", userName: "alexys_lozada", want: "alexys_lozada"},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			got := cleanUserName(tt.userName)
			if got != tt.want {
				t.Fatalf("Got %s, Want %s", got, tt.want)
			}
		})
	}
}
