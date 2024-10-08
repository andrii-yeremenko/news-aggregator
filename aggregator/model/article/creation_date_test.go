package article_test

import (
	"news-aggregator/aggregator/model/article"
	"testing"

	"time"
)

func TestCreationDate_HumanReadableString(t *testing.T) {
	cd := article.CreationDate(time.Date(2024, time.June, 5, 10, 0, 0, 0, time.UTC))

	expectedOutput := "05 Jun 24 10:00 UTC"

	output := cd.HumanReadableString()

	if output != expectedOutput {
		t.Errorf("HumanReadableString() output is incorrect, got: %s, want: %s", output, expectedOutput)
	}
}
