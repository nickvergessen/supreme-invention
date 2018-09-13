package main

import (
	"log"
	"math/rand"

	"github.com/goml/gobrain"
)

var (
	emptyGame = []float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5,
		0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5,
		0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5,
		0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5,
		0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5,
		0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5}
	winOptions = [42][][]int{}
)

func addWinPattern(pattern []int) {
	for i := 0; i < 4; i++ {
		winOptions[pattern[i]] = append(winOptions[pattern[i]], pattern)
	}
}

func main() {
	// set the random seed to 0
	rand.Seed(0)

	// Winning conditions
	patterns1 := [][][]float64{}

	// Horizontal
	for i := 0; i < 7; i++ {
		for j := 0; j < 3; j++ {
			addWinPattern([]int{
				i + 7*(j+0),
				i + 7*(j+1),
				i + 7*(j+2),
				i + 7*(j+3)})

			tempGame := make([]float64, len(emptyGame))
			copy(tempGame, emptyGame)

			tempGame[i+7*(j+0)] = 1
			tempGame[i+7*(j+1)] = 1
			tempGame[i+7*(j+2)] = 1
			tempGame[i+7*(j+3)] = 1
			patterns1 = append(patterns1, [][]float64{tempGame, {1}})

			tempGame[i+7*(j+0)] = 0
			tempGame[i+7*(j+1)] = 0
			tempGame[i+7*(j+2)] = 0
			tempGame[i+7*(j+3)] = 0
			patterns1 = append(patterns1, [][]float64{tempGame, {0}})
		}
	}

	// Vertical
	for i := 0; i < 4; i++ {
		for j := 0; j < 6; j++ {
			addWinPattern([]int{
				i + 0 + 7*j,
				i + 1 + 7*j,
				i + 2 + 7*j,
				i + 3 + 7*j})

			tempGame := make([]float64, len(emptyGame))
			copy(tempGame, emptyGame)

			tempGame[i+0+7*j] = 1
			tempGame[i+1+7*j] = 1
			tempGame[i+2+7*j] = 1
			tempGame[i+3+7*j] = 1
			patterns1 = append(patterns1, [][]float64{tempGame, {1}})

			tempGame[i+0+7*j] = 0
			tempGame[i+1+7*j] = 0
			tempGame[i+2+7*j] = 0
			tempGame[i+3+7*j] = 0
			patterns1 = append(patterns1, [][]float64{tempGame, {0}})
		}
	}

	// top left to bottom right \
	for i := 0; i < 4; i++ {
		for j := 0; j < 3; j++ {
			addWinPattern([]int{
				i + 0 + 7*(j+0),
				i + 1 + 7*(j+1),
				i + 2 + 7*(j+2),
				i + 3 + 7*(j+3)})

			tempGame := make([]float64, len(emptyGame))
			copy(tempGame, emptyGame)

			tempGame[i+0+7*(j+0)] = 1
			tempGame[i+1+7*(j+1)] = 1
			tempGame[i+2+7*(j+2)] = 1
			tempGame[i+3+7*(j+3)] = 1
			patterns1 = append(patterns1, [][]float64{tempGame, {1}})

			tempGame[i+0+7*(j+0)] = 0
			tempGame[i+1+7*(j+1)] = 0
			tempGame[i+2+7*(j+2)] = 0
			tempGame[i+3+7*(j+3)] = 0
			patterns1 = append(patterns1, [][]float64{tempGame, {0}})
		}
	}

	// top right to bottom left /
	for i := 0; i < 4; i++ {
		for j := 0; j < 3; j++ {
			addWinPattern([]int{
				i + 3 + 7*(j+0),
				i + 2 + 7*(j+1),
				i + 1 + 7*(j+2),
				i + 0 + 7*(j+3)})

			tempGame := make([]float64, len(emptyGame))
			copy(tempGame, emptyGame)

			tempGame[i+3+7*(j+0)] = 1
			tempGame[i+2+7*(j+1)] = 1
			tempGame[i+1+7*(j+2)] = 1
			tempGame[i+0+7*(j+3)] = 1
			patterns1 = append(patterns1, [][]float64{tempGame, {1}})

			tempGame[i+3+7*(j+0)] = 0
			tempGame[i+2+7*(j+1)] = 0
			tempGame[i+1+7*(j+2)] = 0
			tempGame[i+0+7*(j+3)] = 0
			patterns1 = append(patterns1, [][]float64{tempGame, {0}})
		}
	}

	// instantiate the Feed Forward
	ff1 := &gobrain.FeedForward{}
	ff2 := &gobrain.FeedForward{}

	// initialize the Neural Network;
	// the networks structure will contain:
	// 42 inputs, 2 hidden nodes and 1 output.
	ff1.Init(42, 42, 1)
	ff2.Init(42, 42, 1)

	// train the network using the XOR patterns
	// the training will run for 1000 epochs
	// the learning rate is set to 0.6 and the momentum factor to 0.4
	// use true in the last parameter to receive reports about the learning error
	ff1.Train(patterns1, 1000, 0.6, 0.4, true)
	ff2.Train(patterns1, 1000, 0.6, 0.4, true)

	ff1.SetContexts(1, nil)
	ff2.SetContexts(1, nil)

	var stats [3]int
	for i := 0; i < 1; i++ {
		winner := playAGame(ff1, ff2)
		if winner == 1 {
			stats[1]++
		} else if winner == 0 {
			stats[2]++
		} else {
			stats[0]++
		}

		log.Printf("Draw: %d, Player X: %d, Player O: %d", stats[0], stats[1], stats[2])
	}

	log.Printf("Draw: %d, Player X: %d, Player O: %d", stats[0], stats[1], stats[2])
}

