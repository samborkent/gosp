package gsp

type SampleProcessor[F Frame[T], T Type] interface {
	// This function gets called for each sample of the input buffer. It should contain the signal processing logic.
	// For stateful algorithms such as filters, the processor is responsible for keeping state.
	Process(inputSample F) (outputSample F)
}

type BufferProcessor[F Frame[T], T Type] interface {
	// This function gets called for each buffer of an input signal. It should contain the signal processing logic.
	// For stateful algorithms such as filters, the processor is responsible for keeping state.
	// The pipeline will assure that the processor is fed a continuous signal by sending a zero buffer in case no input was provided within the sample rate clock interval.
	ProcessBuffer(outputBuffer, inputBuffer []F)
}
