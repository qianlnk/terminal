package terminal

import (
	"fmt"
	"strings"
	"sync"
)

const (
	SYS_ASCII_TAB        = 9
	SYS_ASCII_LF         = 10
	SYS_ASCII_BACK_SPACE = 127
	SYS_ASCII_ETX        = 3
)

const (
	SYS_MAX_HISTORY = 100
)

type Terminal struct {
	prompt       string
	historyList  []string
	historyIndex int
	history      bool
	sysCmdList   []string
	sysCmdMaxLen int
	echo         bool
	cursorX      int
	cursorY      int
	lock         sync.Mutex
}

func NewTerminal(prompt string) *Terminal {
	return &Terminal{
		prompt:       prompt,
		historyIndex: 0,
		history:      true,
		echo:         true,
		cursorX:      0,
		cursorY:      len(prompt),
	}
}

func (tm *Terminal) SetSystemCommand(cmdlist []string) {
	tm.lock.Lock()
	defer tm.lock.Unlock()
	tm.sysCmdList = append(tm.sysCmdList, cmdlist...)
	for _, cmd := range cmdlist {
		if len(cmd) > tm.sysCmdMaxLen {
			tm.sysCmdMaxLen = len(cmd)
		}
	}
}

func (tm *Terminal) History(addToHistory bool) {
	tm.lock.Lock()
	defer tm.lock.Unlock()
	tm.history = addToHistory
}

func (tm *Terminal) Echo(echo bool) {
	tm.lock.Lock()
	defer tm.lock.Unlock()
	tm.echo = echo
}

func (tm *Terminal) addXY(x, y int) {
	tm.lock.Lock()
	defer tm.lock.Unlock()
	tm.cursorX += x
	tm.cursorY += y
}

func (tm *Terminal) cursorMoveUp(x int) {
	tm.lock.Lock()
	defer tm.lock.Unlock()
	if x <= 0 {
		return
	}
	if tm.cursorX <= 0 {
		return
	}
	if x >= tm.cursorX {
		fmt.Printf("\033[%dA", tm.cursorX)
		tm.cursorX = 0
	} else {
		fmt.Printf("\033[%dA", x)
		tm.cursorX -= x
	}
}

func (tm *Terminal) cursorMoveDown(x int) {
	tm.lock.Lock()
	defer tm.lock.Unlock()
	if x <= 0 {
		return
	}
	fmt.Printf("\033[%dB", x)
	tm.cursorX += x

}

func (tm *Terminal) cursorMoveRight(y int) {
	tm.lock.Lock()
	defer tm.lock.Unlock()
	if y <= 0 {
		return
	}
	fmt.Printf("\033[%dC", y)
	tm.cursorY += y
}

func (tm *Terminal) cursorMoveLeft(y int) {
	tm.lock.Lock()
	defer tm.lock.Unlock()
	if y <= 0 {
		return
	}
	if tm.cursorY <= 0 {
		return
	}
	if y >= tm.cursorY {
		fmt.Printf("\033[%dD", tm.cursorY)
		tm.cursorY = 0
	} else {
		fmt.Printf("\033[%dD", y)
		tm.cursorY -= y
	}
}

func (tm *Terminal) cursorMoveTo(x, y int) {
	if x >= tm.cursorX {
		tm.cursorMoveDown(x - tm.cursorX)
	} else {
		tm.cursorMoveUp(tm.cursorX - x)
	}

	if y >= tm.cursorY {
		tm.cursorMoveRight(y - tm.cursorY)
	} else {
		tm.cursorMoveLeft(tm.cursorY - y)
	}

	tm.cursorX = x
	tm.cursorY = y
}

func (tm *Terminal) getXY() (int, int) {
	tm.lock.Lock()
	defer tm.lock.Unlock()
	return tm.cursorX, tm.cursorY
}

