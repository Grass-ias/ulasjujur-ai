import pandas as pd
from datasets import Dataset
from transformers import AutoTokenizer, AutoModelForSequenceClassification, Trainer, TrainingArguments
import torch
import os

print("--- MEMULAI SCRIPT TRAINING MOCK ---")

# 1. BACA DATASET & AMBIL SAMPLE 50 BARIS AJA
print("1. Membaca PRDECT-ID Dataset.csv...")
df_full = pd.read_csv("PRDECT-ID Dataset.csv")
df_sample = df_full.head(50).copy() # Cukup 50 baris buat test run

# 2. MAPPING LABEL KE ANGKA
print("2. Mapping Sentimen dan Emosi...")
sentiment_map = {"Negative": 0, "Positive": 1}
df_sample["label_sentiment"] = df_sample["Sentiment"].map(sentiment_map)

emotion_map = {"Happy": 0, "Sadness": 1, "Fear": 2, "Love": 3, "Anger": 4}
df_sample["label_emotion"] = df_sample["Emotion"].map(emotion_map)

# 3. UBAH PANDAS KE HUGGINGFACE DATASET
hf_dataset = Dataset.from_pandas(df_sample)

# 4. TOKENISASI TEKS
print("3. Memuat Tokenizer IndoBERT...")
tokenizer = AutoTokenizer.from_pretrained("indobenchmark/indobert-base-p1")

def tokenize_function(examples):
    return tokenizer(examples["Customer Review"], padding="max_length", truncation=True, max_length=128)

tokenized_datasets = hf_dataset.map(tokenize_function, batched=True)

# Pisahkan dataset untuk 2 model dengan nama kolom label yang standar (wajib bernama "labels")
ds_sentiment = tokenized_datasets.rename_column("label_sentiment", "labels")
ds_emotion = tokenized_datasets.rename_column("label_emotion", "labels")

# ==========================================
# TRAINING MODEL A: SENTIMEN (2 KELAS)
# ==========================================
print("\n--- MULAI TRAINING MODEL SENTIMEN ---")
model_sent = AutoModelForSequenceClassification.from_pretrained("indobenchmark/indobert-base-p1", num_labels=2)

training_args_sent = TrainingArguments(
    output_dir="./results_sent",
    num_train_epochs=1,     # 1 epoch aja karena cuma test run
    per_device_train_batch_size=8,
    logging_steps=5,
    use_cpu=True # Paksa pakai CPU aja biar laptop aman pas ngetes
)

trainer_sent = Trainer(
    model=model_sent,
    args=training_args_sent,
    train_dataset=ds_sentiment,
)

trainer_sent.train()

# SAVE MODEL SENTIMEN KE LOKAL
print("Menyimpan Model Sentimen ke folder './model_sentimen_saved'...")
model_sent.save_pretrained("./model_sentimen_saved")
tokenizer.save_pretrained("./model_sentimen_saved")

# ==========================================
# TRAINING MODEL B: EMOSI (5 KELAS)
# ==========================================
print("\n--- MULAI TRAINING MODEL EMOSI ---")
model_emo = AutoModelForSequenceClassification.from_pretrained("indobenchmark/indobert-base-p1", num_labels=5)

training_args_emo = TrainingArguments(
    output_dir="./results_emo",
    num_train_epochs=1,     # 1 epoch aja
    per_device_train_batch_size=8,
    logging_steps=5,
    use_cpu=True
)

trainer_emo = Trainer(
    model=model_emo,
    args=training_args_emo,
    train_dataset=ds_emotion,
)

trainer_emo.train()

# SAVE MODEL EMOSI KE LOKAL
print("Menyimpan Model Emosi ke folder './model_emosi_saved'...")
model_emo.save_pretrained("./model_emosi_saved")

print("\n--- TEST RUN SELESAI! MODEL BERHASIL DISIMPAN ---")