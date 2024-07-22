document.addEventListener('DOMContentLoaded', () => {
    const board = document.getElementById('board');
    const startGameButton = document.getElementById('start-game');
    let gamePaused = true;

    document.getElementById('start-game').addEventListener('click', startGame);

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

    async function startGame() {
        try {
            const response = await fetch(startUrl, {
                method: 'GET',
                headers: {
                    'X-CSRFToken': csrfToken
                }
            });

            const data = await response.json();
            updateBoard(data.board);
            gamePaused=false;
            board.classList.remove('disabled');
            board.style.display = 'table'; // Display the board when the game starts
            startGameButton.textContent = 'Start new game'; // Change button text to "New game"

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
            } else {
                updateBoard(data.board);
            }

            setTimeout(async () => {
                await compMove();
            }, Math.random() * 1000 + 1000); // Random delay between 1 and 2 seconds

        } catch (error) {
            console.error('Error making move:', error);
            alert('An error occurred while making the move.');
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
                body: JSON.stringify('isCompMove')
            });

            const data = await response.json();
            if (data.error) {
                alert(data.error);
            } else {
                updateBoard(data.board);
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
});
