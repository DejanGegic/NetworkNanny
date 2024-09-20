# 📑 Table of Contents
**Go to [QUICK START GUIDE](#)**

- [❓ What is this?](#-what-is-this)
- [🌟 Features](#-features)
- [🏃‍♂️ Quick-start guide](#️-quick-start-guide)
- [🏔️ Environment parameters guide](#️-environment-parameters-guide)
- [🛠️ Build it yourself](#️-build-it-yourself)
- [📈🪲 Future improvements, known bugs, etc.](#-future-improvements-known-bugs-etc)

# ❓ What is this?

It's a configurable, performant, and stable file server with a built-in rate limiter, written in Go.

⬆️[Go back to top](#📑-table-of-contents)

# 🌟 Features

## ⚙️ Configurability

**Only 1 file config**  
Pretty much all parameters can be configured in the `.env` file.
`.env` contains ample comments that'll get you up and running, but we'll dive into details later in this document.

## 🚀 Performance

### Sub 1ms response time

On my personal machine, times of under 100µs (0.1ms) are not uncommon.
This is only true when using **BadgerDB**.  

For **Redis** (I imagine keysdb should be a bit faster) response time is usually around 1ms, and never over 2.5ms. Which I recon is because there's another layer of network latency between this server and the database.

All in all, very performant.

## 🗄️ Multiple database options

**Support for BadgerDB, Redis, and KeysDB**  
For now, this project supports Redis or KeysDB for in-memory storage, or BadgerDB as an integrated on-disk database.
Both can be configured in the `.env` file.

## 💻🖥️ Multi-platform support

This project comes pre-compiled for **Windows** and **Linux**. As well as **ARM** and **x86** architectures for both **64** and **32** bit systems.  
Likewise, Mac binaries also exist but due to my lack of access to a machine running macOS, feel free to test yourself and report back. In theory, there shouldn't be any issues.  

Basically, it will run on almost anything.  
<details>
<summary> Additional binaries </summary>
If you need support for niche systems like Solaris, OpenBSD, FreeBSD, Android, etc., feel free to open an issue, or compile it yourself from source code.
</details>

## 🍼 Simple to use

Drop in a binary with the .env file and you're good to go.  
Need to change settings? Just edit the `.env` and restart the server.  

Want to use it standalone?  
With another server pointing to it?  
Integrate it with your existing projects, from custom APIs to WordPress?  
It Just works!

## ➕ Extendable

If you need to integrate it with another DB that's not supported yet, everything's set up for you to code it up. Just implement `DbInterface` form `db/interface.go` with your favorite DB.

## 🐹 Written in Go

What's there not to love?

⬆️[Go back to top](#📑-table-of-contents)

# 🏃‍♂️ Quick-start guide

1. Download the binary and `.env.temp` for your system from "Resources" on GitHub
2. Rename `.env.temp` to `.env` and place it in the same directory as the binary
3. Create a `./public` directory and place all static files in it
4. Change the parameters in `.env` as you please, and start the server by running the binary
5. You're done! Enjoy your rate-limited, ultra-performant file server 😊

⬆️[Go back to top](#📑-table-of-contents)

# 🏔️ Environment parameters guide

| Parameter         | Description                                                                              | Default   | Constraints       |
| ----------------- | ---------------------------------------------------------------------------------------- | --------- | ----------------- |
| PORT              | Port on which the server will listen                                                     | 8000      | int               |
| PUBLIC_DIR        | Path to the directory where static files are stored                                      | /public   | string            |
| ALLOW_BROWSING    | Allow users to browse files without having the link to the specific file.                | false     | "true"            |
| DB_TYPE           | Database type. "redis" is also used for "keysdb" and other compatible DBs                | "badger"  | "badger", "redis" |
| WINDOW            | Time window in SECONDS in which requests will be counted                                 | 600       | int               |
| REQUEST_LIMIT     | Permitted number of requests during a window                                             | 100       | int               |
| LIMITER_NAME      | Name of the limiter, mostly used for logging, or when there are multiple limiters active | 10minutes | string            |
| PERMABAN_TRESHOLD | Like "REQUEST_LIMIT", but only applies to permaban.                                      | 10        | int               |
| PERMABAN_TIME     | How long Permaban will last, in MINUTES. Usually way longer than WINDOW                  | 1440      | int               |
| DB_LOCATION       | Location of your badger database. Only applicable if your database type is `badger`      | badger    | string            |
| REDIS_HOST        | Host address of Redis or Redis-compatible DB.                                            | localhost | string            |
| REDIS_PORT        | Port of Redis or Redis-compatible DB                                                     | 6379      | string            |
| REDIS_PASSWORD    | Password of Redis or Redis-compatible DB.                                                | ""        | string            |

⬆️[Go back to top](#📑-table-of-contents)

# 🛠️ Build it yourself

## For current platform

Download the source code and run `go mod tidy && go build -o server main.go` from the project root.
That's pretty much it. You've built a binary for your current platform.

## For all/different platforms

There's a shell script `build.sh` that will build binary for all supported platforms automatically.
If you wish to use it, simply run `./build.sh` on a linux machine.

### Adding new platforms

If there's a need to compile for additional platforms automatically, simply add it to the `platforms` variable in `build.sh`.
Please do your due diligence and test it.

⬆️[Go back to top](#📑-table-of-contents)

# 📈🪲 Future improvements, known bugs, etc.

## Upcoming features

*The features listed will be added only if interest for this project starts developing, at least slowly.*  

- Official Docker image

## Bugs

No bugs found. If you discover any unusual behavior, please open an issue.
