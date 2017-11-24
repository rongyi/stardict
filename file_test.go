package stardict

import (
	"fmt"
	"testing"
)

func TestFileNew(t *testing.T) {
	i := NewInfo("./testdata/stardict-HanYuChengYuCiDian-new_colors-2.4.2/HanYuChengYuCiDian-new_colors.ifo")
	if i == nil {
		t.Fatalf("%s\n", "NewInfo fail")
	}
}

func TestIndexNew(t *testing.T) {
	idx, err := NewIndex("./testdata/stardict-langdao-ec-gb-2.4.2/langdao-ec-gb.idx")
	if err != nil {
		t.Fatalf("%s\n", "NewIndex get nil Index")
	}
	fmt.Println(len(idx.wordDict))
	fmt.Println(len(idx.wordLst))
}

func TestDictionary(t *testing.T) {
	info := NewInfo("./testdata/stardict-HanYuChengYuCiDian-new_colors-2.4.2/HanYuChengYuCiDian-new_colors.ifo")
	if info == nil {
		t.Fatalf("%s\n", "NewInfo fail")
	}
	idx, err := NewIndex("./testdata/stardict-HanYuChengYuCiDian-new_colors-2.4.2/HanYuChengYuCiDian-new_colors.idx")
	if err != nil {
		t.Fatalf("%s\n", "NewIndex get nil Index")
	}
	idx.Parse()

	d, err := NewDictionary(info, idx, "./testdata/stardict-HanYuChengYuCiDian-new_colors-2.4.2/HanYuChengYuCiDian-new_colors.dict.dz")
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
	info := NewInfo("./testdata/stardict-langdao-ec-gb-2.4.2/langdao-ec-gb.ifo")
	if info == nil {
		t.Fatalf("%s\n", "NewInfo fail")
	}
	idx, err := NewIndex("./testdata/stardict-langdao-ec-gb-2.4.2/langdao-ec-gb.idx")
	if err != nil {
		t.Fatalf("%s\n", "NewIndex get nil Index")
	}
	idx.Parse()

	d, err := NewDictionary(info, idx, "./testdata/stardict-langdao-ec-gb-2.4.2/langdao-ec-gb.dict.dz")
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
