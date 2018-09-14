package parser

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"sync"
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

// on fail return nil
func newInfo(ir io.Reader) (*Info, error) {
	i := &Info{
	// File: dirname,
	}
	c, err := ioutil.ReadAll(ir)
	if err != nil {
		return nil, err
	}
	i.Content = string(c)
	d, err := parseInfo(i.Content)
	if err != nil {
		return nil, err
	}
	i.Dict = d

	return i, nil
}

func parseInfo(content string) (map[string]string, error) {
	lines := strings.Split(content, "\n")
	if len(lines) < 1 {
		return nil, errors.New("content empty")
	}

	ret := make(map[string]string)
	ret["header"] = lines[0]
	for _, l := range lines[1:] {
		if l == "" {
			continue
		}
		secs := strings.SplitN(l, "=", 2)
		if len(secs) != 2 {
			return nil, errors.New("key value pair fail")
		}
		key := strings.Trim(secs[0], "\n ")
		value := strings.Trim(secs[1], "\n ")
		ret[key] = value
	}
	return ret, nil
}

func (i *Info) String() string {
	var ret []string
	for key := range i.Dict {
		cur := fmt.Sprintf("%s: %s", key, i.Dict[key])
		ret = append(ret, cur)
	}

	return strings.Join(ret, "\n")
}

// Word represent the dictionary unit: word
type Word struct {
	w      string // the word to be searched
	offset uint32 // start position at dic file
	size   uint32 // len
	index  uint32 // index serial number
}

// Index reprent the idx file
type Index struct {
	content   []byte
	offset    int
	index     uint32
	indexBits uint32
	wordDict  map[string][]Word
	wordLst   []Word
	parsed    bool
}

// newIndex create a new Index with idx file
func newIndex(ir io.Reader) (*Index, error) {
	c, err := ioutil.ReadAll(ir)
	if err != nil {
		return nil, err
	}

	// from doc:
	// If the version is "3.0.0" and "idxoffsetbits=64", word_data_offset will
	// be 64-bits unsigned number in network byte order. Otherwise it will be
	// 32-bits.
	// word_data_size should be 32-bits unsigned number in network byte order.

	// The dictionary I downloaded is all version 2.4, so here the indexBits is
	// hardcoded to 32, If you need to parse 3.0 or higher, read the documentation
	// above, and rewrite.
	idx := &Index{
		content:   c,
		offset:    0,
		index:     0,
		indexBits: 32,
		wordDict:  make(map[string][]Word),
	}
	idx.parse()

	return idx, nil
}

func (idx *Index) nextWord() (string, error) {
	if idx.offset == len(idx.content) {
		return "", errorEOF
	}
	// In order to make StarDict work on different platforms, these numbers
	// must be in network byte order.

	// format:
	// word_str;  // a utf-8 string terminated by '\0'.
	// word_data_offset;  // word data's offset in .dict file
	// word_data_size;  // word data's total size in .dict file
	end := bytes.IndexByte(idx.content[idx.offset:], '\000')
	end += idx.offset
	// 1. word_str;  // a utf-8 string terminated by '\0'.
	wordStr := string(idx.content[idx.offset:end])

	newWord := Word{
		w: wordStr,
	}

	idx.offset = end + 1
	// 2. word_data_offset;  // word data's offset in .dict file
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
		newWord.offset = wOffset
	} else {
		return "", errorBits
	}

	// word_data_size;  // word data's total size in .dict file
	// word_data_size should be 32-bits unsigned number in network byte order.
	var wSize uint32
	sizeByte := idx.content[idx.offset : idx.offset+4]
	r := bytes.NewReader(sizeByte)
	binary.Read(r, binary.BigEndian, &wSize)
	newWord.size = wSize

	idx.offset += 4

	newWord.index = idx.index
	idx.index++

	// update the cache
	idx.wordLst = append(idx.wordLst, newWord)
	idx.wordDict[wordStr] = append(idx.wordDict[wordStr], newWord)

	return wordStr, nil
}

func (idx *Index) parse() {
	if idx.parsed {
		return
	}

	for _, err := idx.nextWord(); err == nil; _, err = idx.nextWord() {
	}
	idx.parsed = true
}

type Dictionary struct {
	info    *Info
	index   *Index
	content []byte
	offset  uint32
}

func mkchan(ifo, idx, dict io.Reader, d *Dictionary) <-chan error {
	ret := make(chan error, 3)
	defer close(ret)
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		info, err := newInfo(ifo)
		if err != nil {
			ret <- err
			return
		}
		d.info = info
	}()

	go func() {
		defer wg.Done()
		index, err := newIndex(idx)
		if err != nil {
			ret <- err
			return
		}
		d.index = index
	}()

	go func() {
		defer wg.Done()
		raw, err := ioutil.ReadAll(dict)
		if err != nil {
			ret <- err
			return
		}
		content, err := Gunzip(raw)
		if err != nil {
			ret <- err
			return
		}
		d.content = content
	}()

	wg.Wait()

	return ret
}

