const type = {
    0: 'dot',
    1: 'wall',
    2: 'ghost_lair',
    3: 'power-up',
    4: 'empty',
    5: 'pacman'
}


let boards = [],
    width,
    copyBoards = [],
    start,
    pacman = new Pacman()

let score = 0;
let lives = 3;
async function createBoard(mapInd = 0) {
    let map = []
    let resp = await fetch("maps.json")
    let res = await resp.json()
    map = res[mapInd]['layout']
    width = res[mapInd]['width']
    let ghostSpawns = -2
    for (let i = 0; i < map.length; i++) {
        const square = document.createElement('div')
        square.classList.add(type[map[i]])
        if (map[i] == 5) {
            pacman.pos = i
            start = i
            square.className = `pacman${pacman.className}`
        }
        if (map[i] == 2) {
            if (ghostSpawns < 0) {
                ghostSpawns++
            } else if (ghostSpawns < 4) {
                ghosts[ghostSpawns].start = i
                ghosts[ghostSpawns].currentPos = i
                square.classList.add(ghosts[ghostSpawns].className)
                square.classList.add('ghost')
                ghostSpawns++
            }
        }
        boards.push(square)
        copyBoards.push(square.className)
    }
}

async function playGames(grid) {
    await createBoard(0)
    grid.setAttribute("style", `width:${width * 20}px;`)
    for (let i = 0; i < boards.length; i++) {
        grid.appendChild(boards[i])
    }
    animatePacMan()
    ghosts.forEach(moveGhost)
}
document.addEventListener('DOMContentLoaded', () => {
    playGames(document.getElementById('grid'))
})

document.addEventListener('keyup', event => movePacman(event))
function resetBoard() {
    let point = document.querySelectorAll('.dot').length + document.querySelectorAll('.power-up').length
    if (point == 0) {
        clearInterval(move)
        ghosts.forEach(ghost => clearInterval(ghost.timer))
        for (let i = 0; i < boards.length; i++) {
            console.log()
            boards[i].className = copyBoards[i]
            pacman.pos = start
        }
    }
}

function updateScore(score) {
    document.getElementById('score').innerText = score
}

function updateLives(lives) {
    document.getElementById('lives').innerText = lives
}

let move = setInterval(() => { }, 12)
let waiting = ""
let current = ""
let keyCodes = {
    L: 37,
    R: 39,
    U: 38,
    D: 40
}
let clear = false
function movePacman(e) {
    console.log(pacman.pos)
    let key = e.keyCode
    if (clear) {

    }
    if (key != keyCodes[current]) {
        console.log(key, keyCodes[current])
        waiting = e
    }
    switch (key) {
        case 37:
        case 81:
            if (canMove('L')) {
                clearInterval(move)
                move = setInterval(() => {
                    if (canMove(waiting)) {
                        clearInterval(move)
                        movePacman(waiting)
                    } else {
                        if (canMove('L')) {
                            current = 'L'
                            boards[pacman.pos].className = 'empty'
                            pacman.pos = pacman.pos - 1
                            boards[pacman.pos].className == 'dot' ? score += 10 : score = score
                            updateScore(score)
                            pacman.className = 'L'
                            boards[pacman.pos].className = 'pacmanL'
                            resetBoard()
                        }
                    }
                }, 100)
            }
            break;
        case 39:
        case 68:
            if (canMove('R')) {
                current = 'R'
                clearInterval(move)
                move = setInterval(() => {
                    if (canMove(waiting)) {
                        clearInterval(move)
                        movePacman(waiting)
                    } else {
                        if (canMove('R')) {
                            boards[pacman.pos].className = 'empty'
                            pacman.pos = pacman.pos + 1
                            boards[pacman.pos].className == 'dot' ? score += 10 : score = score
                            updateScore(score)
                            pacman.className = 'R'
                            boards[pacman.pos].className = 'pacmanR'
                            resetBoard()
                        }
                    }
                }, 100)
            }
            break;
        case 38:
        case 90:
            if (canMove('U')) {
                current = 'U'
                clearInterval(move)
                move = setInterval(() => {
                    if (canMove(waiting)) {
                        clearInterval(move)
                        movePacman(waiting)
                    } else {
                        if (canMove('U')) {
                            boards[pacman.pos].className = 'empty'
                            pacman.pos = pacman.pos - width
                            boards[pacman.pos].className == 'dot' ? score += 10 : score = score
                            updateScore(score)
                            pacman.className = 'U'
                            boards[pacman.pos].className = 'pacmanU'
                            resetBoard()
                        }
                    }
                }, 100)
            }
            break;
        case 40:
        case 83:
            if (canMove('D')) {
                current = 'D'
                clearInterval(move)
                move = setInterval(() => {
                    if (canMove(keyCodes[waiting.keyCode])) {
                        clearInterval(move)
                        movePacman(waiting)
                    } else {
                        if (canMove('D')) {
                            boards[pacman.pos].className = 'empty'
                            pacman.pos = pacman.pos + width
                            boards[pacman.pos].className == 'dot' ? score += 10 : score = score
                            updateScore(score)
                            pacman.className = 'D'
                            boards[pacman.pos].className = 'pacmanD'
                            resetBoard()
                        }
                    }
                }, 100)
            }
            break;
    }
}