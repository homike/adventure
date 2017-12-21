package util

import "math/rand"

type RandomWeight struct {
	IDs     []int
	Weights []int

	Total int
	Range []int
}

func NewRandom(IDs, Weights []int) *RandomWeight {
	r := &RandomWeight{
		IDs:     IDs,
		Weights: Weights,
		Total:   0,
		Range:   []int{},
	}

	r.Total = 0
	for i := 0; i < len(r.IDs); i++ {
		r.Total += r.Weights[i]
		r.Range[i] = r.Total
	}

	return r
}

func (r *RandomWeight) GetRandomNum() int {
	randIndex := RandNum(0, r.Total)
	for i := 0; i < len(r.Range); i++ {
		if r.Range[i] > randIndex {
			return r.Range[i]
		}
	}
	return r.IDs[0]
}

func RandNum(base int, n int) int {
	if n <= 0 {
		return base
	}
	return base + rand.Intn(n)
}
