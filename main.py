from dotenv import load_dotenv
from logic.chatgpt import chat_with_chatgpt
import logging
from bottle import route, run ,request

@route('/chat')
def chat():
    sentence = request.query.sentence
    logging.info(f"Received sentence: {sentence}")
    res = chat_with_chatgpt(sentence)
    logging.info(f"Response: {res}")
    return {"message": res}

@route('/ping')
def ping():
    return "pong"


if __name__ == "__main__":
    load_dotenv()
    logging.basicConfig(level=logging.INFO)
    logging.info("Starting server")
    run(host="0.0.0.0", port=5000, debug=False)
