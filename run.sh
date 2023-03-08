git pull -r

engine_pid=$(ps -ef | awk '/python main.py/{print $2}')

client_pid=$(ps -ef | awk '/chatgpt-bot/{print $2}')

cd engine || exit

echo "kill engine..."
kill $engine_pid
echo "kill engine done"

source venv/bin/activate

pip install -r requirements.txt --upgrade

nohup python main.py -c ../config.yaml > engine.log &

echo "run engine success"

cd ../client || exit

./build.sh

nohup ./chatgpt-bot -c ../config.yaml > client.log &

echo "run client success."

echo "kill client.."

kill $client_pid

echo "kill client done"

echo "Reboot success."

