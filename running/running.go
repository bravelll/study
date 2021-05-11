package running

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"runtime"
	"unsafe"

	"k8s.io/klog/v2"
)

func Trace() (string, int, string) {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return "?", 0, "?"
	}

	fn := runtime.FuncForPC(pc)
	return file, line, fn.Name()
}

func FileExist(path string) bool {
	_, err := os.Lstat(path)
	return !os.IsNotExist(err)
}

func PrintStack() {
	var buf [1024]byte
	n := runtime.Stack(buf[:], false)
	fmt.Printf("==> %s\n", string(buf[:n]))
}

func Recover() {
	ev := recover()
	if ev != nil {
		switch ev.(type) {
		case runtime.Error: // 运行时错误
			klog.Errorln("panic:", ev)
			PrintStack()
		default: // 非运行时错误
			klog.Errorln("error:", ev)
		}
		klog.Flush()
	}
}

func LoadEnviron() (map[string]string, error) {
	envs := make(map[string]string)
	buf, err := ioutil.ReadFile("/proc/1/environ")
	if err != nil {
		klog.Errorln(err)
		return envs, err
	}
	for _, line := range bytes.Split(buf, []byte{0}) {
		values := bytes.SplitN(line, []byte{'='}, 2)
		if len(values) == 2 {
			envs[string(values[0])] = string(values[1])
		}
	}
	return envs, nil
}

func StringsAdd(slice []string, item string) []string {
	for _, el := range slice {
		if el == item {
			return slice
		}
	}
	slice = append(slice, item)
	return slice
}

func StringsDel(slice []string, item string) []string {
	for i := 0; i < len(slice); i++ {
		if item == slice[i] {
			if i < len(slice)-1 {
				slice = append(slice[:i], slice[i+1:]...)
			} else {
				slice = slice[:i]
			}
			i--
		}
	}
	return slice
}

func String2Bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
