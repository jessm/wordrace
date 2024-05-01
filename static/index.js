////////////////////////////
// WEBSOCKET HANDLING STUFF
////////////////////////////

let getSocketUrl = () => {
  let url = window.location.href;
  if (url.includes('https')) {
    url = url.replace('https', 'wss')
  } else {
    url = url.replace('http', 'ws')
  }

  return url + "connect";
}

let ws = new WebSocket(getSocketUrl())
ws.onopen = () => {
  console.log('connected')
  let el = document.getElementById('status-emoji')
  el.textContent = '🟢'

  // doStart()
}
ws.onclose = (err) => {
  console.log('Socket Close:', JSON.stringify(err, ["message", "arguments", "type", "name"]));
  let el = document.getElementById('status-emoji')
  el.textContent = '🔴'
}
ws.onmessage = (ev) => {
  console.log('Received', ev.data)
  let msg = JSON.parse(ev.data)
  let handleFuncs = {
    'joined': handleJoin,
    'setup': handleSetup,
    'err': handleErr,
    'foundWord': handleFoundWord,
  }
  handleFuncs[msg["kind"]](msg)
}

const sendCmd = (cmd) => {
  console.log('sending', cmd)
  ws.send(JSON.stringify(cmd))
}

////////////////////////////
// VISUALS
////////////////////////////

let createLetterEl = (letter = 'A') => {
  let el = document.createElement('div')
  el.className = 'letter-button'
  el.textContent = letter.toUpperCase()
  return el
}

////////////////////////////
// GAME LOGIC
////////////////////////////

let gState = {
  userName: '',
  doneSetup: false,
  words: [],
  wordsFound: [],

  letterEls: [],
}

const handleJoin = (msg) => {
  if (gState.userName == '') {
    gState.userName = msg["data"]
  }
}

const doStart = () => {
  sendCmd({
    'kind': 'start',
    'from': gState.userName,
  })
}

const doReset = () => {
  sendCmd({
    'kind': 'reset',
  })
}

const handleSetup = (msg) => {
  if (gState.setup) {
    return
  }
  gState.setup = true
  let wordEl = document.getElementById("words")
  for (let i = 0; i < msg["data"]["counts"].length; i++) {
    let wordLen = msg["data"]["counts"][i]
    let el = document.createElement("div")
    el.className = "word"
    for (let j = 0; j < wordLen; j++) {
      let letterEl = document.createElement("div")
      letterEl.className = "empty-letter"
      letterEl.textContent = ' '
      el.appendChild(letterEl)
    }
    gState.words.push(el)
    gState.wordsFound.push(false)
    wordEl.appendChild(el)
  }

  let lettersEl = document.getElementById("letters")
  let letters = msg["data"]['letters'].toUpperCase()
  for (let i = 0; i < letters.length; i++) {
    let letterEl = createLetterEl(letters[i])
    letterEl.onclick = letterOnClick;
    lettersEl.appendChild(letterEl)
    gState.letterEls.push(letterEl)
  }
}

const doTryWord = () => {
  el = document.getElementById('inputWord')
  let word = ''
  for (let i = 0; i < el.children.length; i++) {
    word += el.children[i].textContent
  }
  word = word.toLowerCase()
  
  clearOnClick()

  sendCmd({
    'kind': 'tryWord',
    'from': gState.userName,
    'data': word,
  })
}

const handleErr = (msg) => {
  let msgElement = document.getElementById('serverMessage')
  msgElement.textContent = msg["data"]
  msgElement.style.opacity = 100
  setTimeout(() => {
    msgElement.style.opacity = 0
  }, 2000)
}

const handleFoundWord = (msg) => {
  let pos = msg['data']['pos']
  if (gState.wordsFound[pos]) {
    return
  }
  let scoreEl = null
  let letterColorClass = ''
  if (msg['data']['player'] == gState.userName) {
    letterColorClass = 'your-color'
    scoreEl = document.getElementById('yourScore')
  } else {
    letterColorClass = 'their-color'
    scoreEl = document.getElementById('theirScore')
  }
  gState.wordsFound[pos] = true
  let el = gState.words[pos]
  el.innerHTML = ''

  let word = msg['data']['word'].toUpperCase()
  for (let i = 0; i < word.length; i++) {
    letterEl = document.createElement('div')
    letterEl.classList.add('found-letter', letterColorClass)
    letterEl.textContent = word[i]
    el.appendChild(letterEl)
  }

  // el.textContent = msg['data']['word'].toUpperCase()
  let score = parseInt(scoreEl.textContent)
  score += msg['data']['word'].length
  scoreEl.textContent = score
}

////////////////////////////
// GAME LOGIC HELPERS
////////////////////////////

let getLetterPos = (letter) => {
  for (let i = 0; i < gState.letterEls.length; i++) {
    if (gState.letterEls[i].textContent == letter) {
      return i
    }
  }
  return -1
}

let letterOnClick = (ev) => {
  let letter = ev.target.textContent
  let pos = getLetterPos(letter)
  let letterEl = gState.letterEls[pos]

  // --- all the stuff to revert when done with letter
  letterEl.classList.remove('letter-button')
  letterEl.classList.add('empty-letter-button')
  letterEl.onclick = null
  // --- end stuff

  let inputWordEl = document.getElementById('inputWord')

  inputWordEl.appendChild(createLetterEl(ev.target.textContent))
}

let clearOnClick = () => {
  let inputWordEl = document.getElementById('inputWord')
  inputWordEl.innerHTML = ''

  for (let i = 0; i < gState.letterEls.length; i++) {
    let el = gState.letterEls[i]
    el.classList.remove('empty-letter-button')
    el.classList.add('letter-button')
    el.onclick = letterOnClick
  }
}

let shuffleOnClick = () => {
  let idx = gState.letterEls.length
  while (idx != 0) {
    let randIdx = Math.floor(Math.random() * idx)
    idx--
    [gState.letterEls[idx], gState.letterEls[randIdx]] = [gState.letterEls[randIdx], gState.letterEls[idx]]
  }

  let lettersEl = document.getElementById('letters')
  lettersEl.innerHTML = ''
  for (let i = 0; i < gState.letterEls.length; i++) {
    lettersEl.appendChild(gState.letterEls[i])
  }
}