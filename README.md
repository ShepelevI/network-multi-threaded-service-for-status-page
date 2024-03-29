# Сетевой многопоточный сервис для StatusPage.

Данный сервис принимает запросы по сети (используется симулятор данных для эмитации получения данных из сети) и возвращает данные о состоянии систем.
Выводятся данные о состоянии работы систем: SMS, MMS, Email, Voice Call, Billing, Support, Истории инцидентов.
При каждом запуске симулятор генерирует новые данные, часть из которых доступна в виде файлов (сервис ищет их в директории симулятора и сканирует), а часть данных доступна через API.
Результаты выводятся на web страницу сайта под названием StatusPage.

Использовались следующие библиотеки:

1. https://github.com/biter777/countries - для работы с кодами стран в формате ISO 3166-1 alpha-2;
2. https://github.com/go-chi/chi - в роли маршрутизатора;
3. https://github.com/gorilla/mux - в роли маршрутизатора и диспетчера запросов для сопоставления входящих запросов с их соответствующим обработчиком.

Для запуска сначала перейдите в каталог симулятора и запустите симулятор с помощью команды:
```shell
cd .\simulator\
go run main.go
```
Затем запустите основной сервис:
```shell
go run .\cmd\main.go
```

Страница результатов в web/index.html. Откройте, чтобы увидеть результат.
