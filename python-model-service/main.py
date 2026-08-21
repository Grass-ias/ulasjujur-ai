from fastapi import FastAPI
from pydantic import BaseModel
from typing import List
import torch
from transformers import AutoTokenizer, AutoModelForSequenceClassification

app = FastAPI()

print("Memuat Tokenizer dan 2 Model ke Memori...")
tokenizer = AutoTokenizer.from_pretrained("./model_sentimen_saved")
model_sentiment = AutoModelForSequenceClassification.from_pretrained("./model_sentimen_saved")
model_emotion = AutoModelForSequenceClassification.from_pretrained("./model_emosi_saved")
print("Model Siap Melayani Request!")

map_sent_reverse = {0: "negative", 1: "positive"}
map_emo_business = {
    0: "Puas", 
    1: "Kecewa / Menyesal", 
    2: "Was-was / Khawatir", 
    3: "Sangat Puas / Loyal", 
    4: "Frustrasi / Komplain"
}

class ReviewInput(BaseModel):
    id: int
    text: str

class BatchRequest(BaseModel):
    reviews: List[ReviewInput]

@app.post("/predict-batch")
def predict_batch(payload: BatchRequest):
    texts = [r.text for r in payload.reviews]
    inputs = tokenizer(texts, padding=True, truncation=True, return_tensors="pt", max_length=128)
    results = []
    
    with torch.no_grad():
        # Prediksi Sentimen pakai Softmax buat dapet probabilitas asli
        out_sent = model_sentiment(**inputs)
        probs_sent = torch.nn.functional.softmax(out_sent.logits, dim=-1)
        preds_sent = torch.argmax(probs_sent, dim=-1).tolist()
        
        # Prediksi Emosi
        out_emo = model_emotion(**inputs)
        preds_emo = torch.argmax(out_emo.logits, dim=-1).tolist()

    # # PROSES JAHIT HASILNYA
    # for i, r in enumerate(payload.reviews):
    #     # Ambil angka confidence dari kelas yang menang
    #     conf_score = round(probs_sent[i][preds_sent[i]].item(), 2)
        
    #     # LOGIC MIXED SENTIMENT (Opsi 2 yang kita sepakati)
    #     # Jika AI yakinnya di bawah 75%, kita tandain ini sebagai ulasan abu-abu
    #     is_mixed_flag = False
    #     if conf_score < 0.75:
    #         is_mixed_flag = True

    #     results.append({
    #         "id": r.id,
    #         "sentiment": map_sent_reverse.get(preds_sent[i], "unknown"),
    #         "emotion": map_emo_business.get(preds_emo[i], "unknown"),
    #         "confidence": conf_score,
    #         "is_mixed": is_mixed_flag # Fitur baru buat ditangkap UI!
    #     })

    # return {"results": results}
# PROSES JAHIT HASILNYA
    for i, r in enumerate(payload.reviews):
        # Ambil angka confidence dari kelas yang menang
        conf_score = round(probs_sent[i][preds_sent[i]].item(), 2)

        results.append({
            "id": r.id,
            "sentiment": map_sent_reverse.get(preds_sent[i], "unknown"),
            "emotion": map_emo_business.get(preds_emo[i], "unknown"),
            "confidence": conf_score
            # is_mixed KITA BUANG! Bye-bye beban!
        })

    return {"results": results}