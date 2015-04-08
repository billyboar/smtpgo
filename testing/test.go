package main

import (
	"fmt"

	"github.com/billyboar/smtpgo"
)

func main() {
	err := smtpgo.StartServer("", "SMTP go", "testhost", ":25")
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
