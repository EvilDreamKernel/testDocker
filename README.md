# Тестовое задание

**Основная задача:** <br />
Разработать Java Service, который принимает файл при помощи cURL утилиты. Также нужно разработать sidecar Golang сервис, который отдаёт какие-либо метрики в Prometheus формате. Метрики должны отражать полезные данные для мониторинга основной активности сервиса.

## GetOverHere (Java-app)
----
Java-приложение, которое принимает файл в POST запросе, в виде multipart/form-data. Файл ожидается в параметре **file**, при этом URL Path __/upload__. <br />
Пример работы с приложением через cURL:
```
curl http://localhost:4567/upload -XPOST --form file=@file.txt
```

**Пример ответа (JSON format):**
```
{"message":"File uploaded successfully"}
```
При загрузке файла создается директория **files** рядом с запущенным приложением. В неё помещается файл, загруженный через POST запрос.

### Как работает приложение
----
Приложение слушает порт 4567 (Порт по-умолчанию, который устанавливается фреймворком Sparkjava). Для обслуживания POST запросов использовался framework Sparkjava, при этом не использовались стандартные функции пакетов com.sun, т.к. Oracle не рекомендовала их к использованию.<br />
Приложение обрабатывает POST запросы к URL Path /upload, также есть ограничения для объемов данных, загружаемых таким способом. Ограничения описаны в коде. Можно было бы вынести параметры в отдельный конфигурационный файл, но я решил что в рамках тестового задания и это будет нормально.

### Сборка Java приложения
----
Для сборки использовался Maven. В файле pom.xml определены зависимости, плагины, которые в дальнейшем понадобятся для сборки. В директорию target выводятся транслированные файлы классов, а также готовый jar архив. Вместе с ним добавляются и библиотеки описанные в pom.xml (рядом в директорию libs)

### Docker `GetOverHere`
----
Для того, чтобы уменьшить размер конечного образа Docker было принято решение воспользоваться multi-stage сборкой образа, а именно: была сделана отдельная стадия для сборки, конечные артефакты из которого передаются в следующий docker-образ. <br />
После чего приложение начинает слушать порт 4567 и отдаёт его сервисам внутри своей сети. <br />

## Blaze (Golang-App)
----
Golang-приложение, которое проверяет директорию, в которую Java-приложение загружает файлы, после чего отдаёт определенные метрики для анализа и визуализации сторонними сервисами (Prometheus, Grafana, etc.). <br />
**Пример получения метрик:**
```
curl http://127.0.0.1:9145/metrics
```

**Пример ответа:**
```
# HELP directory_elements_number Elements number in directory
# TYPE directory_elements_number gauge
directory_elements_number{path="files/"} 2
# HELP file_age_unix Last modification time of file
# TYPE file_age_unix gauge
file_age_unix{path="files/file2.txt"} 1.639341396e+09
file_age_unix{path="files/file3.txt"} 1.639343649e+09
# HELP file_size_bytes File size in bytes
# TYPE file_size_bytes gauge
file_size_bytes{path="files/file2.txt"} 12351
file_size_bytes{path="files/file3.txt"} 12351
```
Если директория отсутствует - то никаких метрик возвращено не будет при запросе.
Как только директория появится - то появится одна из метрик:
```
# HELP directory_elements_number Elements number in directory
# TYPE directory_elements_number gauge
directory_elements_number{path="files/"} 0
```

### Возвращаемые метрики
---
**directory_elements_number** - количество элементов в прослушиваемой директории. (Файлов) <br />
**file_age_unix** - Время последней модификации файла. Т.к. сложно было решить, что есть новый файл, а изменение файла вручную не подразумевается, значит время последней модификации файла можно считать временем создания/загрузки файла. <br />
**file_size_bytes** - Размер каждого из файлов в байтах.

### Конфигурация
---
Данное приложение можно конфигурировать, передавая параметры запуска в аргументах shell'а. <br />
**Пример запуска со всеми аргументами:**
```
./Blaze --fs.monitor-directory "files/" --web.listen-address ":9145" --web.metrics-path "/metrics"
```
**Описание параметров:** <br />
_fs.monitor-directory_: Полный или относительный путь до директории, которую нужно мониторить. <br />
_web.listen-address_: Адрес с портом (или просто порт), который необходимо слушать, для приема запроса метрик <br />
_web.metrics-path_: URL Path, который будет слушаться, чтобы отдавать метрики по запросу.

### Сборка Golang приложения.
---
Сборка приложения происходит при помощи стандартной команды:
```
go build
```
Зависимости при этом описаны в go.mod файле, а сами пакеты находятся в директории vendor.

### Docker Blaze
---
В целом можно было реализовать go build также через multi-stage в докере, но в данном случае в этом не было сильной необходимости, т.к. размер образа не сильно подрос. <br />
Конфигурационные параметры передаются в виде переменных окружений, которые также можно задавать снаружи:
```
ENV LISTEN_ADDRESS=":9145"
ENV METRICS_PATH="/metrics"
ENV MONITOR_DIRECTORY="files/"
```

## docker-compose
---
Для быстрого запуска данных утилит использовался docker-compose, которые может определять параметры для запуска blaze-app, определяет общий volume, который разделяли оба приложения, а также порты, которые expose'ились наружу.
```
version: "3.9"
services:
  blaze-app:
    build:
      context: ./Blaze
    ports:
      - "9145:9145"
    volumes:
      - "./files:/app/files"
  
  java-app:
    build:
      context: ./GetOverHere
    ports:
      - "4567:4567"
    volumes:
      - "./files:/app/files"
```
Также из кода видно, что предварительно образы собираются из dockerfile, находящихся в директориях, определенных в параметрах context.
