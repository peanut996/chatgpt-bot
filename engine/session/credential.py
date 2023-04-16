from revChatGPT.V3 import Chatbot


class Credential:
    def __init__(self, api_key=None):
        self.api_key = api_key
        self.chat_gpt_bot = Chatbot(api_key=api_key)

    @staticmethod
    def parse(credential_str: str):
        credential = credential_str.split(":")
        length = len(credential)
        if length == 1:
            return Credential(api_key=credential[0])
        if length != 1:
            raise Exception("token format error")
