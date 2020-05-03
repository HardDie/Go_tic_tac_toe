package pool

import (
	"errors"
	"fmt"
	"tic_tac_toe/database"
	"tic_tac_toe/game"
)

/**
 * Constants
 */

const firstWeight = 3
const constStudyWinWeight = 3
const constStudyLoseWeight = 1

/**
 * Types
 */

type step struct {
	field string
	step  int
}

type Pool struct {
	activeSteps [2][]step
	db          *database.Database
}

/**
 * Methods
 */

func New(db *database.Database) *Pool {
	pool := Pool{}
	pool.db = db
	return &pool
}

func stepWithRotation(pos, rotation int) int {
	var row int
	var line int

	row = pos % game.FieldWidth
	line = pos / game.FieldWidth

	switch rotation {
	case 0:
		return pos
	case 90:
		return (game.FieldHeight - 1 - line) + row*game.FieldWidth
	case 180:
		return (game.FieldWidth - 1 - row) + (game.FieldHeight-1-line)*game.FieldWidth
	case 270:
		return line + (game.FieldWidth-1-row)*game.FieldWidth
	default:
		panic("Noooooo")
	}
}

func (p *Pool) GetStep(gm game.Game, player game.PlayerType) (int, error) {
	var curField string
	var curStep int = -1

	rotation := 0
	fields := gm.FieldToAllVariants()

	for i, tmpField := range fields {
		tmpStep, err := p.db.GetStep(tmpField)
		if err == nil {
			switch i {
			case 0, 1:
				rotation = 0
			case 2, 3:
				rotation = 90
			case 4, 5:
				rotation = 180
			case 6, 7:
				rotation = 270
			default:
				return 0, errors.New("Wrong rotation")
			}
			curField = tmpField
			curStep = tmpStep
			break
		}
	}

	// Create element if not exist
	if curStep == -1 {
		curField = fields[0]
		if err := p.db.CreateField(curField, firstWeight); err != nil {
			return 0, err
		}

		tmpStep, err := p.db.GetStep(curField)
		if err != nil {
			return 0, err
		}
		curStep = tmpStep
	}

	switch player {
	case game.PlayerX:
		p.activeSteps[0] = append(p.activeSteps[0], step{curField, curStep})
	case game.PlayerO:
		p.activeSteps[1] = append(p.activeSteps[1], step{curField, curStep})
	}

	return stepWithRotation(curStep, rotation), nil
}

func (p *Pool) DoWin(player game.PlayerType) {
	loserID := -1
	switch player {
	case game.PlayerX:
		for _, val := range p.activeSteps[0] {
			err := p.db.ChangeValue(val.field, val.step, constStudyWinWeight)
			if err != nil {
				panic(err)
			}
			fmt.Println("Good", p.db.Pool[val.field], val.step)
		}
		loserID = 1
	case game.PlayerO:
		for _, val := range p.activeSteps[1] {
			err := p.db.ChangeValue(val.field, val.step, constStudyWinWeight)
			if err != nil {
				panic(err)
			}
			fmt.Println("Good", p.db.Pool[val.field], val.step)
		}
		loserID = 0
	default:
		panic("ERROR")
	}

	for _, val := range p.activeSteps[loserID] {
		err := p.db.ChangeValue(val.field, val.step, -constStudyLoseWeight)
		if err != nil {
			panic(err)
		}
		fmt.Println("Bad", p.db.Pool[val.field], val.step)
	}

	p.db.IncreaseCountGames()
}

func (p *Pool) DoDraw() {
	for _, val := range p.activeSteps[0] {
		err := p.db.ChangeValue(val.field, val.step, constStudyWinWeight)
		if err != nil {
			panic(err)
		}
		fmt.Println("Good", p.db.Pool[val.field], val.step)
	}
	for _, val := range p.activeSteps[1] {
		err := p.db.ChangeValue(val.field, val.step, constStudyWinWeight)
		if err != nil {
			panic(err)
		}
		fmt.Println("Good", p.db.Pool[val.field], val.step)
	}

	p.db.IncreaseCountGames()
}
