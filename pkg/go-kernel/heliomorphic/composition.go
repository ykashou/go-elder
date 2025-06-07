package heliomorphic

import "math/cmplx"

type FunctionComposer struct {
	Functions map[string]func(complex128) complex128
	Cache     map[string]complex128
}

func NewFunctionComposer() *FunctionComposer {
	return &FunctionComposer{
		Functions: make(map[string]func(complex128) complex128),
		Cache:     make(map[string]complex128),
	}
}

func (fc *FunctionComposer) RegisterFunction(name string, f func(complex128) complex128) {
	fc.Functions[name] = f
}

func (fc *FunctionComposer) Compose(f1Name, f2Name string) func(complex128) complex128 {
	f1 := fc.Functions[f1Name]
	f2 := fc.Functions[f2Name]
	
	return func(z complex128) complex128 {
		return f1(f2(z))
	}
}

func (fc *FunctionComposer) ChainRule(f1Name, f2Name string, z complex128) complex128 {
	f1 := fc.Functions[f1Name]
	f2 := fc.Functions[f2Name]
	
	h := 1e-8
	
	f2_z := f2(z)
	df1_df2 := (f1(f2_z+complex(h, 0)) - f1(f2_z-complex(h, 0))) / complex(2*h, 0)
	df2_dz := (f2(z+complex(h, 0)) - f2(z-complex(h, 0))) / complex(2*h, 0)
	
	return df1_df2 * df2_dz
}
