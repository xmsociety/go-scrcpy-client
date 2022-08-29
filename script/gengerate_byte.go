package script

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// 二进制文件硬编码
func generateByteCode(path string) {
	codeString := "package script\nvar StaticContent = []byte{"
	f, _ := os.Create("./outbut.go")
	bytes, _ := ioutil.ReadFile(filepath.Clean(path))
	for _, b := range bytes {
		r := strconv.FormatUint(uint64(b), 16)
		codeString += fmt.Sprintf("0x%02v, ", r)
	}
	codeString += "}"
	_, err := f.WriteString(codeString)
	if err != nil {
		fmt.Println(err)
	}
}
