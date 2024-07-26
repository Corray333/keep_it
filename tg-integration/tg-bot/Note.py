class Note:
    def __init__(self, content, date, original):
        self.content = content
        self.source = "tg"
        self.icon = {
            "type": "icon",
            "icon": "mingcute:telegram-fill"
        }
        self.original = original

        self.created_at = date




