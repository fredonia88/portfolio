# Generated by Django 5.0.7 on 2024-12-09 14:49

from django.db import migrations, models


class Migration(migrations.Migration):

    dependencies = [
        ('web', '0003_alter_tictactoeresult_difficulty_level'),
    ]

    operations = [
        migrations.AlterField(
            model_name='tictactoeresult',
            name='winner',
            field=models.CharField(max_length=3),
        ),
    ]
