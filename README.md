# UlasJujur - AI Microservice (Layer 1)

Ini adalah Lapis 1 dari arsitektur UlasJujur. Bertugas menerima batch teks ulasan dari Lapis 2 (Golang) dan mengembalikan prediksi Sentimen serta Emosi (berbasis bahasa bisnis) menggunakan model IndoBERT fine-tuned.

## Cara Setup & Run di Local

1. **Clone Repo dan Siapkan Venv**
   ```bash
   git clone <https://github.com/Grass-ias/ulasjujur-ai>
   cd ulasjujur-ai
   python -m venv venv

2. **Aktifkan Venv**

    Windows: .\venv\Scripts\activate
    Mac/Linux: source venv/bin/activate

3. **Install Dependencies**

    pip install -r requirements.txt

*(Catatan: Pastikan folder model model_sentimen_saved dan model_emosi_saved sudah diekstrak di root folder ini. Folder model TIDAK dipush ke GitHub karena ukurannya besar).*

4. **Run Server**

    uvicorn main:app --reload

*Akses dokumentasi Swagger UI di: http://127.0.0.1:8000/docs*