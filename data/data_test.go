package data

import (
	"testing"
)

func TestIsValidData(t *testing.T) {
	// 97 to 122
	testdata := []RelativeFrequency{
		{123: [][][]int{{{10, 10, 10, 10}}, {{10, 10, 10, 10, 97}}}},
		{96: [][][]int{{{10, 10, 10, 10}}, {{10, 10, 10, 10, 97}}}},
		{97: [][][]int{{{10, 10, 10}}, {{10, 10, 10, 10, 97}}}},
		{97: [][][]int{{{10, 10, 10, 10}}, {{10, 10, 10, 10}}}},
		{97: [][][]int{{{10, 10, 10, 10}}, {{100, 10, 10, 10, 97}}}},
		{97: [][][]int{{{10, 10, 10, 10}}, {{10, 100, 10, 10, 97}}}},
		{97: [][][]int{{{10, 10, 10, 10}}, {{10, 10, 100, 10, 97}}}},
		{97: [][][]int{{{10, 10, 10, 10}}, {{10, 10, 10, 100, 97}}}},
	}

	for i, v := range testdata {
		if err := IsValidData(v); err == nil {
			t.Errorf("\n%d: should be an error occur\n", i)
		}
	}
}

func TestData(t *testing.T) {
	err := IsValidData(Data)
	if err != nil {
		t.Errorf("\nerror: %v\n", err)
	}
}
