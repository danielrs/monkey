package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/danielrs/monkey/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	if len(os.Args) == 2 {
		// Tries to read from file.
		data, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		repl.Run(string(data), os.Stdout)
	} else {
		// Runs REPL.
		fmt.Printf("Hello %s!\n", user.Username)
		fmt.Printf("This is the Monkey programming language!\n")
		fmt.Printf("Feel free to type in commands\n")
		repl.Start(os.Stdin, os.Stdout)
	}

}
