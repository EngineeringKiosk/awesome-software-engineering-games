# Contribution Guide

## How to add a new game to the list

Please feel free to open a pull request if you know of a game that fits the list and is not mentioned here.  
If you're in doubt whether a particular game is a good fit for the list, **don't open an issue, but create a pull request right away** because that's easier to handle. Thanks! :smiley:

### Criteria to be accepted

Each game on the list should meet at least one of the following criteria:

- It is a game :smiley:
- You need to write some form of code/logic/procedure during the gameplay
- It is programmable (you can control the game with code)
- It is highly known/appreciated in the software engineering community (e.g., for its logical thinking gameplay)
- is at least 3 months old
- it has the key theme of *technology*, *startups*, *business* or *software engineering* in the gameplay

We do understand that complex criteria for such a list are hard.
See this list as opinionated and curated.
We will not merge everything right away.
Depending on the game and the argumentation, we may discuss it first in the Pull Request if this is a good fit.

Games can be removed from the list if a criterion is no longer met.

#### What is the theme of *technology*, *startups*, *business* or *software engineering* in this context?

When it deals with things like ...

- A simulator game of running a startup (with all the ups and downs): Mimicking the (funny) real world
- A game where the game flow is primarily to automate activities away (with some form of logic, code, etc.)

This list is incomplete and should not limit you from adding your proposal to the list.
Instead, it provides you with background on the thought process during the criteria creation.

### Format

:warning: **The main [`README.md`](/README.md) is a rendered version of the data. Do not edit it manually.**

To add a new game, please create a `.yml` file in the [`/games`](/games) directory like `/games/<game-name>.yml`. 
Feel free to check out a few other YAML files in that directory to see what it should look like.

| Field                        | Type               | Description                                                                                |
|------------------------------|--------------------|--------------------------------------------------------------------------------------------|
| name                         | string             | Name of the game                                                                           |
| website                      | string (url)       | Website (full URL) of the game                                                             |
| steamID                      | integer            | ID of the game in the [Steam Store](https://store.steampowered.com/)                       |
| repository                   | string (url)       | Git repository (Github, Gitlab, etc.) if the source code is publicly available               |
| programmable                 | boolean            | Can the game be controlled by code? (e.g., programmable interface for competitions)          |

Finally, create a pull request with all your changes. 

**Thanks for contributing!** :tada:
