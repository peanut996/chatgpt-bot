import argparse
import logging
import os
import traceback

import yaml
from OpenAIAuth import Error as OpenAIError
from revChatGPT.typings import Error as ChatGPTError
from quart import Quart, request

from session.session import Session

app = Quart(__name__)
session: Session


@app.route('/chat')
async def chat():
    sentence = request.args.get("sentence")
    user_id = request.args.get("user_id")
    model = request.args.get("model")
    try:
        res = await session.chat_with_chatgpt(sentence, user_id=user_id, model=model)
        return {"message": res}
    except OpenAIError as e:
        logging.error(
            "[Engine] chat gpt engine get open api error: status: {}, details: {}".format(e.status_code, e.details))
        return {"detail": e.details, "code": e.status_code}
    except ChatGPTError as e:
        logging.error("[Engine] chat gpt engine get chat gpt error: {}".format(e.message))
        return {"detail": e.message, "code": e.code}
    except Exception as e:
        logging.error(f"[Engine] chat gpt engine get error: {traceback.format_exc()}")
        return {"detail": str(e) if len(str(e)) != 0 else "Internal Server Error", "code": 500}


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
    app.logger.setLevel(logging.WARNING)


    config = load_config()
    session = Session(config=config)
    port = config['engine']['port']
    debug = config['engine'].get('debug', False)

    app.run(host="0.0.0.0", port=port, debug=debug)


if __name__ == "__main__":
    main()
