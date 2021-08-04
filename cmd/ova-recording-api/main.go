package main

import (
	"fmt"
	"github.com/enescakir/emoji"
)

func main() {
	fmt.Printf("hello world %v\n", emoji.WavingHand.Tone(emoji.Light))
}