func NewDictionary(ifo, idx, dict io.Reader) (*Dictionary, error) {
	d := &Dictionary{
		offset: 0,
	}
	for err := range mkchan(ifo, idx, dict, d) {
		if err != nil {
			log.Fatal(err)
		}
	}
	return d, nil
}

func (d *Dictionary) isSameTypeSequence() bool {
	_, ok := d.info.Dict["sametypesequence"]
	return ok
}

// GetWord get the meaning of word
func (d *Dictionary) GetWord(word string) []map[uint8][]byte {
	index, ok := d.index.wordDict[word]
	if !ok {
		return nil
	}
	var ret []map[uint8][]byte
	for _, curWord := range index {
		d.offset = curWord.offset
		if d.isSameTypeSequence() {
			// set offset to this word meaning
			curValue := d.getWordSameSequence(curWord)
			ret = append(ret, curValue)
		} else {
			curValue := d.getWordNonSameSequence(curWord)
			ret = append(ret, curValue)
		}
	}
	return ret
}

func (d *Dictionary) GetFormatedMeaning(word string) []string {
	ret := []string{}
	ms := d.GetWord(word)
	for _, m := range ms {
		for k, v := range m {
			if k == byte('h') {
				txt, err := Unhtml(v)
				if err != nil {
					continue
				}
				lines := strings.Split(txt, "\n")
				ret = append(ret, lines...)
			} else {
				lines := strings.Split(string(v), "\n")
				ret = append(ret, lines...)
			}
		}
	}
	return ret
}

func (d *Dictionary) getWordNonSameSequence(word Word) map[uint8][]byte {
	ret := make(map[uint8][]byte)
	var readSize uint32
	startOffset := d.offset
	for readSize < word.size {
		typeByte := d.content[d.offset : d.offset+1]
		r := bytes.NewReader(typeByte)
		var c uint8
		binary.Read(r, binary.BigEndian, &c)
		// pass type byte
		d.offset++

		if strings.Index("mlgtxykwhnr", string(c)) >= 0 {
			end := bytes.IndexByte(d.content[d.offset:], '\000')
			end += int(d.offset)
			value := d.content[d.offset:end]
			d.offset = uint32(end) + 1
			ret[c] = value
		} else {
			sizeByte := d.content[d.offset : d.offset+4]
			r := bytes.NewReader(sizeByte)
			var s uint32
			binary.Read(r, binary.BigEndian, &s)
			d.offset += 4

			value := d.getEntryFieldSize(s)
			ret[c] = value
		}

		readSize = d.offset - startOffset
	}

	return ret
}

func (d *Dictionary) getWordSameSequence(word Word) map[uint8][]byte {
	ret := make(map[uint8][]byte)
	sametypesequence := d.info.Dict["sametypesequence"]

	startOffset := d.offset
	for i, c := range []byte(sametypesequence) {
		if strings.Index("mlgtxykwhnr", string(c)) >= 0 {
			// The first data entry for each word will have a terminating '\0', but
			// the second entry will not have a terminating '\0'.  The omissions of
			// the type chars and of the last field's size information are the
			// optimizations required by the "sametypesequence" option described
			// above.

			// last one
			if i == len(sametypesequence)-1 {
				value := d.getEntryFieldSize(word.size - (d.offset - startOffset))
				ret[c] = value
			} else {
				end := bytes.IndexByte(d.content[d.offset:], '\000')
				end += int(d.offset)
				value := d.content[d.offset:end]
				d.offset = uint32(end) + 1
				ret[c] = value
			}
		} else if strings.Index("WP", string(c)) >= 0 {
			// The data begins with a network byte-ordered guint32 to identify the wav
			// file's size, immediately followed by the file's content.

			// last one
			if i == len(sametypesequence)-1 {
				ret[c] = d.getEntryFieldSize(word.size - (d.offset - startOffset))
			} else {
				sizeByte := d.content[d.offset : d.offset+4]
				r := bytes.NewReader(sizeByte)
				var s uint32
				binary.Read(r, binary.BigEndian, &s)
				d.offset += 4

				value := d.getEntryFieldSize(s)
				ret[c] = value
			}
		}
	}
	return ret
}

func (d *Dictionary) getEntryFieldSize(size uint32) []byte {
	value := d.content[d.offset : d.offset+size]
	d.offset += size

	return value
}
