package main

type cost struct {
	day   int
	value float64
}

func getDayCosts(costs []cost, day int) []float64 {
	// ?
	ans := []float64{}
	for i := range costs {
		if costs[i].day == day {

			ans = append(ans, costs[i].value)
		}
	}
	return ans
}
