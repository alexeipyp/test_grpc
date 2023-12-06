# Конфигурационный файл - формат всех полей


## Содержание
- [Секция параметра окружения ("env")](#env)
- [Секция конфигурации планировщика ("scheduler")](#scheduler)
    - [Размер очереди ("queue_size")](#queue_size)
    - [Размер пула обработчиков ("worker_pool_size")](#worker_pool_size)
- [Секция конфигурации кэша ("cache")](#cache)
    - [Время жизни кэша ("lifetime")](#lifetime)
- [Секция конфигурации подлючения к внешнему HTTP серверу ("external_http_connection")](#external_http_connection)
    - [Адрес сервера ("host")](#http_host)
    - [Порт сервера ("port")](#http_port)
    - [Логин для аутентификации ("username")](#username)
    - [Пароль для аутентификации ("password")](#password)
    - [Время ожидания ответа от сервера ("timeout")](#timeout)
- [Секция конфигурации gRPC сервера ("grpc")](#grpc)
    - [Адрес сервера ("host")](#grpchost)
    - [Порт сервера ("port")](#grpcport)
- [Секция конфигурации лог-файлов ("log")](#log)
    - [Имя лог-файла: уровни warning, error ("log_filename")](#log_filename)
    - [Имя лог-файла с отладочной информацией ("debug_log_filename")](#debug_log_filename)
    - [Имя лог-файла вывода всего трафика по gRPC ("grpc_trace_log_filename")](#grpc_trace_log_filename)
    - [Имя лог-файла вывода всего трафика по HTTP ("http_trace_log_filename")](#http_trace_log_filename)

<div id="env"></a>

## Секция параметра окружения ("env")

Может принимать одно из двух строковых значений:
- "dev"
- "prod"



<div id="scheduler"></a>

## Секция конфигурации планировщика ("scheduler")


<div id="queue_size"></a>

### Размер очереди ("queue_size")

Является целым, положительным числом.


<div id="worker_pool_size"></a>

### Размер пула обработчиков ("worker_pool_size")

Является целым, положительным числом.



<div id="cache"></a>

## Секция конфигурации кэша ("cache")


<div id="lifetime"></a>

### Время жизни кэша ("lifetime")

Является строковым значением, где:
- Строка заканчивается одним из суффиксов:
    - "ns", обозначая наносекунды
    - "ms", обозначая милисекунды
    - "s", обозначая секунды
    - "m", обозначая минуты
    - "h", обозначая часы
- Суффиксу предшествует числовое значение количества временных единиц, обозначаемых суффиксом



<div id="external_http_connection"></a>

## Секция конфигурации подлючения к внешнему HTTP серверу ("external_http_connection")


<div id="http_host"></a>

### Адрес сервера ("host")

Является строковым непустым значением


<div id="http_port"></a>

### Порт сервера ("port")

Является целым, положительным числом


<div id="username"></a>

### Логин для аутентификации ("username")

Является строковым непустым значением


<div id="password"></a>

### Пароль для аутентификации ("password")

Является строковым значением


<div id="timeout"></a>

### Время ожидания ответа от сервера ("timeout")

Является строковым значением, где:
- Строка заканчивается одним из суффиксов:
    - "ns", обозначая наносекунды
    - "ms", обозначая милисекунды
    - "s", обозначая секунды
    - "m", обозначая минуты
    - "h", обозначая часы
- Суффиксу предшествует числовое значение количества временных единиц, обозначаемых суффиксом


<div id="grpc"></a>

## Секция конфигурации gRPC сервера ("grpc")


<div id="grpchost"></a>

### Адрес сервера ("host")

Является строковым непустым значением

<div id="grpcpost"></a>

### Порт сервера ("port")

Является целым, положительным числом

<div id="log"></a>

## Секция конфигурации лог-файлов ("log")


<div id="log_filename"></a>

### Имя лог-файла: уровни warning, error ("log_filename")

Является строковым непустым значением


<div id="debug_log_filename"></a>

### Имя лог-файла с отладочной информацией ("debug_log_filename")

Является строковым непустым значением


<div id="grpc_trace_log_filename"></a>

### Имя лог-файла вывода всего трафика по gRPC ("grpc_trace_log_filename")

Является строковым непустым значением


<div id="http_trace_log_filename"></a>

### Имя лог-файла вывода всего трафика по HTTP ("http_trace_log_filename")

Является строковым непустым значением