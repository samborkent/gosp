package gsp

import (
	"context"
	"runtime"
)

type SamplePipeline[F Frame[T], T Type] struct {
	processors []SampleProcessor[F, T]

	input, output chan F

	initialized bool
}

func NewSamplePipeline[F Frame[T], T Type](processors ...SampleProcessor[F, T]) *SamplePipeline[F, T] {
	if len(processors) == 0 {
		return &SamplePipeline[F, T]{}
	}

	pipeline := &SamplePipeline[F, T]{
		processors:  processors,
		input:       make(chan F, 1),
		output:      make(chan F, 1),
		initialized: true,
	}

	ctx, cancel := context.WithCancel(context.Background())

	runtime.AddCleanup(pipeline, func(_ int) {
		cancel()
		close(pipeline.input)
		close(pipeline.output)
	}, 0)

	go pipeline.run(ctx)

	return pipeline
}

func (p *SamplePipeline[F, T]) Loop() <-chan F {
	return p.output
}

func (p *SamplePipeline[F, T]) ReadSample() F {
	return <-p.output
}

func (p *SamplePipeline[F, T]) WriteSample(s F) {
	p.input <- s
}

func (p *SamplePipeline[F, T]) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case input := <-p.input:
			p.output <- p.process(input)
		}
	}
}

func (p *SamplePipeline[F, T]) process(sample F) F {
	for _, processor := range p.processors {
		sample = processor.Process(sample)
	}

	return sample
}
