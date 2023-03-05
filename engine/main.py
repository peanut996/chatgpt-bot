import argparse
import logging
import os
import traceback

import yaml
from OpenAIAuth import Error as OpenAIError
from flask import Flask, request

from session.session import Session

app = Flask(__name__)
session: Session


@app.route('/chat')
async def chat():
    sentence = request.args.get("sentence")
    user_id = request.args.get("user_id")
    logging.getLogger("app").info(f"[Engine] chat gpt engine get request: from {user_id}: {sentence} ")
    try:
        res = await session.chat_with_chatgpt(sentence, user_id=user_id)
        logging.getLogger("app").info(f"[Engine] chat gpt engine get response: to {user_id}: {res}")
        return {"message": res}
    except OpenAIError as e:
        logging.error(
            "[Engine] chat gpt engine get open api error: status: {}, details: {}".format(e.status_code, e.details))
        return {"detail": e.details, "code": e.status_code}
    except Exception as e:
        logging.error(f"[Engine] chat gpt engine get error: {traceback.format_exc()}")
        return {"detail": e.args}


@app.route('/ping')
def ping():
    return "pong"


def load_yaml(path):
    with open(path, "r", encoding="utf-8") as f:
        return yaml.load(f, Loader=yaml.FullLoader)


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
    logging.getLogger('werkzeug').setLevel(logging.ERROR)
    logging.getLogger("app").setLevel(logging.INFO)

    config = load_config()
    session = Session(config=config)
    port = config['engine']['port']
    debug = config['engine'].get('debug', False)
    logging.info("Starting server")
    app.run(host="0.0.0.0", port=port, debug=debug)


if __name__ == "__main__":
    main()
