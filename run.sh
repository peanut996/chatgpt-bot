git pull -r

pids=$(ps -ef | awk '/config.yaml/{print $2}')

cd engine

source venv/bin/activate

pip install -r requirements.txt --upgrade

nohup python main.py -c ../config.yaml > engine.log &

echo "run engine success"

cd ../client

./build.sh

nohup ./chatgpt-bot -c ../config.yaml > client.log &

echo "run client success."

for pid in $pids; do
  echo "Killing process $pid..."
  kill $pid
done

echo "kill client and engine done"

