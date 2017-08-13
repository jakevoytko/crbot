CRBot
-----

CRBot is a Discord bot that you can teach to respond to user-defined
commands. CRBot is short for "Call and Response Bot". Commands are taught using
`?learn <call> <response>`, and can retrieve the response with `?call`.

It works really well for media retrieval; the Discord unfurler will inline
images, videos, etc.

This was initially a weekend project to see if my friends would like it. Now
that they haven't banned the bot (or me), I'm in the process of turning this
into something that's maintainable.

Example
--------

```
Jake: ?help
crbot: Type ?help for this message, ?list to list all commands, or
       ?help <command> to get help for a particular command.

Jake: ?learn bearshrug ʅʕ•ᴥ•ʔʃ
crbot: Learned about bearshrug

One of my friends: What's the weather tomorrow?
Jake: ?bearshrug
crbot: ʅʕ•ᴥ•ʔʃ
```

Prerequisites
---------------

* Redis instance, standard port, no password, DB 0
* A bazel build, running at least 0.5.3
* Add `secret.json` with single key, `bot_token`

Running
--------

`go run *.go`

You have to run this with a bot account. You can register a bot account for
free [at the Discord site](https://discordapp.com/developers/docs/intro). Take
the bot ID that you are given, and
visit
[https://discordapp.com/oauth2/authorize?&client_id={$bot_id}](https://discordapp.com/oauth2/authorize?&client_id={bot-id}) in
your browser, where {$bot_id} is replaced with the bot ID that you are given
from the Discord developers site.

Before sending PR
-------------------

`bazel run :go_default_test` from working directory.

If your PR adds features, please add tests in system_test.go that cover your new use cases.

More reading
-------------

**Blog posts**:
- [Writing a Discord bot, and techniques for writing effective small programs](https://www.bitlog.com/index.php/2017/03/31/techniques-for-effectively-growing-small-programs/)
- [My friends trolled each other with my Discord bot, and how we fixed it](https://www.bitlog.com/index.php/2017/05/31/my-friends-trolled-each-other-with-my-discord-bot-and-how-we-fixed-it/)