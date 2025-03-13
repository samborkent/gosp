package gsp

import (
	"context"
	"time"
)

// TODO: consider making channels buffers
// TODO: add back buffer processor support
// TODO: add parallel processing support based on available threads for buffer processing
// TODO: add SIMD support if possible?

type SamplePipeline[F Frame[T], T Type] struct {
	sampleRate int

	clock                  *time.Ticker
	input, bufferedInput   chan F
	output, bufferedOutput chan F

	circularInputBuffer, circularOutputBuffer *circularBuffer[F, F, T]

	overflowCount uint64

	processors []SampleProcessor[F, T]

	done chan struct{}
}

// When the buffer size is set to zero, the pipeline will operate in sample mode.
// This means it will implement SampleReader and SampleWriter.
// If the buffer size is greater than zero, it will implement Reader and Writer.
func NewSamplePipeline[F Frame[T], T Type](sampleRate int, processors ...SampleProcessor[F, T]) (*SamplePipeline[F, T], error) {
	var p SamplePipeline[F, T]

	p.sampleRate = sampleRate
	p.processors = processors
	p.done = make(chan struct{})

	p.clock = time.NewTicker(time.Duration(float64(time.Second) / float64(sampleRate)))

	const circularBufferLength = 3

	p.input = make(chan F)
	p.bufferedInput = make(chan F, circularBufferLength)
	p.output = make(chan F)
	p.bufferedOutput = make(chan F, circularBufferLength)

	p.circularInputBuffer = newCircularBuffer[F, F, T](p.input, p.bufferedInput)
	p.circularOutputBuffer = newCircularBuffer[F, F, T](p.output, p.bufferedOutput)

	return &p, nil
}

func (p *SamplePipeline[F, T]) ReadSample() (F, error) {
	return <-p.bufferedOutput, nil
}

func (p *SamplePipeline[F, T]) WriteSample(s F) error {
	p.input <- s
	return nil
}

func (p *SamplePipeline[F, T]) OverflowCount() uint64 {
	return p.overflowCount + p.circularInputBuffer.OverflowCount() + p.circularOutputBuffer.OverflowCount()
}

func (p *SamplePipeline[F, T]) Start(ctx context.Context) {
	go p.circularInputBuffer.Run(ctx)
	go p.circularOutputBuffer.Run(ctx)

	// Flush processors
	p.processSample(*new(F))

whileLoop:
	for {
		select {
		case <-ctx.Done():
			break whileLoop
		case input := <-p.bufferedInput:
			// If input comes before the clock tick, we still have to wait for the clock tick to synchronize
			<-p.clock.C

			select {
			case p.output <- p.processSample(input):
			default:
				p.overflowCount++
			}
		case <-p.clock.C:
			// In case no input comes in due time, send a zero sample
			select {
			case p.output <- p.processSample(*new(F)):
			default:
				p.overflowCount++
			}
		}
	}

	close(p.done)
}

func (p *SamplePipeline[F, T]) Close() error {
	p.clock.Stop()
	<-p.done

	close(p.input)
	close(p.bufferedInput)
	close(p.output)
	close(p.bufferedOutput)

	return nil
}

func (p *SamplePipeline[F, T]) processSample(sample F) F {
	for _, processor := range p.processors {
		sample = processor.Process(sample)
	}

	return sample
}
