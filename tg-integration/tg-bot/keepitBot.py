import json
from Note import Note
import telebot
from dotenv import load_dotenv
import os
from telebot import types

load_dotenv()
TOKEN = os.getenv('TG_API_TOKEN')
bot = telebot.TeleBot(TOKEN)


@bot.message_handler(content_types=['text', 'photo', 'video'])
def get_message(message):
    formatted_json = json.dumps(message.json, ensure_ascii=False, indent=4)
    print(formatted_json)
    content = []
    rich_text = {

    }

    # Формирование plain_text

    if "text" in message.json.keys():
        rich_text['plain_text'] = message.text
    elif "caption" in message.json.keys():
        rich_text['plain_text'] = message.caption

    # ФОРМИРОВАНИЕ META
    meta = []
    if "entities" not in message.json.keys() and "caption_entities" not in message.json.keys():  # Если оформления нет, мета пустая
        rich_text['meta'] = meta
    elif "entities" in message.json.keys():  # Формирование мета информации для текста БЕЗ ВЛОЖЕНИЙ
        for ent in message.json['entities']:
            meta.append({
                "style": [ent['type']],
                "length": ent['length'],
                "offset": ent['offset']
            })
    elif "caption_entities" in message.json.keys():
        for cap_ent in message.json['caption_entities']:
            meta.append({
                "style": [cap_ent['type']],
                "length": cap_ent['length'],
                "offset": cap_ent['offset']
            })

    if meta:  # Объединение применяемых к одному блоку текста стилей в массив
        for block_index in range(len(meta) - 1):
            if meta[block_index]['length'] == meta[block_index + 1]['length'] and meta[block_index]['offset'] == \
                    meta[block_index + 1]['offset']:
                meta[block_index + 1]['style'].append(meta[block_index]['style'][0])
                meta[block_index] = {}
        clean_meta = [meta_item for meta_item in meta if meta_item]  # Очистка массива от пустых значений
        rich_text['meta'] = clean_meta

    # ФОРМИРОВАНИЕ МЕТА ЗАКОНЧЕНО

    # ФОРМИРОВАНИЕ CONTENT

    content = [{"rich_text": rich_text}, {"type": "p"}]

    if "photo" in message.json.keys():
        file_path = bot.get_file(message.json['photo'][3]['file_id']).file_path
        file_link = f"https://api.telegram.org/file/bot{TOKEN}/{file_path}"
        photo_block = {
            "caption": "",
            "src": file_link,
            "type": "img"
        }
        content.append(photo_block)

    if "video" in message.json.keys():
        file_path = bot.get_file(message.json['video']['file_id']).file_path
        file_link = f"https://api.telegram.org/file/bot{TOKEN}/{file_path}"
        video_block = {
            "caption": "",
            "src": file_link,
            "type": "video"
        }
        content.append(video_block)

    # ФОРМИРОВАНИЕ CONTENT ЗАКОНЧЕНо

    # ФОРМИРОВАНИЕ original
    original = {}
    if "forward_origin" in message.json.keys():
        original = {
            "type": message.forward_origin.type,
            "created_at": message.forward_origin.date,
        }
        if original['type'] == "channel":
            original['title'] = message.forward_origin.chat.title
            original["link"] = f"https://t.me/{message.forward_origin.chat.username}"
        elif original['type'] == "user":
            original['title'] = (message.forward_origin.sender_user.first_name +
                                 " " + message.forward_origin.sender_user.last_name)
            original['link'] = f"https://t.me/{message.forward_origin.sender_user.username}"

    date = message.json['date']
    new_note = Note(content, date, original)
    bot.send_message(message.from_user.id, json.dumps(new_note.__dict__, ensure_ascii=False, indent=4))


bot.infinity_polling()
