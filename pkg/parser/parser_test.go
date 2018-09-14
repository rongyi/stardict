package parser

import (
	"fmt"
	"testing"
	"os"
)

func TestFileNew(t *testing.T) {
	f, err := os.Open("./testdata/stardict-HanYuChengYuCiDian-new_colors-2.4.2/HanYuChengYuCiDian-new_colors.ifo")
	if err != nil {
		t.Fatalf("%s\n", "newInfo fail")
	}
	defer f.Close()
	_, err = newInfo(f)
	if err != nil {
		t.Fatalf("%s\n", "newInfo fail")
	}
}

func TestIndexNew(t *testing.T) {
	f, err := os.Open("./testdata/stardict-langdao-ec-gb-2.4.2/langdao-ec-gb.idx")
	if err != nil {
		t.Fatalf("%s\n", "newIndex get nil Index")
	}
	defer f.Close()
	idx, err := newIndex(f)
	if err != nil {
		t.Fatalf("%s\n", "newIndex get nil Index")
	}
	fmt.Println(len(idx.wordDict))
	fmt.Println(len(idx.wordLst))
}

func TestDictionary(t *testing.T) {
	f1, err := os.Open("./testdata/stardict-HanYuChengYuCiDian-new_colors-2.4.2/HanYuChengYuCiDian-new_colors.ifo")
	if err != nil {
		t.Fatalf("%s\n", "fail to create new dictionary")
	}
	defer f1.Close()
	f2, err := os.Open("./testdata/stardict-HanYuChengYuCiDian-new_colors-2.4.2/HanYuChengYuCiDian-new_colors.idx")
	if err != nil {
		t.Fatalf("%s\n", "fail to create new dictionary")
	}
	defer f2.Close()
	f3, err := os.Open("./testdata/stardict-HanYuChengYuCiDian-new_colors-2.4.2/HanYuChengYuCiDian-new_colors.dict.dz")
	if err != nil {
		t.Fatalf("%s\n", "fail to create new dictionary")
	}
	defer f3.Close()
	d, err := NewDictionary(f1, f2, f3)
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
	vs := d.GetFormatedMeaning("堆金积玉")
	for _, v := range vs {
		fmt.Println(v)
	}
}

func TestNonSequenceDictionary(t *testing.T) {
	f1, err := os.Open("./testdata/stardict-langdao-ec-gb-2.4.2/langdao-ec-gb.ifo")
	if err != nil {
		t.Fatalf("%s\n", "fail to create new dictionary")
	}
	defer f1.Close()

	f2, err := os.Open("./testdata/stardict-langdao-ec-gb-2.4.2/langdao-ec-gb.idx")
	if err != nil {
		t.Fatalf("%s\n", "fail to create new dictionary")
	}
	defer f2.Close()

	f3, err := os.Open("./testdata/stardict-langdao-ec-gb-2.4.2/langdao-ec-gb.dict.dz")
	if err != nil {
		t.Fatalf("%s\n", "fail to create new dictionary")
	}
	defer f3.Close()

	d, err := NewDictionary(f1, f2, f3)
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

func TestUnhtml(t *testing.T) {
	input := []byte(`<font size=5 color=black>堆金积玉<br></font><font size=3 color=green>duī  jīn  jī  yù<br> <br></font><font size=3 color=blue>【解释】金玉多得可以堆积起来。形容聚敛的财富极多。<br></font><font size=3 color=black>【出处】唐·李贺《嘲少年》诗：“长金积玉夸豪毅。”<br></font><font size=3 color=brown>【示例】<br></font><font size=3 color=gray>【拼音码】djjy<br></font><font size=3 color=blue>【近义词】堆金叠玉、腰缠万贯<br></font><font size=3 color=red>【反义词】<br></font><font size=3 color=black>【歇后语】<br></font><font size=3 color=lightgrey>【灯谜面】<br></font><font size=3 color=green>【用法】联合式；作谓语；形容财富充裕<br></font><font size=3 color=purple>【英文】amass a fortune<br></font><font size=3 color=black>【故事】</font>`)
	txt, err := Unhtml(input)
	if err != nil {
		t.Fatal("test fail of unhtml")
	}

	fmt.Println(txt)
}