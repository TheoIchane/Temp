const type = {
    0:'dot' ,
    1:'wall',
    2:'ghost_lair' ,
    3:'power-up',
    4:'empty',
    5: 'pacman'
}


let boards = [],
copyBoards = [],
start,
pacman = {
    className:'L',
    state:0,
    pos:0
}
document.addEventListener('DOMContentLoaded',() => {
    let score = 0;
    let lives = 3;
    let width 
    const grid = document.getElementById('grid');
    function createBoard(mapInd = 0) {
        let map = []
        fetch("maps.json")
            .then((res) => res.json())
            .then((res) => {
                map = res[mapInd]['layout']
                width = res[mapInd]['width']
                grid.setAttribute("style",`width:${width * 20}px;`)
                for (let i = 0; i < map.length; i++) {
                    const square = document.createElement('div')
                    square.classList.add(type[map[i]])
                    if (map[i] == 5) {
                        pacman.pos = i
                        start = i
                        square.className = `pacman${pacman.className}`
                    }
                    grid.appendChild(square)
                    boards.push(square)
                    copyBoards.push(square.className)
                }
                animatePacMan()
            })
            .catch((e) => console.log(console.error(e)))
    }
    createBoard(0)
    let move = setInterval(() => {},12)
    let waiting = ""
    let current = ""
    let keyCodes = {
        L:37,
        R:39,
        U:38,
        D:40
    }
    let clear = false
    updateLives(lives)
    function movePacman(e) {
        let key = e.keyCode
        if (clear) {

        }
        if (key != keyCodes[current]) {
            console.log(key,keyCodes[current])
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
                }}},100)  
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
                }}},100) 
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
                }}},100)
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
                }}},100)
                }
            break;
        }
    }
    document.addEventListener('keydown',event => movePacman(event))
    function canMove(direction) {
        switch (direction) {
            case 'L':
                if (pacman.pos % width == 0) {
                    boards[pacman.pos].className = 'empty'
                    pacman.pos = pacman.pos + width - 1
                    boards[pacman.pos].className = 'pacman'
                    return canMove(direction)
                }
                if (boards[pacman.pos - 1].className == 'wall' || boards[pacman.pos - 1].className == 'ghost_lair') {
                    return false
                }
                return true
            case 'R':
                if (pacman.pos % width == width - 1) {
                    boards[pacman.pos].className = 'empty'
                    pacman.pos = pacman.pos - width + 1
                    boards[pacman.pos].className = 'pacman'
                    return canMove(direction)
                }
                if (boards[pacman.pos + 1].className == 'wall' || boards[pacman.pos + 1].className == 'ghost_lair') {
                    return false
                }
                return true
            case 'U':
                if (pacman.pos - width < 0 || boards[pacman.pos - width].className == 'wall' || boards[pacman.pos - width].className == 'ghost_lair') {
                    return false
                }
                return true
            case 'D':
                if (pacman.pos + width >= boards.length || boards[pacman.pos + width].className == 'wall' || boards[pacman.pos + width].className == 'ghost_lair') {
                    return false
                }
                return true
        }
    }

    function resetBoard() {
        let point = document.querySelectorAll('.dot').length + document.querySelectorAll('.power-up').length
        if (point == 0) {
            clearInterval(move)
            for (let i = 0; i < boards.length; i++) {  
                console.log()
                boards[i].className = copyBoards[i]
                pacman.pos = start
            }
        }
    }
})

const fps = 15

function animatePacMan() {
    setTimeout(() => {
        requestAnimationFrame(animatePacMan);
    },1000 / fps);
    if (pacman.state == 3) {
        pacman.state = 0
    } else {
        pacman.state++
    }
    boards[pacman.pos].style.backgroundPositionX = `${pacman.state * -20}px`
}
function updateScore(score) {
    document.getElementById('score').innerText = score
}

function updateLives(lives) {
    document.getElementById('lives').innerText = lives
}

function moveGhost(color) {
    
}