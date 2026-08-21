cara run 
1. DOWNLOAD Docker dulu jika belum ada : 
link : https://www.docker.com/products/docker-desktop/

2. install docker dan jalankan docker
3. buka terminal di folder root yang sama dengan file docker-compose.yml
4. jalankan command berikut : docker compose up --build
5. Tunggu prosesnya hingga semua container berjalan seperti Golang, MySQL, dan Python. Tekan CTRL + C untuk quit.
6. Jalankan command : docker compose up. Untuk menjalankan ulang.
7. Buka terminal baru untuk melakukan testing.
8. Pindah ke folder ke sample-data folder.
9. Jalankan command berikut untuk mengirim data csv ke server golang yang menyala. 
    curl -X POST http://localhost:8080/api/upload -F "file=@test.csv"
10. Ganti dengan data yang lebih besar 
    curl -X POST http://localhost:8080/api/upload -F "file=@ai_csv_Augmented.csv"
11. Jika output keluar berupa JSON, testing berhasil.