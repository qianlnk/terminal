package terminal

import (
	"fmt"
	"testing"
)

//func TestUserPassword(t *testing.T) {
//	for {
//		//getch()
//		fmt.Printf("user:")
//		user := GetInput(true, false)
//		if user == "exit" {
//			break
//		}
//		fmt.Printf("\npassword:")
//		password := GetInput(false, false)
//		fmt.Println(user, password)
//	}
//}

//func TestCmd(t *testing.T) {
//	for {
//		fmt.Printf("> ")
//		cmd := GetInput(true, true)
//		fmt.Println("\ncmd = ", cmd, "  len = ", len(cmd))
//	}
//}

func TestTab(t *testing.T) {
	cmdList := []string{"hosts", "services", "connect", "connects", "relase", "log", "start", "stop", "restart", "delete", "exit", "test2", "test1", "test3", "tes41", "test5", "test6", "tes71", "test8", "test9", "test10", "test11", "test12", "test13", "test14", "test15"}
	term := NewTerminal("> ")
	term.SetSystemCommand(cmdList)
	term.History(true)
	term.Echo(true)
	for {
		fmt.Printf("> ")
		cmd := term.getInput()
		fmt.Println()
		//fmt.Println("\ncmd = ", cmd, "  len = ", len(cmd))
		if cmd == "exit" {
			break
		}
	}
}
