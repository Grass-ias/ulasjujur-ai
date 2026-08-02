from fastapi import FastAPI
from pydantic import BaseModel
from typing import List
import random

app = FastAPI()

# Definisi struktur Input dari Golang
class ReviewInput(BaseModel):
    id: int
    text: str

class BatchRequest(BaseModel):
    reviews: List[ReviewInput]

# Endpoint untuk batch inference
@app.post("/predict-batch")
def predict_batch(payload: BatchRequest):
    results = []

    # LOGIC DUMMY: Akan diganti model IndoBERT di Minggu 2
    for r in payload.reviews:
        dummy_sentiment = random.choice(["positive", "negative"])
        dummy_emotion = random.choice(["bahagia", "sedih", "marah", "kecewa", "takut"])
        dummy_confidence = round(random.uniform(0.75, 0.99), 2)

        results.append({
            "id": r.id,
            "sentiment": dummy_sentiment,
            "emotion": dummy_emotion,
            "confidence": dummy_confidence
        })

    return {"results": results}