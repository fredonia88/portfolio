package main 

import (
	"strings"
	"strconv"
)


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

type chessPiece interface {
	fullName() string
	name() string
	color() string
	getRow() int 
	getCol() int
	getHasMoved() (error)
	setHasMoved()
	updatePosition(newRow, newCol int)
	isValidMove(to *square, cg *chessGame) (bool, bool, bool, error)
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

func (b *basePiece) getRow() int {
	return b.row
}

func (b *basePiece) getCol() int {
	return b.col
}

func (b *basePiece) getHasMoved() (err error) {
	boolStr := strings.Split(b.piece, "-")[2]
	hasMoved, err := strconv.ParseBool(boolStr)
	if err != nil {
		err = newChessError(errHasMovedConversion, "Error converting '%s' to bool: %v\n", boolStr, err)
	}
	if hasMoved {
		err = newChessError(errHasMoved, "%d has already moved and is ineligible to castle", b.fullName())
	}
	return
}

func (b *basePiece) setHasMoved() {
	b.piece = b.piece[:len(b.piece)-1] + "1"
}

func (b *basePiece) updatePosition(newRow, newCol int) {
	b.row = newRow
	b.col = newCol
}

// helper function to retrieve piece constructor from pieces factory 
func getPieceConstructor(piece string) (constructor func(string, int, int) chessPiece, err error) {
	constructor, found := pieces[piece] 
	if !found {
		err = newChessError(errMissingPiece, "%s not found in pieces factory", piece)
	}
	return
}