from datetime import datetime
from .credential import Credential


class UserSession:

    def __init__(self, user_id, conversation_id=None, parent_id=None, credential: Credential = None):
        self.user_id = user_id
        self.conversation_id = conversation_id
        self.parent_id = parent_id
        self.last_time = datetime.now()
        self.credential = credential

    def update(self, parent_id=None, conversation_id=None):
        self.last_time = datetime.now()
        self.parent_id = parent_id
        self.conversation_id = conversation_id or self.conversation_id
