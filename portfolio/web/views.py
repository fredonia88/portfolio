import json
from .models import TicTacToeResult
from .tic_tac_toe import TicTacToe, MoveIsTaken
from django.views import View
from django.shortcuts import render
from django.http import JsonResponse
from django.db.models import Sum

class HomeView(View):
    template_name = 'home.html'

    def get(self, request):

        return render(request, self.template_name)

class ContactView(View):
    template_name = 'contact.html'

    def get(self, request):
        context = {
            'title': 'Contact'
        }

        return render(request, self.template_name, context)

class AboutView(View):
    template_name = 'about.html'

    def get(self, request):
        context = {
            'title': 'About'
        }

        return render(request, self.template_name, context)

class ErcotView(View):
    template_name = 'ercot.html'

    def get(self, request):
        context = {
            'title': 'Ercot'
        }

        return render(request, self.template_name, context)

class BlsView(View):
    template_name = 'bls.html'

    def get(self, request):
        context = {
            'title': 'Bls'
        }

        return render(request, self.template_name, context)

class TicTacToeView(View):
    template_name = 'tictactoe.html'
    model = TicTacToeResult

    def get(self, request):

        games_played = TicTacToeResult.objects.count()
        user_wins = TicTacToeResult.objects.filter(winner='X').count()
        comp_wins = TicTacToeResult.objects.filter(winner='O').count()
        ties = TicTacToeResult.objects.filter(winner='Tie').count()
        user_win_rate = str(0.0 if games_played == 0 else round((user_wins / games_played) * 100, 1)) + '%'
        comp_win_rate = str(0.0 if games_played == 0 else round((comp_wins / games_played) * 100, 1)) + '%'

        context = {
            'title': 'Tic Tac Toe',
            'games_played': games_played,
            'user_wins': user_wins,
            'comp_wins': comp_wins,
            'ties': ties,
            'user_win_rate': user_win_rate,
            'comp_win_rate': comp_win_rate
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
            game = TicTacToe(board=board)

            body = json.loads(request.body)
            if body == 'compMove':
                game.comp_move()
            else:
                row = int(body['row'])
                col = int(body['col'])
                game.user_move(row, col)

            winner = game.victory_for(game.board)
            if winner: 
                TicTacToeResult.objects.create(
                    winner = winner
                )

            request.session['board'] = game.board
            return JsonResponse({'board': game.board, 'winner': winner})
        except MoveIsTaken as e:
            response = {'error': str(e)}
            return JsonResponse(response, status=400)
        