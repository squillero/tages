Tages
=====

[![License: GPL](https://img.shields.io/badge/license-gpl--3.0-green.svg)](https://opensource.org/licenses/GPL-3.0)
![Language: go](https://img.shields.io/badge/language-go-blue.svg)
![](https://www.google-analytics.com/collect?v=1&t=pageview&tid=UA-28094298-5&cid=4f34399f-f437-4f67-9390-61c649f9b8b2&dl=https%3A%2F%2Fgithub.com%2Fsquillero%2Ftages%2F)

The [iterated prisoner's dilemma](https://en.wikipedia.org/wiki/Prisoner%27s_dilemma) is a famous model of cooperation and conflict in game theory. Its origin can be traced back to the Cold War, and [countless strategies](/strategies.md) for playing it have been proposed so far, either designed by hand or automatically generated by computers. In the 2000s, scholars started focusing on *adaptive players*, that is, players able to classify their opponent's behavior and adopt an effective counter-strategy.

Tages pushes the idea of adaptation one step further: it builds a model of the current adversary from scratch, without relying on any pre-defined archetypes, and tweaks it as the game develops using an evolutionary algorithm; at the same time, it exploits the model to lead the game into the most favorable continuation. Models are compact non-deterministic finite state machines; they are extremely efficient in predicting opponents' replies, without being completely correct by necessity. Experimental results show that such player is able to win several one-to-one games against strong opponents taken from the literature, and that it consistently prevails in round-robin tournaments of different sizes. See the article *Exploiting Evolutionary Modeling to Prevail in Iterated Prisoner's Dilemma Tournaments* (DOI: [10.1109/TCIAIG.2015.2439061](https://dx.doi.org/10.1109/TCIAIG.2015.2439061)) for more details.

> This project marks Squillero's first attempt to learn [Go](http://golang.org/), and the overall quality of the code reflects it.

**Copyright © 2015 Giovanni Squillero.**

Tages is free software: you can redistribute it and/or modify it under the terms of the [GNU General Public License](http://www.gnu.org/licenses/) as published by the *Free Software Foundation*, either [version 3](https://opensource.org/licenses/GPL-3.0) of the License, or (at your option) any later version.
