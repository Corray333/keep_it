from transformers import BertModel, BertTokenizer
import numpy as np

# Загрузка предобученной модели и токенизатора BERT
tokenizer = BertTokenizer.from_pretrained('bert-base-uncased')
model = BertModel.from_pretrained('bert-base-uncased')

def get_text_vector(text):
    inputs = tokenizer(text, return_tensors='pt', max_length=512, truncation=True, padding='max_length')
    outputs = model(**inputs)
    # Берем среднее по токенам, чтобы получить вектор на уровне текста
    return outputs.last_hidden_state.mean(dim=1).detach().numpy()

def transform_to_2d(vector):
    # Простой метод для приведения вектора к двумерному пространству
    # Например, делаем первую и вторую координаты просто элементами вектора
    # Это не будет настоящей проекцией, но обеспечит двумерность
    return np.array([vector[0][0], vector[0][1]])  # Берем первые два элемента вектора

# Печать координат для каждого текста
while True:
    # Ввод текста
    text = input("Введите текст (или 'exit' для выхода): ")
    if text.lower() == 'exit':
        break
    
    # Получение векторного представления текста
    vector = get_text_vector(text)
    
    # Преобразование вектора в двумерные координаты
    reduced_vector = transform_to_2d(vector)
    
    # Выводим координаты текста
    print(f"Текст: {text}")
    print(f"Координаты: {reduced_vector}")
