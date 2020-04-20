package main

import (
	"fmt"
	"tic_tac_toe/game"
	"tic_tac_toe/pool"
)

func main() {
	gg := game.NewGame()
	pl := pool.NewPool()
	players := [...]game.PlayerType{game.PlayerX, game.PlayerO}

	if err := pl.ReadData("brain.bin"); err != nil {
		panic(err)
	}

	index := 0

	for {
		ret, err := pl.GetStep(*gg, players[index])
		if err != nil {
			panic(err)
		}

		err = gg.MakeStep(players[index], ret) //nolint
		if err != nil {
			panic(err)
		}

		gg.Draw()

		if res, val := gg.CheckWin(); res {
			switch val {
			case game.PlayerNone:
				fmt.Println("Draw")
				pl.DoDraw()
			case game.PlayerX:
				fmt.Println("Win: X")
				pl.DoWin(game.PlayerX)
			case game.PlayerO:
				fmt.Println("Win: O")
				pl.DoWin(game.PlayerO)
			}
			break
		}

		switch index {
		case 0:
			index = 1
		case 1:
			index = 0
		}
	}

	if err := pl.WriteData("brain.bin"); err != nil {
		panic(err)
	}

	fmt.Println("Games:", pl.GameCounts)
	fmt.Println("Variants:", len(pl.Pool))
}
