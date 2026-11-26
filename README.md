This is a simple GT7 emulator allowing for easy testing and debugging of a system integrating with the websocket feature.

The emulator is a device in client mode: it opens the WebSocket connection to the host. The host is the management software that you test.

All incoming and outgoing messages are saved in a sqlite db and sent to the frontend.

The db is whipped on every boot.

# How to
go run cmd/main.go --device-id {device-id} --token {token} --target wss://{endpoint}

Open http://localhost:8080 for the frontend.

# Heartbeat
The emulator sends a `Heartbeat` request at the start of each connection and every 60 seconds after that.

If the host gives `HBI` in its `Heartbeat` response, the emulator uses that interval in seconds.

The emulator answers a `Heartbeat` request of the host. While the host sends such requests in each interval, the emulator sends none.

If the host sends no `Heartbeat` request or response for two intervals and 10 seconds, the emulator closes the connection. Then it connects again.

# Event registration
The emulator sends events to a registered host only. The host registers with `RegisterEvent`, for example with the event name `App.*`.

An event stays in the db with the status `queued` until the host is registered to it. The emulator sends the queued events after the `RegisterEvent` response.

A registration is valid for one connection. After a new connection, the host must send `RegisterEvent` again.

# Actions
The frontend has three actions. Each action sends the event `App.CardIdent` with the id that you enter.

| Action | Endpoint | Event |
| --- | --- | --- |
| Card read | `POST /api/actions/card-read` | Card at the reader `INTERNAL`. The id is the `UID`. |
| Barcode | `POST /api/actions/barcode` | Barcode at the reader `BARCODE1`. The id is the barcode data. |
| PIN | `POST /api/actions/pin` | PIN at the device. The id is the PIN and has digits only. |

The body of each endpoint is a JSON object, for example `{"id":"1234"}`.

# Requests of the host
The emulator answers these requests with state 0:

- `Heartbeat`
- `System.SetHostInfo`
- `GetDeviceInfo`
- `RegisterEvent` and `UnregisterEvent`
- `App.StartUnlockProcess`. The emulator then sends the event `App.UnlockProcessResult` with `AccessGranted` true.
- `App.StartDenyProcess`. The emulator sends no event after the response. See "Deny process".

The emulator answers all other requests with state 255 (unknown command).

# Deny process
When the host sends `App.StartDenyProcess`, the emulator does these steps:

1. It answers with state 0 and an empty `Data` object. If a parameter has a wrong type, it answers with state 2.
2. It shows the deny screen with the `DisplayText` of the request. The frontend shows this screen in the line "Device screen".
3. It goes back to the idle screen after `Timeout_ms`. If the request has no `Timeout_ms`, the time is 3 seconds.

The deny completes the identification. A later `App.StartUnlockProcess` does not get the id of the denied card, barcode or PIN.

# Test
go test ./...
