package tim2

import (
	"fmt"
	"strings"
	"testing"
)

func TestShortQuery(t *testing.T) {
	// interms := Tokenize("cong hoa xahoi chu hi")
	// interms := Tokenize("cộng hòa xã    \thội\t chủ nghĩa Việt Nam.độc \nlập")
	interms := Tokenize("Launch HN: Fogbender (YC W22) – B2B support software designed for customer teams (fogbender.com)")
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
		{"thang loi.", "thang loi thang-loi"},
		{`Georg Cantor (3 tháng 3 [lịch cũ 19 tháng m2] năm 1845 – 6 tháng 1 năm 1918) là một nhà toán học người Đức, được biết đến nhiều nhất với tư cách cha
đẻ của lý thuyết tập hợp, một lý thuyết đã trở thành một lý thuyết nền tảng trong toán học. Cantor đã cho thấy tầm quan trọng của quan hệ song ánh giữa
các phần tử của hai tập hợp, định nghĩa các tập vô hạn và các tập sắp tốt, và chứng minh rằng các số thực là "đông đúc" hơn các số tự nhiên.
Trên thực tế, phương pháp chứng minh định lý này của Cantor ngụ ý sự tồn tại "vô hạn các tập vô hạn".
Ông định nghĩa bản số và số thứ tự và phép tính về chúng.
Sự nghiệp toán học vĩ đại của ông nhận được sự quan tâm lớn về mặt triết học,
nhờ đó khiến ông càng được biết đến nhiều hơn.`,
			"thang-3 duc hai-tap cach han-va han-cac nghiep-toan nghia-cac minh nhien-tren su-ton m2-nam toan-hoc thanh-mot tam ve-chung tu-nhien tinh vi-dai hop thuc-te lon hoc-nho duoc thay-tam hop-dinh la-dong cua-ong ngu dai 1918-la cua-ly ong-nhan cach-cha ly-thuyet sap-tot lon-ve nguoi-duc tap-vo cua-cantor phan phan-tu song-anh han 1845-6 va-so nho-do la-mot tu-cach va-phep nguoi trong phap-chung ve-mat 1845 den tap song thu-tu tro-thanh dinh georg lich cha triet khien-ong tang cua-quan so-tu lich-cu han-ong nghiep mot-ly quan-trong cang-duoc nhat-voi khien da-tro hai cac-so tinh-ve quan-tam nha thay cantor-da cho ban-so den-nhieu quan thuc hon-cac mot-nha hoc-cantor he-song nhan mat sap rang rang-cac so-va te-phuong phep-tinh nen-tang do-khien nhieu tu-cua hon cantor-3 thang nha-toan 19-thang tap-hop cantor thuyet-nen giua va-chung nay-cua nghia-ban voi trong-toan phep su-nghiep thanh nhieu-hon hoc thuyet-tap phuong-phap tu-va chung-su cang thang-m2 duoc-biet biet-den so-thuc tren duoc-su tam-lon nen chung-minh phap cantor-ngu nam mot thuyet tot tren-thuc ngu-y tai-vo toan cha-de tang-trong giua-cac cac-tap chung minh-rang thuc-la dinh-ly ly-nay dai-cua thuyet-da anh-giua dinh-nghia ton thu hoc-vi triet-hoc nhieu-nhat tot-va dong dong-duc nhien ton-tai ban vo-han 1918 hoc-nguoi de-cua trong-cua quan-he duc-duoc ong cu-19 cho-thay ong-cang tam-quan cac-phan minh-dinh su-quan nho nhat tro va-cac duc-hon ong-dinh nhan-duoc nam-1845 biet nghia georg-cantor cua-hai phuong voi-tu hop-mot tap-sap tai so-thu mat-triet nam-1918 da-cho thang-1"},
		{`可口的 1 _ - 123 accompaniment nnn Trụ sở: (Tầng 6), tòa nhà Kailash,
ngõ 92 Trần Thái Tông, Di3u@gmail.com Phường Dịch Vọng Hậu, Quận Cầu Giấy, Hà Nội (84)123123211 dieu
https://translate.google.com/?hl=vi&sl=en&tl=vi&text=concrete&op=translate
pneumonoultramicroscopicsilicovolcanoconiosis viet nam. hello viet.nam
va co gmail com`, "dich viet-nam 84123123211 nnn toa-nha thai-tong nam so-tang ngo dieu-https translate concrete-op va-co gmail-com tang 92-tran cau ha-noi vong hau-quan quan cau-giay translate.google.com nha tong di3u-gmail.com noi-84 nam-hello noi hl-vi concrete di3u@gmail.com tang-6 nha-kailash kailash vi-sl text-concrete viet tru tru-so tong-di3u https toa text nnn-tru quan-cau hello 84-123123211 tl-vi viet.nam-va phuong hau giay giay-ha thai gmail.com-phuong en-tl 123 kailash-ngo tran-thai phuong-dich accompaniment ngo-92 tran hello-viet.nam viet.nam 123123211-dieu dieu sl-en dich-vong vi-text op-translate di3u vong-hau 123123211 co-gmail pneumonoultramicroscopicsilicovolcanoconiosis"},
		{"a - b là số gì, nhóm máu axfO và AB+ 2020 ;-; 24", "nhom nhom-mau mau-axfo axfo va-ab ab-2020 2020 la-so so-gi gi-nhom mau axfo-va"},
		{"tôi muốn truy cập trang google.com nhưng không vào dc", "toi muon-truy trang nhung vao-dc khong google.com muon truy truy-cap cap-trang cap nhung-khong khong-vao vao toi-muon"},
		{"cộng hoa xã hội chủ nghĩa Việt. Nam", "cong hoa hoi nam chu nghia viet cong-hoa xa-hoi hoi-chu hoa-xa chu-nghia nghia-viet viet-nam"},
		{`phạm kiều
 thanh`, "pham kieu thanh pham-kieu kieu-thanh"},
	}
	for _, tc := range testCases {
		out := Tokenize(tc.in)
		if len(out) != len(strings.Split(tc.out, " ")) {
			t.Errorf("Len should be eq expect %d, got %d", len(out), len(strings.Split(tc.out, " ")))
		}
		outM := map[string]bool{}
		for _, term := range out {
			outM[term] = true
		}

		for _, term := range strings.Split(tc.out, " ") {
			term = strings.Replace(term, "-", " ", -1)
			if !outM[term] {
				t.Errorf("MISSING term [%s] for tc %d", term, len(term))
			}
		}

		for k := range outM {
			found := false
			for _, term := range strings.Split(tc.out, " ") {
				term = strings.Replace(term, "-", " ", -1)
				if term == k {
					found = true
					break
				}
			}

			if !found {
				t.Errorf("REDURANT term [%s]", k)
			}
		}
	}
}

func TestFindPersonalPhoneNumber(t *testing.T) {
	str := "1 _ - 123 accompaniment hi, 123 phone number is +(84) 0974-304-123., 0974403666 end. +84 2473.021.368"
	phones := findPersonalPhoneNumber(str)

	out := []string{"840974304123", "0974403666", "842473021368"}

	if len(phones) != len(out) {
		t.Errorf("Len should be eq expect %d, got %d", len(out), len(phones))
	}

	for i := range out {
		if out[i] != phones[i] {
			t.Errorf("Should be eq expect %s, got %s", out[i], phones[i])
		}
	}
}
