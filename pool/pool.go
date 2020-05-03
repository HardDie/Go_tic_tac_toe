package pool

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"tic_tac_toe/game"
	"time"
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

type variants [game.FieldWidth * game.FieldHeight]uint64

type step struct {
	field string
	step  int
}

type Pool struct {
	GameCounts  uint64
	Width       uint16
	Height      uint16
	Pool        map[string]*variants
	activeSteps [2][]step
}

/**
 * Methods
 */

func NewPool() *Pool {
	pool := Pool{}
	pool.Pool = make(map[string]*variants)
	pool.Width = game.FieldWidth
	pool.Height = game.FieldHeight
	rand.Seed(time.Now().UnixNano())
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
	var field string
	rotation := 0
	fields := gm.FieldToAllVariants()

	for i, val := range fields {
		if _, ok := p.Pool[val]; ok {
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
			field = val
			break
		}
	}

	// Init element if not exist
	if len(field) == 0 {
		field = fields[0]
		p.Pool[field] = &variants{}
		for i, val := range field {
			if val == ' ' {
				p.Pool[field][i] = firstWeight
			}
		}
	}

	// Calculate max random value
	maxValue := uint64(0)
	for _, val := range p.Pool[field] {
		maxValue += val
	}

	// Choose step
	randValue := rand.Uint64() % maxValue
	tmpValue := uint64(0)
	for i, val := range p.Pool[field] {
		if (randValue >= tmpValue) &&
			(randValue < tmpValue+val) {

			switch player {
			case game.PlayerX:
				p.activeSteps[0] = append(p.activeSteps[0], step{field, i})
			case game.PlayerO:
				p.activeSteps[1] = append(p.activeSteps[1], step{field, i})
			}

			return stepWithRotation(i, rotation), nil
		}
		tmpValue += val
	}

	return 0, errors.New("Can't find step")
}

func (p Pool) WriteData(filename string) error {
	dat, err := json.Marshal(p)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, dat, 0644) // nolint: gosec
	if err != nil {
		return err
	}

	return nil
}

func (p *Pool) ReadData(filename string) error {
	dat, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// Skip reading if file not exist
			return nil
		}
		return err
	}

	err = json.Unmarshal(dat, p)
	if err != nil {
		return err
	}

	if game.FieldWidth != p.Width ||
		game.FieldHeight != p.Height {
		return errors.New("Field size diffirent")
	}

	return nil
}

func (p *Pool) DoWin(player game.PlayerType) {
	loserID := -1
	switch player {
	case game.PlayerX:
		for _, val := range p.activeSteps[0] {
			p.Pool[val.field][val.step] += constStudyWinWeight
			fmt.Println("Good", p.Pool[val.field], val.step)
		}
		loserID = 1
	case game.PlayerO:
		for _, val := range p.activeSteps[1] {
			p.Pool[val.field][val.step] += constStudyWinWeight
			fmt.Println("Good", p.Pool[val.field], val.step)
		}
		loserID = 0
	default:
		panic("ERROR")
	}

	for _, val := range p.activeSteps[loserID] {
		if p.Pool[val.field][val.step] > constStudyLoseWeight {
			p.Pool[val.field][val.step] -= constStudyLoseWeight
			fmt.Println("Bad", p.Pool[val.field], val.step)
		} else {
			p.Pool[val.field][val.step] = 1
			fmt.Println("Bad", p.Pool[val.field], val.step)
		}
	}

	p.GameCounts++
}

func (p *Pool) DoDraw() {
	for _, val := range p.activeSteps[0] {
		p.Pool[val.field][val.step] += constStudyWinWeight
		fmt.Println("Good", p.Pool[val.field], val.step)
	}
	for _, val := range p.activeSteps[1] {
		p.Pool[val.field][val.step] += constStudyWinWeight
		fmt.Println("Good", p.Pool[val.field], val.step)
	}

	p.GameCounts++
}
