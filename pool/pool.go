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
	activeSteps []step
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

func (p *Pool) GetStep(gm game.Game) (int, error) {
	field := gm.Field2String()

	// Init element if not exist
	if _, ok := p.Pool[field]; !ok {
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

			p.activeSteps = append(p.activeSteps, step{field, i})

			return i, nil
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

	err = ioutil.WriteFile(filename, dat, 0644)
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

func (p *Pool) DoWin() {
	for _, val := range p.activeSteps {
		p.Pool[val.field][val.step] += constStudyWinWeight
	}

	p.GameCounts++
}

func (p *Pool) DoLose() {
	for _, val := range p.activeSteps {
		if p.Pool[val.field][val.step] > constStudyLoseWeight {
			p.Pool[val.field][val.step] -= constStudyLoseWeight
		} else {
			p.Pool[val.field][val.step] = 1
		}
	}

	p.GameCounts++
}

func (p Pool) Print() {
	fmt.Println("GameCount:", p.GameCounts)
	fmt.Println("Width:", p.Width)
	fmt.Println("Height:", p.Height)
	for key, val := range p.Pool {
		fmt.Println(key, val)
	}
}
