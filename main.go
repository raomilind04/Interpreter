package main

import (
    "fmt"
    "os"
    "os/user"
    "interpreter/repl"
)

func main() {
    user, err := user.Current()
    if err != nil {
        panic(err)
    }
    fmt.Printf("Hello %s! \n", user.Username)
    fmt.Printf("REPL Started\n")
    repl.Start(os.Stdin, os.Stdout)
}
