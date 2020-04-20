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

	for _, poolEl := range pl {
		if err := poolEl.ReadData("brain.bin"); err != nil {
			panic(err)
		}
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
				pl[0].DoWin()
				pl[0].DoWinTwo(pl[1])
				index = 0
			case game.PlayerX:
				fmt.Println("Win: X")
				pl[0].DoWin()
				pl[0].DoLose(pl[1])
				index = 0
			case game.PlayerO:
				fmt.Println("Win: O")
				pl[1].DoWin()
				pl[1].DoLose(pl[0])
				index = 1
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

	if err := pl[index].WriteData("brain.bin"); err != nil {
		panic(err)
	}

	fmt.Println("Games:", pl[index].GameCounts)
	fmt.Println("Variants:", len(pl[index].Pool))
}
