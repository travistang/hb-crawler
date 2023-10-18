package hiking_buddies

import "testing"

func TestParsePointString(t *testing.T) {
	got, _ := parsePointString("789 Points")
	want := 789

	if *got != want {
		t.Errorf("got %d, wanted %d", *got, want)
	}
}
