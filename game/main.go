package game

import (
	"errors"
	"fmt"
)

/**
 * Constants
 */

const FieldWidth = 3
const FieldHeight = 3

/**
 * Types
 */

type cellType rune

const (
	cell_Empty cellType = ' '
	cell_X              = 'X'
	cell_O              = 'O'
)

type PlayerType int

const (
	Player_None PlayerType = iota
	Player_X
	Player_O
)

type Game struct {
	cells [FieldWidth * FieldHeight]cellType
}

/**
 * Methods
 */

func NewGame() *Game {
	game := Game{}
	game.Reset()
	return &game
}

func (game *Game) Reset() {
	for i := 0; i < FieldWidth*FieldHeight; i++ {
		game.cells[i] = cell_Empty
	}
}

func (game Game) Draw() {
	for i := 0; i < FieldWidth*2+1; i++ {
		fmt.Printf("-")
	}
	fmt.Println()

	for line := 0; line < FieldHeight; line++ {
		for row := 0; row < FieldWidth; row++ {
			fmt.Printf("|%c", game.cells[line*FieldHeight+row])
		}
		fmt.Println("|")

		for i := 0; i < FieldWidth*2+1; i++ {
			fmt.Printf("-")
		}
		fmt.Println()
	}
}

func (game *Game) MakeStep(player PlayerType, pos int) error {
	if pos < 0 || pos >= (FieldWidth*FieldHeight) {
		return errors.New("Invalid position value")
	}

	if game.cells[pos] != cell_Empty {
		return errors.New("Field is already busy")
	}

	switch player {
	case Player_X:
		game.cells[pos] = cell_X
	case Player_O:
		game.cells[pos] = cell_O
	default:
		return errors.New("Wrong player type")
	}

	return nil
}

func getPlayerType(cell cellType) PlayerType {
	switch cell {
	case cell_X:
		return Player_X
	case cell_O:
		return Player_O
	}
	return Player_None
}

func (game Game) CheckWin() (bool, PlayerType) {
	var flag bool
	// Lines
	for line := 0; line < FieldHeight; line++ {
		if game.cells[0+line*FieldWidth] == cell_Empty {
			continue
		}

		flag = true
		for row := 0; row < FieldWidth-1; row++ {
			if game.cells[row+line*FieldWidth] != game.cells[(row+1)+line*FieldWidth] {
				flag = false
				break
			}
		}

		if flag == false {
			continue
		}

		return true, getPlayerType(game.cells[0+line*FieldWidth])
	}

	// Rows
	for row := 0; row < FieldWidth; row++ {
		if game.cells[row+0*FieldWidth] == cell_Empty {
			continue
		}

		flag = true
		for line := 0; line < FieldHeight-1; line++ {
			if game.cells[row+line*FieldWidth] != game.cells[row+(line+1)*FieldWidth] {
				flag = false
				break
			}
		}

		if flag == false {
			continue
		}

		return true, getPlayerType(game.cells[row+0*FieldWidth])
	}

	// Diagonals
	if FieldWidth == FieldHeight {
		// From top left to bottom right
		if game.cells[0] != cell_Empty {
			flag = true
			for i := 0; i < FieldWidth-1; i++ {
				if game.cells[i+i*FieldWidth] != game.cells[(i+1)+(i+1)*FieldWidth] {
					flag = false
					break
				}
			}
		}

		if flag == true {
			return true, getPlayerType(game.cells[0])
		}

		// From top right to bottom left
		if game.cells[FieldWidth-1] != cell_Empty {
			flag = true
			for i := 0; i < FieldWidth-1; i++ {
				if game.cells[(FieldWidth-1-i)+i*FieldWidth] != game.cells[(FieldWidth-2-i)+(i+1)*FieldWidth] {
					flag = false
					break
				}
			}
		}

		if flag == true {
			return true, getPlayerType(game.cells[FieldWidth-1])
		}
	}

	// Draw
	flag = true
	for _, val := range game.cells {
		if val == cell_Empty {
			flag = false
			break
		}
	}
	if flag == true {
		return true, Player_None
	}

	return false, Player_None
}
