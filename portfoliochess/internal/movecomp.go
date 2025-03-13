package main

// function to assign point values to pieces and game
func (cg *chessGame) evalGame() (score int, err error) {
	score = 0
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			sq, sqErr := cg.getSquare(row, col)
			if err != nil {
				err = sqErr
				return
			}
			if sq.cp != nil {
				var m int
				if sq.cp.color() == "W" {
					m = 1
				} else {
					m = -1
				}
				switch sq.cp.name() {
				case "pn":
					score += 1 * m
				case "kt", "bp":
					score += 3 * m
				case "rk":
					score += 5 * m
				case "qn": 
					score += 9 * m
				case "kg":
					score += 100 * m
				}
			}
		}
	}
	return
}

// create new mini max method to determine computer moves

func (cg *chessGame) minimax(cg *chessGame, depth int, maxPlayer bool) (eval int, err error) {
	
}