---
layout: default
---

#### Пререквизиты
1. Linux
2. Установленный Tmux


####  Примеры разворачивания сети c различными типами консенсуса и криптографии (3 майнера)
- Raft + ГОСТ криптография: `devchain/raft_gost`
- Raft + NIST криптография: `devchain/raft`
- Proof-of-Work + ГОСТ криптография: `devchain/pow_gost`
-  Proof-of-Work + NIST криптография: `devchain/pow`
- Proof-of-Authority (Clique) + ГОСТ криптография: `devchain/clique_gost`
-  Proof-of-Authority (Clique) + NIST криптография: `devchain/clique`

#### Перед запуском сети необходимо ее инициализировать:

```sh
./init_chain.sh
```

#### Запустить 3 майнера + бутнода в Tmux:
```sh
./run_chain.sh
```
Подключаться к Tmux окну можно командой:
```sh
tmux a -t node1
```
Проверить открытые окна Tmux
```sh
tmux ls
```
#### Отсановить сеть:
```sh
./kill_chain.sh
```
#### Очистить бд узлов сети:
```sh
./clean.sh
```

[back](./)
