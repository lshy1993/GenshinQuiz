package main

import (
	"genshin-quiz/config"
	"genshin-quiz/internal/console"
)

func main() {
	a := config.NewApp()
	c := console.NewCommandHandler(a)
}
