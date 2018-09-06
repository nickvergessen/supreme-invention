package main

import (
	"log"
	"math/rand"

	"github.com/goml/gobrain"
)

var (
	winOptions = [][][]int{
		{{1, 2}, {4, 8}, {3, 6}},         // 0
		{{0, 2}, {4, 7}},                 // 1
		{{0, 1}, {4, 6}, {5, 8}},         // 2
		{{0, 6}, {4, 5}},                 // 3
		{{0, 8}, {1, 7}, {2, 6}, {3, 5}}, // 4
		{{2, 8}, {3, 4}},                 // 5
		{{0, 3}, {2, 4}, {7, 8}},         // 6
		{{6, 8}, {1, 4}},                 // 7
		{{6, 7}, {0, 4}, {2, 5}}}         // 8
)

func main() {
	// set the random seed to 0
	rand.Seed(0)

	// Winning conditions
	patterns1 := [][][]float64{
		{{0.5, 0.5, 1,
			0.5, 0.5, 1,
			0.5, 0.5, 1}, {1}},
		{{0.5, 1, 0.5,
			0.5, 1, 0.5,
			0.5, 1, 0.5}, {1}},
		{{1, 0.5, 0.5,
			1, 0.5, 0.5,
			1, 0.5, 0.5}, {1}},
		{{1, 1, 1,
			0.5, 0.5, 0.5,
			0.5, 0.5, 0.5}, {1}},
		{{0.5, 0.5, 0.5,
			1, 1, 1,
			0.5, 0.5, 0.5}, {1}},
		{{0.5, 0.5, 0.5,
			0.5, 0.5, 0.5,
			1, 1, 1}, {1}},
		{{1, 0.5, 0.5,
			0.5, 1, 0.5,
			0.5, 0.5, 1}, {1}},
		{{0.5, 0.5, 1,
			0.5, 1, 0.5,
			1, 0.5, 0.5}, {1}},
		{{0.5, 0.5, 0,
			0.5, 0.5, 0,
			0.5, 0.5, 0}, {0}},
		{{0.5, 0, 0.5,
			0.5, 0, 0.5,
			0.5, 0, 0.5}, {0}},
		{{0, 0.5, 0.5,
			0, 0.5, 0.5,
			0, 0.5, 0.5}, {0}},
		{{0, 0, 0,
			0.5, 0.5, 0.5,
			0.5, 0.5, 0.5}, {0}},
		{{0.5, 0.5, 0.5,
			0, 0, 0,
			0.5, 0.5, 0.5}, {0}},
		{{0.5, 0.5, 0.5,
			0.5, 0.5, 0.5,
			0, 0, 0}, {0}},
		{{0, 0.5, 0.5,
			0.5, 0, 0.5,
			0.5, 0.5, 0}, {0}},
		{{0.5, 0.5, 0,
			0.5, 0, 0.5,
			0, 0.5, 0.5}, {0}},
	}

	// traps
	// patterns2 := [][][]float64{
	// 	{{1, 0, 0,
	// 		0, 2, 0,
	// 		0, 0, 1}, {1}},
	// 	{{0, 0, 1,
	// 		0, 2, 0,
	// 		1, 0, 0}, {1}},
	// }

	// instantiate the Feed Forward
	ff1 := &gobrain.FeedForward{}
	ff2 := &gobrain.FeedForward{}

	// initialize the Neural Network;
	// the networks structure will contain:
	// 2 inputs, 2 hidden nodes and 1 output.
	ff1.Init(9, 2, 1)
	ff2.Init(9, 2, 1)

	// train the network using the XOR patterns
	// the training will run for 1000 epochs
	// the learning rate is set to 0.6 and the momentum factor to 0.4
	// use true in the last parameter to receive reports about the learning error
	ff1.Train(patterns1, 1000, 0.6, 0.4, true)
	// ff1.Train(patterns2, 1000, 0.6, 0.4, true)
	ff2.Train(patterns1, 1000, 0.6, 0.4, true)

	ff1.SetContexts(1, nil)
	ff2.SetContexts(1, nil)

	var stats [3]int
	for i := 0; i < 100; i++ {
		winner := playAGame(ff1, ff2)
		if winner == 1 {
			stats[1]++
		} else if winner == 0 {
			stats[2]++
		} else {
			stats[0]++
		}
	}

	log.Println(stats)
}

func playAGame(ff1, ff2 *gobrain.FeedForward) float64 {
	game := []float64{0.5, 0.5, 0.5,
		0.5, 0.5, 0.5,
		0.5, 0.5, 0.5}

	for i := 1; i < 10; i++ {
		player := i % 2
		var field float64
		if player == 1 {
			field = ask(ff1, player, game)
		} else {
			field = ask(ff2, player, game)
		}
		game[int(field)] = float64(player)
		if checkWinner(player, game, winOptions[int(field)]) {
			// log.Printf("\n\n#\n# Player %d wins\n#", player)

			playerOneWon := 1.0
			if player == 0 {
				playerOneWon = 0
			}

			ff1.Train([][][]float64{
				{game, {playerOneWon}},
			}, 200, 0.6, 0.4, false)
			ff2.Train([][][]float64{
				{game, {playerOneWon}},
			}, 200, 0.6, 0.4, false)

			printField(game)
			return playerOneWon
		}

		if i == 9 {
			// log.Printf("\n\n#\n# Draw\n#")
			ff1.Train([][][]float64{
				{game, {0}},
			}, 200, 0.6, 0.4, false)

			ff2.Train([][][]float64{
				{game, {0}},
			}, 200, 0.6, 0.4, false)
		}
	}

	printField(game)
	return 0.5
}

func checkWinner(player int, game []float64, options [][]int) bool {
OPTIONS:
	for _, fields := range options {
		for _, field := range fields {
			if game[field] != float64(player) {
				continue OPTIONS
			}
		}
		return true
	}
	return false
}

func ask(ff *gobrain.FeedForward, player int, game []float64) float64 {

	field := 0.0
	chance := 0.0
	if player == 0 {
		chance = 1.0
	}

	for i := 0; i < len(game); i++ {
		if game[i] != 0.5 {
			continue
		}

		// if player == 2 {
		// 	log.Printf("Player %d picking %d", player, i)
		// 	return float64(i)
		// }

		tempGame := make([]float64, len(game))
		copy(tempGame, game)
		tempGame[i] = float64(player)

		result := ff.Update(tempGame)
		if player == 1 {
			// log.Printf("Player %d field %d rate %f", player, i, result[0])
			if result[0] > chance {
				chance = result[0]
				field = float64(i)
			}
		} else if player == 0 {
			// log.Printf("Player %d field %d rate %f", player, i, result[0])
			if result[0] < chance {
				chance = result[0]
				field = float64(i)
			}
		}
	}

	// log.Printf("Player %d picking %.0f", player, field)
	return field
}

func printField(game []float64) {
	// log.Printf("\n\n%s|%s|%s\n-+-+-\n%s|%s|%s\n-+-+-\n%s|%s|%s\n\n",
	// 	printPlayer(game[0]), printPlayer(game[1]), printPlayer(game[2]),
	// 	printPlayer(game[3]), printPlayer(game[4]), printPlayer(game[5]),
	// 	printPlayer(game[6]), printPlayer(game[7]), printPlayer(game[8]))
}

func printPlayer(player float64) string {
	if player == 0 {
		return "O"
	} else if player == 1 {
		return "X"
	}

	return " "
}
