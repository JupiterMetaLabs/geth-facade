module custom-backend-example

go 1.22.6

require github.com/JupiterMetaLabs/geth-facade v0.0.0

require github.com/gorilla/websocket v1.5.3 // indirect

replace github.com/JupiterMetaLabs/geth-facade => ../../
