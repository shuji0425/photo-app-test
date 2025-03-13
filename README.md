# 機能
写真のアップロードと表示をメインに行うアプリ

## 環境
* DB       : PostgreSQL
* Backend  : Go言語(Ginフレームワーク)
* Frontend : TypeScript(React)

## 環境構築
CPUの容量がある人
```bash
# ymlファイルが2つのため
docker-compose -f docker-compose.yml up -d
```

CPUの容量がない人は以下参照（メモリ8GGB にLightroomれると重くなりました)

### PCのCPUがないためDBのみDockerを使用
**シェルスクリプトに権限を与える**
```bash
# start.shはファイル名
chmod +x start.sh
```
```bash
# 起動
./start.sh
```

### initializationブランチ
新しく環境を作るのが面倒なので、React + Go環境のデフォルトを配置
```bash
docker-compose.db.yml
start.sh
```
上記2ファイルはCPU次第では不要