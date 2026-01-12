package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
)

type TransitionCount map[string]int

type keyProbability struct {
	countName string
	nCount    int
}

func main() {
	lambda0 := math.Pow(200, -1)
	lambda1 := math.Pow(400, -1)
	mu := math.Pow(250, -1)

	ET0 := 1 / lambda0
	ET1 := 1 / lambda1
	ETmu := 1 / mu

	P10 := mu / (lambda0 + mu)
	P01 := mu / (lambda1 + mu)
	P1T := lambda0 / (lambda0 + mu)
	P0T := lambda1 / (lambda1 + mu)

	transitionCount := make(TransitionCount)
	transitionCount["P10"] = 0
	transitionCount["P01"] = 0
	transitionCount["P1T"] = 0
	transitionCount["P0T"] = 0

	counter := 100000

	fmt.Printf("Parameters:\n")
	fmt.Printf("ET0: %.2f, ET1: %.2f, ETmu: %.2f\n", ET0, ET1, ETmu)
	fmt.Printf("P10 (1→0): %.4f\n", P10)
	fmt.Printf("P01 (0→1): %.4f\n", P01)
	fmt.Printf("P1T (1→T): %.4f\n", P1T)
	fmt.Printf("P0T (0→T): %.4f\n", P0T)
	fmt.Println()

	n1 := runSimulations(counter, ET0, ET1, P10, P01, P1T, P0T, lambda0, lambda1, mu, transitionCount)

	displayProbabilities(n1, counter)

	totalTransitions := transitionCount["P10"] + transitionCount["P01"] + transitionCount["P1T"] + transitionCount["P0T"]
	fmt.Printf("\nTransition Counts:\n")
	fmt.Printf("P10 (1→0): %d\n", transitionCount["P10"])
	fmt.Printf("P01 (0→1): %d\n", transitionCount["P01"])
	fmt.Printf("P1T (1→T): %d\n", transitionCount["P1T"])
	fmt.Printf("P0T (0→T): %d\n", transitionCount["P0T"])
	fmt.Printf("Total transitions: %d\n", totalTransitions)

	fmt.Printf("\nEmpirical Transition Probabilities:\n")
	from1 := transitionCount["P10"] + transitionCount["P1T"]
	if from1 > 0 {
		p10_sim := float64(transitionCount["P10"]) / float64(from1)
		p1T_sim := float64(transitionCount["P1T"]) / float64(from1)
		fmt.Printf("P(1→0): %.6f (expected: %.4f)\n", p10_sim, P10)
		fmt.Printf("P(1→T): %.6f (expected: %.4f)\n", p1T_sim, P1T)
	}

	from0 := transitionCount["P01"] + transitionCount["P0T"]
	if from0 > 0 {
		p01_sim := float64(transitionCount["P01"]) / float64(from0)
		p0T_sim := float64(transitionCount["P0T"]) / float64(from0)
		fmt.Printf("P(0→1): %.6f (expected: %.4f)\n", p01_sim, P01)
		fmt.Printf("P(0→T): %.6f (expected: %.4f)\n", p0T_sim, P0T)
	}
}

func runSimulations(
	counter int,
	ET0, ET1, P10, P01, P1T, P0T, lambda0, lambda1, mu float64,
	transitionCount TransitionCount,
) map[int]int {
	freq := make(map[int]int)
	wifi_visit_case := make(map[keyProbability]int)

	for i := 0; i < counter; i++ {
		initial := initialstate(ET0, ET1)
		countOfOnes, path := wifi_visit(initial, P10, P01, P1T, P0T, transitionCount)
		freq[countOfOnes]++
		switch path[0] {
		case 1:
			if path[len(path)-1] == 1 {
				wifi_visit_case[keyProbability{countName: "case1", nCount: countOfOnes}]++
			} else {
				wifi_visit_case[keyProbability{countName: "case2", nCount: countOfOnes}]++
			}
		case 0:
			if path[len(path)-1] == 1 {
				wifi_visit_case[keyProbability{countName: "case3", nCount: countOfOnes}]++
			} else {
				wifi_visit_case[keyProbability{countName: "case4", nCount: countOfOnes}]++
			}
		}
	}

	maxN := 0
	for n := range freq {
		if n > maxN {
			maxN = n
		}
	}

	cases := []struct {
		name string
		fn   func(float64, float64, float64, int) float64
	}{
		{"case1", case1},
		{"case2", case2},
		{"case3", case3},
		{"case4", case4},
	}

	analysisSum := make(map[int]float64)
	simulationSum := make(map[int]float64)

	for _, c := range cases {
		fmt.Printf("\n%s (Analysis vs Simulation): n1 from 0 to %d\n", c.name, maxN)
		fmt.Printf("%-4s %-15s %-15s %-12s\n", "n1", "Analysis P", "Sim P", "%error")
		fmt.Println(strings.Repeat("-", 55))

		for n := 0; n <= maxN; n++ {
			analysis := c.fn(lambda0, lambda1, mu, n)
			key := keyProbability{countName: c.name, nCount: n}
			simulation := float64(wifi_visit_case[key]) / float64(counter)

			analysisSum[n] += analysis
			simulationSum[n] += simulation

			var percentErrorStr string
			if simulation > 1e-10 {
				percentError := math.Abs(simulation-analysis) / simulation * 100
				percentErrorStr = fmt.Sprintf("%.4f%%", percentError)
			} else if analysis > 1e-10 {
				percentErrorStr = ">100% (sim=0)"
			} else {
				percentErrorStr = "0.00%"
			}

			fmt.Printf("%-4d %-15.8f %-15.8f %s\n", n, analysis, simulation, percentErrorStr)
		}
	}

	fmt.Println("\n" + strings.Repeat("=", 65))
	fmt.Println("SUM OF ALL CASES FOR EACH n1")
	fmt.Printf("%-4s %-15s %-15s %-12s\n", "n1", "Analysis Sum", "Sim Sum", "%error")
	fmt.Println(strings.Repeat("-", 55))
	for n := 0; n <= maxN; n++ {
		analysisTotal := analysisSum[n]
		simulationTotal := simulationSum[n]

		var percentErrorStr string
		if simulationTotal > 1e-10 {
			percentError := math.Abs(simulationTotal-analysisTotal) / simulationTotal * 100
			percentErrorStr = fmt.Sprintf("%.4f%%", percentError)
		} else if analysisTotal > 1e-10 {
			percentErrorStr = ">100% (sim=0)"
		} else {
			percentErrorStr = "0.00%"
		}

		fmt.Printf("%-4d %-15.8f %-15.8f %s\n", n, analysisTotal, simulationTotal, percentErrorStr)
	}
	fmt.Println(strings.Repeat("=", 65))

	return freq
}

