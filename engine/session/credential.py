import asyncio
import logging

from revChatGPT.V3 import Chatbot as ChatGPTBot


class Credential:
    def __init__(self, email, password, api_key, conversation_id=None, verbose=False, loop=None):
        self.email = email
        self.password = password
        self.conversation_id = conversation_id
        self.api_key = api_key
        self.lock = asyncio.Lock()
        self.verbose = verbose
        logging.info("[Credential] init: {}".format(email))
        self.chat_gpt_bot = ChatGPTBot(api_key)

    def set_verbose(self, verbose):
        self.verbose = verbose
        self.chat_gpt_bot.verbose = verbose

    @staticmethod
    def parse(credential_str: str):
        credential = credential_str.split(":")
        length = len(credential)
        if length != 3:
            raise Exception("token format error")
        return Credential(credential[0], credential[1], credential[2])
