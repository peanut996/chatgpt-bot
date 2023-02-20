from dotenv import load_dotenv
from logic.chatgpt import async_chat_with_chatgpt
import logging
import os
from flask import Flask, request


app = Flask(__name__)


@app.route('/chat')
async def chat():
    sentence = request.args.get("sentence")
    logging.info(f"[Engine] Request: {sentence}")
    try:
        res = await async_chat_with_chatgpt(sentence)
        logging.info(f"[Engine] Response: {res}")
        return {"message": res}
    except Exception as e:
        logging.error(f"[Engine] Error: {e}")
        return {"message": str(e)}

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

def checkCookie():
    if os.path.exists("cookie.json") is False:
        logging.error("cookie.json not found")
        exit(1)

if __name__ == "__main__":
    load_dotenv()
    logging.basicConfig(level=logging.INFO)
    logging.info("Starting server")
    app.run(host="0.0.0.0", port=5000, debug=False)
    # run(host="0.0.0.0",server='asyncio', port=5000, debug=False)
