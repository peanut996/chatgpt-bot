import logging
import random
import time
from typing import List

import OpenAIAuth

from .credential import Credential


class Session:
    def __init__(self, config):
        self.used_chatgpt_credentials_indexes = []
        try:
            self.chatgpt_credentials: List[Credential] = list(map(Credential.parse, config["engine"]["chatgpt"]["tokens"]))
        except OpenAIAuth.Error as e:
            logging.error("Init Credential Error: status: {}, details: {}".format(e.status_code, e.details))
            raise e
        self.chat_gpt_bot = None
        self.edge_gpt_bot = None
        self.user_to_credential = dict()
        self.user_to_last_chat_time = dict()
        self.verbose = config['engine'].get('debug', False)
        for c in self.chatgpt_credentials:
            c.set_verbose(self.verbose)

    def _get_random_chat_gpt_credential(self):
        length_range = len(self.chatgpt_credentials) - 1
        index = random.randint(0, length_range)
        while index in self.used_chatgpt_credentials_indexes:
            if len(self.used_chatgpt_credentials_indexes) == len(self.chatgpt_credentials):
                self.used_chatgpt_credentials_indexes = []
            index = random.randint(0, length_range)

        self.used_chatgpt_credentials_indexes.append(index)
        return self.chatgpt_credentials[index]

    def _generate_chat_gpt_bot(self, user_id=None) -> Credential:
        if user_id is None:
            credential = self._get_random_chat_gpt_credential()
            return credential
        else:
            if user_id not in self.user_to_last_chat_time:
                self.user_to_last_chat_time[user_id] = time.time()
                credential = self._get_random_chat_gpt_credential()
                self.user_to_credential[user_id] = credential
                return credential
            else:
                self.user_to_last_chat_time[user_id] = time.time()
                if time.time() - self.user_to_last_chat_time[user_id] > 60 * 5:
                    credential = self._get_random_chat_gpt_credential()
                    self.user_to_credential[user_id] = credential
                    return credential
                else:
                    return self.user_to_credential[user_id]

    async def chat_with_chatgpt(self, sentence: str, user_id=None) -> str:
        bot = self._generate_chat_gpt_bot(user_id=user_id)
        logging.info("ChatGPTBot using token: {}".format(bot.email))
        async with bot.lock:
            try:
                bot.refresh_token()
                res = ""
                prev_text = ""
                for data in bot.chat_gpt_bot.ask(sentence):
                    message = data["message"][len(prev_text):]
                    res += message
                    prev_text = data["message"]
                if len(res) == 0:
                    raise Exception("empty response")
                return res
            except Exception as e:
                logging.error("ChatGPTBot error: {}".format(e))
                raise e
