
# UlasJujur Backend Services - Panduan Instalasi & Testing Lokal

Dokumen ini berisi panduan langkah demi langkah untuk melakukan *setup* lingkungan pengembangan lokal dan panduan pengujian API UlasJujur menggunakan *dataset* CSV berskala besar.

Arsitektur sistem ini menggunakan tiga *microservices* utama yang berjalan di dalam Docker:

1. **Lapis 1 (AI Service):** Python FastAPI & PyTorch (IndoBERT)
2. **Lapis 2 (Main API):** Golang Backend
3. **Database:** MySQL (Kamus Isu Kategori)

---

## Tahap 1: Instalasi Docker Desktop (Prasyarat Utama)

Untuk menjalankan seluruh sistem hanya dengan satu perintah, komputer harus memiliki Docker Desktop yang sudah aktif.

1. Unduh *installer* Docker Desktop untuk Windows/Mac melalui situs resmi: [https://www.docker.com/products/docker-desktop/](https://www.docker.com/products/docker-desktop/)
2. Lakukan instalasi seperti biasa. Pastikan opsi **Use WSL 2 instead of Hyper-V** (untuk pengguna Windows) tetap tercentang.
3. *Restart* komputer jika diminta oleh sistem.
4. Buka aplikasi **Docker Desktop** dari menu Start.
5. Tunggu beberapa saat hingga indikator di sudut kiri bawah aplikasi berubah menjadi hijau dengan status **"Engine running"**.

---

## Tahap 2: Menyalakan Sistem (Orkestrasi)

Setelah Docker menyala di latar belakang, kita bisa langsung membangun dan menyalakan ketiga mesin *backend* tersebut.

1. Buka terminal (Command Prompt / PowerShell / Git Bash).
2. Arahkan terminal ke dalam folder utama proyek ini.
3. Jalankan perintah berikut untuk mengunduh *library* dan menyalakan *server*:

```bash
docker compose up --build

```

4. Proses instalasi awal (khususnya mengunduh model AI) akan memakan waktu sekitar 10-15 menit tergantung kecepatan internet.
5. Sistem siap digunakan jika terminal sudah menampilkan pesan berikut:
* `Container ulasjujur-mysql Healthy`
* `Berhasil memuat 162 kata kunci ke dalam memory cache`
* `Server Golang berjalan di port 8080`



---

## Tahap 3: Panduan Testing Data Besar (Ribuan Baris)

Sistem ini telah dikonfigurasi untuk menangani *file* CSV berskala besar (seperti `ai_csv_PRDECT_Tokopedia_Balanced_8850.csv` yang berisi ribuan baris). Pengujian dapat dilakukan menggunakan **cURL** (via terminal) atau **Postman**.

### Opsi A: Menggunakan cURL (Via Terminal Baru)

Buka tab terminal baru, arahkan ke folder `sample-data`, lalu jalankan:

```bash
curl -X POST http://localhost:8080/api/upload -F "file=@ai_csv_PRDECT_Tokopedia_Balanced_8850.csv"

```

### Opsi B: Menggunakan Postman / API Tester

1. Buka aplikasi Postman.
2. Buat *request* baru dengan metode **POST**.
3. Masukkan URL: `http://localhost:8080/api/upload`
4. Pilih tab **Body**, lalu pilih opsi **form-data**.
5. Pada kolom *Key*, ketik `file` (ubah tipe input di sebelah kanan dari *Text* menjadi *File*).
6. Pada kolom *Value*, unggah file `ai_csv_PRDECT_Tokopedia_Balanced_8850.csv`.
7. Klik **Send**.

### ⚠️ Catatan Penting Saat Testing Data Besar

* Proses analisis untuk lebih dari 8.000 baris menggunakan CPU biasa akan memakan waktu **beberapa menit**.
* Terminal atau Postman akan terlihat dalam status *loading* / *waiting for response*. Ini adalah perilaku yang **normal**.
* Sistem Golang sudah diberikan batas toleransi (*timeout*) yang tinggi agar tidak memutus koneksi di tengah jalan saat model AI sedang bekerja membedah sentimen dan emosi satu per satu.
* Harap **tidak menutup terminal utama** atau menekan *Cancel* pada Postman selama proses ini berlangsung.

---

## Tahap 4: Mematikan Sistem

Jika pengujian sudah selesai dan Anda ingin mematikan sistem agar tidak memakan RAM komputer:

1. Buka kembali terminal tempat Docker berjalan.
2. Tekan tombol `Ctrl + C` secara bersamaan.
3. Tunggu hingga semua kontainer berstatus *Stopped*.

---

Semoga berhasil dengan presentasi dan pengujian *end-to-end* di depan PM-mu! Ada bagian lain dari proyek ini yang perlu kita rapikan sebelum kamu mendemonstrasikannya?