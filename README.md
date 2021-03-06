[![Go](https://github.com/hihoak/auto-standup/actions/workflows/go.yaml/badge.svg)](https://github.com/hihoak/auto-standup/actions/workflows/go.yaml)
[![codecov](https://codecov.io/gh/hihoak/auto-standup/branch/main/graph/badge.svg?token=1G0M0FKKWS)](https://codecov.io/gh/hihoak/auto-standup)

# standup-generator
Automates your standup report

## Как пользоваться?
1. установить утилитку одним из удобных способов
    
* `go install github.com/hihoak/auto-standup`
* собрать локально, скрипт положит бинарь в $GOPATH/bin/auto-standup `make build_{arm|amd|win}`
* скачать бинарь с [последнего релиза](https://github.com/hihoak/auto-standup/releases/latest)

2. Настройка. Тут два варианта:
* Создать конфигурационный файл `.standup.yaml`. По умолчанию утилитка ищет его по пути `$HOME/.standup.yaml` или передавать флагом `--config-path "/path/to/my/config.yaml"`
```yaml
# Обязательные параметры
# Логин и пароль в Jira. Используется для получения введенных тобой тикетов
username: admin
password: Mysecuresuperadminpass!

# Необязательные параметры
# Параметр, который добавит в отчет estimated time каждого тикета и так же суммарное кол-во запланированного времени. Может быть подключен через флаг --estimated-time команды
# include_estimated_time: true # default false
# Параметр, который добавит в отчет сумму залогированного времени за прошедший рабочий день. Так же может быть добавлен подключен через флаг --log-time команды
# include_logged_time: true # default false
# Параметр, в котором задаются пользователи, чья активность засчитывается при авто-нахождении тикетов за прошедший рабочий день (указывай username через запятую без пробелов)
# eligible_users_histories: gitlab # default "gitlab,{.username}"
# Параметр, в котором можно указать проекты в Jira, тикеты из этих проектов не будут включены в отчет
# exclude_jira_projects: retest # default "retest"
```
* Передавать логин и пароль через флаги команды пример `auto-standup -u "admin" -p "admin"`

3. Запускаем и радуемся или не радуемся, если все развалилось. Вот примеры команд для твоего удобства:
* `auto-standup -t "RE-2000,RE-3000,RE-4000"`
Максимально автоматизируем составление отчета, с помощью флага `-t` перечисляем через запятую тикеты, которые планируем сделать, а тикеты, которые сделал будут определены автоматически. Логика определения тикета: достаются все тикеты за промежуток времени с текущего времени и за прошлый день для указанного пользователя, если запускаешь в выходные или в понедельник, то будут взяты тикеты с пятницы -> тикеты фильтруются по активности от пользователей, которых считаем валидными и по проектам в Jira (см. п.2 "Создание конфигурационного файла")
* `auto-standup -d "RE-1000,RE-1500" -t "RE-2000,RE-3000,RE-4000"`
Минимум автоматизации. С помощью флага `-d` перечисляем сделанные тикеты

4. Пример работы программы
```bash
auto-standup create --log-time --estimated-time -t "RE-1000,RE-2000,RETEST-528"
```
```text
**Что вы делали с прошлого опроса?**
* [RE-6862](https://jira.example.com/browse/RE-6862) - Сделать прикольную фичу [log: 4h]
* [RE-6760](https://jira.example.com/browse/RE-6760) - Продумать архитектурное решение [log: 2h]
* [RE-5977](https://jira.example.com/browse/RE-5977) -  Зачинить баг в программе [log: 30m]
*Суммарно залогировано времени: 6h 30m*
**Что вы будете делать до следующего опроса?**
* [RE-1000](https://jira.example.com/browse/RE-1000) - Проблема с выкачиванием image из registry [no estimate]
* [RE-2000](https://jira.example.com/browse/RE-2000) - Корневой базовый образ [no estimate]
* [RETEST-528](https://jira.example.comru/browse/RETEST-528) - for test [3h 20m]
*Суммарно запланировано времени: 3h 20m*
```

## Контакты и полезные ссылки

* Борда с тикетами, куда вы можете добавить фичу, чтобы ее реализовали или взять любой из тикетов и помочь развитию утилитки [![Trello board](https://upload.wikimedia.org/wikipedia/en/8/8c/Trello_logo.svg)](https://trello.com/b/OxH7R79n/auto-standup-board)

### Author
<table>
<tr>
  <td align="right"><a href="https://github.com/hihoak"><img src="https://github.com/hihoak.png" width="100px;" alt=""/><br /><sub><b>Mihaylov Artem</b></sub></a></td>
</tr>
</table>

[![Telegram url](https://icons.iconarchive.com/icons/alecive/flatwoken/48/Apps-Telegram-icon.png)](https://t.me/ez_buckets)

### Contributors
Пока, что совсем никого :(
