package rest

import "fmt"

const header = `
package main

import "github.com/gofiber/fiber/v2"

`

func makeHandlerHeader(methodName string) string {
	return fmt.Sprintf("func %s(ctx *fiber.Ctx) error {\n")
}

func declareApiPrefix(prefix string) string {
	return fmt.Sprintf("\tg := mainApi.Group(%q)\n", prefix)
}

func attatchHandler(verb, url, methodName string) string {
	return fmt.Sprintf("\tg.%s(%q, %s)", verb, url, methodName)
}

func declareServerListen(port int) string {
	return fmt.Sprintf("\tmainApi.Listen(\":%d\")\n", port)
}

const funcMainHeader = `
func main() {
	mainApi := fiber.New()
`

const handlerFooter = `
}
`

const funcMainFooter = `
}
`
