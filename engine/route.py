import argparse
import logging
import os
import traceback
from OpenAIAuth import Error as OpenAIError
from revChatGPT.typings import Error as ChatGPTError
from quart import Quart, request, Response, stream_with_context

from engine.session.session import Session

app = Quart(__name__)
session: Session


@app.route('/chat', methods=["GET"])
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


@app.route('/chat-stream', methods=["GET"])
async def chat():
    sentence = request.args.get("sentence")
    user_id = request.args.get("user_id")
    model = request.args.get("model")
    try:
        async def generate():
            async for message in session.chat_stream_with_chatgpt(sentence, user_id=user_id, model=model):
                yield message

        return Response(stream_with_context(generate()), content_type='text/plain')
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


def set_session(s: Session):
    global session
    session = s
