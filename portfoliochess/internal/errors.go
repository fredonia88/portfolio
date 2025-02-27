package main

import (
	"fmt"
	"errors"
	"runtime"
	"strings"
	"log"
)

var (
	errCollision = errors.New("COLLISION")
	errOutOfRange = errors.New("OUT_OF_RANGE")
	errEmptySquare = errors.New("EMPTY_SQUARE")
	errOccupiedSquare = errors.New("OCCUPIED_SQUARE")
	errInvalidMove = errors.New("INVALID_MOVE")
	errHasMovedConversion = errors.New("HAS_MOVED_CONVERSION")
	errMissingPiece = errors.New("MISSING_PIECE")
)

type chessError struct {
	details error
	message string
	file string
	line int
}

func handleError(err error) {
	if err != nil {
		log.Fatal("Error ", err)
	}
}

func (e *chessError) Error() string {
	return fmt.Sprintf("[%v]: %s (%s: %d)", e.details, e.message, e.file, e.line)
}

func newChessError(details error, message string, args ...interface{}) error {
	_, file, line, _ := runtime.Caller(1)
	subdirs := strings.Split(file, "/")
	file = subdirs[len(subdirs)-1]

	return &chessError{
		details,
		fmt.Sprintf(message, args...),
		file,
		line,
	}
}



// from move.go

// collision error
// empty square
// square is occupied with player's piece
// invalid piece
// fatal error, move cannot be made

// pawn must move forward
// pawn must move one space at a time, or two from starting position
// pawn cannot move laterally more than one square
// pawn cannot move laterally unless capturing a piece or en passant

// rook must move along a single row or column
// rook or king has already moved while attempting to castle

// knight must move vertically and horiztonally by 2 and 1 squares

// bishop must move diagonally

// king must move one square at a time, unless castling


// from piece.go

// error converting string bool 
// piece not found in factory


// game.go

// piece not found in factory 
// getSquare error
// King is in check
// getSquare is out of range


// main.go
// getSquare error

