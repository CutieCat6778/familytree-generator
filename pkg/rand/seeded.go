package rand

import (
	"math/rand"
	"sync"
)


type SeededRandom struct {
	rng *rand.Rand
	mu  sync.Mutex
}


func New(seed int64) *SeededRandom {
	return &SeededRandom{
		rng: rand.New(rand.NewSource(seed)),
	}
}


func (r *SeededRandom) Int() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rng.Int()
}


func (r *SeededRandom) Intn(n int) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rng.Intn(n)
}


func (r *SeededRandom) Float64() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rng.Float64()
}


func (r *SeededRandom) NormFloat64() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rng.NormFloat64()
}


func (r *SeededRandom) IntRange(min, max int) int {
	if min >= max {
		return min
	}
	return min + r.Intn(max-min+1)
}


func (r *SeededRandom) Float64Range(min, max float64) float64 {
	return min + r.Float64()*(max-min)
}


func (r *SeededRandom) Bool() bool {
	return r.Float64() < 0.5
}


func (r *SeededRandom) Chance(probability float64) bool {
	return r.Float64() < probability
}


func Choice[T any](r *SeededRandom, slice []T) T {
	return slice[r.Intn(len(slice))]
}


func (r *SeededRandom) WeightedChoice(weights []float64) int {
	total := 0.0
	for _, w := range weights {
		total += w
	}

	threshold := r.Float64() * total
	cumulative := 0.0

	for i, w := range weights {
		cumulative += w
		if threshold < cumulative {
			return i
		}
	}

	return len(weights) - 1
}


func Shuffle[T any](r *SeededRandom, slice []T) {
	for i := len(slice) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}


func (r *SeededRandom) NormalDistribution(mean, stddev float64) float64 {
	return mean + r.NormFloat64()*stddev
}
