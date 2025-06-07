package heliomorphic

type ConvolutionKernel struct {
	Kernel []complex128
	Size   int
}

func (ck *ConvolutionKernel) Convolve(input []complex128) []complex128 {
	output := make([]complex128, len(input))
	for i := range input {
		sum := complex(0, 0)
		for j := range ck.Kernel {
			if i-j >= 0 && i-j < len(input) {
				sum += ck.Kernel[j] * input[i-j]
			}
		}
		output[i] = sum
	}
	return output
}
