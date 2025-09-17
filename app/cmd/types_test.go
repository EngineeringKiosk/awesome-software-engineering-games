package cmd

import (
	"testing"
)

func TestGameInformation_GetReleaseDate(t *testing.T) {
	tests := []struct {
		name string
		date string
		want string
	}{
		{
			name: "Valid date 1",
			date: "Aug 23, 2018",
			want: "Thursday, 23 August 2018",
		},
		{
			name: "Valid date 2",
			date: "Oct 25, 2006",
			want: "Wednesday, 25 October 2006",
		},
		{
			name: "Invalid date",
			date: "not a date",
			want: "",
		},
		{
			name: "Empty date",
			date: "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := GameInformation{
				ReleaseDate: ReleaseDate{
					Date: tt.date,
				},
			}
			got := g.GetReleaseDate()
			if got != tt.want {
				t.Errorf("GetReleaseDate() = %q, want %q", got, tt.want)
			}
		})
	}
}
