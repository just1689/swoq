# swoq
<a href="https://github.com/just1689/swoq/releases"><img src="https://img.shields.io/badge/version-alpha-blue" /></a>&nbsp;

Scaleable Websockets over Queues


The goal of this project is to decouple websocket front-ends from the worker nodes that carry out server-side work. This allows shaping and throttling for workers.

Currently the project supports the following:
- Adding a Gorilla websocket client to the default http mux.
- Connecting to NATs
- Creating workers


Missing functionality:
- Closing a client queue subscriber when the client disconnects from the websocket server
- Allowing clients to indicate they are a previous client reconnecting
