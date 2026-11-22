async function refresh() {
  try {
    const [status, events] = await Promise.all([
      fetch('/api/status').then(r => r.json()),
      fetch('/api/events').then(r => r.json())
    ]);
    document.getElementById('status').textContent =
      status.connected ? 'WebSocket connected' : 'WebSocket disconnected — events will be queued';
    const tbody = document.getElementById('events');
    tbody.replaceChildren();
    for (const e of events) {
      const tr = document.createElement('tr');
      for (const value of [e.id, e.created_at, e.direction, e.status, JSON.stringify(e.payload)]) {
        const td = document.createElement('td');
        td.textContent = value;
        tr.appendChild(td);
      }
      tbody.appendChild(tr);
    }
  } catch (err) {
    document.getElementById('message').textContent = String(err);
  }
}

document.getElementById('trigger').onclick = async () => {
  const res = await fetch('/api/trigger', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: document.getElementById('payload').value
  });
  document.getElementById('message').textContent =
    res.ok ? 'Event queued' : await res.text();
  refresh();
};

refresh();

setInterval(refresh, 2000);

