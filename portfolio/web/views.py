import json
from .models import (
    Home, 
    About,
    Contact,
    Projects
)
from .tic_tac_toe import TicTacToe, MoveIsTaken
from django.views import View
from django.shortcuts import render
from django.http import JsonResponse

class HomeView(View):
    model = Home
    template_name = 'home.html'

    def get(self, request):

        return render(request, self.template_name)

class ContactView(View):
    model = Contact
    template_name = 'contact.html'

    def get(self, request):
        context = {
            'title': 'Contact'
        }

        return render(request, self.template_name, context)

class AboutView(View):
    model = About
    template_name = 'about.html'

    def get(self, request):
        context = {
            'title': 'About'
        }

        return render(request, self.template_name, context)

class ErcotView(View):
    model = Projects
    template_name = 'ercot.html'

    def get(self, request):
        context = {
            'title': 'Ercot'
        }

        return render(request, self.template_name, context)

class BlsView(View):
    model = Projects
    template_name = 'bls.html'

    def get(self, request):
        context = {
            'title': 'Bls'
        }

        return render(request, self.template_name, context)

class TicTacToeView(View):
    template_name = 'tictactoe.html'

    def get(self, request):
        context = {
            'title': 'Tic Tac Toe'
        }

        return render(request, self.template_name, context)

class TicTacToeBoardView(View):
    template_name = 'tictactoe.html'

    def get(self, request):
        game = TicTacToe()
        request.session['board'] = game.board
        return JsonResponse({'board': game.board})
    
    def post(self, request):
        try:
            board = request.session.get('board')
            if board is None:
                raise ValueError('Game not found')
            
            game = TicTacToe(board=board)
            body = json.loads(request.body)
            if body == 'compMove':
                response = game.comp_move
            else:
                row = int(body['row'])
                col = int(body['col'])
                response = game.user_move(row, col)

            winner = game.victory_for()
            request.session['board'] = game.board
            return JsonResponse({'board': game.board, 'winner': winner})
        except MoveIsTaken as e:
            response = {'error': str(e)}
            return JsonResponse(response, status=400)
        