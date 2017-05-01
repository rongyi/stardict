package stardict

import (
	"fmt"
	"io/ioutil"
	"strings"
)

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
