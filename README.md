# Kbot
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/matvrus/kbot/blob/main/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/matvrus/kbot)](https://goreportcard.com/report/github.com/matvrus/kbot)
[![Release](https://img.shields.io/github/release/matvrus/kbot.svg?style=flat-square)](https://github.com/matvrus/kbot/releases/latest)

🤖 
KBot - Telegram-бот на мові програмування Golang
KBot є Telegram-ботом, написаним на мові програмування Golang, який може відповідати на текстові повідомлення та згодом буде виконувати різноманітні команди.
## Посилання на бот

 ✉️ https://t.me/matvrus_bot
## Встановленя


Щоб встановити KBot, вам потрібно спочатку встановити Golang та налаштувати своє середовище розробки. Якщо ви використовуєте Codespaces, то все необхідне вже налаштовано.

```bash
 git clone https://github.com/matvrus/kbot
```
перейдіть у папку з ботом

```bash
 cd kbot
```
виконайте

    go get

потім

    go build -ldflags "-X="github.com/matvrus/kbot/cmd.appVersion=v1.0.1


після цього можете запускати бот
     
    ./kbot start

тепер можна взаємодіяти з ботом

нижче наведено приклади текстових повідомлень, які Kbot розуміє:

    /start hello 
    - Kbot відповість з "Hello, I'm Kbot <Ім'я відправника та версія>!".

Імпортуйте необхідні бібліотеки та встановіть бібліотеку gopkg.in/telebot.v3 за допомогою go get. Отримайте токен бота зі змінної середовища.
## свій бот
 `TELE_TOKEN`. Ви можете отримати свій токен від [@BotFather](https://t.me/BotFather).


## Usage/Examples

🤖To run Kbot, set the TELE_TOKEN environment variable to your Telegram bot token, and then run the following command

./kbot start

## 📝Ліцензія

[MIT](https://choosealicense.com/licenses/mit/)

