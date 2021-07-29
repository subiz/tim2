package tim

import (
	"fmt"
	"strings"
	"testing"
)

func TestLongQuery(t *testing.T) {
	// interms := Tokenize("cong hoa xahoi chu hi")
	interms := Tokenize("cộng hòa xã    \thội\t chủ nghĩa Việt Nam.độc \nlập")
	fmt.Printf("%#v\n", interms)
	var terms []string
	if len(interms) > 5 {
		biwords := make([]string, 0)
		for _, term := range interms {
			if strings.Contains(term, " ") {
				biwords = append(biwords, term)
			}
		}
		terms = make([]string, 0)
		for i := 0; i < 2 && i < len(biwords); i++ {
			terms = append(terms, biwords[i])
		}
		if len(terms) < 2 {
			for i := 0; i < 4-len(terms); i++ {
				terms = append(terms, interms[i])
			}
		}
	} else {
		terms = interms
	}
	fmt.Println(terms)
	t.Error("TRUE")
}

func TestTokenize(t *testing.T) {
	testCases := []struct {
		in  string
		out string
	}{
		{"cộng hoa xã hội chủ nghĩa Việt. Nam", "cong hoa xa hoi nam chu nghia viet cong-hoa xa-hoi hoi-chu hoa-xa chu-nghia nghia-viet viet-nam"},
		{`phạm kiều
 thanh`, "pham kieu thanh pham-kieu kieu-thanh"},
	}
	for _, tc := range testCases {
		out := Tokenize(tc.in)
		if len(out) != len(strings.Split(tc.out, " ")) {
			fmt.Println(out)
			t.Errorf("Len should be eq for tc %s", tc.in)
		}
		outM := map[string]bool{}
		for _, term := range out {
			outM[term] = true
		}

		for _, term := range strings.Split(tc.out, " ") {
			if !outM[term] {
				fmt.Println("OUT FOR ", tc.in, out)
				t.Errorf("MISSING term %s for tc %s", term, tc.in)
			}
		}
	}
}
