package processors

import (
	"context"
	"github.com/gamenotover/mini_processor/consts"
	"github.com/gamenotover/mini_processor/model"
	"github.com/gamenotover/mini_processor/processor"
)

type nameProcessor struct {
}

func NewNameProcessor() processor.Processor {
	return &nameProcessor{}
}

func (n nameProcessor) Tag() consts.ProcessorTag {
	return consts.NameProcessor
}

func (n nameProcessor) Wanted() []consts.ProcessorTag {
	return nil
}

func (n nameProcessor) Process(ctx context.Context, tag2val map[consts.ProcessorTag]*model.BasicInfo) *model.BasicInfo {
	name := "wuyuhang"
	return &model.BasicInfo{UserInfo: &model.UserInfo{
		Name: name,
	}}
}
