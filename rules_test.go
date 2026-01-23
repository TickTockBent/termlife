package main

import (
	"reflect"
	"testing"
)

func TestParseRuleConway(t *testing.T) {
	rule, err := ParseRule("B3/S23")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(rule.Birth, []int{3}) {
		t.Errorf("expected birth [3], got %v", rule.Birth)
	}
	if !reflect.DeepEqual(rule.Survival, []int{2, 3}) {
		t.Errorf("expected survival [2,3], got %v", rule.Survival)
	}
}

func TestParseRuleHighLife(t *testing.T) {
	rule, err := ParseRule("B36/S23")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(rule.Birth, []int{3, 6}) {
		t.Errorf("expected birth [3,6], got %v", rule.Birth)
	}
	if !reflect.DeepEqual(rule.Survival, []int{2, 3}) {
		t.Errorf("expected survival [2,3], got %v", rule.Survival)
	}
}

func TestParseRuleCaseInsensitive(t *testing.T) {
	testCases := []string{"b3/s23", "B3/S23", "b3/S23", "B3/s23"}

	for _, tc := range testCases {
		rule, err := ParseRule(tc)
		if err != nil {
			t.Errorf("ParseRule(%q) unexpected error: %v", tc, err)
			continue
		}
		if !reflect.DeepEqual(rule.Birth, []int{3}) {
			t.Errorf("ParseRule(%q) birth: expected [3], got %v", tc, rule.Birth)
		}
	}
}

func TestParseRuleEmptySurvival(t *testing.T) {
	// Seeds rule: B2/S
	rule, err := ParseRule("B2/S")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(rule.Birth, []int{2}) {
		t.Errorf("expected birth [2], got %v", rule.Birth)
	}
	if len(rule.Survival) != 0 {
		t.Errorf("expected empty survival, got %v", rule.Survival)
	}
}

func TestParseRuleInvalidFormat(t *testing.T) {
	testCases := []string{
		"B3S23",      // Missing /
		"3/23",       // Missing B and S
		"B3/S23/X1",  // Too many parts
		"",           // Empty
	}

	for _, tc := range testCases {
		_, err := ParseRule(tc)
		if err == nil {
			t.Errorf("ParseRule(%q) expected error, got nil", tc)
		}
	}
}

func TestParseRuleInvalidDigit(t *testing.T) {
	testCases := []string{
		"B9/S23", // 9 is invalid (max neighbors is 8)
		"Ba/S23", // 'a' is not a digit
	}

	for _, tc := range testCases {
		_, err := ParseRule(tc)
		if err == nil {
			t.Errorf("ParseRule(%q) expected error, got nil", tc)
		}
	}
}

func TestRuleBirths(t *testing.T) {
	rule := ConwayRule()

	testCases := []struct {
		neighbors int
		expected  bool
	}{
		{0, false},
		{1, false},
		{2, false},
		{3, true}, // Birth on 3
		{4, false},
		{5, false},
		{6, false},
		{7, false},
		{8, false},
	}

	for _, tc := range testCases {
		if result := rule.Births(tc.neighbors); result != tc.expected {
			t.Errorf("Births(%d): expected %v, got %v", tc.neighbors, tc.expected, result)
		}
	}
}

func TestRuleSurvives(t *testing.T) {
	rule := ConwayRule()

	testCases := []struct {
		neighbors int
		expected  bool
	}{
		{0, false},
		{1, false},
		{2, true}, // Survive on 2
		{3, true}, // Survive on 3
		{4, false},
		{5, false},
		{6, false},
		{7, false},
		{8, false},
	}

	for _, tc := range testCases {
		if result := rule.Survives(tc.neighbors); result != tc.expected {
			t.Errorf("Survives(%d): expected %v, got %v", tc.neighbors, tc.expected, result)
		}
	}
}

func TestRuleString(t *testing.T) {
	testCases := []struct {
		rule     *Rule
		expected string
	}{
		{ConwayRule(), "B3/S23"},
		{RuleHighLife, "B36/S23"},
		{RuleSeeds, "B2/S"},
	}

	for _, tc := range testCases {
		if result := tc.rule.String(); result != tc.expected {
			t.Errorf("expected %q, got %q", tc.expected, result)
		}
	}
}

func TestConwayRule(t *testing.T) {
	rule := ConwayRule()

	if !reflect.DeepEqual(rule.Birth, []int{3}) {
		t.Errorf("ConwayRule birth: expected [3], got %v", rule.Birth)
	}
	if !reflect.DeepEqual(rule.Survival, []int{2, 3}) {
		t.Errorf("ConwayRule survival: expected [2,3], got %v", rule.Survival)
	}
}
