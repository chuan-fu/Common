package filex

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/chuan-fu/Common/util"
	"github.com/chuan-fu/Common/zlog"
)

func CheckExist(path string) bool {
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

// 读取文件
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// 逐行读取文件
func ReadFileLine(path string) (<-chan []byte, error) {
	fi, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	ch := make(chan []byte)

	go func(f *os.File, c chan []byte) {
		defer util.DeferFunc()
		defer f.Close()
		defer close(c)

		br := bufio.NewReader(f)
		for {
			line, _, err2 := br.ReadLine()
			if err2 == io.EOF {
				break
			}
			c <- line
		}
		close(ch)
	}(fi, ch)

	return ch, nil
}

// 覆盖式写入
// 文件不存在则创建
func WriteFile(path string, data []byte) error {
	return ioutil.WriteFile(path, data, os.ModePerm)
}

// 追加式写入
// 文件不存在则创建
func AppendFile(path, data string) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return nil
}

// 创建目录
func Mkdir(path string) error {
	return os.Mkdir(path, os.ModePerm)
}

// 创建目录
func MkdirAll(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// 删除文件
func Remove(path string) error {
	return os.Remove(path)
}

// 遍历读取目录
func ReadPath(paths ...string) <-chan string {
	ch := make(chan string)
	for k := range paths {
		if !strings.HasSuffix(paths[k], "/") { // 添加后缀/
			paths[k] += "/"
		}
		if _, err := os.Stat(paths[k]); err != nil { // 查询目录是否存在，如果不存在，则报错
			panic(err)
		}
	}

	go func(pathList []string, c chan<- string) {
		defer util.DeferFunc()
		defer close(c)
		for k := range pathList {
			readDir(pathList[k], c)
		}
	}(paths, ch)

	return ch
}

func readDir(path string, c chan<- string) {
	rds, err := ioutil.ReadDir(path)
	if err != nil {
		log.Error(err)
		return
	}
	for _, fi := range rds {
		if fi.IsDir() {
			p2 := fmt.Sprintf(`%s%s/`, path, fi.Name())
			c <- p2
			readDir(p2, c)
		} else {
			c <- path + fi.Name()
		}
	}
}
