package stardict

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

var (
	errorEOF  = errors.New("end of index, no next word")
	errorBits = errors.New("offset bits only support 64 or 32")
	errorRead = errors.New("read file error")
	errorGzip = errors.New("gunzip fail")
)

// Info indicate the stardict ifo file
type Info struct {
	File    string
	Content string
	Dict    map[string]string
}

// NewInfo create a new Info struct
func NewInfo(dirname string) *Info {
	i := &Info{
		File: dirname,
	}
	c, err := ioutil.ReadFile(i.File)
	if err != nil {
		return nil
	}
	i.Content = string(c)
	lines := strings.Split(i.Content, "\n")
	i.Dict = make(map[string]string)
	if len(lines) < 1 {
		return nil
	}
	i.Dict["header"] = lines[0]
	for _, l := range lines[1:] {
		if l == "" {
			continue
		}
		secs := strings.Split(l, "=")
		if len(secs) != 2 {
			return nil
		}
		key := strings.Trim(secs[0], "\n ")
		value := strings.Trim(secs[1], "\n ")
		i.Dict[key] = value
	}

	return i
}

func (i *Info) String() string {
	var ret []string
	for key := range i.Dict {
		cur := fmt.Sprintf("%s: %s", key, i.Dict[key])
		ret = append(ret, cur)
	}

	return strings.Join(ret, "\n")
}

type Word struct {
	w      string // the word
	offset uint32 // start position
	size   uint32 // size
	index  uint32 // index serial number
}

type Index struct {
	content   []byte
	offset    int
	index     uint32
	indexBits uint32
	wordDict  map[string][]Word
	wordLst   []Word
}

func NewIndex(filename string) (*Index, error) {
	c, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	idx := &Index{
		content:   c,
		offset:    0,
		index:     0,
		indexBits: 32,
		wordDict:  make(map[string][]Word),
	}

	return idx, nil
}

func (idx *Index) NextWord() (string, error) {
	if idx.offset == len(idx.content) {
		return "", errorEOF
	}
	// format:
	// word_str;  // a utf-8 string terminated by '\0'.
	// word_data_offset;  // word data's offset in .dict file
	// word_data_size;  // word data's total size in .dict file
	end := bytes.IndexByte(idx.content[idx.offset:], '\000')
	end += idx.offset
	wordStr := string(idx.content[idx.offset:end])
	fmt.Println(wordStr)

	newWord := Word{
		w: wordStr,
	}

	idx.offset = end + 1
	if idx.indexBits == 64 {
		var wOffset uint64
		offByte := idx.content[idx.offset : idx.offset+8]
		r := bytes.NewReader(offByte)
		binary.Read(r, binary.BigEndian, &wOffset)
		idx.offset += 8
		newWord.offset = uint32(wOffset)
	} else if idx.indexBits == 32 {
		var wOffset uint32
		offByte := idx.content[idx.offset : idx.offset+4]
		r := bytes.NewReader(offByte)
		binary.Read(r, binary.BigEndian, &wOffset)
		idx.offset += 4
		fmt.Printf("offset reading: %d\n", wOffset)
		newWord.offset = wOffset
	} else {
		return "", errorBits
	}
	var wSize uint32
	sizeByte := idx.content[idx.offset : idx.offset+4]
	r := bytes.NewReader(sizeByte)
	binary.Read(r, binary.BigEndian, &wSize)
	fmt.Printf("size reading: %d\n", wSize)
	newWord.size = wSize

	idx.offset += 4

	newWord.index = idx.index
	idx.index++

	// update the cache
	idx.wordLst = append(idx.wordLst, newWord)
	idx.wordDict[wordStr] = append(idx.wordDict[wordStr], newWord)

	return wordStr, nil
}

type Dictionary struct {
	info  *Info
	index *Index
	content []byte
	offset uint32
}


func NewDictionary(i *Info, idx *Index, filename string) (*Dictionary, error) {
	d := &Dictionary{
		info: i,
		index: idx,
		offset: 0,
	}
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errorRead
	}
	content, err := Gunzip(raw)
	if err != nil {
		return nil, errorGzip
	}
	d.content = content

	return d, nil
}
