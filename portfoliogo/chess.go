package main

import (
	"fmt"
	"strings"
	"math"
	"errors"
	"log"
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
Keep track of whose turn it is, available moves, and gamee-nding conditions.

Opponent Logic (Optional):

Implement an AI opponent (start with random moves, then explore minimax or similar algorithms).

Communication Layer:

Expose the game logic via APIs (REST or gRPC) or a CLI for testing.
*/

type chessBoard struct {
	board [8][8]chessPiece
	player string
	wCaptured []chessPiece
	bCaptured []chessPiece
}

func (cb *chessBoard) newGame(player string) {
	cb.player = player
	majorMinor := []string{"rk","kt","bp","qn","kg","bp","kt","rk"}
	for i := 0; i < 8; i++ {

		if constructor, found := pieces[majorMinor[i]]; !found {
			err := errors.New("piece type 'pn' not found in map")
			log.Fatal("Error:", err)
		} else {
			cb.board[0][i] = constructor(majorMinor[i] + "-B", 0, i)
			cb.board[7][i] = constructor(majorMinor[i] + "-W", 7, i)
		}

		if constructor, found := pieces["pn"]; !found {
			err := errors.New("piece type 'pn' not found in map")
			log.Fatal("Error:", err)
		} else {
			cb.board[1][i] = constructor("pn-B", 1, i)
			cb.board[6][i] = constructor("pn-W", 6, i)
		}
	}
}

func (cb *chessBoard) loadGame(game [8][8]string, player string, wCaptured []string, bCaptured []string) {
	cb.player = player
	for row := range game {
		for col := range game[row] {
			piece := strings.TrimSpace(strings.Split(game[row][col], "-")[0])

			if piece == "e" {
				cb.board[row][col] = nil
				continue
			}
			
			constructor, found := pieces[piece] 
			if !found {
				err := fmt.Errorf("Piece not found in pieces map %s", game[row][col])
				log.Fatal("Error:", err)
			}

			cb.board[row][col] = constructor(game[row][col], row, col)
		}
	}
	// load captured pieces
	for _, p := range wCaptured {
		piece := strings.TrimSpace(strings.Split(p, "-")[0])
		constructor, _ := getPieceConstructor(piece)
		cb.wCaptured = append(cb.wCaptured, constructor(p, -1, -1))
	}

	for _, p := range bCaptured {
		piece := strings.TrimSpace(strings.Split(p, "-")[0])
		constructor, _ := getPieceConstructor(piece)
		cb.bCaptured = append(cb.bCaptured, constructor(p, -1, -1))
	}
}

func (cb *chessBoard) getPiece(pos position) chessPiece {
	return cb.board[pos.row][pos.col]
}

func (cb *chessBoard) getPosition(row, col int) (position, error) {
	if row > 7 || row < 0 || col > 7 || col < 0 {
		return position{}, fmt.Errorf("From position is out of range: %d %d", row, col)
	}
	return cb.board[row][col].Position(), nil
}

func (cb *chessBoard) displayBoard() {
	for row := 0; row < 8; row++ {
		line := make([]string, 8)
		for col := 0; col < 8; col++ {
			if cb.board[row][col] == nil {
				line[col] = " e- "
			} else {
				line[col] = cb.board[row][col].fullName()
			}
		}
		fmt.Println(line)
	}
}

func (cb *chessBoard) displayBoardPositions() {
	for row := 0; row < 8; row++ {
		line := make([]position, 8)
		for col := 0; col < 8; col++ {
			if cb.board[row][col] == nil {
				line[col] = position{-1, -1}
			} else {
				line[col] = cb.board[row][col].Position()
			}
		}
		fmt.Println(line)
	}
}

func (cb *chessBoard) displayCapturedPieces() {
	var w []string
	var b []string
	for _, piece := range cb.wCaptured {
		w = append(w, piece.fullName())
	}
	for _, piece := range cb.bCaptured {
		b = append(b, piece.fullName())
	}
	fmt.Println("White pieces captured:", w)
	fmt.Println("Black pieces captured:", b)
}

type position struct {
	row, col int
}

type chessPiece interface {
	fullName() string
	Name() string
	Color() string
	Position() position
	SetPosition(newPos position)
	makeMove(to position, cb *chessBoard)
	isValidMove(to position, cb *chessBoard) (bool, error)
}

type basePiece struct {
	piece string
	pos position
}

func (b *basePiece) fullName() string {
	return b.piece
}

func (b *basePiece) Name() string {
	return strings.Split(b.piece, "-")[0]
}

func (b *basePiece) Color() string {
	return strings.Split(b.piece, "-")[1]
}

func (b *basePiece) Position() position {
	return b.pos
}

func (b *basePiece) SetPosition(newPos position) {
	b.pos = newPos
}

func (b *basePiece) makeMove(to position, cb *chessBoard) {
	toPiece := cb.getPiece(to)

	// does the basePiece belong to the player?
	if (cb.player == "User" && b.Color() == "B") || (cb.player == "Computer" && b.Color() == "W") {
		err := fmt.Errorf("%s cannot move %s pieces", cb.player, b.Color())
		log.Fatal("Error:", err)
	}

	// is the to position occupied by the current player's piece?
	if toPiece != nil && toPiece.Color() == b.Color() {
		err := fmt.Errorf("To position is occupied with: %s", toPiece.fullName())
		log.Fatal("Error:", err)
	}

	// get the implementation of chessPiece and call its specific isValidMove 
	fromPiece := cb.getPiece(b.pos)
	if canMove, err := fromPiece.isValidMove(to, cb); canMove {
		// update the board
		cb.board[b.pos.row][b.pos.col] = nil
		cb.board[to.row][to.col] = fromPiece
		fromPiece.SetPosition(to)
	} else if err != nil {
		log.Fatal("Error:", err)
	}
}

type pawn struct {
	*basePiece
}

type rook struct {
	basePiece
}

type knight struct {
	basePiece
}

type bishop struct {
	basePiece
}

type queen struct {
	basePiece
}

type king struct {
	basePiece
}

var pieces = map[string]func(string, int, int) chessPiece {
	"pn": func(piece string, row, col int) chessPiece {
		return &pawn{&basePiece{piece, position{row, col}}}
	},
	"rk": func(piece string, row, col int) chessPiece {
		return &rook{basePiece{piece, position{row, col}}}
	},
	"kt": func(piece string, row, col int) chessPiece {
		return &knight{basePiece{piece, position{row, col}}}
	},
	"bp": func(piece string, row, col int) chessPiece {
		return &bishop{basePiece{piece, position{row, col}}}
	},
	"qn": func(piece string, row, col int) chessPiece {
		return &queen{basePiece{piece, position{row, col}}}
	},
	"kg": func(piece string, row, col int) chessPiece {
		return &king{basePiece{piece, position{row, col}}}
	},
}

// helper function to retrieve piece constructor from pieces factory 
func getPieceConstructor(piece string) (func(string, int, int) chessPiece, error) {
	constructor, found := pieces[piece] 
	if !found {
		err := fmt.Errorf("%s not found in pieces factory", piece)
		log.Fatal("Error:", err)
	}
	return constructor, nil
}

// each piece will have its own isValidMove to evaluate its specific rules
func (p *pawn) isValidMove(to position, cb *chessBoard) (bool, error) {
	// implement isValidMove for other pieces
	// evaluate edge cases
	//		DONE: pawn: moving two spaces off start line
	// 		pawn: promotion
	//		pawn: en passant
	// 		DONE: pawn: taking opponent's piece diagonally
	//		rook and king: castling
	// 		all but knight: collisions
	// keep track of captured pieces

	// ensure pawn only moves forward
	if (p.Color() == "W" && to.row >= p.pos.row) || (p.Color() == "B" && to.row <= p.pos.row) {
		err := fmt.Errorf("%s can only move forward", p.fullName())
		return false, err
	}

	// ensure pawn only moves one row at a time, but allow two rows off starting row
	advance := int(math.Abs(float64(to.row - p.pos.row)))
	if advance > 1 {
		if advance == 2 && (p.pos.row == 1 || p.pos.row == 6) {
			// pawn is allowed to move two spaces forward off starting row
		} else {
			err := fmt.Errorf("%s can only move one space at a time, or two from starting row", p.fullName())
			return false, err
		}
	}

	// ensure pawn move is in the same column, unless capturing an opponent's piece or en passant
	lateral := int(math.Abs(float64(to.col - p.pos.col)))
	toPiece := cb.getPiece(to)
	if lateral > 0 {
		// pawn is allowed to capture opponent's pieces
		if toPiece != nil && toPiece.Color() != p.Color() {
			toPiece.SetPosition(position{row: -1, col: -1})
			if toPiece.Color() == "W" {
				cb.wCaptured = append(cb.wCaptured, toPiece)
			} else {
				cb.bCaptured = append(cb.bCaptured, toPiece)
			}
		//} else if {
			// evaluate en passant
			// need to keep track of previous move
		} else {
			err := fmt.Errorf("%s cannot move laterally unless capturing or using en passant", p.fullName())
			return false, err
		}
	}

	// allow pawn promotion

	return true, nil
}

func (p *rook) isValidMove(to position, cb *chessBoard) (bool, error) {
	return true, nil
}

func (p *knight) isValidMove(to position, cb *chessBoard) (bool, error) {
	return true, nil
}

func (p *bishop) isValidMove(to position, cb *chessBoard) (bool, error) {
	return true, nil
}

func (p *queen) isValidMove(to position, cb *chessBoard) (bool, error) {
	return true, nil
}

func (p *king) isValidMove(to position, cb *chessBoard) (bool, error) {
	return true, nil
}


func main() {
	// create vars and initialize board
	var board [8][8]chessPiece
	var player string
	var wCaptured []chessPiece
	var bCaptured []chessPiece
	cb := chessBoard{board, player, wCaptured, bCaptured}

	// parse vars from front end:
	//		mainCreateNewGame
	//		mainBoard
	//		mainPlayer (move)
	//		mainFrom
	//		mainTo
	//		mainWhiteCaptured
	//		mainBlackCaptured

	// create vars and load the game
	mainCreateNewGame := false
	mainBoard := [8][8]string{
		{"rk-B", "kt-B", "bp-B", "qn-B", "kg-B", "bp-B", "kt-B", "rk-B"},
		{"pn-B", "pn-B", "pn-B", "pn-B", " e-", "pn-B", "pn-B", "pn-B"},
		{" e- ", " e- ", " e- ", " e- ", " e- ", " e- ", " e- ", " e- "},
		{" e- ", " e- ", " e- ", " e- ", " e- ", " e- ", " e- ", " e- "},
		{" e- ", " e- ", " e- ", " e- ", "pn-B", " e- ", " e- ", " e- "},
		{" e- ", " e- ", " e- ", "pn-W", " e- ", " e- ", " e- ", " e- "},
		{"pn-W", "pn-W", "pn-W", " e- ", "pn-W", "pn-W", "pn-W", "pn-W"},
		{"rk-W", "kt-W", "bp-W", "qn-W", "kg-W", "bp-W", "kt-W", "rk-W"},
	}
	mainPlayer := "Computer" // User or Computer
	mainFrom := []int{4, 4}
	mainTo := []int{5, 3}
	mainWCaptured := []string{}
	mainBCaptured := []string{}
	
	if mainCreateNewGame {
		cb.newGame("User") // new games always start with User
	} else {
		cb.loadGame(mainBoard, mainPlayer, mainWCaptured, mainBCaptured)
	}

	fromPos, err := cb.getPosition(mainFrom[0], mainFrom[1])
	if err != nil {
		log.Fatal("Error:", err)
	}

	toPos, err := cb.getPosition(mainTo[0], mainTo[1])
	if err != nil {
		log.Fatal("Error:", err)
	}
	
	piece := cb.getPiece(fromPos)
	piece.makeMove(toPos, &cb)
	cb.displayBoard()
	//cb.displayBoardPositions()
	cb.displayCapturedPieces()
}
