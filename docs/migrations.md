#### Миграции
* Создаст файл миграции - там пишем нужный sql код и выполняем `make goose cmd="up"` для запуска миграции
```shell
make goose cmd="create some_migration sql"
```