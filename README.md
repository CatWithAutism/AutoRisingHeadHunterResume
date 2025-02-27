# Автоматическое поднятие резюме на hh.ru

С помощью этой утилиты можно автоматически поднимать резюме на hh.ru каждые четыре часа. 

Бесплатная альтернатива для HH PRO

По желанию есть информирование о действиях в Telegram.

Проще всего запустить бота в докер контейнере

- TG_BOT_TOKEN - токен можно получить у BotFather @BotFather (по желанию)
- TG_ADMIN_ID - свой Telegram ID можно взять у @getmyid_bot (по желанию)
- HH_USERNAME - логин от hh.ru
- HH_PASSWORD - пароль от hh.ru
- HH_UPDATE_HOURS - раз во сколько часов обновляем резюме(рекомендуется минимум 4)

```
version: '3.8'

services:
  headhunter-auto-rising:
    image: catwithautism/headhunter-auto-rising:latest
    container_name: hh-auto-riser
    restart: unless-stopped
    environment:
      - TG_BOT_TOKEN=${TG_BOT_TOKEN}
      - TG_ADMIN_ID=${TG_ADMIN_ID}
      - HH_USERNAME=${HH_USERNAME}
      - HH_PASSWORD=${HH_PASSWORD}
      - HH_UPDATE_HOURS=${HH_UPDATE_HOURS}

```

Копируем в compose.yaml, заполняем поля

```bash
docker-compose up -d
```