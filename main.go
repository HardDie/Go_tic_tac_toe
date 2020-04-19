package main

import (
	"fmt"
	"tic_tac_toe/game"
	"tic_tac_toe/pool"
)

func main() {
	gg := game.NewGame()
	pl := [...]*pool.Pool{pool.NewPool(), pool.NewPool()}
	players := [...]game.PlayerType{game.PlayerX, game.PlayerO}

	if err := pl[0].ReadData("brain_X.bin"); err != nil {
		panic(err)
	}
	if err := pl[1].ReadData("brain_O.bin"); err != nil {
		panic(err)
	}

	index := 0

	for {
		ret, err := pl[index].GetStep(*gg)
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
				pl[0].DoLose()
				pl[1].DoLose()
			case game.PlayerX:
				fmt.Println("Win: X")
				pl[0].DoWin()
				pl[1].DoLose()
			case game.PlayerO:
				fmt.Println("Win: O")
				pl[0].DoLose()
				pl[1].DoWin()
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

	if err := pl[0].WriteData("brain_X.bin"); err != nil {
		panic(err)
	}
	if err := pl[1].WriteData("brain_O.bin"); err != nil {
		panic(err)
	}

	fmt.Println("Games:", pl[index].GameCounts)
	fmt.Println("Variants:", len(pl[index].Pool))
}
