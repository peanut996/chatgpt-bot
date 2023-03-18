dir_path=$(cd "$(dirname "$0")" && pwd)

export ACCESS_TOKEN=""
export PUID=""
export GIN_MODE=release
export CHATGPT_BASE_URL="http://127.0.0.1:8080/api/"


echo "##################################"
echo "###### RUN proxy server ##########"
echo "##################################"


cd "$dir_path/../proxy"
go build -ldflags '-w -s' -o proxy .
chmod +x ./proxy

proxy_pid=$(ps -ef | awk '/.\/proxy/{print $2}')

kill $proxy_pid

nohup ./proxy > /tmp/proxy.log &

echo "start proxy on port 8080 done"

echo "##################################"
echo "###### RUN engine server #########"
echo "##################################"

git pull -r

engine_pid=$(ps -ef | awk '/python main.py/{print $2}')

client_pid=$(ps -ef | awk '/chatgpt-bot/{print $2}')

cd "$dir_path/../engine" || exit

echo "kill engine..."

kill $engine_pid

echo "kill engine done"

source venv/bin/activate

pip install -r requirements.txt --upgrade

nohup python main.py -c ../etc/config.yaml > /tmp/engine.log &

echo "run engine success"

echo "##################################"
echo "###### RUN client ################"
echo "##################################"

cd "$dir_path/../client" || exit

./build.sh

nohup ./chatgpt-bot -c ../etc/config.yaml > client.log &

echo "run client success."

echo "kill client.."

kill $client_pid

echo "kill client done"

echo "Reboot success."