func playAGame(ff1, ff2 *gobrain.FeedForward) float64 {
	game := make([]float64, len(emptyGame))
	copy(game, emptyGame)

	// TODO debugging the NN
	game[41] = 0
	game[34] = 0
	game[27] = 0
	game[40] = 1
	game[33] = 1
	game[26] = 1

	for i := 1; i <= len(emptyGame); i++ {
		player := i % 2
		var field int
		if player == 1 {
			field = ask(ff1, player, game)
		} else {
			field = ask(ff2, player, game)
		}

		// Make the move
		game[field] = float64(player)
		printField(game)

		if checkWinner(player, game, field) {
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

		if i == len(emptyGame) {
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

func checkWinner(player int, game []float64, move int) bool {
	options := winOptions[move]
OPTIONS:
	for _, fields := range options {
		for _, field := range fields {
			if move != field && game[field] != float64(player) {
				continue OPTIONS
			}
		}
		return true
	}
	return false
}

func isValidMove(game []float64, move int) bool {
	return move >= 0 && move < len(game) && (move+7 >= len(game) || game[move+7] != 0.5) && game[move] == 0.5
}

func ask(ff *gobrain.FeedForward, player int, game []float64) int {

	field := 0
	chance := 0.0
	if player == 0 {
		chance = 1.0
	}

	if player == 1 {
		printField(game)
	}

	for i := 0; i < len(game); i++ {
		if !isValidMove(game, i) {
			continue
		}

		tempGame := make([]float64, len(game))
		copy(tempGame, game)
		tempGame[i] = float64(player)

		// TODO this should be learned, not hardcoded
		// if checkWinner(player, tempGame, i) {
		// 	return i
		// }
		// if checkWinner(player+1%2, tempGame, i) {
		// 	return i
		// }

		result := ff.Update(tempGame)
		if player == 1 {
			log.Printf("Player %d field %d rate %f", player, i, result[0])
			if result[0] > chance {
				chance = result[0]
				field = i
			}
		} else if player == 0 {
			// log.Printf("Player %d field %d rate %f", player, i, result[0])
			if result[0] < chance {
				chance = result[0]
				field = i
			}
		}
	}

	// log.Printf("Player %d picking %.0f", player, field)
	return field
}

func printField(game []float64) {
	log.Printf("\n\n%s|%s|%s|%s|%s|%s|%s\n-+-+-+-+-+-+-\n%s|%s|%s|%s|%s|%s|%s\n-+-+-+-+-+-+-\n%s|%s|%s|%s|%s|%s|%s\n-+-+-+-+-+-+-\n%s|%s|%s|%s|%s|%s|%s\n-+-+-+-+-+-+-\n%s|%s|%s|%s|%s|%s|%s\n-+-+-+-+-+-+-\n%s|%s|%s|%s|%s|%s|%s\n\n",
		printPlayer(game[0]), printPlayer(game[1]), printPlayer(game[2]), printPlayer(game[3]), printPlayer(game[4]), printPlayer(game[5]), printPlayer(game[6]),
		printPlayer(game[7]), printPlayer(game[8]), printPlayer(game[9]), printPlayer(game[10]), printPlayer(game[11]), printPlayer(game[12]), printPlayer(game[13]),
		printPlayer(game[14]), printPlayer(game[15]), printPlayer(game[16]), printPlayer(game[17]), printPlayer(game[18]), printPlayer(game[19]), printPlayer(game[20]),
		printPlayer(game[21]), printPlayer(game[22]), printPlayer(game[23]), printPlayer(game[24]), printPlayer(game[25]), printPlayer(game[26]), printPlayer(game[27]),
		printPlayer(game[28]), printPlayer(game[29]), printPlayer(game[30]), printPlayer(game[31]), printPlayer(game[32]), printPlayer(game[33]), printPlayer(game[34]),
		printPlayer(game[35]), printPlayer(game[36]), printPlayer(game[37]), printPlayer(game[38]), printPlayer(game[39]), printPlayer(game[40]), printPlayer(game[41]),
	)
}

func printPlayer(player float64) string {
	if player == 0 {
		return "O"
	} else if player == 1 {
		return "X"
	}

	return " "
}
