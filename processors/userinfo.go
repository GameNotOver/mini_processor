package processors

import (
	"context"
	"fmt"
	"github.com/gamenotover/mini_processor/consts"
	"github.com/gamenotover/mini_processor/model"
	"github.com/gamenotover/mini_processor/processor"
)

type userInfoProcessor struct {
}

func NewUserInfoProcessor() processor.Processor {
	return &userInfoProcessor{}
}

func (s userInfoProcessor) Tag() consts.ProcessorTag {
	return consts.UserInfoProcessor
}

func (s userInfoProcessor) Wanted() []consts.ProcessorTag {
	return []consts.ProcessorTag{consts.NameProcessor, consts.AgeProcessor, consts.GenderProcessor}
}

func (s userInfoProcessor) Process(ctx context.Context, tag2val map[consts.ProcessorTag]*model.BasicInfo) *model.BasicInfo {
	name := tag2val[consts.NameProcessor].UserInfo.Name
	age := tag2val[consts.AgeProcessor].UserInfo.Age
	gender := tag2val[consts.GenderProcessor].UserInfo.Gender
	fmt.Printf("name: %v\n", *tag2val[consts.NameProcessor].UserInfo)
	fmt.Printf("age: %v\n", *tag2val[consts.AgeProcessor].UserInfo)
	fmt.Printf("gender: %v\n", *tag2val[consts.GenderProcessor].UserInfo)
	return &model.BasicInfo{UserInfo: &model.UserInfo{
		Name:   name,
		Age:    age,
		Gender: gender,
	}}
}
