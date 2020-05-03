package main

import (
	"fmt"
	"tic_tac_toe/database"
	"tic_tac_toe/game"
	"tic_tac_toe/pool"
)

func playGameAI(db *database.Database) error {
	gg := game.NewGame()
	pl := pool.New(db)
	players := [...]game.PlayerType{game.PlayerX, game.PlayerO}

	index := 0

	for {
		ret, err := pl.GetStep(*gg, players[index])
		if err != nil {
			return err
		}

		err = gg.MakeStep(players[index], ret) //nolint
		if err != nil {
			return err
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

	return nil
}

func main() {
	// Read database
	db := database.New(game.FieldWidth, game.FieldHeight)
	if err := db.ReadData("brain.bin"); err != nil {
		panic(err)
	}

	if err := playGameAI(db); err != nil {
		panic(err)
	}

	if err := db.WriteData("brain.bin"); err != nil {
		panic(err)
	}

	fmt.Println("Games:", db.CountPlayedGames)
	fmt.Println("Variants:", len(db.Pool))
}
