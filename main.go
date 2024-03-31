package main

import (
	"context"
	"github.com/gamenotover/mini_processor/processor"
	"github.com/gamenotover/mini_processor/processors"
)

func init() {
	processor.RegisterFn(processors.NewAgeProcessor)
	processor.RegisterFn(processors.NewGenderProcessor)
	processor.RegisterFn(processors.NewNameProcessor)
	processor.RegisterFn(processors.NewUserInfoProcessor)
	processor.AssertValid()
}

func main() {
	processor.Run(context.Background())
}
