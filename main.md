___
## main()
this function is the first function that runs and is responsible for calling other functions and doing stuff...

#### Inputs:
```
// no inputs
```

#### Outputs:
```
// no outputs
```

### Main Operation
1. Read & Parse JSON Config File
	![[OpenParseJSON.png]]
2. Establish Multiple Websocket Connections
	[[#connectExchange()]]
	[[#routeSubscribe()]]
	![[establishConnections.png]]
3. Receive & Save Responses Asynchronously
	[[#routeResponse()]]
	![[listenAsynchronously.png]]
4. Listen for Graceful Shutdown
	[[#gracefulShutdown()]]
	![[listenGracefulShutdown.png]]
___
## gracefulShutdown()
this function is the first function that runs and is responsible for calling other functions and doing stuff...

#### Inputs:
```
// []*websocket.Conn
```

#### Outputs:
```
// no outputs
```

### Operation
5. Create a Channel & Listen for Shutdown Signals
	![[shutdownChannel.png]]
6. Gracefully Close Connections
	![[closeConnections.png]]

___
## connectExchange()
this function is responsible for making the connection to the exchange server through HTTPS and websockets.

#### Inputs:
```
7. Struct [[#ExchangeConfig]]
```

#### Outputs:
```
8. *websocket.Conn
9. error
```

### Operation
10. Make Connection & Return it. 
    ![[makeConnection.png]]


___
## routeSubscribe()
this function is responsible for sending the subscription message to the exchange so that we are receiving real time market time.

#### Inputs:
```
11. Struct [[#ExchangeConfig]]
12. *websocket.Conn
```

#### Outputs:
```
13. error
```

### Operation
14. Parse Streams into array map
	![[ParseStreams.png]]
15. Loop Through Streams & Send Message
	![[LoopStreams.png]]
___
## routeResponse()
This function is responsible for routing incoming responses to the correlating and correct data handler.

#### Inputs:
```
16. Struct [[#ExchangeConfig]]
17. *websocket.Conn
```

#### Outputs:
```
18. error
```

### Operation
19. Set Luh Defer
	![[LuhDefer.png]]
20. Loop Through Incoming Messages
	[[coinex_handler]]
	[[binance_handler]]
	[[bybit_handler]]
	[[bitfinex handler]]
	![[loopIncoming.png]]

___
