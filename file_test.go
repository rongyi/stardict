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

func TestUnhtml(t *testing.T) {
	input := []byte(`<font size=5 color=black>堆金积玉<br></font><font size=3 color=green>duī  jīn  jī  yù<br> <br></font><font size=3 color=blue>【解释】金玉多得可以堆积起来。形容聚敛的财富极多。<br></font><font size=3 color=black>【出处】唐·李贺《嘲少年》诗：“长金积玉夸豪毅。”<br></font><font size=3 color=brown>【示例】<br></font><font size=3 color=gray>【拼音码】djjy<br></font><font size=3 color=blue>【近义词】堆金叠玉、腰缠万贯<br></font><font size=3 color=red>【反义词】<br></font><font size=3 color=black>【歇后语】<br></font><font size=3 color=lightgrey>【灯谜面】<br></font><font size=3 color=green>【用法】联合式；作谓语；形容财富充裕<br></font><font size=3 color=purple>【英文】amass a fortune<br></font><font size=3 color=black>【故事】</font>`)
	txt, err := Unhtml(input)
	if err != nil {
		t.Fatal("test fail of unhtml")
	}

	fmt.Println(txt)
}
