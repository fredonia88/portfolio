package main 

import (
	"fmt"
	"log"
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
	hasMoved() bool
	updatePosition(newRow, newCol int)
	isValidMove(to *square, cg *chessGame) (bool, bool, bool, bool, error)
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

// helper function to retrieve piece constructor from pieces factory 
func getPieceConstructor(piece string) (func(string, int, int) chessPiece, error) {
	constructor, found := pieces[piece] 
	if !found {
		err := fmt.Errorf("%s not found in pieces factory", piece)
		log.Fatal("Error:", err)
	}
	return constructor, nil
}