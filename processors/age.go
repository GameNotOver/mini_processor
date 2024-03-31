package processors

import (
	"context"
	"github.com/gamenotover/mini_processor/consts"
	"github.com/gamenotover/mini_processor/model"
	"github.com/gamenotover/mini_processor/processor"
)

type ageProcessor struct {
}

func NewAgeProcessor() processor.Processor {
	return &ageProcessor{}
}

func (n ageProcessor) Tag() consts.ProcessorTag {
	return consts.AgeProcessor
}

func (n ageProcessor) Wanted() []consts.ProcessorTag {
	return nil
}

func (n ageProcessor) Process(ctx context.Context, tag2val map[consts.ProcessorTag]*model.BasicInfo) *model.BasicInfo {
	age := 24
	return &model.BasicInfo{UserInfo: &model.UserInfo{
		Age: age,
	}}
}
