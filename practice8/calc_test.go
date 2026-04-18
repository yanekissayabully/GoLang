package main

import (
	"testing"
)

func TestSubtractTableDriven(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"both pos", 10, 3, 7},
		{"pos - zero", 5, 0, 5},
		{"neg - pos", -5, 3, -8},
		{"both neg", -7, -2, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Subtract(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("subtract(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestDivide(t *testing.T) {
	tests := []struct {
		name      string
		a, b      int
		want      int
		expectErr bool
	}{
		{"uspeh", 10, 2, 5, false},
		{"delenie na nol", 10, 0, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Divide(tt.a, tt.b)
			if tt.expectErr && err == nil {
				t.Errorf("dolzhna bit error no ne bilo")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("kakayata error: %v", err)
			}
			if got != tt.want {
				t.Errorf("div(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}