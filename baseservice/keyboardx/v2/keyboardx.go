package keyboardx

import (
	"bytes"
	"fmt"

	"github.com/chuan-fu/Common/baseservice/keyboardx"

	"github.com/chuan-fu/Common/baseservice/colorx"
	"github.com/chuan-fu/Common/zlog"
	"github.com/eiannone/keyboard"
)

var (
	CmdExit     = []string{"q", "quit", "exit"}
	CmdExitDesc = "退出"

	CmdHelp     = []string{"h", "help"}
	CmdHelpDesc = "功能介绍"
)

const (
	fastTaskHistory = 1 << iota
	fastTaskHelp
	fastTaskExit
)

const (
	CmdHistory        = "history"
	CmdHistoryDesc    = "历史记录"
	DefaultBufferSize = 10
	space             = " "
)

type CheckInHistoryFunc func(string) (string, bool)

type KeyboardXServ interface {
	TaskService

	SetPrefix(s string) KeyboardXServ           // 前缀 默认使用《➜ ## 》
	SetColor(color int) KeyboardXServ           // 关键词颜色 默认红色
	SetBufferSize(bufferSize int) KeyboardXServ // 输入缓冲区 默认为10，基本不用修改
	SetEntryCmd(t Task) KeyboardXServ           // 空回车逻辑处理
	SetCheckInHistory(f CheckInHistoryFunc) KeyboardXServ

	ResetTaskService(f TaskService) KeyboardXServ

	InitCmdList(list []string) KeyboardXServ // 重置CmdList

	AddHistory() KeyboardXServ
	AddHelp() KeyboardXServ
	AddExit() KeyboardXServ

	Run() error
}

type keyboardX struct {
	TaskService // 任务服务

	prefix     string // 前缀 默认使用《➜ ## 》
	color      int    // 关键词颜色 默认红色
	bufferSize int    // 输入缓冲区 默认为10，基本不用修改

	emptyCmd Task // 空回车逻辑处理

	runTaskFunc RunTaskFunc // 执行任务 默认baseservice/keyboardx/v2/task.go:41

	checkInHistoryFunc CheckInHistoryFunc // 校验是否加入历史
	cmdListSvc         CmdListServ        // 历史指令服务

	fastTask int
}

func NewKeyboardX() KeyboardXServ {
	return &keyboardX{
		prefix:      colorx.GreenBluePrefix,
		color:       colorx.WordRed,
		bufferSize:  DefaultBufferSize,
		TaskService: NewTaskService(),
		runTaskFunc: runTask,
	}
}

func (k *keyboardX) FmtPrefix() {
	fmt.Print(k.prefix)
}

func (k *keyboardX) ResetPrefix(s string) {
	fmt.Print("\033[2K\r", k.prefix, s)
}

func (k *keyboardX) SetPrefix(s string) KeyboardXServ {
	k.prefix = s
	return k
}

func (k *keyboardX) SetColor(color int) KeyboardXServ {
	k.color = color
	return k
}

func (k *keyboardX) SetBufferSize(bufferSize int) KeyboardXServ {
	k.bufferSize = bufferSize
	return k
}

func (k *keyboardX) SetEntryCmd(t Task) KeyboardXServ {
	k.emptyCmd = t
	return k
}

func (k *keyboardX) SetRunTaskFunc(f RunTaskFunc) KeyboardXServ {
	k.runTaskFunc = f
	return k
}

func (k *keyboardX) SetCheckInHistory(f CheckInHistoryFunc) KeyboardXServ {
	k.checkInHistoryFunc = f
	return k
}

func (k *keyboardX) ResetTaskService(taskSvc TaskService) KeyboardXServ {
	k.TaskService = taskSvc
	return k
}

func (k *keyboardX) InitCmdList(list []string) KeyboardXServ {
	if k.cmdListSvc == nil {
		k.cmdListSvc = newCmdListServ(list)
		return k
	}
	k.cmdListSvc.InitCmdList(list)
	return k
}

