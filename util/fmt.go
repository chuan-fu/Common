package util

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

const (
	GrepReplace = `%c[1;31m%s%c[0m` // 颜色替换 0x1B[0;31m %s 0x1B[0m
	RedGrep     = 0x1B              // 红色
)

/*
颜色 背景
30	40	黑色
31	41	红色
32	42	绿色
33	43	黄色
34	44	蓝色
35	45	紫红色
36	46	青蓝色
37	47	白色

0	终端默认设置
1	高亮显示
4	使用下划线
5	闪烁
7	反白显示
8	不可见

0x1B是标记，[开始定义颜色，1代表高亮，40代表黑色背景，32代表绿色前景，0代表恢复默认颜色。

fmt.Printf("%c[1;37;41m 内容 %c[0m", 0x1B, 0x1B)
*/
var (
	BlackArrow     = fmt.Sprintf(`%c[1;30m➜ %c[0m`, 0x1B, 0x1B) // 黑色
	RedArrow       = fmt.Sprintf(`%c[1;31m➜ %c[0m`, 0x1B, 0x1B) // 红色
	GreenArrow     = fmt.Sprintf(`%c[1;32m➜ %c[0m`, 0x1B, 0x1B) // 绿色
	YellowArrow    = fmt.Sprintf(`%c[1;33m➜ %c[0m`, 0x1B, 0x1B) // 黄色
	BlueArrow      = fmt.Sprintf(`%c[1;34m➜ %c[0m`, 0x1B, 0x1B) // 蓝色
	PurpleRedArrow = fmt.Sprintf(`%c[1;35m➜ %c[0m`, 0x1B, 0x1B) // 紫红色
	GreenBlueArrow = fmt.Sprintf(`%c[1;36m➜ %c[0m`, 0x1B, 0x1B) // 青蓝色
	WhiteArrow     = fmt.Sprintf(`%c[1;37m➜ %c[0m`, 0x1B, 0x1B) // 白色

	DefaultPrefix   = "-> # "
	RedBluePrefix   = fmt.Sprintf(`%c[1;31m➜ %c[0m%c[1;34m# %c[0m`, 0x1B, 0x1B, 0x1B, 0x1B)
	GreenBluePrefix = fmt.Sprintf(`%c[1;32m➜ %c[0m%c[1;34m# %c[0m`, 0x1B, 0x1B, 0x1B, 0x1B)
)

type fmtGrain struct {
	i        int
	str      string
	spaceNum int
}

type FmtList struct {
	list []fmtGrain
	i    int
}

func NewFmtList() *FmtList {
	return &FmtList{
		i:    1,
		list: make([]fmtGrain, 0),
	}
}

func (p *FmtList) Add(str string) (i int) {
	p.list = append(p.list, fmtGrain{
		i:   p.i,
		str: str,
	})
	i = p.i
	p.i++
	return
}

func (p *FmtList) InitSpace() {
	length := p.GetLen(p.i - 1)
	for k := range p.list {
		v := &p.list[k]
		v.spaceNum = length - p.GetLen(v.i)
	}
}

func (p *FmtList) String() string {
	p.InitSpace()
	sBy := bytes.Buffer{}
	for k := range p.list {
		v := &p.list[k]
		sBy.WriteString(fmt.Sprintf("%d)", v.i))
		for i := 0; i <= v.spaceNum; i++ {
			sBy.WriteByte(' ')
		}
		sBy.WriteString(v.str)
		sBy.WriteString("\n")
	}
	return sBy.String()
}

func (p *FmtList) GetLen(i int) int {
	return len(strconv.Itoa(i))
}

func (p *FmtList) Len() int {
	return p.i - 1
}

func FmtColor(str, key string, grep int) string {
	if str == "" || key == "" {
		return str
	}
	if count := strings.Count(str, key); count > 0 {
		l := make([]interface{}, 0, 3*count)
		for i := 0; i < count; i++ {
			l = append(l, grep, key, grep)
		}
		return fmt.Sprintf(strings.ReplaceAll(str, key, GrepReplace), l...)
	}
	return str
}
