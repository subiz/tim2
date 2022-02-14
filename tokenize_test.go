package tim2

import (
	"fmt"
	"testing"
	"strings"
)

func TestShortQuery(t *testing.T) {
	// interms := Tokenize("cong hoa xahoi chu hi")
	// interms := Tokenize("cộng hòa xã    \thội\t chủ nghĩa Việt Nam.độc \nlập")
	interms := Tokenize("vào -  thanh + hoa")
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

func TestLongQuery(t *testing.T) {
	// interms := Tokenize("cong hoa xahoi chu hi")
	// interms := Tokenize("cộng hòa xã    \thội\t chủ nghĩa Việt Nam.độc \nlập")
	interms := Tokenize("vào file xyz.txt là đc<script>console.log('1')</script> lasjdfl;k asjlkfdj alksjfdlkasj dkldsjalkfjlaksdjflkasjdlkfjaslkfjlkalskdjflkajsflkjasdfasdjf asldkfj")
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

func TestTokenizeWiki(t *testing.T) {
	str := `Georg Cantor (3 tháng 3 [lịch cũ 19 tháng 2] năm 1845 – 6 tháng 1 năm 1918) là một nhà toán học người Đức, được biết đến nhiều nhất với tư cách cha đẻ của lý thuyết tập hợp, một lý thuyết đã trở thành một lý thuyết nền tảng trong toán học. Cantor đã cho thấy tầm quan trọng của quan hệ song ánh giữa các phần tử của hai tập hợp, định nghĩa các tập vô hạn và các tập sắp tốt, và chứng minh rằng các số thực là "đông đúc" hơn các số tự nhiên. Trên thực tế, phương pháp chứng minh định lý này của Cantor ngụ ý sự tồn tại "vô hạn các tập vô hạn". Ông định nghĩa bản số và số thứ tự và phép tính về chúng. Sự nghiệp toán học vĩ đại của ông nhận được sự quan tâm lớn về mặt triết học, nhờ đó khiến ông càng được biết đến nhiều hơn.`
	literals := Tokenize(str)

	fmt.Printf("%#v", literals)
}

func TestTokenizeLiteral(t *testing.T) {
	str := `可口的 1 _ - 123 accompaniment nnn Trụ sở: (Tầng 6), tòa nhà Kailash,
		ngõ 92 Trần Thái Tông, Di3u@gmail.com Phường Dịch Vọng Hậu, Quận Cầu Giấy, Hà Nội (84)123123211 dieu
		https://translate.google.com/?hl=vi&sl=en&tl=vi&text=concrete&op=translate
		pneumonoultramicroscopicsilicovolcanoconiosis viet nam. hello viet.nam
		va co gmail com`
	literals := Tokenize(str)

	fmt.Printf("%#v", literals)
	return
	fmt.Printf("%#v", Tokenize(`cong, hoa, xa, hoi, 'chu nghia" viet "`))

	fmt.Println("")
	fmt.Printf("%#v", Tokenize("cong hoa xa"))
	fmt.Printf("%#v", Tokenize("cong hoa xa hoi chu nghia viet nam"))
	if true {
		t.Error("TRUE")
	}
}

func TestFindPersonalPhoneNumber(t *testing.T) {
	t.Skip()
	str := "1 _ - 123 accompaniment hi, 123 phone number is +(84) 0974-304-123., 0974403666 end. +84 2473.021.368"
	phones := findPersonalPhoneNumber(str)
	fmt.Printf("%#v", phones)
	t.Error("TRUE")
}