func (tm *Terminal) getInput() string {
	var cmd []byte
	var leftCmd, rightCmd []byte
	tabYet := false
	tabIndex := -1
	tm.cursorX = 0
	tabLine := 0
	tm.cursorY = len(tm.prompt)
	var sameCmdList []string
	for {
		buf, parse := getch()

		if buf == SYS_ASCII_TAB && tabYet == true {
			tabYet = true
			tabIndex++
		} else if tabYet == true {
			if buf != SYS_ASCII_LF {
				tabYet = false
			}
			tabIndex = -1
			sameCmdList = nil
			_, tabY := tm.getXY()
			tm.cursorMoveLeft(tabY)
			for i := 0; i < tabLine; i++ {
				tm.cursorMoveDown(1)
				cleanLine()
			}
			tm.cursorMoveUp(tabLine)
			tm.cursorMoveRight(tabY)
		}
		if buf == SYS_ASCII_LF {
			if tabYet == true {
				tabYet = false
				fmt.Printf(" ")
				leftCmd = append(leftCmd, byte(' '))
				tm.addXY(0, 1)
				continue
			} else {
				break
			}
		}
		switch buf {
		case SYS_UP:
			if !tm.echo || tm.historyIndex == 0 {
				break
			}
			if tm.historyIndex > 0 {
				tm.historyIndex--
				tm.cursorMoveLeft(len(leftCmd))
				cleanCell(len(leftCmd) + len(rightCmd))
				tm.addXY(0, len(leftCmd)+len(rightCmd))
				tm.cursorMoveLeft(len(leftCmd) + len(rightCmd))
				leftCmd = []byte(tm.historyList[tm.historyIndex])
				rightCmd = nil
				fmt.Printf("%s", string(leftCmd))
				tm.addXY(0, len(leftCmd))
			}
			break
		case SYS_DOWN:
			if !tm.echo {
				break
			}
			if tm.historyIndex < len(tm.historyList) {
				tm.historyIndex++
				tm.cursorMoveLeft(len(leftCmd))
				cleanCell(len(leftCmd) + len(rightCmd))
				tm.addXY(0, len(leftCmd)+len(rightCmd))
				tm.cursorMoveLeft(len(leftCmd) + len(rightCmd))
				if tm.historyIndex >= len(tm.historyList) {
					leftCmd = nil
					rightCmd = nil
				} else {
					leftCmd = []byte(tm.historyList[tm.historyIndex])
					rightCmd = nil
				}
				fmt.Printf("%s", string(leftCmd))
				tm.addXY(0, len(leftCmd))
			}
			break
		case SYS_LEFT:
			if !tm.echo {
				break
			}
			if len(leftCmd) > 0 {
				tm.cursorMoveLeft(1)
				var tmpRight []byte
				tmpRight = append(tmpRight, leftCmd[len(leftCmd)-1])
				rightCmd = append(tmpRight, rightCmd...)
				if len(leftCmd) > 1 {
					leftCmd = leftCmd[0 : len(leftCmd)-1]
				} else {
					leftCmd = nil
				}
			}
			break
		case SYS_RIGHT:
			if !tm.echo {
				break
			}
			if len(rightCmd) > 0 {
				tm.cursorMoveRight(1)
				leftCmd = append(leftCmd, rightCmd[0])
				rightCmd = rightCmd[1:]
			}
			break
		case SYS_PARSE:
			if tm.echo {
				fmt.Printf("%s%s", parse, string(rightCmd))
			} else {
				for i := 0; i < len(parse)+len(rightCmd); i++ {
					fmt.Printf("*")
				}
			}
			tm.addXY(0, len(parse))
			if len(rightCmd) != 0 {
				tm.cursorMoveLeft(len(rightCmd))
			}
			leftCmd = append(leftCmd, []byte(parse)...)
			break
		case SYS_ASCII_TAB:
			if !tm.echo {
				break
			}
			if len(tm.sysCmdList) > 0 {
				if tabYet == false {
					for _, cmd := range tm.sysCmdList {
						if len(leftCmd) > len(cmd) {
							continue
						}
						if string(leftCmd) == string([]byte(cmd)[0:len(leftCmd)]) {
							sameCmdList = append(sameCmdList, cmd)
						}
					}

					if len(sameCmdList) == 1 {
						tm.cursorMoveLeft(len(leftCmd))
						cleanCell(len(leftCmd))
						tm.addXY(0, len(leftCmd))
						tm.cursorMoveLeft(len(leftCmd))
						leftCmd = []byte(sameCmdList[0])
						fmt.Printf("%s", string(leftCmd))
						tm.addXY(0, len(leftCmd))
					} else if len(sameCmdList) > 1 {
						x, y := tm.getXY()
						for i, cmd := range sameCmdList {
							if i%3 == 0 {
								fmt.Println()
								tm.cursorX++
								tm.cursorY = 0
							}
							fmt.Printf("%-30s ", cmd)
							tm.cursorY += 31
						}
						tabLine = tm.cursorX - x
						tm.cursorMoveUp(tm.cursorX - x)
						tm.cursorMoveLeft(tm.cursorY - y)
					}
					tabYet = true
				} else {
					if len(sameCmdList) > 1 {
						x, y := tm.getXY()
						if tabIndex == len(sameCmdList) {
							tabIndex = 0
						}
						for i, cmd := range sameCmdList {
							if i%3 == 0 {
								fmt.Println()
								tm.cursorX++
								tm.cursorY = 0
							}
							if i == tabIndex {
								fmt.Printf("\033[31;47m%-30s\033[0m ", cmd)
							} else {
								fmt.Printf("%-30s ", cmd)
							}
							tm.cursorY += 31
						}
						tm.cursorMoveUp(tm.cursorX - x)
						tm.cursorMoveLeft(tm.cursorY - y)

						tm.cursorMoveLeft(len(leftCmd))
						cleanCell(len(leftCmd))
						tm.addXY(0, len(leftCmd))
						tm.cursorMoveLeft(len(leftCmd))
						leftCmd = []byte(sameCmdList[tabIndex])
						fmt.Printf("%s", string(leftCmd))
						tm.addXY(0, len(leftCmd))
					}
				}
			}
			break
		case SYS_ASCII_BACK_SPACE:
			if len(leftCmd) > 0 {
				if len(leftCmd) > 1 {
					leftCmd = leftCmd[0 : len(leftCmd)-1]
				} else {
					leftCmd = nil
				}
				tm.cursorMoveLeft(1)
				fmt.Printf("%s%c", string(rightCmd), ' ')
				tm.addXY(0, len(rightCmd)+1)
				tm.cursorMoveLeft(len(rightCmd) + 1)
			}
			break
		case SYS_ASCII_ETX:
			return ""
		default:
			if tm.echo {
				fmt.Printf("%c%s", buf, string(rightCmd))
			} else {
				fmt.Printf("*")
			}
			tm.cursorY++
			tm.addXY(0, len(rightCmd))
			if len(rightCmd) != 0 {
				tm.cursorMoveLeft(len(rightCmd))
			}
			leftCmd = append(leftCmd, byte(buf))
		}
	}
	cmd = append(cmd, leftCmd...)
	cmd = append(cmd, rightCmd...)
	if tm.history {
		if len(strings.Trim(string(cmd), " ")) != 0 {
			if len(tm.historyList) == 0 {
				tm.historyList = append(tm.historyList, strings.Trim(string(cmd), " "))
			} else if tm.historyList[len(tm.historyList)-1] != strings.Trim(string(cmd), " ") {
				if len(tm.historyList) >= SYS_MAX_HISTORY {
					tm.historyList = tm.historyList[1:]
				}
				tm.historyList = append(tm.historyList, strings.Trim(string(cmd), " "))
			}
		}
		tm.historyIndex = len(tm.historyList)
	}
	return strings.Trim(string(cmd), " ")
}

func cleanCell(num int) {
	for i := 0; i < num; i++ {
		fmt.Printf(" ")
	}
}

func cleanLine() {
	fmt.Printf("\033[K")
}

func (tm *Terminal) GetCommand() string {
	tm.History(true)
	tm.Echo(true)
	return tm.getInput()
}

func (tm *Terminal) GetUser() string {
	tm.History(false)
	tm.Echo(true)
	return tm.getInput()
}

func (tm *Terminal) GetPassword() string {
	tm.History(false)
	tm.Echo(false)
	return tm.getInput()
}
