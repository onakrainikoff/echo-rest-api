# echo-rest-api
В данном проекте показан пример создания REST API на Go с использованием Echo framework.

### В проекте использованы
#### HTTP
- `github.com/labstack/echo` - для создания rest сервиса; использованы: _router_, _data binding_ и _data rendering_, _logger middleware_

#### Валидация
- `gopkg.in/go-playground/validator.v9` - для _data validation_  

#### База данных
- `database/sql` - для работы с запросами и транзакционностью
- `github.com/lib/pq` - в качестве СУБД использовалась postgresql
- `github.com/rubenv/sql-migrate` - миграционная тулза для sql
- `docker postgres image` - для запуска postgresql
#### Конфигурация
- `flag` - для передачи параметров при запуске
- `github.com/jinzhu/configor` - для загрузки конфигурации из yaml

#### Логгирование
- `github.com/sirupsen/logrus` - для логов приложения 

#### Тестирование
- `testing` - для написания тестов
- `github.com/golang/mock/gomock` - для генерации моков
- `github.com/stretchr/testify/assert` - исключительно для _assert_ в тестах

#### Зависимости
`github.com/golang/dep` - менеджер зависимостей

### Для запуска
- `make init` - установит необходимые тулзы (dep, sql-migrate, mockgen), затем через `dep` установит зависимости проекта
- `make run`  - запустит приложение, при этом запустив docker контейнер с БД
- `make test` - запустит приложение, при этом запустив docker контейнер с БД
