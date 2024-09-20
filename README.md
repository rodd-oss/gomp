# Tomb Mates

Процесс написания игры можно посмотреть на [YouTube](https://www.youtube.com/watch?v=jMqC_VUEAgs). Версия кода на момент записи видео доступна в [этом коммите](https://github.com/jilio/tomb_mates/commit/a6823c442f40ba63f5ecd04db2814d47b5e6b150).

![Скриншот](./screenshot.png 'Скриншотес')

## Запустить сервер

```
cd cmd/server
go run *go
```

## Запустить игру

###### Не будет работать без запущенного сервера

```
cd cmd/client
go run main.go
```

## Спрайты

[DungeonTilesetII](https://0x72.itch.io/dungeontileset-ii)

---

## Команды

```
protoc --go_out=. *.proto
sh svg2icns.sh icon.svg
```
