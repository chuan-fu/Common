package keyboardx

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/chuan-fu/Common/baseservice/cast"

	"github.com/chuan-fu/Common/baseservice/colorx"

	"github.com/chuan-fu/Common/util"

	"github.com/chuan-fu/Common/zlog"
	"github.com/eiannone/keyboard"
)

const (
	space = " "
)

type Task func(s string) (isEnd bool, err error)

func Exit(s string) (isEnd bool, err error) {
	return true, nil
}

func KeyboardX(f Task, opts ...Option) error {
	c := buildConfig(opts)

	checkTask := func(s string) Task {
		handle := func(task Task) Task {
			return func(s string) (isEnd bool, err error) {
				if c.preHandle != nil {
					c.preHandle(s)
				}
				defer func() {
					if c.postHandle != nil {
						c.postHandle(s)
					}
				}()
				return task(s)
			}
		}
		if c.taskSvc != nil {
			if task := c.taskSvc.Match(s); task != nil {
				return handle(task)
			}
		}
		return handle(f)
	}

	keysEvents, err := keyboard.GetKeys(c.bufferSize)
	if err != nil {
		log.Error(err)
		return err
	}

	str := bytes.Buffer{}
	var isEnd bool

	cmd := newCommandHistory(c.cmdList, c.color)
	newTab := newTabSvc()

	fmt.Print(c.prefix)

	for {
		event := <-keysEvents
		if event.Err != nil {
			log.Error(event.Err)
		}

		switch event.Key {
		case keyboard.KeyEnter:
			fmt.Println()
			newTab.close()
			if str.Len() > 0 {
				s := str.String()
				str.Reset()

				// 历史命令 history XXX
				isHistory := c.needHistory && strings.HasPrefix(s, CmdHistory)
				if isHistory {
					cmd.History(s)
					fmt.Print(c.prefix)
				}

				// 校验是否加入历史数据
				if c.checkInHistoryHandle != nil {
					if s2, ok := c.checkInHistoryHandle(s); ok {
						cmd.add(s2)
					}
				} else {
					cmd.add(s)
				}

				if isHistory { // 如果历史触发，则不触发后续任务
					continue
				}

				if task := checkTask(s); task != nil {
					isEnd, err = task(s)
				} else {
					fmt.Println("---未搜索到任务---")
				}
			} else { // 空回车处理
				if c.emptyEnter != nil {
					isEnd, err = c.emptyEnter("")
				}
			}
			if err != nil {
				log.Error(err)
			}
			if isEnd {
				return nil
			}
			fmt.Print(c.prefix)
		case keyboard.KeySpace:
			str.WriteString(space)
			fmt.Print(space)
		case keyboard.KeyCtrlC, keyboard.KeyCtrlX, keyboard.KeyCtrlZ:
			fmt.Println()
			fmt.Println()
			return nil
		case keyboard.KeyArrowUp, keyboard.KeyArrowDown, keyboard.KeyPgup, keyboard.KeyPgdn:
			s := cmd.setIndex(func() int {
				if event.Key == keyboard.KeyArrowUp || event.Key == keyboard.KeyPgup {
					return -1
				}
				return 1
			}())
			str.Reset()
			str.WriteString(s)
			fmt.Print("\033[2K\r", c.prefix, s)
		case keyboard.KeyArrowLeft,
			keyboard.KeyArrowRight,
			keyboard.KeyCtrlE,
			keyboard.KeyCtrlO,
			keyboard.KeyEsc,
			keyboard.KeyBackspace,
			keyboard.KeyF1,
			keyboard.KeyF2,
			keyboard.KeyF3,
			keyboard.KeyF4,
			keyboard.KeyF5,
			keyboard.KeyF6,
			keyboard.KeyF7,
			keyboard.KeyF8,
			keyboard.KeyF9,
			keyboard.KeyF10,
			keyboard.KeyF11,
			keyboard.KeyF12,
			keyboard.KeyInsert,
			keyboard.KeyDelete,
			keyboard.KeyHome,
			keyboard.KeyEnd:
			// 特殊字符，不处理

		case keyboard.KeyBackspace2:
			newTab.close()
			if str.Len() == 0 {
				break
			}
			sR := []rune(str.String())
			s := string(sR[:len(sR)-1])
			str.Reset()
			str.WriteString(s)
			fmt.Print("\033[2K\r", c.prefix, s)
		case keyboard.KeyTab:
			s, index := newTab.run(str.String())
			s2, show := cmd.find(s, index)
			if s2 == "" {
				break
			}
			str.Reset()
			str.WriteString(s2)
			fmt.Print("\033[2K\r", c.prefix, show)
		default:
			s := string(event.Rune)
			str.WriteString(s)
			fmt.Print(s)
		}
	}
}

