# CASHFLOW-BE


## Prasyarat

Pastikan Anda telah menginstal:
- [Go](https://golang.org/) (versi 1.21 atau lebih tinggi)
- [Docker](https://www.docker.com/) dan Docker Compose (untuk development dan deployment)
- [Make](https://www.gnu.org/software/make/) (opsional, untuk menjalankan Makefile commands)

## Setup Projek Lokal

### 1. Clone Repository
```bash
git clone https://github.com/kenziehh/cashflow-be
cd cashflow-be
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Setup Environment Variables

Salin file `.env.example` menjadi `.env`:
```bash
cp .env.example .env
```

Kemudian edit file `.env` dan sesuaikan nilai-nilai konfigurasi sesuai kebutuhan Anda:
```bash
nano .env
# atau gunakan editor favorit Anda
```



### 4. Jalankan Development Environment

Gunakan `docker-compose.dev.yml` untuk development yang mendukung hot-reload:
```bash
docker-compose -f docker-compose.dev.yml up
```

**Fitur Development Mode:**
- ✅ Hot-reload otomatis saat code berubah
- ✅ Volume mounting untuk development
- ✅ Debug logging enabled
- ✅ Database dan services pendukung lainnya

### 5. Jalankan di Background
```bash
docker-compose -f docker-compose.dev.yml up -d
```

## Deployment Production dengan Docker

### Prasyarat Deployment

Pastikan Anda telah:
1. Menginstal Docker dan Docker Compose di server production
2. Menyiapkan file `.env` dengan konfigurasi production
3. Membuild image atau memiliki akses ke image registry

### Cara Deploy Production

#### 1. Persiapan File Environment

Pastikan file `.env` sudah dikonfigurasi dengan benar untuk production:
```bash
cp .env.example .env
# Edit .env dengan konfigurasi production.
nano .env
```

#### 2. Build dan Jalankan Container Production

Jalankan perintah berikut untuk memulai deployment:
```bash
docker-compose -f docker-compose.prod.yml up -d
```

**Penjelasan flag:**
- `-f docker-compose.prod.yml`: Menggunakan file konfigurasi production
- `up`: Membuat dan menjalankan container
- `-d`: Menjalankan container di background (detached mode)

#### 3. Verifikasi Deployment

Periksa status container:
```bash
docker-compose -f docker-compose.prod.yml ps
```


