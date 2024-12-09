from django.db import models

class TicTacToeResult(models.Model):
    difficulty_level = models.CharField(max_length=12)
    winner = models.CharField(max_length=3)
    created_at = models.DateTimeField(auto_now_add=True)