type tabSvc struct {
	isTab    bool
	tabStr   string
	tabIndex int
}

func newTabSvc() *tabSvc {
	return &tabSvc{}
}

func (t *tabSvc) close() {
	if t.isTab {
		t.isTab = false
	}
}

func (t *tabSvc) run(s string) (string, int) {
	if !t.isTab {
		t.init(s)
	} else {
		t.tabIndex++
	}
	return t.tabStr, t.tabIndex
}

func (t *tabSvc) init(s string) {
	t.isTab = true
	t.tabStr = s
	t.tabIndex = 0
}

type commandHistory struct {
	commandList []string
	sum, index  int
	color       int
}

func newCommandHistory(cmdList []string, color int) *commandHistory {
	return &commandHistory{
		commandList: cmdList,
		sum:         len(cmdList),
		index:       len(cmdList),
		color:       color,
	}
}

func (c *commandHistory) setIndex(i int) string {
	c.index += i
	if c.index < 0 {
		c.index = 0
	}
	if c.index > c.sum {
		c.index = c.sum
	}

	if c.index == c.sum {
		return ""
	}
	return c.commandList[c.index]
}

func (c *commandHistory) add(s string) {
	if c.sum > 0 {
		if c.commandList[c.sum-1] == s {
			c.index = c.sum
			return
		}
	}
	c.commandList = append(c.commandList, s)
	c.sum++
	c.index = c.sum
}

func (c *commandHistory) find(s string, index int) (original, show string) {
	for i := c.sum - 1; i >= 0; i-- {
		if strings.HasPrefix(c.commandList[i], s) {
			if index == 0 {
				original = c.commandList[i]
				show = fmt.Sprintf("%s%s", colorx.Sprint(c.color, s), original[len(s):])
				return
			}
			index--
		}
	}
	return
}

func (c *commandHistory) History(key string) {
	key = strings.TrimSpace(strings.TrimPrefix(key, CmdHistory))

	fmtList := util.NewFmtList()
	for k, v := range c.commandList {
		if key == "" {
			fmtList.Add(c.commandList[k])
			continue
		}
		if strings.Contains(v, key) {
			fmtList.Add(colorx.KeywordsSprintf(c.color, v, key))
		}
	}
	fmt.Println(fmtList.String())
}

func Close() error {
	return keyboard.Close()
}

func CheckStrList(list []string) (resp string) {
	switch len(list) {
	case 0:
		fmt.Println()
	case 1:
		fmt.Printf("%s\n\n", list[0])
		resp = list[0]
	default:
		f := util.NewFmtList()
		for k := range list {
			f.Add(list[k])
		}
		fmt.Print(f.String())

		_ = KeyboardX(
			func(s string) (isEnd bool, err error) {
				index, err2 := cast.ToInt(s)
				if err2 == nil {
					index -= 1
					if len(list) > index && index >= 0 {
						resp = list[index]
						return true, nil
					}
				}
				fmt.Println("---下标有误---")
				return true, nil
			},
			WithEmptyEnter(func(s string) (isEnd bool, err error) {
				resp = list[0]
				return true, nil
			}),
			WithPrefix(colorx.PurpleRedArrow),
		)
	}
	return
}
