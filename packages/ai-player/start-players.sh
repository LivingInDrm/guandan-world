#!/bin/bash

ROOM_CODE=${1:-""}
SERVER_URL=${2:-"http://localhost:8080"}
PLAYER_COUNT=${3:-3}

if [ -z "$ROOM_CODE" ]; then
  echo "Usage: ./start-players.sh <room-code> [server-url] [player-count]"
  echo "Example: ./start-players.sh 7185"
  echo "         ./start-players.sh 7185 http://localhost:8080 3"
  exit 1
fi

echo "Starting $PLAYER_COUNT AI players for room $ROOM_CODE..."

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

for i in $(seq 1 $PLAYER_COUNT); do
  echo "Starting AI Player $i..."
  npx tsx "$SCRIPT_DIR/src/index.ts" \
    --server "$SERVER_URL" \
    --room-code "$ROOM_CODE" \
    --auto-register \
    --nickname "AI-$i" &
done

echo "All $PLAYER_COUNT AI players started. Press Ctrl+C to stop all."

trap "echo 'Stopping all AI players...'; kill 0" SIGINT SIGTERM

wait
