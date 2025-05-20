package main

import (

	"github.com/walterfan/llm-agent-go/cmd"

)



func main() {

	cmd.Execute()
}

func init() {

	if err := cmd.InitConfig(); err != nil {
		panic(err)
	}
}