package main

import (
	"fmt"
	"strings"
	"math"
	"math/rand"
	//"errors"
	"log"
	"strconv"
)
/*
Board Representation:

Decide how to represent the chessboard (e.g., 2D array, 1D array, or FEN strings).
Include a mapping for pieces (e.g., P for white pawn, p for black pawn, etc.).

Game Rules:

Implement the rules for legal moves (start with pawns, then other pieces).
Handle special cases like castling, en passant, and pawn promotion.

Move Validation:

Write functions to check if a move is valid.
Include checks for checks and checkmates.

Game State Management:

Track the current state of the board.
Keep track of whose turn it is, available moves, and game-ending conditions.

Opponent Logic (Optional):

Implement an AI opponent (start with random moves, then explore minimax or similar algorithms).

Communication Layer:

Expose the game logic via APIs (REST or gRPC) or a CLI for testing.
*/


/* ----- variables ----- */
// var moveFrom *square
// var moveTo *square
var moveFromPrior *square
var moveToPrior *square
var pieces = map[string]func(string, int, int) chessPiece {
	"pn": func(piece string, row, col int) chessPiece {
		return &pawn{&basePiece{piece, row, col}}
	},
	"rk": func(piece string, row, col int) chessPiece {
		return &rook{&basePiece{piece, row, col}}
	},
	"kt": func(piece string, row, col int) chessPiece {
		return &knight{&basePiece{piece, row, col}}
	},
	"bp": func(piece string, row, col int) chessPiece {
		return &bishop{&basePiece{piece, row, col}}
	},
	"qn": func(piece string, row, col int) chessPiece {
		return &queen{&basePiece{piece, row, col}}
	},
	"kg": func(piece string, row, col int) chessPiece {
		return &king{&basePiece{piece, row, col}}
	},
} 


/* ----- types ----- */

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

type basePiece struct {
	piece string
	row int // rank
	col int // file
}

type pawn struct {
	*basePiece
}

type rook struct {
	*basePiece
}

type knight struct {
	*basePiece
}

type bishop struct {
	*basePiece
}

type queen struct {
	*basePiece
}

type king struct {
	*basePiece
}


/* ----- interfaces ----- */

type chessPiece interface {
	fullName() string
	name() string
	color() string
	hasMoved() bool
	updatePosition(newRow, newCol int)
	isValidMove(to *square, cg *chessGame) (bool, bool, bool, error)
}


/* ----- methods ----- */

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

