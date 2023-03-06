import asyncio
import time
import logging

from revChatGPT.V1 import Chatbot as ChatGPTBot


class Credential:
    def __init__(self, email, password, conversation_id=None, verbose=False):
        self.email = email
        self.password = password
        self.conversation_id = conversation_id
        loop = asyncio.get_event_loop()
        self.lock = asyncio.Lock(loop=loop)
        self.verbose = verbose
        self.last_update_time = time.time()
        self.chat_gpt_bot = ChatGPTBot(config={
            'email': email,
            'password': password,
            'verbose': verbose
        }, conversation_id=conversation_id)

    def set_verbose(self, verbose):
        self.verbose = verbose
        self.chat_gpt_bot.verbose = verbose

    def refresh_token(self):
        if time.time() - self.last_update_time > 60 * 30:
            self.chat_gpt_bot = ChatGPTBot(config={
                'email': self.email,
                'password': self.password,
                'verbose': self.verbose
            }, conversation_id=self.conversation_id)
            self.last_update_time = time.time()
            logging.info("ChatGPTBot token refreshed: {}".format(self.email))

    @staticmethod
    def parse(credential_str: str):
        credential = credential_str.split(":")
        length = len(credential)
        if length != 2 and length != 3:
            raise Exception("token format error")
        if length == 2:
            return Credential(credential[0], credential[1])
        else:
            return Credential(credential[0], credential[1], credential[2])
