# kp2imdb

[English](#english) | [Русский](#russian)

<a name="english"></a>

## Migrate your movie ratings from Kinopoisk to IMDb

kp2imdb is a Go utility that helps you transfer your movie ratings from [Kinopoisk](https://kinopoisk.ru) to [IMDb](https://imdb.com).

### Getting Started

#### Prerequisites

- Go 1.24 or later - [download](https://go.dev/doc/install)
- OMDb API key (get one at [omdbapi.com](https://www.omdbapi.com/apikey.aspx))
- IMDb account with valid session cookie

#### Installation

1. Clone the repository:
   ```
   git clone https://github.com/OutOfStack/kp2imdb.git
   cd kp2imdb
   ```

2. Configure the application:
    - Copy your IMDb session cookie (instructions below)
    - Get an OMDb API key
    - Edit `config.json`:
      ```json
      {
        "locale": "en",
        "omdb_api_key": "your_omdb_api_key",
        "imdb_cookie": "your_imdb_cookie"
      }
      ```

3. Prepare your movie data:
    - Export your Kinopoisk ratings using the provided script
    - Save the exported data as `data.json`

#### How to get your IMDb Cookie

1. Log in to IMDb in your browser
2. Go to browser developer tools (F12) -> Network tab
3. Refresh the page
4. Find any request to imdb.com
5. In the request headers, copy the entire `cookie` value

### Exporting Kinopoisk Ratings

1. Log in to your Kinopoisk account
2. Go to your rated movies page
3. Go to browser developer tools (F12) -> Console tab
4. Paste the content of `scripts/collect.js` into the console
5. Copy the resulting JSON output
6. Save it to a file named `data.json` in the project directory
> **Important:** Kinopoisk limits each ratings page to display a maximum of 200 movies. If you have rated more than 200 movies, you will need to:
>
> 1. Navigate through each page of your ratings
> 2. Run the `collect.js` script on each page separately
> 3. Combine the outputs into a single `data.json` file, or process each batch separately
> 4. Run the application for each batch of collected movies

### Usage

Run the application:

```
go run .
```

Repeat step `3. Prepare your movie data` for every page with ratings
 

The program will:
1. Read your movie ratings from `data.json`
2. Search for each movie on IMDb using the OMDb API
3. Rate each movie on your IMDb account
4. Log any issues to the console and to `warnings.json`

### Troubleshooting & Error Handling

The application records all issues into console and in `warnings.json` with detailed error messages. Here are common issues and their solutions:

- **Movie Not Found**: Some movies may not be found due to title differences or limited database entries. Check `warnings.json` for the specific error and try manually searching and rating these films on IMDb.

- **Rating Failed**: If the rating update fails, you'll see "IMDb API error" messages. Try refreshing your IMDb cookie in the config file, as session cookies expire periodically.

- **API Limits**: The free tier of OMDb API has daily usage limits. If you see "OMDb API limit exceeded" errors, wait until the next day to continue processing.

- **Potentially Incorrect Matches**: The application will flag movies where the title match is uncertain. These entries will still be rated but marked in `warnings.json` for you to verify manually.

- **Missing Ratings**: Movies without ratings in your Kinopoisk export will be skipped. Add these to your IMDb "Check-in" list manually.

---

<a name="russian"></a>

## Перенос оценок фильмов с Кинопоиска на IMDb

kp2imdb - утилита на Go, которая помогает перенести ваши оценки фильмов с [Кинопоиска](https://kinopoisk.ru) на [IMDb](https://imdb.com).

### Начало работы

#### Требования

- Go 1.24 или новее - [скачать](https://go.dev/doc/install)
- Ключ API OMDb (получите на [omdbapi.com](https://www.omdbapi.com/apikey.aspx))
- Аккаунт IMDb с валидной сессионной cookie

#### Установка

1. Клонируйте репозиторий:
   ```
   git clone https://github.com/OutOfStack/kp2imdb.git
   cd kp2imdb
   ```

2. Настройте приложение:
    - Скопируйте cookie сессии IMDb (инструкция ниже)
    - Получите ключ API OMDb
    - Отредактируйте `config.json`:
      ```json
      {
        "locale": "ru",
        "omdb_api_key": "ваш_ключ_omdb_api",
        "imdb_cookie": "ваша_cookie_imdb"
      }
      ```

3. Подготовьте данные о фильмах:
    - Экспортируйте оценки с Кинопоиска с помощью предоставленного скрипта
    - Сохраните экспортированные данные как `data.json`

#### Как получить cookie IMDb

1. Войдите в IMDb в браузере
2. Перейдите в инструменты разработчика бразуера (F12) -> Вкладка Сеть
3. Обновите страницу
4. Найдите любой запрос к imdb.com
5. В заголовках запроса скопируйте полное значение поля `cookie`

### Экспорт оценок с Кинопоиска

1. Войдите в свой аккаунт на Кинопоиске
2. Перейдите на страницу с оценками фильмов
3. Перейдите в инструменты разработчика браузера (F12) -> Вкладка Консоль
4. Вставьте содержимое `scripts/collect.js` в консоль
5. Скопируйте полученный JSON-вывод
6. Сохраните его в файл `data.json` в директории проекта
> **Важно:** Кинопоиск ограничивает отображение на странице максимум 200 фильмов. Если вы оценили более 200 фильмов, вам нужно:
>
> 1. Переходить по каждой странице с вашими оценками
> 2. Запускать скрипт `collect.js` на каждой странице отдельно
> 3. Объединить результаты в один файл `data.json` или обрабатывать каждый набор отдельно
> 4. Запускать приложение для каждого набора собранных фильмов

### Использование

Запустите приложение:

```
go run .
```

Программа выполнит:
1. Чтение оценок фильмов из `data.json`
2. Поиск каждого фильма на IMDb с помощью API OMDb
3. Проставление оценок в вашем аккаунте IMDb
4. Журналирование проблем в консоль и файл `warnings.json`

### Устранение проблем и обработка ошибок

Приложение записывает все проблемы в консоль и в файл `warnings.json` с подробными сообщениями. Вот частые проблемы и их решения:

- **Фильм не найден**: Некоторые фильмы могут не найтись из-за различий в названиях или ограничений базы данных. Проверьте `warnings.json` для уточнения ошибки и попробуйте вручную найти и оценить эти фильмы на IMDb.

- **Ошибка проставления оценки**: Если обновление оценки не удалось, вы увидите сообщение "Ошибка IMDb API". Попробуйте обновить cookie IMDb в конфигурационном файле, так как сессионные cookie периодически истекают.

- **Лимиты API**: Бесплатный тариф OMDb API имеет ограничения на ежедневное использование. Если вы видите ошибку "Превышен лимит запросов к OMDb", подождите до следующего дня, чтобы продолжить обработку.

- **Потенциально неверные совпадения**: Приложение отметит фильмы, где совпадение названия неоднозначно. Эти записи всё равно будут оценены, но помечены в `warnings.json` для ручной проверки.

- **Отсутствующие оценки**: Фильмы без оценок в вашем экспорте с Кинопоиска будут пропущены. Добавьте их в список "Check-in" на IMDb вручную.
