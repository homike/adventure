package util

import "math/rand"

type RandomWeight struct {
	IDs     []int32
	Weights []int32

	Total int32
	Range []int32
}

// 测试提交123123
func NewRandom(IDs, Weights []int32) *RandomWeight {
	r := &RandomWeight{
		IDs:     IDs,
		Weights: Weights,
		Total:   0,
		Range:   []int32{},
	}

	r.Total = 0
	r.Range = make([]int32, len(r.IDs))
	for i := 0; i < len(r.IDs); i++ {
		r.Total += r.Weights[i]
		r.Range[i] = r.Total
	}

	return r
}

func (r *RandomWeight) GetRandomNum() int32 {
	randIndex := RandNum(0, r.Total)
	for i := 0; i < len(r.Range); i++ {
		if r.Range[i] > randIndex {
			return r.Range[i]
		}
	}
	return r.IDs[0]
}

func RandNum(base int32, n int32) int32 {
	if n <= 0 {
		return base
	}
	return base + int32(rand.Intn(int(n)))
}
