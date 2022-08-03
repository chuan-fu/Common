package keyboardx

import (
	"fmt"
	"strings"

	"github.com/chuan-fu/Common/baseservice/colorx"
	"github.com/chuan-fu/Common/util"
)

const (
	findIndexClose = -1
	findIndexStart = 0
)

type CmdListServ interface {
	InitCmdList(list []string)
	CmdList() []string

	SetColor(color int)
	Last() string             // 上一个【上键】
	Next() string             // 下一个【下键】
	Add(s string) CmdListServ // 添加指令

	FindStart(keywords string) (original, show string) // 开始查找，持续查找
	FindClose()                                        // 结束查找

	Find(keywords string, index int) (original, show string) // index为找到的第几个关键词

	HistoryPrefixTask(matchKey string) Task // 历史指令列表 matchKey为匹配指令
}

type cmdListServ struct {
	list       []string
	sum, index int
	color      int

	findIndex   int
	findKeyword string
}

func newCmdListServ(list []string) CmdListServ {
	return &cmdListServ{
		list:      list,
		sum:       len(list),
		index:     len(list),
		color:     colorx.WordRed,
		findIndex: findIndexClose,
	}
}

func (c *cmdListServ) InitCmdList(list []string) {
	c.list = list
	c.sum = len(list)
	c.index = c.sum
}

func (c *cmdListServ) CmdList() []string {
	return c.list
}

func (c *cmdListServ) SetColor(color int) {
	c.color = color
}

// 上一个【上键】
func (c *cmdListServ) Last() string {
	if c.index > 0 {
		c.index -= 1
	}
	return c.list[c.index]
}

// 下一个【下键】
func (c *cmdListServ) Next() string {
	c.index += 1
	if c.index >= c.sum {
		c.index = c.sum
		return ""
	}
	return c.list[c.index]
}

func (c *cmdListServ) Add(s string) CmdListServ {
	if c.sum > 0 { // 历史最后命令 = s，则不写入
		if c.list[c.sum-1] == s {
			c.index = c.sum
			return c
		}
	}
	c.list = append(c.list, s)
	c.sum++
	c.index = c.sum
	return c
}

func (c *cmdListServ) FindStart(keywords string) (original, show string) { // 开始查找，持续查找
	c.findIndex += 1
	if c.findIndex == findIndexStart { // 为0 表示第一次查询
		c.findKeyword = keywords
	}
	return c.Find(c.findKeyword, c.findIndex)
}

func (c *cmdListServ) FindClose() { // 结束查找
	c.findIndex = findIndexClose
}

func (c *cmdListServ) Find(keywords string, index int) (original, show string) {
	for i := c.sum - 1; i >= 0; i-- {
		if strings.HasPrefix(c.list[i], keywords) {
			if index == 0 {
				original = c.list[i]
				show = fmt.Sprintf("%s%s", colorx.Sprint(c.color, keywords), original[len(keywords):])
				return
			}
			index--
		}
	}
	return
}

// matchKey为匹配的key
func (c *cmdListServ) HistoryPrefixTask(matchKey string) Task {
	return NewHandleTask(func(s string) (isEnd bool, err error) {
		keywords := strings.TrimSpace(strings.TrimPrefix(s, matchKey))
		c.history(keywords)
		return
	})
}

// 获取历史命令
func (c *cmdListServ) history(keywords string) {
	fmtList := util.NewFmtList()
	for k, v := range c.list {
		if keywords == "" {
			fmtList.Add(c.list[k])
			continue
		}
		if strings.Contains(v, keywords) {
			fmtList.Add(colorx.KeywordsSprintf(c.color, v, keywords))
		}
	}
	fmt.Println(fmtList.String())
}
