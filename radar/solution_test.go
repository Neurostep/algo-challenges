package radar

import "testing"

func Test_Solution(t *testing.T) {
	tests := []struct {
		name string
		input []string
		expected int
	}{
		{
			name: "allow pass",
			input: []string{"CHARGE: card_country=US&currency=USD&amount=90&ip_country=CA", "ALLOW:amount<100"},
			expected: 1,
		},
		{
			name: "allow pass (equal)",
			input: []string{"CHARGE: card_country=US&currency=USD&amount=100&ip_country=CA", "ALLOW:amount<=100"},
			expected: 1,
		},
		{
			name: "allow no pass (gt)",
			input: []string{"CHARGE: card_country=US&currency=USD&amount=100&ip_country=CA", "ALLOW:amount>100"},
			expected: 0,
		},
		{
			name: "allow pass (equal, gte)",
			input: []string{"CHARGE: card_country=US&currency=USD&amount=100&ip_country=CA", "ALLOW:amount>=100"},
			expected: 1,
		},
		{
			name: "allow no pass",
			input: []string{"CHARGE: card_country=US&currency=USD&amount=110&ip_country=CA", "ALLOW:amount<100"},
			expected: 0,
		},
		{
			name: "block pass",
			input: []string{"CHARGE: card_country=US&currency=USD&amount=90&ip_country=CA", "BLOCK:amount > 100"},
			expected: 1,
		},
		{
			name: "block no pass",
			input: []string{"CHARGE: card_country=US&currency=USD&amount=90&ip_country=CA", "BLOCK:card_country != ip_country"},
			expected: 0,
		},
		{
			name: "complex",
			input: []string{"CHARGE: card_country=US&currency=USD&amount=90&ip_country=CA", "ALLOW:amount<100", "BLOCK:card_country != ip_country AND amount > 100"},
			expected: 1,
		},
		{
			name: "complex (OR)",
			input: []string{"CHARGE: card_country=US&currency=USD&amount=90&ip_country=CA", "ALLOW:amount<100", "BLOCK:card_country != ip_country OR amount > 100"},
			expected: 0,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			result := Solution(tt.input)

			if result != tt.expected {
				t.Errorf("expected %d but got %d instead", tt.expected, result)
			}
		})
	}
}
