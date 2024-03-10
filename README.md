# Crypto Discord Bot

Tracks mentions of tickers for coins and looks up
their current price on [ethplorer](https://ethplorer.io).

## Configuration

Configure tickers like this ([config.yaml](./config.yaml))

```yaml
---
channels_with_coins:
- channel_id: 827974330666450957 # discord channel ID
  coins:
  - ticker: XD # the ticker to listen for in chat
    address: 0x3005003BDA885deE7c74182e5FE336e9E3Df87bB # the address of the coin
```

## Running

If you are not familiar with how discord bots work,
I recommend starting with [their docs](https://discord.com/developers/docs/intro).

Once you have the bot joined to a server, set up your config like above.
Add your token to the ENV.

```bash
export DISCORD_TOKEN=<my_token>
```

Then, build and start the bot.

```bash
make
out/crypto-discord-bot -path=./config.yaml
[INFO] Opening ./config.yaml...
[INFO] Loading ./config.yaml...
[INFO] ...Loaded!
[INFO] Connecting to discord...
[INFO] ...Connected!
[INFO] Listening for messages...
```
