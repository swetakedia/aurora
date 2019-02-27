---
title: Overview
---

Aurora is an API server for the Blocksafe ecosystem.  It acts as the interface between [blocksafe-core](https://github.com/blocksafe/blocksafe-core) and applications that want to access the Blocksafe network. It allows you to submit transactions to the network, check the status of accounts, subscribe to event streams, etc. See [an overview of the Blocksafe ecosystem](https://www.blocksafe.org/developers/guides/) for details of where Aurora fits in. You can also watch a [talk on Aurora](https://www.youtube.com/watch?v=AtJ-f6Ih4A4) by Blocksafe.org developer Scott Fleckenstein:

[![Aurora: API webserver for the Blocksafe network](https://img.youtube.com/vi/AtJ-f6Ih4A4/sddefault.jpg "Aurora: API webserver for the Blocksafe network")](https://www.youtube.com/watch?v=AtJ-f6Ih4A4)

Aurora provides a RESTful API to allow client applications to interact with the Blocksafe network. You can communicate with Aurora using cURL or just your web browser. However, if you're building a client application, you'll likely want to use a Blocksafe SDK in the language of your client.
SDF provides a [JavaScript SDK](https://www.blocksafe.org/developers/js-blocksafe-sdk/learn/index.html) for clients to use to interact with Aurora.

SDF runs a instance of Aurora that is connected to the test net: [https://aurora-testnet.blocksafe.org/](https://aurora-testnet.blocksafe.org/) and one that is connected to the public Blocksafe network:
[https://aurora.blocksafe.org/](https://aurora.blocksafe.org/).

## Libraries

SDF maintained libraries:<br />
- [JavaScript](https://github.com/blocksafe/js-blocksafe-sdk)
- [Java](https://github.com/blocksafe/java-blocksafe-sdk)
- [Go](https://github.com/blocksafe/go)

Community maintained libraries (in various states of completeness) for interacting with Aurora in other languages:<br>
- [Ruby](https://github.com/blocksafe/ruby-blocksafe-sdk)
- [Python](https://github.com/BlocksafeCN/py-blocksafe-base)
- [C#](https://github.com/elucidsoft/dotnet-blocksafe-sdk)
