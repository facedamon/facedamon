package utils

import (
	"math/rand"
	"testing"
)

func TestSortInts(t *testing.T) {
	ints := []interface{}{4, 1, 2, 3}
	Sort(ints, IntComparator)

	for i := 1; i < len(ints); i++ {
		if ints[i-1].(int) > ints[i].(int) {
			t.Errorf("Not sorted!")
		}
	}
}

func TestSortString(t *testing.T) {
	strings := []interface{}{"d", "a", "b", "c"}

	Sort(strings, StringComparator)

	for i := 1; i < len(strings); i++ {
		if strings[i-1].(string) > strings[i].(string) {
			t.Errorf("Not sorted")
		}
	}
}

func TestSortRandon(t *testing.T) {
	ints := []interface{}{}
	for i := 0; i < 10000; i++ {
		ints = append(ints, rand.Int())
	}
	Sort(ints, IntComparator)
	for i := 1; i < len(ints); i++ {
		if ints[i-1].(int) > ints[i].(int) {
			t.Errorf("Not sorted!")
		}
	}
}

func BenchmarkGoSortRandom(b *testing.B) {
	b.StopTimer()
	ints := []interface{}{}
	for i := 0; i < 100000; i++ {
		ints = append(ints, rand.Int())
	}
	b.StartTimer()
	Sort(ints, IntComparator)
	b.StopTimer()
}
