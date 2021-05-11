package running

import (
	"flag"
	"log"

	"github.com/stretchr/testify/assert"

	"testing"

	"k8s.io/klog/v2"
)

func TestStringSliceDelete(t *testing.T) {
	log.SetFlags(log.Lshortfile)
	_, _, funcname := Trace()
	klog.Infoln(funcname + " start")
	defer klog.Infoln(funcname + " end")
	flag.Parse()

	s := []string{"aaa", "bbb", "bbb", "ccc", "ddd", "ee"}
	ss := StringsDel(s, "bbb")
	klog.Infoln(ss)
	assert.Equal(t, len(ss), 4)
}
