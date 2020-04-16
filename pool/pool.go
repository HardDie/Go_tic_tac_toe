package pool

import (
	"errors"
	"math/rand"
	"tic_tac_toe/game"
	"time"
)

/**
 * Constants
 */

const firstWeight = 3

/**
 * Types
 */

type variants [game.FieldWidth * game.FieldHeight]uint64

type Pool struct {
	pool map[string]*variants
}

/**
 * Methods
 */

func NewPool() *Pool {
	pool := Pool{}
	pool.pool = make(map[string]*variants)
	rand.Seed(time.Now().UnixNano())
	return &pool
}

func (p *Pool) GetStep(gm game.Game) (int, error) {
	field := gm.Field2String()

	// Init element if not exist
	if _, ok := p.pool[field]; !ok {
		p.pool[field] = &variants{}
		for i, val := range field {
			if val == ' ' {
				p.pool[field][i] = firstWeight
			}
		}
	}

	// Calculate max random value
	max_value := uint64(0)
	for _, val := range p.pool[field] {
		max_value += val
	}

	// Choose step
	rand_value := rand.Uint64() % max_value
	tmp_value := uint64(0)
	for i, val := range p.pool[field] {
		if (rand_value >= tmp_value) &&
			(rand_value < tmp_value+val) {
			return i, nil
		}
		tmp_value += val
	}

	return 0, errors.New("Can't find step")
}
