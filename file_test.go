package stardict
import (
	"testing"
	"fmt"
)

func TestFileNew(t *testing.T) {
	i := NewInfo("/tmp/stardict-HanYuChengYuCiDian-new_colors-2.4.2/HanYuChengYuCiDian-new_colors.ifo")
	if i == nil {
		t.Fatalf("%s\n", "NewInfo fail")
	}
	fmt.Println(i)
}
