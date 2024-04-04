package common

import (
	"crypto/sha1"
	"fmt"
	"io"
	"testing"
)

func TestSha1En(t *testing.T) {
	salt := GetRandomBoth(4)
	data := "super11" + salt
	fmt.Printf("salt:%s\n", salt)
	t1 := sha1.New() ///产生一个散列值得方式
	_, _ = io.WriteString(t1, data)
	fmt.Printf("%x\n", t1.Sum(nil))
}