func (cg *chessGame) inCheck(moveTo *square) (bool, error, chessPiece) { 
	for row := 0; row < 8; row++ {
		for col := 0; col < 8; col++ {
			moveFrom, err := cg.getSquare(row, col)
			if err != nil {
				fmt.Errorf("There was an error")
			}
			if moveFrom.cp != nil && (moveFrom.cp.color() != moveTo.cp.color()) {
				canMove, _, _, err := moveFrom.cp.isValidMove(moveTo, cg)
				return canMove, err, moveFrom.cp
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

func (cg *chessGame) makeMove(moveFrom, moveTo *square) {

	// does the moveFrom square contain a chessPiece?
	if moveFrom.cp == nil {
		err := fmt.Errorf("This square is empty. Choose a square with your piece")
		log.Fatal("Error:", err)
	}

	// does the basePiece belong to the player?
	if (cg.player == "User" && moveFrom.cp.color() == "B") || (cg.player == "Computer" && moveFrom.cp.color() == "W") {
		err := fmt.Errorf("%s cannot move %s color pieces", cg.player, moveFrom.cp.color())
		log.Fatal("Error:", err)
	}

	// is the moveTo square occupied by the current player's piece?
	if moveTo.cp != nil && moveTo.cp.color() == moveFrom.cp.color() {
		moveToName := moveTo.cp.name()
		moveFromName := moveFrom.cp.name()
		// raise error if the move doesn't look like a castle
		if !((moveToName == "kg" && moveFromName == "rk") || (moveToName == "rk" && moveFromName == "kg")) {
			err := fmt.Errorf("moveTo square is occupied with: %s", moveTo.cp.fullName())
			log.Fatal("Error:", err)
		}
	}

	// get the implementation of chessPiece and call its specific isValidMove 
	if canMove, promotePawn, castleRook, err := moveFrom.cp.isValidMove(moveTo, cg); canMove {
		// update the chessPiece's position, first check castleRook
		if castleRook {
			if moveFrom.col == 0 {
				moveFrom.cp.updatePosition(moveFrom.row, 3)
				sq, _ := cg.getSquare(moveTo.row, 3)
				sq.cp = moveFrom.cp
				moveFrom.cp = nil
			} else {
				moveFrom.cp.updatePosition(moveFrom.row, 5)
				sq, _ := cg.getSquare(moveTo.row, 5)
				sq.cp = moveFrom.cp
				moveFrom.cp = nil
			}
			moveFrom.cp = nil
		} else {
			moveFrom.cp.updatePosition(moveTo.row, moveTo.col)
			// update the squares
			moveTo.cp = moveFrom.cp
			moveFrom.cp = nil
			
			// promote the pawn, if indicated
			// if pawn reaches opposing last row, player can promote pawn to a: queen, bishop, knight, rook
			// the pawn must be replaced
			// the replacing piece does not have to be a captured piece, so you can have multiple queens after promoting
			if promotePawn {
				// will need a way to pause the game and have player select the new piece
				majorMinor := []string{"rk","kt","bp","qn","kg"}
				n := rand.Intn(5)
				constructor := pieces[majorMinor[n]]
				cg.board[moveTo.row][moveTo.col] = square{moveTo.row, moveTo.col, constructor(majorMinor[n] + "-" + moveTo.cp.color(), moveTo.row, moveTo.col)}
			}
		}

	} else if err != nil {
		log.Fatal("Error:", err)
	}
}

func (b *basePiece) fullName() string {
	return b.piece
}

func (b *basePiece) name() string {
	return strings.Split(b.piece, "-")[0]
}

func (b *basePiece) color() string {
	return strings.Split(b.piece, "-")[1]
}

func (b *basePiece) opponentColor() string {
	if strings.Split(b.piece, "-")[1] == "W" {
		return "B"
	} else {
		return "W"
	}
}

func (b *basePiece) hasMoved() bool {
	boolStr := strings.Split(b.piece, "-")[2]
	hasMoved, err := strconv.ParseBool(boolStr)
		if err != nil {
			err := fmt.Errorf("Error converting '%s' to bool: %v\n", boolStr, err)
			log.Fatal("Error:", err)
		}
	return hasMoved
}

func (b *basePiece) updatePosition(newRow, newCol int) {
	b.row = newRow
	b.col = newCol
}

// each piece will have its own isValidMove to evaluate its specific rules
func (p *pawn) isValidMove(moveTo *square, cg *chessGame) (bool, bool, bool, error) {
	// implement isValidMove for other pieces
	// evaluate edge cases
	//		DONE: pawn: moving two spaces off start line
	// 		HALF DONE: pawn: promotion
	//		DONE: pawn: en passant
	// 		DONE: pawn: taking opponent's piece diagonally
	//		rook and king: castling
	// 		all but knight: collisions
	// DONE: keep track of captured pieces

	// ensure pawn only moves forward
	if (p.color() == "W" && moveTo.row >= p.row) || (p.color() == "B" && moveTo.row <= p.row) {
		err := fmt.Errorf("%s can only move forward", p.fullName())
		return false, false, false, err
	}

	// ensure pawn only moves one row at a time, but allow two rows off starting row
	advance := int(math.Abs(float64(moveTo.row - p.row)))
	if advance > 1 {
		if advance == 2 && (p.row == 1 || p.row == 6) {
			// pawn is allowed to move two spaces forward off starting row
		} else {
			err := fmt.Errorf("%s can only move one space at a time, or two from starting row", p.fullName())
			return false, false, false, err
		}
	}

	// ensure pawn move is in the same column, unless capturing an opponent's piece or en passant
	lateral := int(math.Abs(float64(moveTo.col - p.col)))
	if lateral == 1 {
		// pawn is allowed to capture opponent's pieces
		if moveTo.cp != nil && moveTo.cp.color() != p.color() {
			moveTo.cp.updatePosition(-1, -1)
			if moveTo.cp.color() == "W" {
				cg.wCaptured = append(cg.wCaptured, moveTo.cp)
			} else {
				cg.bCaptured = append(cg.bCaptured, moveTo.cp)
			}
		// evaluate en passant
		} else if (moveToPrior.cp != nil &&
			moveToPrior.cp.name() == "pn" && 
			int(math.Abs(float64(moveToPrior.row - moveFromPrior.row))) == 2 &&
			moveToPrior.row == p.row &&
			moveToPrior.col == moveTo.col) {
			moveToPrior.cp.updatePosition(-1, -1)
			if moveToPrior.cp.color() == "W" {
				cg.wCaptured = append(cg.wCaptured, moveToPrior.cp)
			} else {
				cg.bCaptured = append(cg.bCaptured, moveToPrior.cp)
			}
			moveToPrior.cp = nil	
		} else {
			err := fmt.Errorf("%s cannot move laterally unless capturing or using en passant", p.fullName())
			return false, false, false, err
		}
	} else if lateral > 1 {
		err := fmt.Errorf("%s cannot move laterally more than one square", p.fullName())
		return false, false, false, err
	}

	// check for pawn promotion
	promote := false
	if p.color() == "W" && moveTo.row == 0 || p.color() == "B" && moveTo.row == 7 {
		promote = true
	}

	return true, promote, false, nil
}

func (r *rook) isValidMove(moveTo *square, cg *chessGame) (bool, bool, bool, error) {
	verticalMove := r.col == moveTo.col && r.row != moveTo.row
	horizontalMove := r.col != moveTo.col && r.row == moveTo.row

	// check if move is along a single axis and not diagonal
	if !(verticalMove != horizontalMove) { // xor logic
		err := fmt.Errorf("%s can only move along a single row or column", r.fullName())
		return false, false, false, err
	}
	
	// check for collision
	if verticalMove {
		for i := r.row + 1; i < moveTo.row; i++ {
			if cg.board[i][moveTo.col].cp != nil {
				msg := fmt.Sprintf("%s collides at (%d, %d)", r.fullName(), i, r.col)
				err := fmt.Errorf(msg)
				return false, false, false, err
			}
		}
	}
	if horizontalMove {
		fmt.Println("in horizontalMove")
		for i := r.col + 1; i < moveTo.col; i++ {
			if cg.board[moveTo.row][i].cp != nil {
				msg := fmt.Sprintf("%s collides at (%d, %d)", r.fullName(), i, r.col)
				err := fmt.Errorf(msg)
				return false, false, false, err
			}
		}
	}

	// check if castling
	//		DONE: rook or king cannot have moved
	//		DONE: spaces between pieces must be unoccupied
	//		king cannot be in check
	//		the spaces king moves to cannot be under attack
	//		the spaces the king moves through cannot be under attack
	if horizontalMove && moveTo.cp.name() == "kg" && !moveTo.cp.hasMoved() && !r.hasMoved() {
		// evaluate if any starting, thru or ending sqaures would be in check while castling
		squaresToEval := []*square{}
		if r.col == 0 {
			for col := 4; col > 1; col-- {
				sq, err := cg.getSquare(r.row, col)
				if err != nil {
					fmt.Errorf("There's an error")
				}
				squaresToEval = append(squaresToEval, sq)
			}
		} else {
			for col := 4; col < 7; col++ {
				sq, err := cg.getSquare(r.row, col)
				if err != nil {
					fmt.Errorf("There's an error")
				}
				squaresToEval = append(squaresToEval, sq)
			}
		}
		fmt.Println("I'm going to check if you can castle!")
		for s := range squaresToEval {
			// if canCheck, err, cp := cg.inCheck(squaresToEval[s]); canCheck {
			if 1 == 0 {
				// fmt.Println(canCheck)
				// fmt.Println("This castle doesn't work")
				// fmt.Println(cp.fullName())
				fmt.Println(squaresToEval[s])
				// return false, false, false, err
			}
		}
		// the castle is valid, move the king
		if r.col == 0 {
			moveTo.cp.updatePosition(moveTo.row, 2)
			sq, _ := cg.getSquare(moveTo.row, 2)
			sq.cp = moveTo.cp
			moveTo.cp = nil
		} else {
			moveTo.cp.updatePosition(moveTo.row, 6)
			sq, _ := cg.getSquare(moveTo.row, 6)
			sq.cp = moveTo.cp
			moveTo.cp = nil
		}
	}

	return true, false, true, nil
}

func (k *knight) isValidMove(moveTo *square, cg *chessGame) (bool, bool, bool, error) {
	return true, false, false, nil
}

func (b *bishop) isValidMove(moveTo *square, cg *chessGame) (bool, bool, bool, error) {
	return true, false, false, nil
}

func (q *queen) isValidMove(moveTo *square, cg *chessGame) (bool, bool, bool, error) {
	return true, false, false, nil
}

func (k *king) isValidMove(moveTo *square, cg *chessGame) (bool, bool, bool, error) {
	return true, false, false, nil
}


/* ----- functions ----- */

// helper function to retrieve piece constructor from pieces factory 
func getPieceConstructor(piece string) (func(string, int, int) chessPiece, error) {
	constructor, found := pieces[piece] 
	if !found {
		err := fmt.Errorf("%s not found in pieces factory", piece)
		log.Fatal("Error:", err)
	}
	return constructor, nil
}

func main() {
	// create vars and initialize board and game
	var board [8][8]square
	var player string
	var wCaptured []chessPiece
	var bCaptured []chessPiece
	var err error
	
	cg := &chessGame{board, player, wCaptured, bCaptured}

	// parse vars from front end and load the game
	mCreateNewGame := false
	mBoard := [8][8]string{
		{"rk-B-0", "  e-  ", "  e-  ", "  e-  ", "kg-B-0", "bp-B-1", "kt-B-1", "rk-B-0"},
		{"pn-B-0", "pn-B-0", "pn-B-0", "qn-B-0", "pn-B-0", "pn-B-0", "  e-  ", "pn-B-0"},
		{"kt-B-0", "  e-  ", "  e-  ", "pn-B-0", "bp-B-0", "  e-  ", "  e-  ", "  e-  "},
		{"  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "pn-B-1", "  e-  "},
		{"qn-W-1", "  e-  ", "pn-W-1", "pn-W-1", "  e-  ", "  e-  ", "  e-  ", "  e-  "},
		{"  e-  ", "pn-W-0", "  e-  ", "  e-  ", "bp-W-1", "bp-W-0", "kt-W-0", "  e-  "},
		{"pn-W-0", "  e-  ", "  e-  ", "  e-  ", "pn-W-0", "pn-W-0", "pn-W-0", "pn-W-0"},
		{"rk-W-0", "kt-W-0", "  e-  ", "  e-  ", "kg-W-0", "  e-  ", "  e-  ", "rk-W-0"},
	}
	/*
		{"rk-B-0", "kt-B-0", "bp-B-0", "qn-B-0", "kg-B-0", "bp-B-0", "kt-B-0", "rk-B-0"},
		{"pn-B-0", "pn-B-0", "pn-B-0", "qn-B-0", "pn-B-0", "pn-B-0", "pn-B-0", "pn-B-0"},
		{"  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  "},
		{"  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  "},
		{"  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  "},
		{"  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  ", "  e-  "},
		{"pn-W-0", "pn-W-0", "pn-W-0", "pn-W-0", "pn-W-0", "pn-W-0", "pn-W-0", "pn-W-0"},
		{"rk-W-0", "kt-W-0", "bp-W-0", "qn-W-0", "kg-W-0", "bp-W-0", "kt-W-0", "rk-W-0"},
	*/
	mPlayer := "User" // User or Computer
	mMoveFrom := []int{7, 7}
	mMoveTo := []int{7, 4}
	mMoveFromPrior := [2]int{7, 2}
	mMoveToPrior := [2]int{5, 4}
	mWCaptured := []string{}
	mBCaptured := []string{}
	
	if mCreateNewGame {
		cg.newGame()
		fmt.Println()
		fmt.Println("New game:")
		cg.displayBoard()
	} else {
		cg.loadGame(mBoard, mPlayer, mWCaptured, mBCaptured)
		fmt.Println()
		fmt.Println("Game loaded:")
		cg.displayBoard()

		moveFromPrior, err = cg.getSquare(mMoveFromPrior[0], mMoveFromPrior[1])
		if err != nil {
			log.Fatal("Error:", err)
		}

		moveToPrior, err = cg.getSquare(mMoveToPrior[0], mMoveToPrior[1])
		if err != nil {
			log.Fatal("Error:", err)
		}

		fromSquare, err := cg.getSquare(mMoveFrom[0], mMoveFrom[1])
		if err != nil {
			log.Fatal("Error:", err)
		}

		toSquare, err := cg.getSquare(mMoveTo[0], mMoveTo[1])
		if err != nil {
			log.Fatal("Error:", err)
		}

		fmt.Println()
		msg := fmt.Sprintf("Attempting move from (%d, %d) to (%d, %d)", fromSquare.row, fromSquare.col, toSquare.row, toSquare.col)
		fmt.Println(msg)
		cg.makeMove(fromSquare, toSquare)
		fmt.Println()
		fmt.Println("Board after move:")
		cg.displayBoard()
		fmt.Println()
		cg.displayCapturedPieces()
		//cg.displayBoardPositions()
	}
}
