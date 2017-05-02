package stardict

import (
	"fmt"
	"testing"
)

func TestFileNew(t *testing.T) {
	i := NewInfo("/tmp/stardict-HanYuChengYuCiDian-new_colors-2.4.2/HanYuChengYuCiDian-new_colors.ifo")
	if i == nil {
		t.Fatalf("%s\n", "NewInfo fail")
	}
	// fmt.Println(i)
}

func TestIndexNew(t *testing.T) {
	idx, err := NewIndex("/tmp/stardict-HanYuChengYuCiDian-new_colors-2.4.2/HanYuChengYuCiDian-new_colors.idx")
	if err != nil {
		t.Fatalf("%s\n", "NewIndex get nil Index")
	}
	for w, err := idx.NextWord(); err == nil ; w, err = idx.NextWord() {
		fmt.Println(w)
	}
	fmt.Println(len(idx.wordDict))
	fmt.Println(len(idx.wordLst))
	// input := "堆金积玉"
	// w, ok := idx.wordDict[input]
	// if !ok {
	// 	t.Fatalf("%s\n", "fail to get word: 堆金积玉")
	// }
	// fmt.Println(w)
}

func TestDictionary(t *testing.T) {
	info := NewInfo("/tmp/stardict-HanYuChengYuCiDian-new_colors-2.4.2/HanYuChengYuCiDian-new_colors.ifo")
	if info == nil {
		t.Fatalf("%s\n", "NewInfo fail")
	}
	idx, err := NewIndex("/tmp/stardict-HanYuChengYuCiDian-new_colors-2.4.2/HanYuChengYuCiDian-new_colors.idx")
	if err != nil {
		t.Fatalf("%s\n", "NewIndex get nil Index")
	}

	_, err = NewDictionary(info, idx, "/tmp/stardict-HanYuChengYuCiDian-new_colors-2.4.2/HanYuChengYuCiDian-new_colors.dict.dz")
	if err != nil {
		t.Fatalf("%s\n", "fail to create new dictionary")
	}
}
