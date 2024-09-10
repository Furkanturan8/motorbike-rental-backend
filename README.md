# Motorbike Rental Backend

Bu proje, motorbike kiralama sistemi için geliştirilmiş bir backend uygulamasıdır. Kullanıcılar motorbike kiralayabilir, harita üzerinden motorbike'leri görüntüleyebilir ve bluetooth bağlantısı ile motorbike'leri yönetebilirler.

## İçindekiler
1. [Gereksinimler](#gereksinimler)
2. [Kullanılan Teknolojiler](#kullanılan-teknolojiler)
3. [API Endpoint'leri](#api-endpointleri)
   - [Kullanıcı ve Kimlik Doğrulama](#kullanıcı-ve-kimlik-doğrulama)
   - [Admin'in User ile İlgili İşlemleri](#adminin-user-ile-ilgili-işlemleri)
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

### Admin'in User Ile Ilgili Işlemleri

| Method | Endpoint            | Açıklama                              |
|--------|---------------------|---------------------------------------|
| GET   | `/api/users`   | Tüm kullanıcıları getirir.         |
| GET   | `/api/users/:id`    | Belirli bir kullanıcıyı getirir.               |
| POST   | `/api/user/create`  | Yeni bir kullanıcı ekler.                |
| POST    | `/api/user/createAdmin`       | Yeni bir admin ekler. |
| DELETE    | `/api/user/:id`       | Kullanıcıyı siler. |
| PUT   | `/api/user/update/:id`   | Kullanıcı bilgilerini günceller.               |

### Motorbike Işlemleri

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

### Sürüş Işlemleri

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

### Harita Işlemleri

| Method  | Endpoint                               | Açıklama                                 |
|---------|----------------------------------------|------------------------------------------|
| POST    | `/api/map`                             | Yeni bir harita ekler.                   |
| DELETE  | `/api/map/:id`                         | Belirli bir haritayı siler.              |
| GET     | `/api/maps`                            | Tüm haritaları getirir.                  |
| GET     | `/api/maps/:id`                        | Belirli bir haritayı getirir.            |
| GET     | `/api/motorbikes/:motorbikeID/map`     | Motorbike'e ait haritayı getirir.        |
| PUT     | `/api/map/update/:id`                  | Harita bilgilerini günceller.            |
| PUT     | `/api/motorbikes/:motorbikeID/map/update` | Motorbike'e ait haritayı günceller.   |

### Bluetooth Bağlantı Işlemleri

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

