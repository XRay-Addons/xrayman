# XrayMan

Набор сервисов для развёртывания и управления собственной инфраструктурой прокси-серверов на базе XRay-Core.

## Зачем нужен XRAYMAN

Самый простой способ поднять прокси — взять VPS, установить на неё XRay, сгенерировать файл-подписку (конфигурацию клиента) и раздать его тем, кому нужен доступ (друзьям, знакомым, друзьям друзей).

Такой подход быстро перестаёт быть удобным, как только ситуация усложняется:

- Ноды нестабильны — сервер может быть заблокирован, временно недоступен, перезагружен или упасть, и об этом никто не узнает, пока пользователи не пожалуются.
- Нужен выбор нод — некоторые IP-адреса могут быть забанены в отдельных сервисах, поэтому полезно иметь несколько нод и уметь переключаться между ними.
- Нужен мониторинг пользователей — хочется видеть, кто сколько трафика потребляет (сами пользователи остаются анонимными при этом).
- Нужна автономность — сервис должен следить за собой сам: перезапускать упавшие компоненты, исключать недоступные ноды из выдачи, и падение одной ноды не должно приводить к отказу всей системы (отказ менеджера нод или базы данных — отдельный случай, тут ничего не поделать).

XRayMan решает эти задачи, разделяя систему на независимые компоненты: лёгкие ноды-исполнители и центральный менеджер, который следит за ними, управляет пользователями и отдаёт клиентам актуальные подписки.

## Архитектура

Система состоит из двух типов сервисов и базы данных:

- **xray-node** — тонкая обёртка над XRay-Core. Не хранит состояние, управляется только через REST API.

- **xray-nodeman** — центральный менеджер: управляет нодами и пользователями, мониторит доступность нод, отдаёт клиентам подписки, предоставляет админ-панель и метрики. Хранит своё состояние в PostgreSQL.

Соединение между компонентами одностороннее: **nodeman обращается к node**, обратного взаимодействия нет.

```
Клиент (v2rayNG / HAPP / …)
        │
        │  запрос подписки (HTTPS)
        ▼
   xray-nodeman ───────────────► PostgreSQL
        │
        │  REST API (управление, телеметрия, health-check)
        │  соединение защищено cert pinning + JWT
        ▼
   ┌───────────┬───────────┬───────────┐
   │ xray-node │ xray-node │ xray-node │  ...
   └───────────┴───────────┴───────────┘
```

```mermaid
flowchart TD
    Client["Клиент (v2rayNG / HAPP / …)"] -->|"запрос подписки (HTTPS)"| Nodeman["xray-nodeman"]
    Nodeman --> DB[("PostgreSQL")]
    Nodeman -->|"REST API: управление, телеметрия, health-check\n(cert pinning + JWT)"| Node1["xray-node #1"]
    Nodeman --> Node2["xray-node #2"]
    Nodeman --> Node3["xray-node #N"]
```

### Логика работы:

1. Клиент запрашивает у **xray-nodeman** страницу подписки.

2. **xray-nodeman** формирует список доступных на данный момент нод
   и персональные конфигурации пользователя для каждой из них.

3. **xray-nodeman** периодически опрашивает каждую **xray-node**: проверяет доступность, снимает телеметрию (трафик, CPU, RAM), при необходимости перезапускает или временно исключает ноду из выдачи.

4. Пользователи и ноды администрируются через веб-интерфейс **xray-nodeman**.

## Компоненты

### xray-node

Сервис-обёртка над XRay-Core. Не хранит собственное состояние — вся конфигурация приходит от xray-nodeman через REST API. Полная OpenAPI-спецификация есть в репозитории; здесь — только список возможностей:

- запуск XRay с заданным списком пользователей;
- остановка XRay;
- проверка статуса — запущен процесс или нет;
- добавление и удаление пользователей «на лету», без перезапуска;
- отдача телеметрии: трафик по пользователям, загрузка CPU и RAM.

### xray-nodeman

Центральный менеджер нод и пользователей. Периодически опрашивает каждую ноду, чтобы получить телеметрию и проверить доступность — сама нода к nodeman не обращается.

Возможности:

- добавление, удаление, включение и отключение пользователей и нод;

- формирование страницы подписок — эндпоинта, который опрашивают клиенты, чтобы получить список доступных им нод и конфигурации для каждой из них;

- мониторинг состояния нод: недоступные исключаются из выдачи, упавшие — пытается перезапустить;

