let board = []
let result = []
let pacman = false
let width = 0
let height = 0

let classToNum = {
    wall:1,
    dot:0,
    'ghost-lair':2,
    'power-up':3,
    empty:4,
    pacman:5
}
function createGrid() {
    width = Number(document.getElementById('width').value)
    height = Number(document.getElementById('height').value)
    let grid = document.getElementById('grid')
    board.forEach((x) => {grid.removeChild(x),board = board.slice(1)})
    grid.setAttribute("style",`width:${width * 20}px;`)
    for (let i = 0; i < width * height; i++) {
        const square = document.createElement('button')
        square.setAttribute('onclick',`setItem(${i})`)
        square.className = 'empty'
        result.push(classToNum[square.className])
        grid.appendChild(square)
        board.push(square)
    }
    updateResultArray()
}

function selectItem(item) {
    let items = ["dot","wall","power-up","ghost-lair","pacman","empty"]
    items.forEach((x) => {
        if (x !== item) {
            document.getElementById(x).style.opacity = '70%'
        }
    })
    document.getElementById(item).style.opacity = '100%'
    localStorage.setItem('item',item)
}

function setItem(value) {
    let className = localStorage.getItem('item');
    if (className == 'pacman') {
        if (!pacman) {
            board[value].className = className
            pacman = true
        }
    } else {
        if (board[value].className == 'pacman') {
            pacman = false
        }
        board[value].className = className
    }
    updateResultArray()
}

function buildWall() {
    for (let i = 0; i < board.length; i++) {
        if (i < width || i > board.length - width || i % width == width - 1 || i%width == 0 ) {
            board[i].className = 'wall'
        }
    }
    updateResultArray()
}

function resetGrid() {
    board.forEach((x) => x.className = 'empty')
    updateResultArray()
}

function fillWithDot() {
    for (let i = 0; i < board.length; i++) {
        if (board[i].className == 'empty') {
            board[i].className = 'dot'
        }
    }
    updateResultArray()
}

function updateResultArray() {
    for (let i = 0; i < board.length; i++) {
        result[i] = classToNum[board[i].className]
    }
    document.getElementById('result').innerHTML = ''
    for (let i = 0; i < result.length;i++) {
        if (i % width == 0 ) {
            document.getElementById('result').innerHTML += ' <br> '
        }
        document.getElementById('result').innerHTML +=`${result[i]}, `
    }
}
