package main

import (
	"fmt"
	"smtpgo"
)

func main() {
	err := smtpgo.StartServer("", "SMTP go", "testhost", ":25")
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
