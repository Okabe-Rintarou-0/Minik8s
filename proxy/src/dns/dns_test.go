package dns

import (
	"testing"
)

func Test(t *testing.T) {
	dm := New("/home/vectorxj/Desktop/cloud/coredns/hosts")
	if err := dm.DelIfExistEntry("should.not.exist"); err != nil {
		t.Error(err)
	}
	if err := dm.AddEntry("should.exist", "127.0.0.2"); err != nil {
		t.Error(err)
	}
}
