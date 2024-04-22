let words = [];

console.log('working');

let ws = new WebSocket(`wss://${window.location.href}:8080/connect`)
ws.onopen = () => {
  let el = document.getElementById('status');
  el.textContent = 'connected';
}

ws.onmessage = (ev) => {
    console.log('Received', ev.data)
    let msg = JSON.parse(ev.data)
    const p = el.appendChild(document.createElement('p'))
    p.textContent = msg.message
}

const sendMessage = () => {
    content = document.getElementById('input').value
    console.log('sending', content)
    ws.send(JSON.stringify({
        "name": "user",
        "message": content,
    }))
}

