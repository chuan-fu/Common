package colorx

import (
	"fmt"
	"strings"
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

const (
	GrepReplace = `%c[1;31m%s%c[0m` // 颜色替换 0x1B[0;31m %s 0x1B[0m
)

const (
	WordBlack = iota + 30
	WordRed
	WordGreen
	WordYellow
	WordBlue
	WordPurpleRed
	WordGreenBlue
	WordWhite
)

const (
	BgBlack = iota + 40
	BgRed
	BgGreen
	BgYellow
	BgBlue
	BgPurpleRed
	BgGreenBlue
	BgWhite
)

func Print(color int, s string) {
	fmt.Printf("%c[1;%dm%s%c[0m", 0x1B, color, s, 0x1B)
}

func Sprint(color int, s string) string {
	return fmt.Sprintf("%c[1;%dm%s%c[0m", 0x1B, color, s, 0x1B)
}

func Printf(color int, s string, a ...interface{}) {
	fmt.Printf("%c[1;%dm%s%c[0m", 0x1B, color, fmt.Sprintf(s, a...), 0x1B)
}

func Sprintf(color int, s string, a ...interface{}) string {
	return fmt.Sprintf("%c[1;%dm%s%c[0m", 0x1B, color, fmt.Sprintf(s, a...), 0x1B)
}

func Println(color int, s string) {
	fmt.Printf("%c[1;%dm%s%c[0m\n", 0x1B, color, s, 0x1B)
}

func KeywordsSprintf(color int, s, keywords string) string {
	if s == "" || keywords == "" {
		return s
	}
	if count := strings.Count(s, keywords); count > 0 {
		l := make([]interface{}, 0, 3*count)
		for i := 0; i < count; i++ {
			l = append(l, color, keywords, color)
		}
		return fmt.Sprintf(strings.ReplaceAll(s, keywords, GrepReplace), l...)
	}
	return s
}
