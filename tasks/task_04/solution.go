package main

import "math"

type Stats struct {
	Count         int
	Sum, Min, Max int64
}

func Calc(nums []int64) Stats {
	if len(nums) < 2 {
		return Stats{}
	}

	res := Stats{
		Count: len(nums) - 1,
		Min:   math.MaxInt64,
		Max:   math.MinInt64,
	}

	for i := 1; i < len(nums); i++ {
		d := nums[i] - nums[i-1]
		res.Sum += d
		if d < res.Min {
			res.Min = d
		}
		if d > res.Max {
			res.Max = d
		}
	}

	return res
}
