document.addEventListener('DOMContentLoaded', (event) => {
    const board = document.getElementById('board');
    board.addEventListener('click', function(event) {
        const target = event.target;
        if (target.tagName === 'TD') {
            const row = target.parentNode.rowIndex;
            const col = target.cellIndex;
            makeMove(row, col);
        }
    });
});

function makeMove(row, col) {
    console.log('in makeMove')
    fetch(moveUrl, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'X-CSRFToken': csrfToken
        },
        body: JSON.stringify({'row': row, 'col': col})
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            alert(data.error);
        } else {
            updateBoard(data);
        }
    });
}

function updateBoard(data) {
    const board = document.getElementById('board');
    const boardState = data.board;
    for (let i = 0; i < boardState.length; i++) {
        for (let j = 0; j < boardState[i].length; j++) {
            board.rows[i].cells[j].innerText = boardState[i][j];
        }
    }
}