package main 

import (
	"fmt"
	"log"
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

func (cg *chessGame) newGame() {
	majorMinor := []string{"rk","kt","bp","qn","kg","bp","kt","rk"}
	for i := 0; i < 8; i++ {

		if constructor, found := pieces[majorMinor[i]]; !found {
			err := fmt.Errorf("piece type '%s' not found in map", majorMinor[i])
			log.Fatal("Error:", err)
		} else {
			cg.board[0][i] = square{0, i, constructor(majorMinor[i] + "-B-0", 0, i)}
			cg.board[7][i] = square{7, i, constructor(majorMinor[i] + "-W-0", 7, i)}
		}

		if constructor, found := pieces["pn"]; !found {
			err := fmt.Errorf("piece type 'pn' not found in map")
			log.Fatal("Error:", err)
		} else {
			cg.board[1][i] = square{1, i, constructor("pn-B-0", 1, i)}
			cg.board[6][i] = square{6, i, constructor("pn-W-0", 6, i)}
		}

		cg.board[2][i] = square{2, i, nil}
		cg.board[3][i] = square{2, i, nil}
		cg.board[4][i] = square{2, i, nil}
		cg.board[5][i] = square{2, i, nil}
	}
}

func (cg *chessGame) loadGame(mGame [8][8]string, mPlayer string, mWCaptured []string, mBCaptured []string) {
	cg.player = mPlayer
	for row := range mGame {
		for col := range mGame[row] {
			piece := strings.TrimSpace(strings.Split(mGame[row][col], "-")[0])

			if piece == "e" {
				cg.board[row][col] = square{row, col, nil}
				continue
			}
			
			constructor, found := pieces[piece] 
			if !found {
				err := fmt.Errorf("Piece not found in pieces map %s", mGame[row][col])
				log.Fatal("Error:", err)
			}
			
			cg.board[row][col] = square{row, col, constructor(mGame[row][col], row, col)}
		}
	}

	// load captured pieces
	for _, p := range mWCaptured {
		piece := strings.TrimSpace(strings.Split(p, "-")[0])
		constructor, _ := getPieceConstructor(piece)
		cg.wCaptured = append(cg.wCaptured, constructor(p, -1, -1))
	}

	for _, p := range mBCaptured {
		piece := strings.TrimSpace(strings.Split(p, "-")[0])
		constructor, _ := getPieceConstructor(piece)
		cg.bCaptured = append(cg.bCaptured, constructor(p, -1, -1))
	}
}

var visitedSquare = make(map[*square]bool)

func (cg *chessGame) inCheck(moveTo *square) (bool, error, chessPiece) { 
	if visitedSquare[moveTo] {
		return false, nil, nil
	}
	visitedSquare[moveTo] = true

	fmt.Println("evaluating this square in inCheck")
	fmt.Println(moveTo.row, moveTo.col)

	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			moveFrom, err := cg.getSquare(row, col)
			fmt.Println("This is the square that will checking")
			fmt.Println(row, col)
			if err != nil {
				fmt.Errorf("There was an error")
			}
			if moveFrom.cp != nil && (moveFrom.cp.color() != cg.player) {
				fmt.Println("checking isValidMove for this piece %d", moveFrom.cp.fullName())
				fmt.Println(moveTo)
				canMove, _, _, _, err := moveFrom.cp.isValidMove(moveTo, cg)
				if canMove {
					err = fmt.Errorf("King will be in check!")
					return canMove, err, moveFrom.cp
				}
			}
		}
	}
	return false, nil, nil
}

func (cg *chessGame) getSquare(row, col int) (*square, error) {
	if row > 7 || row < 0 || col > 7 || col < 0 {
		return &square{}, fmt.Errorf("Square is out of range: %d %d", row, col)
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