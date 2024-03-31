package processor

import (
	"context"
	"github.com/ahmetb/go-linq/v3"
	"github.com/gamenotover/mini_processor/concurrent"
	"github.com/gamenotover/mini_processor/consts"
	"github.com/gamenotover/mini_processor/model"
	"github.com/gamenotover/mini_processor/sets"
	"sync"
)

type Processor interface {
	Tag() consts.ProcessorTag
	Wanted() []consts.ProcessorTag
	Process(ctx context.Context, tag2val map[consts.ProcessorTag]*model.BasicInfo) *model.BasicInfo
}

var processorFns []func() Processor

func RegisterFn(fn func() Processor) {
	processorFns = append(processorFns, fn)
}

func NewProcessors() []Processor {
	res := make([]Processor, 0, len(processorFns))
	for _, fn := range processorFns {
		res = append(res, fn())
	}
	return res
}

func AssertValid() {
	tags := make([]consts.ProcessorTag, 0, len(processorFns))
	wanted := make([]consts.ProcessorTag, 0, len(processorFns))
	tag2processor := make(map[consts.ProcessorTag]Processor, len(processorFns))
	for _, fn := range processorFns {
		processor := fn()
		tags = append(tags, processor.Tag())
		wanted = append(wanted, processor.Wanted()...)
		tag2processor[processor.Tag()] = processor
	}
	// 检查是否存在根节点
	var rootTags []consts.ProcessorTag
	linq.From(tags).Except(linq.From(wanted)).ToSlice(&rootTags)
	if len(rootTags) == 0 {
		panic("no root tag!")
	}
	// 检查是否成环
	var fn func(pc Processor, tags []consts.ProcessorTag)
	fn = func(pc Processor, tags []consts.ProcessorTag) {
		if linq.From(tags).Contains(pc.Tag()) {
			panic("has cycle!")
		}
		for _, tag := range pc.Wanted() {
			next := tag2processor[tag]
			fn(next, append(tags, pc.Tag()))
		}
	}
	for _, tag := range rootTags {
		fn(tag2processor[tag], []consts.ProcessorTag{})
	}
}

type node struct {
	p      Processor
	wanted []*node
	mutex  sync.Mutex
	once   sync.Once
	res    *model.BasicInfo
}

func (n *node) run(ctx context.Context) {
	async := concurrent.NewAsyncController()
	syncMap := async.NewMap()
	fns := make([]func(), 0, len(n.wanted))
	for _, _wanted := range n.wanted {
		wanted := _wanted
		fns = append(fns, func() {
			syncMap.Store(wanted.p.Tag(), wanted.Run(ctx))
		})
	}
	async.Do(ctx, fns...)
	tag2val := make(map[consts.ProcessorTag]*model.BasicInfo, len(fns))
	syncMap.Range(func(key, value interface{}) {
		tag2val[key.(consts.ProcessorTag)] = value.(*model.BasicInfo)
	})
	n.res = n.p.Process(ctx, tag2val)
}

func (n *node) Run(ctx context.Context) *model.BasicInfo {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.once.Do(func() {
		// 打点日志，或上报埋点
		n.run(ctx)
	})
	return n.res
}

// Run 运行已经注册的Processor
func Run(ctx context.Context) *model.BasicInfo {
	tag2node := make(map[consts.ProcessorTag]*node, len(processorFns))
	tagSets := sets.NewString()
	processors := NewProcessors()
	for _, p := range processors {
		tag2node[p.Tag()] = &node{p: p}
		tagSets.Add(string(p.Tag()))
	}
	for _, p := range processors {
		wanted := make([]*node, 0, len(p.Wanted()))
		for _, tag := range p.Wanted() {
			wanted = append(wanted, tag2node[tag])
			tagSets.Remove(string(tag))
			tag2node[p.Tag()].wanted = wanted
		}
	}
	dp := &node{p: &done{wanted: make([]consts.ProcessorTag, 0, len(tagSets))}}
	for t := range tagSets {
		tag := consts.ProcessorTag(t)
		dp.wanted = append(dp.wanted, tag2node[tag])
	}
	return dp.Run(ctx)
}

type done struct {
	wanted []consts.ProcessorTag
}

func (d done) Tag() consts.ProcessorTag {
	return consts.DoneProcessor
}

func (d done) Wanted() []consts.ProcessorTag {
	return d.wanted
}

func (d done) Process(ctx context.Context, tag2val map[consts.ProcessorTag]*model.BasicInfo) *model.BasicInfo {
	return &model.BasicInfo{}
}
