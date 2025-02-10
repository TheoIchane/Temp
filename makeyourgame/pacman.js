const fps = 15

function animatePacMan() {
    setTimeout(() => {
        requestAnimationFrame(animatePacMan);
    }, 1000 / fps);
    if (pacman.state == 3) {
        pacman.state = 0
    } else {
        pacman.state++
    }
    boards[pacman.pos].style.backgroundPositionX = `${pacman.state * -20}px`
}

class Pacman {
    constructor() {
        this.className = 'L'
        this.state = 0
        this.currentPos = 0
        this.start = 0
    }
}

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

function killPacman() {
    lives--
    if (lives == 0) {
        gameOver()
    } else {
        updateLives()
        boards[pacman.pos].remove('pacman')
        pacman.pos = pacman.start
        boards[pacman.pos].add('pacman')
        ghosts.forEach((ghost) => {
            clearInterval(ghost.timer)
            boards[ghost.currentPos].classList.remove('ghost', ghost.className)
            ghost.currentPos = ghost.start
            boards[ghost.currentPos].classList.add(ghost.className, 'ghost')
        })
    }
}