import os
from revChatGPT.V1 import Chatbot
import asyncio


bot: Chatbot = None

def init_chatgpt():
    global bot
    account = os.getenv('OPEN_AI_ACCOUNT')
    passwd = os.getenv('OPEN_AI_PASSWORD')
    conversation_id = os.getenv('OPEN_AI_CONVERSATION_ID')
    bot = Chatbot(config={
        "email": account,
        "password": passwd
    }, conversation_id=conversation_id)

async def async_chat_with_chatgpt(sentence: str) -> str:
    loop = asyncio.get_event_loop()
    result = await loop.run_in_executor(None, chat_with_chatgpt, sentence)
    return result

def chat_with_chatgpt(sentence: str) -> str:
    if bot is None:
        init_chatgpt()
    res = ""
    prev_text = ""
    for data in bot.ask(sentence):
        message = data["message"][len(prev_text):]
        res += message
        prev_text = data["message"]
    return res
