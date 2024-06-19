# gordon

Gordon is a lightweight protocol for sharing linked documents with a minimal syntax. These documents are transmitted as bytestreams, served over UDP using DTLS.

The protocol is defined in the [mint](https://github.com/vinyl-linux/mint) document DDL and may be found in the [./mint](./mint) directory.

Furthermore, this repository contains a sample implementation in [./sample-app](./sample-app), along with client code and a sample client. (A graphical client, of sorts, can also be found at [github.com/jspc/gordon-browser](https://github.com/jspc/gordon-browser).

The sample-app is deployed to a server at `//gordon.beasts.jspc.pw/` for testing.


## Design Decisions

Gordon prioritises

1. Speed
2. Security
3. Accessibility

(In any particular order)

Every decision we make is made around these. We use binary streams over UDP for speed. We use DTLS for security. We use highly structured, largely plaintext documents to be as accessible as humanly possible.

Think we're doing it wrong, or failing one of those goals? Tell us.


## The Protocol

Documents are served following a pretty common pattern; send a request, receive a response.

A request, as defined in [./mint/requests.mint](./mint/requests.mint), specify a `Verb`, a document `ID`, and an optional hash of Args.

In response, we receive a `Page`.

A `Page` is a structured document, containing specific metadata, and split into sections. `Page`s _also_ can link to other documents in one of two ways:

1. An index link is a shorthand for linking from within text; the argument `[l:0]`, for instance, refers to element 0 in the list of `Links`
2. A relationship is a triple representing how a specific document links to an other; gordon comes with a handful of predicates


## The Encoding

Payloads are encoded to bytestreams using [mint](https://github.com/vinyl-linux/mint). This is a non-describing stream of binary data, with validations and transformations.

These payloads look much the same as thrift, protobuf, and so on.

### Why mint?

Why anything? This whole project is made of daft little projects I've created, and mint encodes at a decent clip, and has reasonably small payloads.


## The Network

Sending and Requesting data is done over UDP using DTLS.


## Licence

BSD 3-Clause License

Copyright (c) 2024, James Condron
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
   list of conditions and the following disclaimer.

2. Redistributions in binary form must reproduce the above copyright notice,
   this list of conditions and the following disclaimer in the documentation
   and/or other materials provided with the distribution.

3. Neither the name of the copyright holder nor the names of its
   contributors may be used to endorse or promote products derived from
   this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
