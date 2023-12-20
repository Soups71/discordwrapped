# Discord Wrapped - A Discord Spellchecker Bot

## Introduction

Discord Wrapped, a bespoke Discord bot, is your ticket to a year-in-review extravaganza, akin to the beloved Spotify Wrapped. This ingenious creation is tailored to offer a comprehensive summary of your Discord server's activity throughout the year.

Unveil the secrets with two simple commands:

    * !wrapped server
    * !wrapped channel

Embark on a journey of insights, whether it's for the entire server or a specific channel:

    * Discover the number of messages sent by each user 📬
    * Unearth the GIF mastery with the count of GIFs shared by each user 🎥
    * Witness the visual tales with the number of images sent by each user 📸

Join the celebration of your server's vibrant moments, meticulously curated by Soups71. Let Discord Wrapped weave the narrative of your year in the Discordverse! 🚀


## Getting Started

To start using Discord Wrapped, follow these steps:

1. Follow the instructions in [installation instructions](INSTALL.md)

2. Download the repository code:
    * `git clone https://github.com/Soups71/discordwrapped.git && cd discordwrapped`

3. Create a file named `config.json` in the root directory of the project. This file will store your bot's secret token. The format of `config.json` should be as follows:

```json
{
    "Token": "MTE***nds",
    "BotPrefix": "!wrapped",
    "DBConn": "mongodb://localhost:27017"
}
```

4. Add Discord Wrapped to your Discord server. For proper functionality, it is recommended to grant the bot Administrator privileges on the server. While restricting permissions might work, it could lead to unexpected behavior.

5. Start the server:
    * `go run cmd/wrapped/main.go <logfile>`



## Current Tasks

The following tasks are in progress or need attention for future improvements:

* Determine the minimal permissions required for the bot's smooth operation.
* Implement thorough error handling, as the current version lacks adequate error management.
* Create a systemd service file and establish a streamlined installation process to automate the setup and execution of the bot.
