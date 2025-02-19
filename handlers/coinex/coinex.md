___
[[coinex/save]]
[[coinex/structs]]
[[coinex/tests]]
[[coinex/demo_server]]
___
okay, the Coinex Message Handler, this file holds the necessary functions for reading and correctly parsing incoming coinex messages.

___
## Main Goal :
The goal is to turn a byte array carrying ASCII characters into a standardized format to save inside a CSV file using buffers for efficiency.

___
## Obstacles :
Coinex data messages are compressed using GZIP, we will have to decompress this in order to parse the messages.

___
## Process :
First we need to detect if the message is GZIPed or not, sometimes coinex messages aren't compressed.

Then we parse the uncompressed byte array into a form our program can easily read so we can identify what kind of message this is.

*We parse the message once more according to the type of response is received.*

After it is parsed, we process the message into a standardized form and send it off to the buffer.