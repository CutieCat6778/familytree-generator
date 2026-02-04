package rand

import (
	"math/rand"
	"sync"
)

// SeededRandom provides a reproducible random number generator
type SeededRandom struct {
	rng *rand.Rand
	mu  sync.Mutex
}

// New creates a new SeededRandom with the given seed
func New(seed int64) *SeededRandom {
	return &SeededRandom{
		rng: rand.New(rand.NewSource(seed)),
	}
}

// Int returns a non-negative pseudo-random int
func (r *SeededRandom) Int() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rng.Int()
}

// Intn returns a non-negative pseudo-random int in [0, n)
func (r *SeededRandom) Intn(n int) int {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rng.Intn(n)
}

// Float64 returns a pseudo-random float64 in [0.0, 1.0)
func (r *SeededRandom) Float64() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rng.Float64()
}

// NormFloat64 returns a normally distributed float64 with mean 0 and stddev 1
func (r *SeededRandom) NormFloat64() float64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rng.NormFloat64()
}

// IntRange returns a pseudo-random int in [min, max]
func (r *SeededRandom) IntRange(min, max int) int {
	if min >= max {
		return min
	}
	return min + r.Intn(max-min+1)
}

// Float64Range returns a pseudo-random float64 in [min, max)
func (r *SeededRandom) Float64Range(min, max float64) float64 {
	return min + r.Float64()*(max-min)
}

// Bool returns a random boolean with 50% probability
func (r *SeededRandom) Bool() bool {
	return r.Float64() < 0.5
}

// Chance returns true with the given probability (0.0 to 1.0)
func (r *SeededRandom) Chance(probability float64) bool {
	return r.Float64() < probability
}

// Choice returns a random element from the slice
func Choice[T any](r *SeededRandom, slice []T) T {
	return slice[r.Intn(len(slice))]
}

// WeightedChoice returns an index based on weights
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

// Shuffle randomly shuffles a slice in place
func Shuffle[T any](r *SeededRandom, slice []T) {
	for i := len(slice) - 1; i > 0; i-- {
		j := r.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

// NormalDistribution returns a value from a normal distribution with given mean and stddev
func (r *SeededRandom) NormalDistribution(mean, stddev float64) float64 {
	return mean + r.NormFloat64()*stddev
}
