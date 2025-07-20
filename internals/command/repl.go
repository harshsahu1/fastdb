package command

import (
	"bufio"
	"fmt"
	"os"
)
const banner = `
███████╗███████╗███████╗████████╗██████╗ ██████╗
██╔════╝██╔══██║██╔════╝╚══██╔══╝██╔══██╗██╔══██╗
█████╗  ███████║███████╗   ██║   ██║  ██║██████╔╝
██╔══╝  ██╔══██║╚════██║   ██║   ██║  ██║██╔══██╗
██║     ██║  ██║███████║   ██║   ██████╔╝██████╔╝
╚═╝     ╚═╝  ╚═╝╚══════╝   ╚═╝   ╚═════╝ ╚═════╝
        ⚡ FastDB — blazing fast KV store
`

func StartREPL(executor *Executor) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println(banner) 
	fmt.Println("🚀 Welcome to FastDB CLI (type 'exit' to quit)")
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		if line == "exit" {
			break
		}

		cmd, err := ParseCommand(line)
		if err != nil || cmd == nil {
			fmt.Println("(error)", err)
			continue
		}

		result, err := executor.ExecuteCommand(cmd)
		if err != nil {
			fmt.Println("(error)", err)
		} else {
			fmt.Println(result)
		}
	}
}
