let words = [];

console.log(window.location.href);
let url = `wss://${window.location.href.replace("https://", "")}connect`
console.log(url);
let ws = new WebSocket(url)
ws.onopen = () => {
  console.log(`connected`);
  let el = document.getElementById('status');
  el.textContent = 'connected';
  el.style.color = 'green';

  ws.send(JSON.stringify({
    "kind": 'join',
    "from": 'user',
    "data": 'user',
  }))
}
ws.onclose = (err) => {
  console.log(JSON.stringify(err, ["message", "arguments", "type", "name"]));
}
ws.onmessage = (ev) => {
    console.log('Received', ev.data)
    let msg = JSON.parse(ev.data)
    switch (msg["kind"]) {
       case setup:
          console.log(msg["data"]);
          break;
       default:
          break;
    }
    
}

const sendStart = () => {
  console.log('sending start');
  ws.send(JSON.stringify({
    "kind": 'start',
    "from": 'user',
  }))
}

const sendMessage = () => {
    content = document.getElementById('input').value
    console.log('sending', content)
    ws.send(JSON.stringify({
        "name": "user",
        "message": content,
    }))
}

