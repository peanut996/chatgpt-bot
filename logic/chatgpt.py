import os
from revChatGPT.V1 import Chatbot


def chat_with_chatgpt(sentence: str) -> str:
    account = os.getenv('OPEN_AI_ACCOUNT')
    passwd = os.getenv('OPEN_AI_PASSWORD')
    conversation_id = os.getenv('OPEN_AI_CONVERSATION_ID')
    chatbot = Chatbot(config={
        "email": account,
        "password": passwd
    }, conversation_id=conversation_id)
    res = ""
    prev_text = ""
    for data in chatbot.ask(sentence):
        message = data["message"][len(prev_text):]
        res += message
        prev_text = data["message"]
    return res
