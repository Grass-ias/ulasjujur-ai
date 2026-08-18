from fastapi import FastAPI
from pydantic import BaseModel
from typing import List
import torch
from transformers import AutoTokenizer, AutoModelForSequenceClassification

app = FastAPI()

print("Memuat Tokenizer dan 2 Model ke Memori...")

# 1. Load Tokenizer & Model dari folder lokal hasil test run
tokenizer = AutoTokenizer.from_pretrained("./model_sentimen_saved")
model_sentiment = AutoModelForSequenceClassification.from_pretrained("./model_sentimen_saved")
model_emotion = AutoModelForSequenceClassification.from_pretrained("./model_emosi_saved")

print("Model Siap Melayani Request!")

# ==========================================
# MAPPING LABEL KE BAHASA BISNIS / SELLER
# ==========================================
map_sent_reverse = {0: "negative", 1: "positive"}

# Translasi dari label asli dataset (Happy, Sad, dll) ke bahasa e-commerce
map_emo_business = {
    0: "Puas",                # Asal: Happy
    1: "Kecewa / Menyesal",   # Asal: Sadness
    2: "Was-was / Khawatir",  # Asal: Fear
    3: "Sangat Puas / Loyal", # Asal: Love
    4: "Frustrasi / Komplain" # Asal: Anger
}

class ReviewInput(BaseModel):
    id: int
    text: str

class BatchRequest(BaseModel):
    reviews: List[ReviewInput]

@app.post("/predict-batch")
def predict_batch(payload: BatchRequest):
    texts = [r.text for r in payload.reviews]

    # Tokenisasi sekaligus (Batch)
    inputs = tokenizer(texts, padding=True, truncation=True, return_tensors="pt", max_length=128)

    results = []
    
    with torch.no_grad():
        # Prediksi Sentimen
        out_sent = model_sentiment(**inputs)
        probs_sent = torch.nn.functional.softmax(out_sent.logits, dim=-1) # Hitung persentase asli
        preds_sent = torch.argmax(probs_sent, dim=-1).tolist()
        
        # Prediksi Emosi
        out_emo = model_emotion(**inputs)
        preds_emo = torch.argmax(out_emo.logits, dim=-1).tolist()

    # PROSES JAHIT HASILNYA
    for i, r in enumerate(payload.reviews):
        # Ambil nilai probabilitas tertinggi sebagai confidence score
        conf_score = round(probs_sent[i][preds_sent[i]].item(), 2) 

        results.append({
            "id": r.id,
            "sentiment": map_sent_reverse.get(preds_sent[i], "unknown"),
            "emotion": map_emo_business.get(preds_emo[i], "unknown"),
            "confidence": conf_score # Sekarang angkanya ASLI dari otak AI!
        })

    return {"results": results}