from django.urls import path
from .views import (
    HomeView, 
    ContactView,
    About,
    TicTacToeView,
    TicTacToeBoard,
    ErcotView,
    BlsView
)

urlpatterns = [
    path('', HomeView.as_view(), name="home"),
    path('contact', ContactView.as_view(), name="contact"),
    path('about', About.as_view(), name="about"),
    path('tictactoe', TicTacToeView.as_view(), name="tictactoe"),
    path('tictactoe/board', TicTacToeBoard.as_view(), name="tictactoe_board"),
    path('ercot', ErcotView.as_view(), name="ercot"),
    path('bls', BlsView.as_view(), name="bls"),
]