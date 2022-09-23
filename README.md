# GSES2 BTC application

## Інструкція з запуску
- Створити .env файл зі своїми api ключами за прикладом файлу .env_sample
- docker-compose build
- docker-compose up

## Запуск лінтера
- локально встановити golangci-lint
-  ```shell
    golangci-lint run
   ```

## Архітектура проекту

![Architecture]([https://drive.google.com/file/d/1DaEDnkuUf2x0Q1IiRUYFrCsVF5KbkCqr/view?usp=sharing](https://lh6.googleusercontent.com/RWAqrZvQUiq_aQ3-6-VBcHuR9noXqRZqm41BzXZDcnWGJ3R25S9B1feNIEK2IfK8hb0=w2400))

### Код написаний за правилами чистої архітектури, що має наступні переваги:
- Такі системи легко розширюються
- Такий код легше тестувати
- Легко змінювати зовнішні залежності при цьому не переписуючи логіку застосунку

### Основна задача - розподіл відповідальності між усіма шарами додатку:
- handler - шар роботи з HTTP

Приймаємо запити й передаємо їх на рівень бізнес логіки. Тому при створенні цього шару впроваджується залежність(dependency injection) із сервісами.
- service - Business logic layer(BLL)

Для роботи цей шар має мати доступ до бд тому впроваджуємо залежність(dependency injection) із шаром доступу до даних
- repository - Data access layer(DAL)

Відповідає за взаємодію з базою даних.

## Безпека
api ключі зберігаються в environment variables що визначаються в .env файлі, який не пушиться у відкриті репозиторії(у файлі docker-compose зазначаеться env_file).
У репозиторії є файл .env_sample, який показує як саме додати потрібні ключі доступу.

Але щоб Ви мали змогу перевірити працездатність застосунку - надам тут свої ключі
```
COINMARKETCAP_API_KEY=55c0acb9-6fa3-40ff-88b6-916b7a2838de
MAILJET_PUBLIC_KEY=16bc99c5fd3923f9e4c7a6985557aa4d
MAILJET_PRIVATE_KEY=f5d009c20e2f9c8cbe687481e17bb536
```

## Використані сервіси
### Coinmarketcap api
Coinmarketcap це найбільш надійне та точне джерело капіталізації криптовалютного ринку, ціноутворення та іншої інформації. Сервіс має безкоштовний api.

### Mailjet
Для відправки емейлів підключено сторонній сервіс, який надає 30 денний free trial період.
