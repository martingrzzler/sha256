package main

import "testing"

func TestNumPaddZero(t *testing.T) {
	testCases := []struct {
		L     int
		Zeros int
	}{{55, 7}, {16, 319}, {154, 239}, {5, 407}}

	for _, tc := range testCases {
		if tc.Zeros != NumPaddZero(tc.L) {
			t.Errorf("For NumPaddZero(%d), Expected %d, Actual %d", tc.L, tc.Zeros, NumPaddZero(tc.L))
		}
	}
}
