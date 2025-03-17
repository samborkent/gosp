package gsp

import (
	"context"
	"runtime"
)

type Pipeline[F Frame[T], T Type] struct {
	processors []BufferProcessor[F, T]

	input, output chan []F

	pool *FramePool[F, T]

	initialized bool
}

func NewPipeline[F Frame[T], T Type](processors ...BufferProcessor[F, T]) *Pipeline[F, T] {
	if len(processors) == 0 {
		return &Pipeline[F, T]{}
	}

	pipeline := &Pipeline[F, T]{
		processors:  processors,
		input:       make(chan []F, 1),
		output:      make(chan []F, 1),
		pool:        NewFramePool[F, T](1024),
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

func (p *Pipeline[F, T]) Loop() <-chan []F {
	return p.output
}

func (p *Pipeline[F, T]) Read(output []F) (int, error) {
	return copy(output, <-p.output), nil
}

func (p *Pipeline[F, T]) Write(input []F) (int, error) {
	p.input <- input
	return len(input), nil
}

func (p *Pipeline[F, T]) run(ctx context.Context) {
	output := []F{}

	for {
		select {
		case <-ctx.Done():
			return
		case input := <-p.input:
			bufPtr := p.pool.Get()
			output = *bufPtr

			if len(input) > len(output) {
				output = make([]F, len(input))
			}

			p.process(output, input)
			p.output <- output[:len(input)]

			clear(output)
			bufPtr = &output
			p.pool.Put(bufPtr)
		}
	}
}

func (p *Pipeline[F, T]) process(output, input []F) {
	for _, processor := range p.processors {
		processor.ProcessBuffer(output, input)
	}
}
