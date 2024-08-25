from django.urls import path
from .views import (
    HomeView, 
    ContactView,
    AboutView,
    TicTacToeView,
    TicTacToeBoardView,
    ErcotView,
    BlsView
)

urlpatterns = [
    path('', HomeView.as_view(), name="home"),
    path('contact', ContactView.as_view(), name="contact"),
    path('about', AboutView.as_view(), name="about"),
    path('tictactoe/', TicTacToeView.as_view(), name="tictactoe"),
    path('tictactoe/board', TicTacToeBoardView.as_view(), name="tictactoe_board"),
    path('ercot', ErcotView.as_view(), name="ercot"),
    path('bls', BlsView.as_view(), name="bls"),
]