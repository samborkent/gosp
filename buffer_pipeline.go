package gsp

import (
	"context"
	"time"
)

// TODO: consider making channels buffers
// TODO: add back buffer processor support
// TODO: add parallel processing support based on available threads for buffer processing
// TODO: add SIMD support if possible?

type BufferPipeline[F Frame[T], T Type] struct {
	sampleRate, bufferSize int

	clock                  *time.Ticker
	input, bufferedInput   chan []F
	output, bufferedOutput chan []F

	circularInputBuffer, circularOutputBuffer *circularBuffer[[]F, F, T]

	overflowCount uint64

	processors []BufferProcessor[F, T]

	done chan struct{}
}

// When the buffer size is set to zero, the pipeline will operate in sample mode.
// This means it will implement SampleReader and SampleWriter.
// If the buffer size is greater than zero, it will implement Reader and Writer.
func NewBufferPipeline[F Frame[T], T Type](sampleRate, bufferSize int, processors ...BufferProcessor[F, T]) (*BufferPipeline[F, T], error) {
	var p BufferPipeline[F, T]

	p.sampleRate = sampleRate
	p.bufferSize = bufferSize
	p.processors = processors
	p.done = make(chan struct{})

	p.clock = time.NewTicker(time.Duration(float64(time.Second) * float64(bufferSize) / float64(sampleRate)))

	const circularBufferLength = 3

	p.input = make(chan []F)
	p.bufferedInput = make(chan []F, circularBufferLength)
	p.output = make(chan []F)
	p.bufferedOutput = make(chan []F, circularBufferLength)

	p.circularInputBuffer = newCircularBuffer[[]F, F, T](p.input, p.bufferedInput)
	p.circularOutputBuffer = newCircularBuffer[[]F, F, T](p.output, p.bufferedOutput)

	return &p, nil
}

func (p *BufferPipeline[F, T]) Read(s []F) (n int, err error) {
	return copy(s, <-p.bufferedOutput), nil
}

func (p *BufferPipeline[F, T]) Write(s []F) (n int, err error) {
	p.input <- s
	return len(s), nil
}

func (p *BufferPipeline[F, T]) OverflowCount() uint64 {
	return p.overflowCount + p.circularInputBuffer.OverflowCount() + p.circularOutputBuffer.OverflowCount()
}

func (p *BufferPipeline[F, T]) Start(ctx context.Context) {
	go p.circularInputBuffer.Run(ctx)
	go p.circularInputBuffer.Run(ctx)

	input := make([]F, p.bufferSize)

	// Flush processors
	p.processBuffer(input, input)

whileLoop:
	for {
		select {
		case <-ctx.Done():
			break whileLoop
		case input = <-p.bufferedInput:
			// If input comes before the clock tick, we still have to wait for the clock tick to synchronize
			<-p.clock.C
		case <-p.clock.C:
			// In case no input comes in due time, send a buffer filled with zeros
			clear(input)
		}

		// TODO: use sync.Pool to avoid allocation
		output := make([]F, len(input))
		p.processBuffer(output, input)

		select {
		case p.output <- output:
		default:
			p.overflowCount++
		}
	}

	close(p.done)
}

func (p *BufferPipeline[F, T]) Close() error {
	p.clock.Stop()
	<-p.done
	close(p.input)
	close(p.output)

	return nil
}

func (p *BufferPipeline[F, T]) processBuffer(output, input []F) {
	for _, processor := range p.processors {
		processor.ProcessBuffer(output, input)
	}
}
