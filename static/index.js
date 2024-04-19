const el = document.body.appendChild(document.createElement("p"))
el.textContent = "testinnnnn";

ws = new WebSocket('ws://localhost:8080/connect')
ws.onopen = () => {
    const p = el.appendChild(document.createElement('p'))
    p.textContent = 'connected';
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

