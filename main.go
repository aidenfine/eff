package main

import (
	"fmt"
	"time"

	"github.com/aidenfine/eff/cmd"
)

func main() {
	start := time.Now()
	cmd.Execute()
	t := time.Now()
	fmt.Println(t.Sub(start))

}
