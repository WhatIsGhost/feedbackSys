package main

import (
	_ "feedbackSys/routers"
	_ "feedbackSys/sysinit"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}