// 使用默认history命令
func (k *keyboardX) AddHistory() KeyboardXServ {
	k.fastTask |= fastTaskHistory
	return k
}

// 使用默认help命令
func (k *keyboardX) AddHelp() KeyboardXServ {
	k.fastTask |= fastTaskHelp
	return k
}

func (k *keyboardX) AddExit() KeyboardXServ {
	k.fastTask |= fastTaskExit
	return k
}

func (k *keyboardX) Init() {
	if k.cmdListSvc == nil {
		k.cmdListSvc = newCmdListServ(nil)
	}
	k.cmdListSvc.SetColor(k.color) // 初始化颜色

	if k.fastTask&fastTaskHistory > 0 {
		k.AddPrefixTask(CmdHistory, CmdHistoryDesc, k.cmdListSvc.HistoryPrefixTask(keyboardx.CmdHistory))
	}
	if k.fastTask&fastTaskHelp > 0 {
		k.AddHelpTasks(CmdHelp, CmdHelpDesc)
	}
	if k.fastTask&fastTaskExit > 0 {
		k.AddFullyTasks(CmdExit, CmdExitDesc, NewExitTask())
	}
}

func (k *keyboardX) Run() error {
	k.Init()

	keysEvents, err := keyboard.GetKeys(k.bufferSize)
	if err != nil {
		log.Error(err)
		return err
	}

	cmd := bytes.Buffer{} // 存指令
	var isEnd bool        // 是否结束

	k.FmtPrefix()
	for {
		event := <-keysEvents
		if event.Err != nil {
			log.Error(event.Err)
		}

		switch event.Key {
		case keyboard.KeyEnter:
			fmt.Println() // 回车键效果
			k.cmdListSvc.FindClose()
			if cmd.Len() == 0 {
				isEnd, err = k.runTaskFunc(k.emptyCmd, "")
			} else {
				s := cmd.String()
				cmd.Reset()
				if task := k.MatchTask(s); task != nil {
					isEnd, err = k.runTaskFunc(task, s)
				}
				SetHistory(k.checkInHistoryFunc, k.cmdListSvc, s)
			}
			if err != nil {
				log.Error(err)
			}
			if isEnd {
				return nil
			}
			k.FmtPrefix()

		case keyboard.KeySpace: // 空格
			cmd.WriteString(space)
			fmt.Print(space)
		case keyboard.KeyCtrlC, keyboard.KeyCtrlX, keyboard.KeyCtrlZ: // 结束任务
			fmt.Println()
			fmt.Println()
			return nil

		case keyboard.KeyArrowUp, keyboard.KeyPgup: // 上键
			s := k.cmdListSvc.Last()
			cmd.Reset()
			cmd.WriteString(s)
			k.ResetPrefix(s)
		case keyboard.KeyArrowDown, keyboard.KeyPgdn: // 下键
			s := k.cmdListSvc.Next()
			cmd.Reset()
			cmd.WriteString(s)
			k.ResetPrefix(s)

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

		case keyboard.KeyBackspace2: // 删除键
			k.cmdListSvc.FindClose()
			if cmd.Len() == 0 {
				continue
			}
			sR := []rune(cmd.String())
			s := string(sR[:len(sR)-1])
			cmd.Reset()
			cmd.WriteString(s)
			k.ResetPrefix(s)

		case keyboard.KeyTab: // 补全键
			s2, show := k.cmdListSvc.FindStart(cmd.String())
			if s2 == "" {
				continue
			}
			cmd.Reset()
			cmd.WriteString(s2)
			k.ResetPrefix(show)
		default:
			s := string(event.Rune)
			cmd.WriteString(s)
			fmt.Print(s)
		}
	}
}

func SetHistory(check CheckInHistoryFunc, cmdList CmdListServ, s string) {
	if check == nil {
		cmdList.Add(s)
		return
	}
	if s2, ok := check(s); ok {
		cmdList.Add(s2)
	}
}
