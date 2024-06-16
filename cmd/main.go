package main

import (
	"parser/internal/parser"
	"parser/internal/webdriver"
)

func main() {

	wd, chrome := webdriver.StartWebDriver()
	parser.OpenCategoryParser(wd)
	webdriver.StopApp(wd, chrome)

}
