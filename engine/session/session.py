import asyncio
import logging
import random

from revChatGPT.V1 import Chatbot as ChatGPTBot


class Session:
    def __init__(self, config):
        self.used_chatgpt_credentials_indexes = []
        self.chatgpt_credentials = list(map(Session.map_token_to_dict, config["engine"]["chatgpt"]["tokens"]))
        self.chat_gpt_bot = None
        self.edge_gpt_bot = None
        self.verbose = config['engine'].get('debug', False)

    def _get_random_chat_gpt_credential(self):
        length_range = len(self.chatgpt_credentials) - 1
        index = random.randint(0, length_range)
        while index in self.used_chatgpt_credentials_indexes:
            if len(self.used_chatgpt_credentials_indexes) == len(self.chatgpt_credentials):
                self.used_chatgpt_credentials_indexes = []
            index = random.randint(0, length_range)

        self.used_chatgpt_credentials_indexes.append(index)
        return self.chatgpt_credentials[index]

    def _generate_chat_gpt_bot(self):
        credential = self._get_random_chat_gpt_credential()
        self.chat_gpt_bot = Session.init_chat_gpt_bot_with_credential(credential, self.verbose)

    def _chat_with_chat_gpt(self, sentence: str) -> str:
        if self.chat_gpt_bot is None:
            self._generate_chat_gpt_bot()
        max_retry = len(self.chatgpt_credentials)
        current = 0
        while current < max_retry:
            try:
                res = ""
                prev_text = ""
                for data in self.chat_gpt_bot.ask(sentence):
                    message = data["message"][len(prev_text):]
                    res += message
                    prev_text = data["message"]
                if len(res) == 0:
                    raise Exception("empty response")
                return res
            except Exception as e:
                logging.error("ChatGPTBot error: {}".format(e))
                self._generate_chat_gpt_bot()
                current += 1
        raise Exception("exceed max retry")

    async def chat_with_chatgpt(self, sentence: str) -> str:
        loop = asyncio.get_event_loop()
        result = await loop.run_in_executor(None, self._chat_with_chat_gpt, sentence)
        return result

    @staticmethod
    def map_token_to_dict(token):
        credential = token.split(":")
        length = len(credential)
        if length != 2 and length != 3:
            raise Exception("token format error")
        if length == 2:
            return {'email': credential[0], 'password': credential[1]}
        else:
            return {'email': credential[0], 'password': credential[1], 'conversation_id': credential[2]}

    @staticmethod
    def init_chat_gpt_bot_with_credential(credential, verbose=False):
        return ChatGPTBot(config={
            **credential,
            "verbose": verbose
        }, conversation_id=credential.get("conversation_id"))
