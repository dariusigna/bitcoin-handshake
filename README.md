# Bitcoin Handshake
 
A simple application that connects to a bitcoin target node and performs a successful handshake.
It uses `btcd` for the ease of testing. 

### Testing
1. Install [btcd](https://github.com/btcsuite/btcd).
2. Run`btcd --configfile ./btcd.conf`
3. Update `test.env` with the desired config values or set the environment variables
    - `TARGET_NODE_ADDRESS` - The address of the node you want to connect to
    - `NETWORK` - The network you want to connect to. (simnet, mainnet, testnet, regtest)
4. Run `go run main.go`

### Build binary for linux
   Run `env GOOS=linux GOARCH=amd64 go build -o bitcoin-handshake-linux main.go`

### What it does?
1. The application starts a simple server with a peer handler
2. The peer handler creates a tcp connection to the target node.
2. The server performs a handshake with the target node.
3. It sleeps for 5 seconds.
4. Disconnects from the target node.


### How the handshake protocol works?
1. The server sends a version message to the target node.
2. The target node responds with his version message.
3. The server sends a verack message to the target node.
4. The target node sends a verack message to the server.

##### If all the above steps are successful and in the correct order, the handshake is complete.


