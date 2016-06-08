# terminal
A custom terminal witch can use up/down press to call history command, and tab to call syscall command.
## use
getch
```golang
ch, paste := getch()  //ch is the keyboard value you've press，paste is the value if you use command+v
```
if you'll create a terminal, you can do like so:
```golang
	systemCmdList := []string{"hosts", "services", "connect", "connects", "relase", "log", "start", "stop", "restart", "delete", "exit"}
	term := NewTerminal("> ")
	term.SetSystemCommand(systemCmdList)
	for {
		fmt.Printf("> ")
		cmd := term.GetCommand()
		fmt.Println()
		//fmt.Println("\ncmd = ", cmd, "  len = ", len(cmd))
		if cmd == "exit" {
			break
		}
	}
```
the system command list will used when you press `tab` keyboard. if you use func `GetCommand`,you can call history command by ⬆️ and ⬇️.
if terminal request user to input username then you can call `term.GetUser()`, if so, it will not record to history list, at the same time, if get password you can call `term.GetPassword()`, it will echo '*' when you input password.
