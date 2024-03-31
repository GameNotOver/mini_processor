package processors

import (
	"context"
	"github.com/gamenotover/mini_processor/consts"
	"github.com/gamenotover/mini_processor/model"
	"github.com/gamenotover/mini_processor/processor"
)

type genderProcessor struct {
}

func NewGenderProcessor() processor.Processor {
	return &genderProcessor{}
}

func (n genderProcessor) Tag() consts.ProcessorTag {
	return consts.GenderProcessor
}

func (n genderProcessor) Wanted() []consts.ProcessorTag {
	return nil
}

func (n genderProcessor) Process(ctx context.Context, tag2val map[consts.ProcessorTag]*model.BasicInfo) *model.BasicInfo {
	gender := consts.MALE
	return &model.BasicInfo{UserInfo: &model.UserInfo{
		Gender: gender,
	}}
}
