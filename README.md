This is a simple GT7 emulator allowing for easy testing and debugging of a system integrating with the websocket feature.

A heartbeat is sent at boot and every 60 seconds thereafter.

All incoming and outgoing messages are saved in a sqlite db and sent to the frontend.

The db is whipped on every boot.

# How to
go run cmd/main.go --device-id {device-id} --token {token} --target wss://{endpoint}
