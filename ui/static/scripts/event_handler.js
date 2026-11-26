async function refresh() {
  try {
    const [status, events] = await Promise.all([
      fetch('/api/status').then(r => r.json()),
      fetch('/api/events').then(r => r.json())
    ]);
    document.getElementById('status').textContent =
      status.connected ? 'WebSocket connected' : 'WebSocket disconnected — events will be queued';
    const subscriptions = status.subscriptions || [];
    document.getElementById('subscription').textContent = subscriptions.length
      ? 'Host registered to: ' + subscriptions.join(', ')
      : 'No host is registered to events — events stay queued until the host sends RegisterEvent';
    const screen = status.screen || {name: 'idle'};
    document.getElementById('screen').textContent = screen.name === 'deny'
      ? 'Device screen (' + screen.device + '): access denied — ' + screen.display_text
      : 'Device screen: idle';
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

for (const form of document.querySelectorAll('form.action')) {
  form.onsubmit = async (submit) => {
    submit.preventDefault();
    const message = document.getElementById('message');
    try {
      const res = await fetch('/api/actions/' + form.dataset.action, {
        method: 'POST',
        headers: {'Content-Type': 'application/json'},
        body: JSON.stringify({id: form.elements.id.value})
      });
      if (!res.ok) {
        message.textContent = await res.text();
      } else {
        const result = await res.json();
        message.textContent = result.subscribed
          ? 'Event ' + result.id + ' queued'
          : 'Event ' + result.id + ' queued — it is sent when the host registers to events';
      }
    } catch (err) {
      message.textContent = String(err);
    }
    refresh();
  };
}

refresh();

setInterval(refresh, 2000);
