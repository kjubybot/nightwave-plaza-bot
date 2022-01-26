# Nightwave Plaza Bot
This simple Discord bot allows to listen to [Nightwave Plaza Radio](https://plaza.one) in voice channels

## Running
You can compile and run it using Go compiler or Docker. Either way `DISCORD_TOKEN` environment variable should be set to Discord bot token

```bash
go build ./cmd/plaza
export DISCORD_TOKEN=token
./plaza
```
or
```bash
docker build -t plaza .
docker run -d -e DISCORD_TOKEN=token plaza
```
