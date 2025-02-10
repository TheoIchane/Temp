class Ghost {
    constructor(className, speed, start = 0) {
        this.className = className
        this.speed = speed
        this.start = start
        this.isScared = false
        this.timer = 0
        this.currentPos = start
    }

}
const ghosts = [
    new Ghost('blinky', 200, 0),
    new Ghost('pinky', 300, 0),
    new Ghost('inky', 400, 0),
    new Ghost('clyde', 500, 0)
]

function moveGhost(ghost) {
    const directions = [-1, 1, width, -width]
    let direction = directions[Math.floor(Math.random() * directions.length)]
    ghost.timer = setInterval(() => {
        let next = ghost.currentPos + direction
        if (!boards[next].classList.contains('wall') && !boards[next].classList.contains('ghost')) {
            boards[ghost.currentPos].classList.remove('ghost', ghost.className)
            ghost.currentPos = next
            boards[ghost.currentPos].classList.add(ghost.className, 'ghost')
        } else {
            direction = directions[Math.floor(Math.random() * directions.length)]
        }
        if (boards[ghost.currentPos].classList.contains('pacman')) {
            killPacman()
            return
        }
    }, ghost.speed)
}
