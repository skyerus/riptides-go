package main

import "github.com/skyerus/riptides-go/pkg/api"

func main()  {
	main := &api.App{}
	main.Initialize()
	main.Run(":80")
}
