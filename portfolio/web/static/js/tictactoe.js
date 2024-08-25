document.addEventListener('DOMContentLoaded', () => {
    const board = document.getElementById('board');
    const startGameButton = document.getElementById('start-game');
    const messageDiv = document.getElementById('message');
    const difficultyTitle = document.getElementById('difficulty-title');
    const difficultyOptions = document.getElementById('difficulty-options');
    const dropdownDifficulty = document.querySelectorAll('.dropdown-difficulty');
    let gamePaused = true; // Pauses event listener

    startGameButton.addEventListener('click', showDifficultyOptions);
    difficultyTitle.addEventListener('click', selectDifficulty);
    document.getElementById('easy').addEventListener('click', () => startGame('Easy'));
    document.getElementById('intermediate').addEventListener('click', () => startGame('Intermediate'));
    document.getElementById('unbeatable').addEventListener('click', () => startGame('Unbeatable'));

    fetchStats('All'); // default for page refresh

    board.addEventListener('click', function(event) {
        if (gamePaused) {
            return;
        }
        const target = event.target;
        if (target.tagName === 'TD') {
            const row = target.parentNode.rowIndex;
            const col = target.cellIndex;
            userMove(row, col);
        }
    });

    dropdownDifficulty.forEach(item => {
        item.addEventListener('click', function(event) {
            event.preventDefault();
            const difficulty = event.target.getAttribute('difficulty-level');
            fetchStats(difficulty);
        });
    });

    function showDifficultyOptions() {
        messageDiv.textContent = '';
        board.style.display = 'none';
        difficultyTitle.style.display = 'block';
        difficultyOptions.style.display = 'block';
    }

    function selectDifficulty() {
        difficultyTitle.style.display = 'none';
        difficultyOptions.style.display = 'none';
    }

    async function startGame(difficulty) {
        fetchStats(difficulty=difficulty);
        selectDifficulty();
        try {
            const response = await fetch(startUrl, {
                method: 'GET',
                headers: {
                    'X-CSRFToken': csrfToken,
                    'Difficulty-Level': difficulty
                }
            });

            const data = await response.json();
            updateBoard(data.board);
            gamePaused=false;
            board.style.display = 'table';
            startGameButton.textContent = 'Start new game';

        } catch (error) {
            console.error('Error starting game:', error);
            alert('An error occurred while starting the game.');
        }
    }

    async function userMove(row, col) {
        try {
            gamePaused=true;
            const response = await fetch(moveUrl, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'X-CSRFToken': csrfToken
                },
                body: JSON.stringify({ 'row': row, 'col': col })
            });

            const data = await response.json();
            if (data.error) {
                alert(data.error);
                gamePaused = false;
            } else {
                updateBoard(data.board);
                if (data.winner) {
                    if (data.winner === 'Tie') {
                        messageDiv.textContent = `${data.winner}!`;
                    } else {
                        messageDiv.textContent = `${data.winner} wins!`;
                    }
                    fetchStats(difficulty=data.difficulty_level);
                    return;
                }

                setTimeout(async () => {
                    await compMove();
                }, 1000);
            }

        } catch (error) {
            console.error('Error making move:', error);
            alert('An error occurred while making the move: ' + error.message);
            gamePaused=false;
        }
    }

    async function compMove() {
        try {
            const response = await fetch(moveUrl, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'X-CSRFToken': csrfToken
                },
                body: JSON.stringify('compMove')
            });

            const data = await response.json();
            if (data.error) {
                alert(data.error);
            } else {
                updateBoard(data.board);
                if (data.winner) {
                    if (data.winner === 'Tie') {
                        messageDiv.textContent = `${data.winner}!`;
                    } else {
                        messageDiv.textContent = `${data.winner} wins!`;
                    }
                    fetchStats(difficulty=data.difficulty_level);
                    return;
                }
            }

        } catch (error) {
            console.error('Error making move:', error);
            alert('An error occurred while making the move.');
        
        } finally {
            gamePaused=false;
        }
    }

    function updateBoard(boardState) {
        board.innerHTML = ''; // Clear the board
        for (let i = 0; i < boardState.length; i++) {
            const row = document.createElement('tr');
            for (let j = 0; j < boardState[i].length; j++) {
                const cell = document.createElement('td');
                cell.innerHTML = boardState[i][j] === '' ? '&nbsp;' : boardState[i][j];
                row.appendChild(cell);
            }
            board.appendChild(row);
        }
    }

    async function fetchStats(difficulty='All') {
        try {
            const response = await fetch(`/tictactoe/?difficulty=${difficulty}`, {
                method: 'GET',
                headers: {
                    'X-CSRFToken': csrfToken,
                    'X-Requested-With': 'XMLHttpRequest',
                },
            });
            
            const data = await response.json();
            if (difficulty === 'All') {
                subtitle = `${difficulty} difficulty levels`;
            } else {
                subtitle = `${difficulty} difficulty level`;
            };
            updateStats(data, subtitle);

        } catch (error) {
            console.error('Error fetching stats:', error);
            alert('An error occurred while fetching the stats.');
        }
    }

    function updateStats(data, subtitle) {
        const userWinRate = parseFloat(data.user_win_rate).toFixed(1);
        const compWinRate = parseFloat(data.comp_win_rate).toFixed(1);
        const tieRate = parseFloat(data.tie_rate).toFixed(1);

        document.querySelector('#metrics').innerHTML = `
            <p style="margin-top: 40px;"><en style="font-style: italic;">${subtitle}</en></p>
            <p>Total games played: ${data.games_played}</p>
            <p>User wins: ${data.user_wins}</p>
            <p>Computer wins: ${data.comp_wins}</p>
            <p>Ties: ${data.ties}</p>
            <p>User win rate: ${userWinRate}%</p>
            <p>Computer win rate: ${compWinRate}%</p>
            <p>Tie rate: ${tieRate}%</p>
        `;
    }
});
