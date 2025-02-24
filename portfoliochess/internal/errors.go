package main

import (
	"fmt"
	"errors"
)

var (
	errCollision = errors.New("Your piece %d collides with %d at (%d, %d)")
)

type chessError struct {
	code error
	message string
}

func (e *chessError) Error() string {
	return fmt.Sprintf("[%s]: %s", e.code, e.message)
}

func newChessError(code error, message string, args ...interface{}) error {
	return &chessError{
		code, 
		fmt.Sprintf(message, args...),
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

