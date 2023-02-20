from EdgeGPT import Chatbot 
import asyncio

bot: Chatbot = None

def init_edgegpt():
    global bot
    bot = Chatbot(cookiePath='./cookie.json')

async def chat_with_edgegpt(sentence: str) -> str:
    if bot is None:
        init_edgegpt()
    return await bot.ask(prompt=sentence)
