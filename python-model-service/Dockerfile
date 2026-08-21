# Menggunakan image Python versi slim agar ukurannya tidak terlalu bengkak
FROM python:3.10-slim

# Menentukan direktori kerja di dalam container
WORKDIR /app

# Meng-copy file requirements terlebih dahulu untuk efisiensi caching Docker
COPY requirements.txt .

# Menginstal semua library yang dibutuhkan (FastAPI, PyTorch, Transformers, dll)
RUN pip install --no-cache-dir -r requirements.txt

# Meng-copy seluruh kode dan folder model AI ke dalam container
COPY . .

# Membuka port 5000
EXPOSE 5000

# Perintah untuk menyalakan server saat container berjalan
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "5000"]