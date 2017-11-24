package stardict

import (
	"fmt"
	"testing"
)

func TestFileNew(t *testing.T) {
	_, err := newInfo("./testdata/stardict-HanYuChengYuCiDian-new_colors-2.4.2/HanYuChengYuCiDian-new_colors.ifo")
	if err != nil {
		t.Fatalf("%s\n", "newInfo fail")
	}
}

func TestIndexNew(t *testing.T) {
	idx, err := newIndex("./testdata/stardict-langdao-ec-gb-2.4.2/langdao-ec-gb.idx")
	if err != nil {
		t.Fatalf("%s\n", "newIndex get nil Index")
	}
	fmt.Println(len(idx.wordDict))
	fmt.Println(len(idx.wordLst))
}

func TestDictionary(t *testing.T) {

	d, err := NewDictionary("./testdata/stardict-HanYuChengYuCiDian-new_colors-2.4.2/HanYuChengYuCiDian-new_colors.ifo",
		"./testdata/stardict-HanYuChengYuCiDian-new_colors-2.4.2/HanYuChengYuCiDian-new_colors.idx",
		"./testdata/stardict-HanYuChengYuCiDian-new_colors-2.4.2/HanYuChengYuCiDian-new_colors.dict.dz")
	if err != nil {
		t.Fatalf("%s\n", "fail to create new dictionary")
	}
	values := d.GetWord("堆金积玉")
	for _, v := range values {
		for k, m := range v {
			fmt.Println(k)
			fmt.Println(string(m))
		}
	}
}

func TestNonSequenceDictionary(t *testing.T) {

	d, err := NewDictionary("./testdata/stardict-langdao-ec-gb-2.4.2/langdao-ec-gb.ifo",
		"./testdata/stardict-langdao-ec-gb-2.4.2/langdao-ec-gb.idx",
		"./testdata/stardict-langdao-ec-gb-2.4.2/langdao-ec-gb.dict.dz")
	if err != nil {
		t.Fatalf("%s\n", "fail to create new dictionary")
	}
	values := d.GetWord("mail")
	if len(values) == 0 {
		t.Fatal("simple mail word not found, parse fail")
	}
	for _, v := range values {
		for _, m := range v {
			fmt.Println(string(m))
		}
	}
}
