
<h1 align="center"> Discord Web Soundboard </h1>

![Discord Web Soundboard](/doc/header_thumb.png)


## Introduction
Simply and quickly connect a Bot to your Discord server, and make him play sounds and musics of your choice by using a Web-interface. 
All your users will be able to play with sounds to emphasize funny discussion/gaming moments !


## Functionalities
* Integrated Authentication with Passport for Discord Account
* Personnalized dashboard with :
	* A global sounds library
	* Possibility to create and order user's categories
	* Drag & Drop sounds from library into your categories, remove and re-order them
* Upload functionality for MP3/OGG.
* Click on a sound, and the bot will join you in your channel (if you're online), and play your sound/music (Impossible to play the sound if you're not here)
* Shared button to stop a sound from the bot (Troll safety first :p)


## Built with
* [React.JS](https://reactjs.org/) : 16.13.1 : React.JS is used to provide a quick one-page web-application.
* [Node.JS](https://nodejs.org) : 13.13.0 : Node.JS provide the back-end functionnalities of the application.
* [Socket.IO](https://socket.io/) : 2.3.0 : Socket.IO allow us to easily and quickly share informations between clients and our application functionalities.
* [Discord.JS](https://discord.js.org) : 12.2.0 : Discord.JS is used for having a Bot playing required sound directly in your Discord server.
* [Express.JS](https://expressjs.com/) : 4.17.1 : Express.JS handles all our requests from the client
* [Formidable](https://www.npmjs.com/package/formidable) : 1.2.2 : Formidable is used to allow users to upload their own sounds
* [ffmpeg-static](https://www.npmjs.com/package/ffmpeg-static) : 4.2.2 : FFMPEG is used to convert **mp3** common files to **ogg** for optimisation
* [Dotenv](https://www.npmjs.com/package/dotenv) : 8.2.0 : Using Dotenv to manage environment variables from a .env file (eg: "password: process.env.DB_PASS")
* [Discordjs/opus](https://github.com/discordjs/opus) : 0.3.2 
* [PostgreSQL](https://www.postgresql.org/) : Database to store sounds and users dashboard
* [Knex.JS](http://knexjs.org/) : SQL query builder
* [Passport](http://www.passportjs.org/) : Using discord authentication with Passport OAuth
* [Sass](https://sass-lang.com/) : SASS preprocessor
* [React-dnd](https://react-dnd.github.io/react-dnd/about) : Drag & Drop for React
* [Redis](https://redis.io/) : In-memory data structure store, used as a database for Cache session (for authentication)


## Local Setup
First of all, ensure you install [Node.JS](https://nodejs.org), and [Git](https://git-scm.com/) if you plan to clone the repository

```bash
# Clone this repository
$ git clone https://github.com/Nadrielle/Discord-Web-Soundboard.git

# Go into the repository ***Server*** folder
$ cd discord-web-soundboard/server

# Install Server dependencies
$ npm install

# Go into the repository ***Client*** folder
$ cd ../client

# Install Client dependencies
$ npm install

# //TO DO : Install PostgreSQL Database and create an empty database

# //TO DO : Install a Redis database
```


You'll need to Create a .env file to store some variables as follow :
```
# Content of the .env file

NODE_ENV=development
PORT=5000
DISCORD_TOKEN=#MyDiscordToken (//TO DO : Add a wiki link to explain)
SOUNDS_FOLDER="./sounds/"
TEMP_FOLDER="../sounds/temp"
OAUTH_CLIENT_ID=#OAuth Client (//TO DO : Add a wiki link to explain)
OAUTH_CLIENT_SECRET=#OAuth Secret (//TO DO : Add a wiki link to explain)
OAUTH_CLIENT_REDIRECT=/auth/ (//TO DO : Add a wiki link to explain)

POSTGRES_DB=#Database name
POSTGRES_USER=#Database user
POSTGRES_PASSWORD=#Database user password
```

You will need one command prombt to run the server, and one to run the react application
Here are the basics, but you can use Nodemon in Development or Concurrently

```
# Run the server script in one command prombt, in the 'discord-web-soundboard/server' folder
$ node server.js
```

```
# Run the react app in another command prombt, in the 'discord-web-soundboard/client' folder
$ npm start
```
The app should launch in your web browser at **localhost:3000/**

![Server and App running](/doc/command_prombts.png)


***To Do : ***
* ***Create a wiki to explain how to create a bot on discord.com***


## Future features
* Discord Login whitelisting (The bot check if the user exist in your server)
* Visual display/history of users playing sounds (We want to know who is spamming this annoying sound =D)
* Auto-disconnect of the bot after a certain time
* Design refinement
* Global Code Refactoring
* More to come ...


## Support
To do


## License
MIT