- просмотр текущего состояния нод и пользователей;

- фронтенд для пользователей — страница, на которой можно создать себе учётную запись и получить инструкцию по установке и настройке клиента XRay;

- админ-панель — список пользователей и нод, их состояние, возможность управления;

- эндпоинт с метриками для экспорта в Prometheus / Grafana;

- хранение состояния в PostgreSQL. Схема базы и миграции (на [goose](https://github.com/pressly/goose)) встроены в само приложение и накатываются автоматически при старте — администратору достаточно один раз указать connection string к пустой базе данных.

## Формат подписки

xray-nodeman отдаёт подписку по субскрипшн-ссылке в формате JSON, содержащей список доступных пользователю нод с их конфигурациями.

Формат совместим с клиентом [HAPP](https://github.com/happ-proxy) — достаточно указать URL подписки в приложении, и оно само подтянет актуальный список нод.

**Happ** используется потому, что он кроссплатформенный, простой для пользователя и позволяет передавать в подписке кастомный роутинг.

## Безопасность

### Соединение nodeman ↔ node

У ноды нет настоящего домена и валидного сертификата, поэтому соединение защищено иначе:

- у каждой ноды есть собственный ключ, публичная часть которого (хэш публичного ключа сертификата) добавляется в конфигурацию nodeman — это **cert pinning**: nodeman принимает соединение с нодой, только если её сертификат совпадает с закреплённым

- дополнительно у пары node/nodeman есть общий secret, которым подписываются **JWT**-токены — по ним nodeman авторизуется перед нодой при каждом запросе к её REST API.

Инициатором соединения всегда выступает nodeman.

### Админ-панель

Админ-панель защищена паролем, но пароль сам по себе не защищает от перехвата трафика. Поэтому админ-панель нужно размещать **за reverse-proxy (nginx или Caddy) с TLS-сертификатом на реальный домен** — не открывайте её напрямую по HTTP.

## Технологии

- **Go** — xray-node и xray-nodeman
- **TypeScript** — фронтенды
  - **Astro** — пользовательская страница
  - **Astro + Vue** — админ-панель.

## Установка

Минимальная рабочая конфигурация состоит из трёх частей:

1. **xray-nodeman** — systemd-сервис или Docker-контейнер;

2. хотя бы одна **xray-node** — systemd-сервис (в Docker не поставляется и не имеет смысла);

3. **PostgreSQL** — локально рядом с nodeman, в Docker или в стороннем хостинге; от администратора требуется только connection string к пустой базе данных.

### xray-node

#### Установка из .deb (Ubuntu, Debian и т.п.)

```sh
# скачать .deb под нужную архитектуру (amd64 или arm64)
curl -fsSL "https://github.com/XRay-Addons/xrayman/releases/latest/download/xray-node-amd64.deb" -o xray-node.deb
# установить
sudo apt install ./xray-node.deb
# посмотреть инструкцию по настройке
sudo apt show xray-node
# запустить
sudo systemctl daemon-reload
sudo systemctl enable --now xray-node
# или перезапустить
sudo systemctl daemon-reload
sudo systemctl restart xray-node
# посмотреть логи
sudo journalctl -u xray-node -n 50 -f
```

#### Сборка вручную

Требования:

- go v1.26.2

```sh

# скачать исходный код
git clone https://github.com/XRay-Addons/xrayman.git
cd xrayman
git submodule update --init --recursive
# установить инструменты сборки
make tools
# собрать
make build_node
# проверить
./build/xray-node/xray-node -h
```

### xray-nodeman

#### Предварительные требования

Нужен PostgreSQL — как отдельное приложение, Docker-образ или сторонний хостинг. Единственное, что требуется передать nodeman — **connection string**:

```
postgresql://username:password@host:port/database_name
```

#### Установка из .deb (Ubuntu, Debian и т.п.)

```sh
# скачать .deb под нужную архитектуру (amd64 или arm64)
curl -fsSL "https://github.com/XRay-Addons/xrayman/releases/latest/download/xray-nodeman-amd64.deb" -o xray-nodeman.deb
# установить
sudo apt install ./xray-nodeman.deb
# посмотреть инструкцию по настройке
sudo apt show xray-nodeman
# запустить
sudo systemctl daemon-reload
sudo systemctl enable --now xray-nodeman
# или перезапустить
sudo systemctl daemon-reload
sudo systemctl restart xray-nodeman
# посмотреть логи
sudo journalctl -u xray-nodeman -n 50 -f
```

После успешного запуска доступны следующие эндпоинты:

- `${ENDPOINT}/u` — пользовательская страница;
- `${ENDPOINT}/adm` — админ-панель;
- `${ENDPOINT}/api` — защищённый авторизацией API;
- `${ENDPOINT}/api/version` — незащищённый статусный эндпоинт;
- `${METRICS_ENDPOINT}/metrics` — метрики для Prometheus.

#### Docker

##### Standalone

Контейнер только с xray-nodeman, ничего лишнего. PostgreSQL и connection string к нему нужно предоставить отдельно (см. выше).

```sh
# скачать docker-compose
curl -fsSL "https://github.com/XRay-Addons/xrayman/releases/latest/download/docker-compose.standalone.tar.gz" | tar -xzf -
# создать .env файл
sudo nano .env
# посмотреть инструкцию по настройке (добавлять переменные в .env, пока не заработает)
docker compose run --rm xray-nodeman --help
# запустить
docker compose up -d
# или перезапустить
docker compose down && docker compose up -d
# посмотреть логи
docker compose logs -n 50 -f xray-nodeman
```

##### All-in-one

Контейнер включает всё необходимое сразу:

- PostgreSQL;
- бэкап PostgreSQL;
- xray-nodeman;
- Prometheus;
- Grafana.

```sh
# скачать docker-compose
curl -fsSL "https://github.com/XRay-Addons/xrayman/releases/latest/download/docker-compose.all-in-one.tar.gz" | tar -xzf -
# создать .env файл
sudo nano .env
# посмотреть инструкцию по настройке (добавлять переменные в .env, пока не заработает)
docker compose run --rm xray-nodeman --help
# запустить
docker compose up -d
# или перезапустить
docker compose down && docker compose up -d
# посмотреть логи
docker compose logs -n 50 -f xray-nodeman
```

#### Сборка вручную

Требования:

- go v1.26.2
- pnpm v10.33.2
- node.js v24.4.0

```sh
# скачать исходный код
git clone https://github.com/XRay-Addons/xrayman.git
cd xrayman
git submodule update --init --recursive
# установить инструменты сборки
make tools
# собрать
make build_nodeman
# проверить
./build/xray-nodeman/xray-nodeman -h
```

### Удобный просмотр логов

Оба сервиса пишут структурированные (zap) логи в JSON. Чтобы читать их в человекочитаемом виде, удобно использовать `jq`.

#### Установка jq

Инструкция по установке: https://jqlang.org/download/

#### Alias для форматирования логов

Откройте файл настроек шелла (`~/.zshrc` на macOS, `~/.bashrc` или аналог на Linux):

```sh
sudo nano ~/.zshrc
```

Добавьте алиас `zl` для форматирования zap-логов:

```sh
# >>> zl - jq setup for go zap log formatting >>>
zl() {
jq -Rr '
def level_color:
    if . == "error" then "\u001b[31m"
    elif . == "warn" then "\u001b[33m"
    elif . == "info" then "\u001b[32m"
    elif . == "debug" then "\u001b[36m"
    else "\u001b[0m"
    end;

    . as $line
    | (try ($line | fromjson) catch null) as $log
    | if ($log | type) == "object" then
      "\($log.ts) \($log.level | level_color)\($log.level | ascii_upcase)\u001b[0m \($log.msg)",
      ($log
        | to_entries[]
        | select(.key != "ts" and .key != "level" and .key != "msg")
        | if (.value | type) == "string" and (.value | contains("\n")) then
            "- \u001b[36m\(.key)\u001b[0m:\n\(.value | split("\n") | map("    " + .) | join("\n"))"
          else
            "- \u001b[36m\(.key)\u001b[0m: \(.value)"
          end
      )
    else
      $line
    end

'
}
# <<< zl setup done <
```

Применить изменения:

```sh
source ~/.zshrc
```

#### Использование алиаса `zl`

Добавьте флаг для вывода «сырых» логов без обёртки (например, `-o cat` для journald, `--no-log-prefix` для docker-compose), затем допишите `| zl` к любой команде, выводящей логи:

```sh
docker compose logs -n 50 -f --no-log-prefix xray-nodeman | zl
# или
docker logs -n 50 -f xray-nodeman | zl
```
