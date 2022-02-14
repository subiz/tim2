package tim2

import (
	"fmt"
	"testing"
)

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
