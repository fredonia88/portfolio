package main 

import (
	"fmt"
	"strings"
)


var moveTo *square
var moveFrom *square
var moveFromPrior *square
var moveToPrior *square

type chessGame struct {
	board [8][8]square
	player string
	wCaptured []chessPiece
	bCaptured []chessPiece
}

type square struct {
	row int // rank
	col int // file
	cp chessPiece
}

func (cg *chessGame) newGame() (err error) {
	majorMinor := []string{"rk","kt","bp","qn","kg","bp","kt","rk"}
	for i := 0; i < 8; i++ {

		majorMinorConstructor, majorMinorErr := getPieceConstructor(majorMinor[i])
		if majorMinorErr != nil {
			err = majorMinorErr
			return
		}
		cg.board[0][i] = square{0, i, majorMinorConstructor(majorMinor[i] + "-B-0", 0, i)}
		cg.board[7][i] = square{7, i, majorMinorConstructor(majorMinor[i] + "-W-0", 7, i)}

		pawnConstructor, pawnErr := getPieceConstructor("pn")
		if pawnErr != nil {
			err = pawnErr
			return
		}
		cg.board[1][i] = square{1, i, pawnConstructor("pn-B-0", 1, i)}
		cg.board[6][i] = square{6, i, pawnConstructor("pn-W-0", 6, i)}

		cg.board[2][i] = square{2, i, nil}
		cg.board[3][i] = square{2, i, nil}
		cg.board[4][i] = square{2, i, nil}
		cg.board[5][i] = square{2, i, nil}
	}

	return
}

func (cg *chessGame) loadGame(mGame [8][8]string, mPlayer string, mWCaptured []string, mBCaptured []string) (err error) {
	cg.player = mPlayer
	for row := range mGame {
		for col := range mGame[row] {
			piece := strings.TrimSpace(strings.Split(mGame[row][col], "-")[0])

			if piece == "e" {
				cg.board[row][col] = square{row, col, nil}
				continue
			}
			
			constructor, conErr := getPieceConstructor(piece)
			if conErr != nil {
				err = conErr
				return
			}
			
			cg.board[row][col] = square{row, col, constructor(mGame[row][col], row, col)}
		}
	}

	// load captured pieces
	for _, p := range mWCaptured {
		piece := strings.TrimSpace(strings.Split(p, "-")[0])
		constructor, conErr := getPieceConstructor(piece)
		if conErr != nil {
			err = conErr
			return
		}
		cg.wCaptured = append(cg.wCaptured, constructor(p, -1, -1))
	}

	for _, p := range mBCaptured {
		piece := strings.TrimSpace(strings.Split(p, "-")[0])
		constructor, conErr := getPieceConstructor(piece)
		if conErr != nil {
			err = conErr
			return
		}
		cg.bCaptured = append(cg.bCaptured, constructor(p, -1, -1))
	}
	
	return 
}

var visitedSquare = make(map[*square]bool)

func (cg *chessGame) inCheck(moveTo *square) (error, chessPiece) { 
	if visitedSquare[moveTo] {
		return nil, nil
	}
	visitedSquare[moveTo] = true

	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			moveFrom, err := cg.getSquare(row, col)
			if err != nil {
				fmt.Errorf("There was an error")
			}
			if moveFrom.cp != nil && (moveFrom.cp.color() != cg.player) {
				_, _, _, err := moveFrom.cp.isValidMove(moveTo, cg)
				if err == nil {
					err = fmt.Errorf("King will be in check!")
					return err, moveFrom.cp
				}
			}
		}
	}
	return nil, nil
}

func (cg *chessGame) getSquare(row, col int) (*square, error) {
	if row > 7 || row < 0 || col > 7 || col < 0 {
		return &square{}, newChessError(errOutOfRange, "Selected square (%d %d) is out of range", row, col)
	}
	return &cg.board[row][col], nil
}

func (cg *chessGame) displayBoard() {
	for row := 0; row < 8; row++ {
		line := make([]string, 8)
		for col := 0; col < 8; col++ {
			if cg.board[row][col].cp == nil {
				line[col] = "|  e-  |"
			} else {
				line[col] = "|" + cg.board[row][col].cp.fullName() + "|"
			}
		}
		fmt.Println(line)
	}
}

// function to check if piece positions align with board
func (cg *chessGame) displayBoardPositions() {
	for row := 0; row < 8; row++ {
		line := make([][2]int, 8)
		for col := 0; col < 8; col++ {
			line[col] = [2]int{cg.board[row][col].row, cg.board[row][col].col}
		}
		fmt.Println(line)
	}
}

func (cg *chessGame) displayCapturedPieces() {
	var w []string
	var b []string
	for _, piece := range cg.wCaptured {
		w = append(w, piece.fullName())
	}
	for _, piece := range cg.bCaptured {
		b = append(b, piece.fullName())
	}
	fmt.Println("White pieces captured:", w)
	fmt.Println("Black pieces captured:", b)
}