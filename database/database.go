package database

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

type Database struct {
	CountPlayedGames uint64
	Pool             map[string]*[]uint64

	// Global variables
	width  uint16
	height uint16
	// Exporting values
	Width  uint16
	Height uint16
}

func New(width, height uint16) *Database {
	db := Database{}
	db.Pool = make(map[string]*[]uint64)
	db.Width, db.width = width, width
	db.Height, db.height = height, height
	rand.Seed(time.Now().UnixNano())
	return &db
}

func (db Database) Serialize() ([]byte, error) {
	dat, err := json.Marshal(db)
	if err != nil {
		return nil, err
	}

	return dat, nil
}

func (db Database) WriteData(filename string) error {
	data, err := db.Serialize()
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, data, 0644) //nolint:gosec
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) Deserialize(data []byte) error {
	err := json.Unmarshal(data, db)
	if err != nil {
		return err
	}

	if db.width != db.Width ||
		db.height != db.Height {
		return errors.New("Field size diffirent")
	}

	return nil
}

func (db *Database) ReadData(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			// Skip reading if file not exist
			return nil
		}
		return err
	}

	err = db.Deserialize(data)
	if err != nil {
		return err
	}

	return nil
}

func (db Database) GetStep(field string) (int, error) {
	if _, ok := db.Pool[field]; !ok {
		return -1, errors.New("Value not exist")
	}

	// Calculate max random value
	maxValue := uint64(0)
	for _, val := range *db.Pool[field] {
		maxValue += val
	}

	if maxValue == 0 {
		return -1, errors.New("Pool is empty")
	}

	// Choose step
	randValue := rand.Uint64() % maxValue
	tmpValue := uint64(0)
	for i, val := range *db.Pool[field] {
		if (randValue >= tmpValue) &&
			(randValue < tmpValue+val) {
			return i, nil
		}
		tmpValue += val
	}

	return -1, errors.New("Can't find step")
}

func (db Database) CreateField(field string, initValue uint64) error {
	if _, ok := db.Pool[field]; ok {
		return errors.New("Element already exist")
	}

	db.Pool[field] = &[]uint64{}

	for _, val := range field {
		if val == ' ' {
			*db.Pool[field] = append(*db.Pool[field], initValue)
		} else {
			*db.Pool[field] = append(*db.Pool[field], 0)
		}
	}

	return nil
}

func (db *Database) ChangeValue(field string, index, diff int) error {
	if _, ok := db.Pool[field]; !ok {
		return errors.New("Element not exist")
	}

	if index < 0 || index >= len(*db.Pool[field]) {
		return errors.New("Wrong index")
	}

	if (*db.Pool[field])[index] == 0 {
		return errors.New("Can't change frozen value")
	}

	if diff < 0 {
		if uint64(diff*-1) > (*db.Pool[field])[index] {
			(*db.Pool[field])[index] = 1
		} else {
			(*db.Pool[field])[index] -= uint64(diff * -1)
		}
	} else if diff > 0 {
		(*db.Pool[field])[index] += uint64(diff)
	}

	return nil
}

func (db *Database) IncreaseCountGames() {
	db.CountPlayedGames++
}
