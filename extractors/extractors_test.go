package extractors

import "testing"

func TestExtractors(t *testing.T) {

	nhentai := NhentaiExtractor()
	testers := []Tester{nhentai}
	for _, tester := range testers {
		tester.Test(t)
	}
}
