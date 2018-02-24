package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/racoonberus/media/cmd/transcodersvc/preset"
)

func main() {
	inFiles:=strings.Split(os.Args[1] , " ")

	var p preset.WebVideo
	out, err := p.Execute(inFiles)
	if nil != err {
		panic(err.Error())
	}

	for _, val := range out {
		fmt.Println(val)
	}
}