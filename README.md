# WhatsApp AutoJoin Bot

**Premise:** My Uni uses WhatsApp for most unofficial communication during coronavirus "study from home". Naturally students share homework/assignments/quiz answers with each other. WhatsApp has a limit of 256 members so you need to join quiz and homework share groups fast before they get filled up. And since I'm not active on WhatsApp, I miss out. Hence this bot.

There is a whitelist with strings that should be in a group title to join them in `handler.go`

Pardon the shoddy code but this was quickly thrown together and I may have ignored best practices for faster dev time.

**Usage:**

- Edit your group whitelist, compile and deploy on a server/raspberry pi with an init service like systemd. (set restart to on and workingdirectory to dir of binary, copy templates folder there too)
- Then visit `http://server-ip:8755`. Note that there is no auth since I'm running this on a local server at home.
- Route `/login` will give you a QR code that you need to scan in order to login with your whatsapp.
- Logs are at `/logs`
- `/reset` in case the bot starts misbehaving. You will have to logout again. (and restart manually incase you havent set that in systemd service file)
- If you use whatsapp web, this bot will stop working temporarily until you quit whatsapp web and it reconnects.
