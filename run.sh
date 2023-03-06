git pull -r

kill $(ps -ef | awk '/config.yaml/{print $2}')

echo "kill client and engine done"

cd engine

source venv/bin/activate

pip install -r requirements.txt --upgrade

nohup python main.py -c ../config.yaml > engine.log &

echo "run engine success"

cd ../client

./build.sh

nohup ./chatgpt-bot -c ../config.yaml > client.log &

echo "run client success."