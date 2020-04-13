package game

import "testing"

func TestReset(t *testing.T) {
	game := NewGame()
	for i := 0; i < FieldWidth*FieldHeight; i++ {
		game.MakeStep(Player_X, i)
	}

	game.Reset()

	for i := 0; i < FieldWidth*FieldHeight; i++ {
		if err := game.MakeStep(Player_X, i); err != nil {
			t.Errorf("Can't set value after reset: %v", err)
		}
	}
}

func TestMakeStep(t *testing.T) {
	game := NewGame()
	if err := game.MakeStep(Player_X, -1); err == nil {
		t.Errorf("Should be error")
	}

	if err := game.MakeStep(Player_X, FieldWidth*FieldHeight); err == nil {
		t.Errorf("Should be error")
	}

	for i := 0; i < FieldWidth*FieldHeight; i++ {
		if err := game.MakeStep(Player_X, i); err != nil {
			t.Errorf("Can't set value: %v", err)
		}
	}
}

func TestCheckWin(t *testing.T) {
	game := NewGame()

	players := [2]PlayerType{Player_X, Player_O}
	for i := 0; i < len(players); i++ {
		// Lines
		for line := 0; line < FieldHeight; line++ {
			game.Reset()
			for row := 0; row < FieldWidth; row++ {
				game.MakeStep(players[i], row+line*FieldWidth)
			}

			if ret, p_type := game.CheckWin(); ret != true {
				t.Errorf("The victory was not marked!")
			} else if p_type != players[i] {
				game.Draw()
				t.Errorf("Wrong victory player type!: %d", p_type)
			}
		}

		// Rows
		for row := 0; row < FieldWidth; row++ {
			game.Reset()
			for line := 0; line < FieldHeight; line++ {
				game.MakeStep(players[i], row+line*FieldWidth)
			}

			if ret, p_type := game.CheckWin(); ret != true {
				t.Errorf("The victory was not marked!")
			} else if p_type != players[i] {
				t.Errorf("Wrong victory player type!")
			}
		}

		// Diagonals
		if FieldWidth != FieldHeight {
			t.Errorf("Can't be diagonale win, field size: %dx%d", FieldWidth, FieldHeight)
		}
		// From left top to right bottom
		game.Reset()
		for j := 0; j < FieldWidth; j++ {
			game.MakeStep(players[i], j+j*FieldWidth)
		}
		if ret, p_type := game.CheckWin(); ret != true {
			t.Errorf("The victory was not marked!")
		} else if p_type != players[i] {
			t.Errorf("Wrong victory player type!")
		}
		// From right top to left bottom
		game.Reset()
		for j := 0; j < FieldWidth; j++ {
			game.MakeStep(players[i], (FieldWidth-1-j)+j*FieldWidth)
		}
		if ret, p_type := game.CheckWin(); ret != true {
			t.Errorf("The victory was not marked!")
		} else if p_type != players[i] {
			t.Errorf("Wrong victory player type!")
		}
	}
}
