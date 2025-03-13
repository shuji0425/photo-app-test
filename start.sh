# DBの起動
docker-compose -f docker-compose.db.yml up -d

# Goの起動
cd ./backend && air &

# Reactの起動
cd ./frontend && npm run dev