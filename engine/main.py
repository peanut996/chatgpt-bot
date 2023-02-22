import argparse
import logging
import os

import yaml
from OpenAIAuth import Error as OpenAIError
from flask import Flask, request

from session.session import Session

app = Flask(__name__)
session = None


@app.route('/chat')
async def chat():
    sentence = request.args.get("sentence")
    logging.info(f"[Engine] chat gpt engine get request: {sentence}")
    try:
        # noinspection PyUnresolvedReferences
        res = await session.chat_with_chatgpt(sentence)
        logging.info(f"[Engine] chat gpt engine get response: {res}")
        return {"message": res}
    except OpenAIError as e:
        logging.error(
            "[Engine] chat gpt engine get open api error: status: {}, details: {}".format(e.status_code, e.details))
        return {"detail": e.details, "code": e.status_code}
    except Exception as e:
        logging.error(f"[Engine] chat gpt engine get error: {str(e)}")
        return {"detail": str(e)}


# @app.route('/bing')
# async def bing_chat():
#     sentence = request.args.get("sentence")
#     logging.info(f"[Engine] Request: {sentence}")
#     try:
#         res = await chat_with_edgegpt(sentence)
#         logging.info(f"[Engine] Response: {res}")
#         return {"message": res}
#     except Exception as e:
#         logging.error(f"[Engine] Error: {e}")
#         return {"message": "Error: " + str(e)}


@app.route('/ping')
def ping():
    return "pong"


def check_cookie():
    if os.path.exists("cookie.json") is False:
        logging.error("cookie.json not found")
        exit(1)


# load yaml with file path
def load_yaml(path):
    with open(path, "r", encoding="utf-8") as f:
        return yaml.load(f, Loader=yaml.FullLoader)


# receive arg from cmd line
def get_config_path():
    if os.getenv("BOT_ENGINE_CONFIG_PATH") is not None:
        return os.getenv("BOT_ENGINE_CONFIG_PATH")
    else:
        parser = argparse.ArgumentParser()
        parser.add_argument("-c", "--config", type=str, default="config.yaml")
        args = parser.parse_args()
        return args.config


def load_config():
    try:
        config_path = get_config_path()
        return load_yaml(config_path)

    except Exception as e:
        logging.error(f"load config error: {e}")
        exit(1)


def main():
    global session
    logging.basicConfig(level=logging.INFO)
    config = load_config()
    session = Session(config=config)
    port = config['engine']['port']
    debug = config['engine'].get('debug', False)
    logging.info("Starting server")
    app.run(host="0.0.0.0", port=port, debug=debug)


if __name__ == "__main__":
    main()
