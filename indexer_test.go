package indexer

import (
	"testing"
)

func TestCompare(t *testing.T) {
	var comp struct {
		loc, rem []*FileInfo
	}

	comp.loc = []*FileInfo{
		&FileInfo{
			Path: "/a",
			Hash: "111",
		},
		&FileInfo{
			Path: "/b",
			Hash: "222",
		},
		&FileInfo{
			Path: "/c",
			Hash: "333",
		},
	}
	comp.rem = []*FileInfo{
		&FileInfo{
			Path: "/a",
			Hash: "111",
		},
		&FileInfo{
			Path: "/b",
			Hash: "223",
		},
		&FileInfo{
			Path: "/d",
			Hash: "444",
		},
	}

	ind, _ := NewIndexer("http://a.com/", "devID", "/")
	ret, err := ind.mergeFiles(comp.rem, comp.loc)
	if err != nil {
		t.Error(err)
		return
	}
	if len(ret.Add) != 1 || ret.Add[0].Path != "/c" {
		t.Errorf("ret.Add: %v", ret.Add)
		return
	}
	if len(ret.Update) != 1 || ret.Update[0].Path != "/b" {
		t.Errorf("ret.Update: %+v", ret.Update)
		return
	}
	if len(ret.Del) != 1 || ret.Del[0].Path != "/d" {
		t.Errorf("ret.Del: %+v", ret.Del)
		return
	}

}
