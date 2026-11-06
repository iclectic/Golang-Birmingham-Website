package rendering

import (
	"strings"
	"testing"
)

func TestInitials(t *testing.T) {
	for in, want := range map[string]string{"Sarah Chen": "SC", "Ibim": "I", "Mary Ann de Silva": "MS", "": ""} {
		if got := Initials(in); got != want {
			t.Errorf("Initials(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestParagraphs(t *testing.T) {
	got := string(Paragraphs("Intro line\ncontinues.\n\nTopics:\n- One <b>\n- Two\n\nEnd."))
	want := "<p>Intro line continues.</p><p>Topics:</p><ul><li>One &lt;b&gt;</li><li>Two</li></ul><p>End.</p>"
	if got != want {
		t.Errorf("Paragraphs:\n got %s\nwant %s", got, want)
	}
	if strings.Contains(got, "<b>") {
		t.Error("HTML was not escaped")
	}
}
