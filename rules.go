package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Rule defines birth and survival conditions for cellular automata
type Rule struct {
	Birth    []int // neighbor counts that cause birth
	Survival []int // neighbor counts that allow survival
}

// ConwayRule returns the classic Conway's Game of Life rule (B3/S23)
func ConwayRule() *Rule {
	return &Rule{
		Birth:    []int{3},
		Survival: []int{2, 3},
	}
}

// ParseRule parses a rule string in B/S notation (e.g., "B3/S23", "B36/S23")
// Accepts case-insensitive input
func ParseRule(s string) (*Rule, error) {
	s = strings.ToUpper(strings.TrimSpace(s))

	// Find B and S sections
	parts := strings.Split(s, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid rule format: expected B.../S..., got %q", s)
	}

	var birthPart, survivalPart string

	for _, part := range parts {
		if strings.HasPrefix(part, "B") {
			birthPart = strings.TrimPrefix(part, "B")
		} else if strings.HasPrefix(part, "S") {
			survivalPart = strings.TrimPrefix(part, "S")
		} else {
			return nil, fmt.Errorf("invalid rule part: expected B or S prefix, got %q", part)
		}
	}

	birth, err := parseDigits(birthPart)
	if err != nil {
		return nil, fmt.Errorf("invalid birth rule: %w", err)
	}

	survival, err := parseDigits(survivalPart)
	if err != nil {
		return nil, fmt.Errorf("invalid survival rule: %w", err)
	}

	return &Rule{
		Birth:    birth,
		Survival: survival,
	}, nil
}

// parseDigits converts a string of digits (e.g., "23") into a slice of ints
func parseDigits(s string) ([]int, error) {
	if s == "" {
		return []int{}, nil
	}

	result := make([]int, 0, len(s))
	for _, ch := range s {
		n, err := strconv.Atoi(string(ch))
		if err != nil {
			return nil, fmt.Errorf("invalid digit %q", ch)
		}
		if n < 0 || n > 8 {
			return nil, fmt.Errorf("neighbor count must be 0-8, got %d", n)
		}
		result = append(result, n)
	}
	return result, nil
}

// Births returns true if a dead cell with n neighbors should become alive
func (r *Rule) Births(n int) bool {
	for _, b := range r.Birth {
		if b == n {
			return true
		}
	}
	return false
}

// Survives returns true if a live cell with n neighbors should stay alive
func (r *Rule) Survives(n int) bool {
	for _, s := range r.Survival {
		if s == n {
			return true
		}
	}
	return false
}

// String returns the rule in B/S notation
func (r *Rule) String() string {
	birthStr := ""
	for _, b := range r.Birth {
		birthStr += strconv.Itoa(b)
	}
	survivalStr := ""
	for _, s := range r.Survival {
		survivalStr += strconv.Itoa(s)
	}
	return fmt.Sprintf("B%s/S%s", birthStr, survivalStr)
}

// Common rule presets
var (
	RuleConway     = &Rule{Birth: []int{3}, Survival: []int{2, 3}}         // B3/S23 - Conway's Game of Life
	RuleHighLife   = &Rule{Birth: []int{3, 6}, Survival: []int{2, 3}}      // B36/S23 - HighLife
	RuleMaze       = &Rule{Birth: []int{3}, Survival: []int{1, 2, 3, 4, 5}} // B3/S12345 - Maze
	RuleReplicator = &Rule{Birth: []int{1, 3, 5, 7}, Survival: []int{1, 3, 5, 7}} // B1357/S1357 - Replicator
	RuleSeeds      = &Rule{Birth: []int{2}, Survival: []int{}}             // B2/S - Seeds
)
