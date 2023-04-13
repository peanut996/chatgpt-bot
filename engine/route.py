import json
import logging
import traceback
import asyncio

from OpenAIAuth import Error as OpenAIError
from quart import Quart, request, make_response
from revChatGPT.typings import Error as ChatGPTError

from session.session import Session

app = Quart(__name__)
session: Session

from dataclasses import dataclass

STREAM_TIMEOUT = 'app_quart_stream_timeout'
STREAM_DONE = 'app_quart_stream_done'


@dataclass
class ServerSentEvent:
    data: str
    event: str = 'event'
    id: int = None
    retry: int = None

    def encode(self) -> bytes:
        if self.data != '[DONE]' and self.data != '[START]' and self.data != '[KEEP]':
            self.data = json.dumps({
                'message': self.data,
            })
        message = f"data: {self.data}"
        if self.event is not None:
            message = f"{message}\nevent: {self.event}"
        if self.id is not None:
            message = f"{message}\nid: {self.id}"
        if self.retry is not None:
            message = f"{message}\nretry: {self.retry}"
        message = f"{message}\r\n\r\n"
        return message.encode('utf-8')

    @staticmethod
    def done_event():
        return ServerSentEvent("[DONE]", event="event")

    @staticmethod
    def start_event():
        return ServerSentEvent("[START]", event="event")

    @staticmethod
    def keep_event():
        return ServerSentEvent("[KEEP]", event="event")


@app.route('/chat', methods=["GET"])
async def chat():
    sentence = request.args.get("sentence")
    user_id = request.args.get("user_id")
    model = request.args.get("model") or 'text-davinci-002-render-sha'
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
async def chat_stream():
    sentence = request.args.get("sentence")
    user_id = request.args.get("user_id")
    model = request.args.get("model") or 'text-davinci-002-render-sha'
    start_time = asyncio.get_event_loop().time()

    async def send_events():
        should_stop = False

        async def put_stream_to_queue(stream, queue):
            async for message in stream:
                if should_stop is True:
                    break
                await queue.put(message)
            await queue.put(STREAM_DONE)

        try:
            yield ServerSentEvent.start_event().encode()

            stream_generator = session.chat_stream_with_chatgpt(sentence, user_id=user_id, model=model)
            queue = asyncio.Queue()
            asyncio.create_task(put_stream_to_queue(stream_generator, queue))

            while True:
                try:
                    message = await asyncio.wait_for(queue.get(), timeout=12)
                except asyncio.TimeoutError:
                    message = STREAM_TIMEOUT

                if queue.empty() and message is STREAM_DONE:
                    break
                if message is STREAM_TIMEOUT:
                    current_time = asyncio.get_event_loop().time()
                    if current_time - start_time > 120:
                        should_stop = True
                        logging.warning("[Engine] chat gpt engine get stream timeout")
                        yield ServerSentEvent("😱 机器人负载过多，请稍后再试(The robot is overwhelmed, please try again later)").encode()
                        break
                    yield ServerSentEvent.keep_event().encode()
                else:
                    yield ServerSentEvent(message).encode()



        except OpenAIError as exception:
            logging.error(
                "[Engine] chat gpt engine get open api error: status: {}, details: {}".format(exception.status_code,
                                                                                              exception.details))
            yield ServerSentEvent(exception.details).encode()
        except ChatGPTError as exception:
            logging.error("[Engine] chat gpt engine get chat gpt error: {}".format(exception.message))
            yield ServerSentEvent(exception.message).encode()
        except Exception as exception:
            logging.error(f"[Engine] chat gpt engine get error: {traceback.format_exc()}")
            msg = str(exception) if len(str(exception)) != 0 else "Internal Server Error"
            yield ServerSentEvent(msg).encode()
        finally:
            yield ServerSentEvent.done_event().encode()

    response = await make_response(
        send_events(),
        {
            'Content-Type': 'text/event-stream',
            'Cache-Control': 'no-cache, no-transform',
            'Transfer-Encoding': 'chunked',
        },
    )
    response.timeout = None
    return response


@app.route('/ping')
def ping():
    return "pong"


def set_session(s: Session):
    global session
    session = s
