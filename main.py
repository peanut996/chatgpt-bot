from dotenv import load_dotenv
from logic.chatgpt import chat_with_chatgpt
from flask import Flask, request, jsonify
import logging

app = Flask(__name__)

app.config['JSON_AS_ASCII'] = False

@app.route("/chat", methods=["GET"])
def chat():
    res = dict()
    sentence = request.values.get("sentence")
    logging.info(f"chat: {sentence}")
    res['message'] = chat_with_chatgpt(sentence=sentence)
    return jsonify(res)


if __name__ == "__main__":
    load_dotenv()
    app.run(host="0.0.0.0", port=5000, debug=True)