func displayProbabilities(n1 map[int]int, total int) {
	maxN := 0
	for n := range n1 {
		if n > maxN {
			maxN = n
		}
	}

	fmt.Printf("\nResults from %d simulations:\n", total)
	fmt.Printf("Frequency count: %v\n", n1)
	fmt.Printf("\nProbabilities:\n")
	fmt.Printf("%-8s %-10s %-15s\n", "n1", "Frequency", "P(n1 = x)")
	fmt.Println(strings.Repeat("-", 35))

	totalProbability := 0.0
	for n := 0; n <= maxN; n++ {
		freq := n1[n]
		probability := float64(freq) / float64(total)
		totalProbability += probability
		fmt.Printf("%-8d %-10d %-15.6f\n", n, freq, probability)
	}

	fmt.Printf("%-8s %-10s %-15.6f\n", "Total", "", totalProbability)
}

func initialstate(ET0, ET1 float64) int {
	U := rand.Float64()
	P0 := ET0 / (ET0 + ET1)
	if U <= P0 {
		return 0
	}
	return 1
}

func wifi_visit(
	initialstate int,
	P10, P01, P1T, P0T float64,
	transitionCount TransitionCount,
) (int, []int) {
	state := initialstate
	countOfOnes := 0
	var wifiVisit []int
	_ = P1T
	_ = P0T

	wifiVisit = append(wifiVisit, state)
	if state == 1 {
		countOfOnes++
	}

	for {
		if state == 0 {
			U := rand.Float64()
			if U <= P01 {
				state = 1
				wifiVisit = append(wifiVisit, state)
				transitionCount["P01"]++
				countOfOnes++
			} else {
				transitionCount["P0T"]++
				break
			}
		} else if state == 1 {
			U := rand.Float64()
			if U <= P10 {
				state = 0
				wifiVisit = append(wifiVisit, state)
				transitionCount["P10"]++
			} else {
				transitionCount["P1T"]++
				break
			}
		}
	}

	return countOfOnes, wifiVisit
}

func case1(lambda0, lambda1, mu float64, n int) float64 {
	if n == 0 {
		return 0.0
	}
	return (lambda0 / (lambda1 + lambda0)) * math.Pow((math.Pow(mu, 2)/((lambda0+mu)*(lambda1+mu))), float64(n-1)) * (lambda0 / (lambda0 + mu))
}

func case2(lambda0, lambda1, mu float64, n int) float64 {
	if n == 0 {
		return 0.0
	}
	return (lambda0 / (lambda1 + lambda0)) * math.Pow((mu/(lambda0+mu)), float64(n)) * math.Pow((mu/(lambda1+mu)), float64(n-1)) * (lambda1 / (lambda1 + mu))
}

func case3(lambda0, lambda1, mu float64, n int) float64 {
	if n == 0 {
		return 0.0
	}
	return (lambda1 / (lambda1 + lambda0)) * math.Pow((mu/(lambda1+mu)), float64(n)) * math.Pow((mu/(lambda0+mu)), float64(n-1)) * (lambda0 / (lambda0 + mu))
}

func case4(lambda0, lambda1, mu float64, n int) float64 {
	return (lambda1 / (lambda1 + lambda0)) * math.Pow((math.Pow(mu, 2)/((lambda1+mu)*(lambda0+mu))), float64(n)) * (lambda1 / (lambda1 + mu))
}
