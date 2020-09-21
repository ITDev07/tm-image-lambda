package main

import (
	"fmt"
	"math"
	"strconv"
)

type scorePair struct {
	PA int
	PB int
}

type scoreResult struct {
	total_per float64
	adj_per float64
	total_score float64
	rounded_total_score float64
}

func decodeScore(score string) [4]scorePair {
	var score_pair [4]scorePair

    for idx, ch := range score {
        map_value, _:= strconv.Atoi(key_map[string(ch)])
		score_pair[idx].PA = map_value / 10 // Get first digit from left
		score_pair[idx].PB = map_value % 10 // Get first digit from right
	}

	return score_pair
}

func compareScores(client_score_pairs, match_score_pairs [4]scorePair) *scoreResult {
	same_variance := 0.0
	total_per := 0.0
	adj_per := 0.0
	total_score := 0.0
	total_variance := 0.0


    for i := 0; i < 4; i++ {
		client_score := client_score_pairs[i]
		match_score := match_score_pairs[i]

        varianceA := math.Abs(float64(client_score.PA-match_score.PA))
		varianceB := math.Abs(float64(client_score.PA-match_score.PB))

		total_variance += varianceA
		total_variance += varianceB

		if varianceA == 0 {
			same_variance += 1
		}

		if varianceB == 0 {
			same_variance += 1
		}
	}

	total_per = total_variance/32.0
	adj_per = same_variance/8.0
	total_score = (1-total_per) + (total_per*adj_per)

	total_result := scoreResult{}
	total_result.total_per = total_per
	total_result.adj_per = adj_per
	total_result.total_score = total_score
	total_result.rounded_total_score = math.Round(total_score*100)

	fmt.Println(total_result)

	return &total_result
}