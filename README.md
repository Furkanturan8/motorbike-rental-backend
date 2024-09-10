# Motorbike Rental Backend

Bu proje, motorbike kiralama sistemi için geliştirilmiş bir backend uygulamasıdır. Kullanıcılar motorbike kiralayabilir, harita üzerinden motorbike'leri görüntüleyebilir ve bluetooth bağlantısı ile motorbike'leri yönetebilirler.

## İçindekiler
1. [Gereksinimler](#gereksinimler)
2. [Kullanılan Teknolojiler](#kullanılan-teknolojiler)
3. [Proje Kurulumu](#proje-kurulumu)
   - [Depoyu Klonlayın](#depoyu-klonlayın)
   - [Bağımlılıkları Yükleyin](#bağımlılıkları-yükleyin)
   - [Çevresel Değişkenleri Yapılandırın](#çevresel-değişkenleri-yapılandırın)
4. [Docker Veritabanı İşlemleri](#docker-veritabanı-işlemleri)
5. [Migration İşlemleri](#migration-işlemleri)
6. [Projenin Derlenmesi ve Çalıştırılması](#projenin-derlenmesi-ve-çalıştırılması)
7. [API Endpoint'leri](#api-endpointleri)
   - [Kullanıcı ve Kimlik Doğrulama](#kullanıcı-ve-kimlik-doğrulama)
   - [Motorbike İşlemleri](#motorbike-işlemleri)
   - [Sürüş İşlemleri](#sürüş-işlemleri)
   - [Harita İşlemleri](#harita-işlemleri)
   - [Bluetooth Bağlantı İşlemleri](#bluetooth-bağlantı-işlemleri)


## Gereksinimler

- GoLang
- Docker
- PostgreSQL (Docker üzerinden çalıştırılacak)

## Kullanılan Teknolojiler

- **Go**: Backend servisi için kullanılan dil.
- **Fiber**: HTTP web framework.
- **GORM**: ORM kütüphanesi.
- **PostgreSQL**: Veritabanı.
- **JWT**: Kimlik doğrulama için JSON Web Token.
- **Docker**: PostgreSQL veritabanını konteyner ortamında çalıştırmak için kullanılan araç.


## Proje Kurulumu

1. **Depoyu klonlayın:**

    ```bash
    git clone https://github.com/Furkanturan8/motorbike-rental-backend.git
    cd motorbike-rental-backend
    ```

2. **Bağımlılıkları yükleyin:**

    ```bash
    go mod tidy
    ```

3. **Çevresel değişkenleri yapılandırın:**

   Projenizle birlikte bir `.env` dosyasına ihtiyacınız olacak. Aşağıdaki bilgileri doldurup `.env` dosyasını oluşturabilirsiniz:

```
DB_DOCKER_CONTAINER=my_postgres_container
DB_USERNAME=your_db_username
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
DB_HOST=localhost
DB_PORT=5432
SERVER_PORT=8080
DSN=postgres://your_db_username:your_db_password@localhost:5432/your_db_name?sslmode=disable
BINARY_NAME=motorbike-rental-backend
```

## Docker Veritabanı İşlemleri

Veritabanı işlemlerini Docker üzerinde gerçekleştirmek için aşağıdaki komutları kullanabilirsiniz.

1. **Diğer Docker Konteynerlerini Durdurmak İçin:**

   ```bash
   make stop_containers
   ```

   Bu komut, çalışan diğer Docker konteynerlerini durdurur.

2. **Veritabanı Konteynerini Oluşturmak İçin:**

   ```bash
   make create_container
   ```

   Bu komut, PostgreSQL veritabanı konteynerini oluşturur ve çalıştırır.

3. **Veritabanı Oluşturmak İçin:**

   ```bash
   make create_db
   ```

   Bu komut, PostgreSQL konteyneri içinde belirtilen kullanıcı ve veritabanı adıyla yeni bir veritabanı oluşturur.

4. **Konteyneri Başlatmak İçin (Daha Önceden Oluşturulduysa):**

   ```bash
   make start_container
   ```

   Bu komut, daha önce oluşturulmuş bir PostgreSQL konteynerini başlatır.

## Migration İşlemleri

Veritabanı migration işlemlerini SQLX kullanarak gerçekleştirebilirsiniz.

1. **Yeni Bir Migration Oluşturmak İçin:**

   ```bash
   make create_migrations
   ```

   Bu komut, `sqlx migrate` komutunu kullanarak yeni bir migration dosyası ekler.

2. **Migration Yüklemek İçin (migrate up):**

   ```bash
   make migrate_up
   ```

   Bu komut, veritabanı üzerinde tanımlı migration'ları çalıştırır.

3. **Migration Geri Almak İçin (migrate down):**

   ```bash
   make migrate_down
   ```

   Bu komut, en son yapılan migration'ı geri alır.

## Projenin Derlenmesi ve Çalıştırılması

1. **Binary Oluşturmak İçin:**

   ```bash
   make build
   ```

   Bu komut, projeyi derleyerek belirtilen binary dosyasını oluşturur.

2. **API Sunucusunu Başlatmak İçin:**

   ```bash
   make start
   ```

   Bu komut, veritabanı konteynerini başlatır ve API sunucusunu çalıştırır.

3. **API Sunucusunu Durdurmak İçin:**

   ```bash
   make stop
   ```

   Bu komut, çalışan API sunucusunu durdurur.

4. **API Sunucusunu Yeniden Başlatmak İçin:**

   ```bash
   make restart
   ```

   Bu komut, önce sunucuyu durdurur, ardından yeniden başlatır.


## API Endpoint'leri

Aşağıda uygulamada kullanılan temel API endpoint'leri verilmiştir.

### Kullanıcı ve Kimlik Doğrulama

| Method | Endpoint            | Açıklama                              |
|--------|---------------------|---------------------------------------|
| POST   | `/api/user/create`   | Yeni bir kullanıcı oluşturur.         |
| POST   | `/api/auth/login`    | Kullanıcı giriş işlemi.               |
| POST   | `/api/auth/refresh`  | JWT token'ını yeniler.                |
| GET    | `/api/user/me`       | Giriş yapmış kullanıcının bilgilerini alır. |
| PUT    | `/api/user/me`       | Giriş yapmış kullanıcının bilgilerini günceller. |
| POST   | `/api/auth/logout`   | Kullanıcı çıkış işlemi.               |

### Motorbike İşlemleri

| Method  | Endpoint                         | Açıklama                                  |
|---------|---------------------------------- |-------------------------------------------|
| POST    | `/api/motorbike`                 | Yeni bir motorbike ekler.                 |
| PUT     | `/api/motorbike/:id`             | Bir motorbike'in bilgilerini günceller.   |
| DELETE  | `/api/motorbike/:id`             | Bir motorbike'i siler.                    |
| GET     | `/api/motorbikes`                | Tüm motorbike'leri getirir.               |
| GET     | `/api/motorbikes/:id`            | Belirli bir motorbike'i getirir.          |
| GET     | `/api/available-motorbikes`      | Kiralanabilir motorbike'leri getirir.     |
| GET     | `/api/maintenance-motorbikes`    | Bakımda olan motorbike'leri getirir.      |
| GET     | `/api/rented-motorbikes`         | Kiralanmış motorbike'leri getirir.        |
| GET     | `/api/motorbike-photos/:id`      | Belirli motorbike'in fotoğraflarını getirir. |

### Sürüş İşlemleri

| Method  | Endpoint                                       | Açıklama                                      |
|---------|------------------------------------------------|-----------------------------------------------|
| POST    | `/api/ride`                                    | Yeni bir sürüş başlatır.                      |
| GET     | `/api/rides`                                   | Tüm sürüşleri getirir.                        |
| GET     | `/api/rides/:id`                               | Belirli bir sürüşü getirir.                   |
| GET     | `/api/rides/user/:userID`                      | Belirli bir kullanıcıya ait sürüşleri getirir.|
| PUT     | `/api/ride/update/:id`                         | Belirli bir sürüşü günceller.                 |
| PUT     | `/api/ride/finish/:id`                         | Bir sürüşü tamamlar.                          |
| DELETE  | `/api/ride/:id`                                | Bir sürüşü siler.                             |
| GET     | `/api/rides/user/:userID/filter?start_time=...`| Tarih aralığına göre kullanıcı sürüşleri getirir.|
| GET     | `/api/motorbike/:bikeID/rides`                 | Belirli bir motorbike'e ait sürüşleri getirir.|

### Harita İşlemleri

| Method  | Endpoint                               | Açıklama                                 |
|---------|----------------------------------------|------------------------------------------|
| POST    | `/api/map`                             | Yeni bir harita ekler.                   |
| DELETE  | `/api/map/:id`                         | Belirli bir haritayı siler.              |
| GET     | `/api/maps`                            | Tüm haritaları getirir.                  |
| GET     | `/api/maps/:id`                        | Belirli bir haritayı getirir.            |
| GET     | `/api/motorbikes/:motorbikeID/map`     | Motorbike'e ait haritayı getirir.        |
| PUT     | `/api/map/update/:id`                  | Harita bilgilerini günceller.            |
| PUT     | `/api/motorbikes/:motorbikeID/map/update` | Motorbike'e ait haritayı günceller.   |

### Bluetooth Bağlantı İşlemleri

| Method  | Endpoint                                    | Açıklama                                  |
|---------|---------------------------------------------|-------------------------------------------|
| POST    | `/api/connection/connect`                   | Motorbike ile bluetooth bağlantısı kurar. |
| POST    | `/api/connection/disconnect/:id`            | Bluetooth bağlantısını keser.             |
| GET     | `/api/connections`                          | Tüm bağlantıları getirir.                 |
| GET     | `/api/connection/:id`                       | Belirli bir bağlantıyı getirir.           |
| DELETE  | `/api/connection/:id`                       | Bir bağlantıyı siler.                     |
| GET     | `/api/connection/motorbike/:motorbikeID`    | Belirli motorbike'e ait bağlantıları getirir. |
| GET     | `/api/connection/user/:userID`              | Belirli kullanıcıya ait bağlantıları getirir. |


---

