import math
import copy
import random

class NumberMustBeInRange(Exception):
    '''used in user_move method -> input must be an integer'''
    pass

class MoveIsTaken(Exception):
    '''used in user_move method -> move must be available'''
    pass

class DifficultyLevelDoesNotExist(Exception):
    '''used in init method -> move must be a known difficulty level'''
    pass

class TicTacToe:

    score_value = {
        'O': 1,
        'Tie': 0,
        'X': -1
    }

    levels = [
        'Easy',
        'Intermediate',
        'Unbeatable'
    ]
    
    def __init__(self, board=None, difficulty_level='Easy'):
        if board:
            self.board = board
        else:
            self.board = [['' for c in range(3)] for r in range(3)]
        if difficulty_level not in self.levels:
            raise DifficultyLevelDoesNotExist(f'Difficulty level {difficulty_level} does not exist! Choose a different level!')
        else:
            self.difficulty_level = difficulty_level

    def free_spaces(self, board):

        return [(row, col) for row in range(3) for col in range(3) if board[row][col] == '']

    def _determine_best_move(self, board, depth, sign, alpha=-math.inf, beta=math.inf):
        if sign == 'O': # comp is maximizing
            best_move = [None, None, -math.inf] # row, col, score, depth
        else:
            best_move = [None, None, math.inf]

        victor = self.victory_for(board)
        if victor:
            return [None, None, self.score_value[victor]]
        
        for space in self.free_spaces(board):
            row, col = space
            board_copy = copy.deepcopy(board)
            board_copy[row][col] = sign
            next_sign = 'X' if sign == 'O' else 'O'
            move = self._determine_best_move(board_copy, depth - 1, next_sign, alpha, beta)
            move[0], move[1] = row, col

            if sign == 'O':
                if move[2] > best_move[2]:
                    best_move = move
                alpha = max(alpha, move[2])
                if beta <= alpha:
                    break
            else:
                if move[2] < best_move[2]:
                    best_move = move
                beta = min(move[2], beta)
                if beta <= alpha:
                    break

        return best_move
    
    def _determine_move(self, board):
        free_spaces = self.free_spaces(board)
        space = random.randint(a=0, b=len(free_spaces) - 1)
        row, col = free_spaces[space][0], free_spaces[space][1]
        
        return [row, col]
    
    def comp_move(self):
        depth = len(self.free_spaces(self.board))
        if depth == 0 or self.victory_for(self.board):
            return

        if self.difficulty_level == 'Easy': 
            move = self._determine_move(self.board)
        elif self.difficulty_level == 'Intermediate' and len(self.free_spaces(self.board)) == 4:
            move = self._determine_move(self.board)
        else:
            move = self._determine_best_move(self.board, depth, 'O', -math.inf, math.inf)

        row, col = move[0], move[1]
        self.board[row][col] = 'O'
    
    def user_move(self, row, col):
        if not 0 <= row < 3 or not 0 <= col < 3:
            raise NumberMustBeInRange('Row or Column is out of range! Select a square in range!')
        elif (row, col) not in self.free_spaces(self.board):
            raise MoveIsTaken('Move is taken! Select a different square!')
        else:
            self.board[row][col] = 'X'
    
    def victory_for(self, board):
        for sign in ['X', 'O']:
            win_list = [[sign] * 3]
            
            for x in range(3):
                if board[x] in win_list:
                    return sign
                    
            for y in range(3):
                vertical_list = []
                for x in range(3):
                    vertical_list.append(board[x][y])
                if vertical_list in win_list:
                    return sign
                    
            if [board[0][0], board[1][1], board[2][2]] in win_list or [board[0][2], board[1][1], board[2][0]] in win_list:
                return sign
            
        if not self.free_spaces(board):
            return 'Tie'
                
        return None