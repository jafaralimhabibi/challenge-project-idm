# Challenges Project

Ini adalah contoh sederhana implementasi dari Echo Golang + MongoDB + Socket.io

Struktur folder yang saya biasa desain sebagai berikut:
|   config.yml
|   go.mod
|   go.sum
|   main.go
|   README.md
|
+---config
|       config.go
|       mongodb.go
|
+---controller
|   +---auth
|   |       auth_controller.go
|   |
|   \---products
|           sampleproduct_controller.go
|
+---lib
|       response.go
|
+---middleware
|       middleware.go
|
+---model
|   +---products
|   |       simpleproduct_model.go
|   |
|   \---users
|           user_model.go
|
+---router
|   +---auth
|   |       auth_router.go
|   |
|   \---products
|           sampleproducts_router.go
|
\---services
    +---scheduler
    \---socket
            socket.go

## REST API dan Socket.io

REST API yang berjalan di project ini adalah sebagai berikut:

User segment (authentication)
1. /auth/login (method: POST): bertujuan untuk menjalankan fungsi login ketika user akan login, dengan body request adalah username dan password, dan response yang dihasilkan adalah token dan object user dengan id dan username berada di dalamnya.

2. /auth/register (method: POST): bertujuan untuk menjalankan fungsi registrasi pelanggan yg akan login, dengan request berupa username dan password, yang nantinya akan memberikan response id dan username.

Product segment (CRUD)
1. /api/product/get-list/all/v1 (method: GET): bertujuan untuk memanggil API dari keseluruhan data product yg tersimpan di collection products di database MongoDB, meliputi _id, name, category, description, price, amount, created_at, dan updated_at.

2. /api/product/view/:id/v1 (method: GET): bertujuan untuk memanggil API dari id spesifik yg ada di collection product, dengan id tertentu.

3. /api/product/create/v1 (method: POST): bertujuan untuk memasukkan data yg dikirimkan melalui API URI ini dengan request body berupa name, category, description, price, dan amount, dimana nilai created_at dan updated_at di generate dari controller langsung.

4. /api/product/update/:id/v1 (method: PUT): bertujuan untuk memutakhirkan document dengan id tertentu yang dimasukkan ke dalam param :id, dengan request body berupa name, category, description, price, amount, dan updated_at yg disimpan dalam bentuk integer.

5. /product/delete/:id/v1 (method: DELETE): bertujuan untuk menghapus data document yang tersimpan di database, dengan param :id sebagai indikator nilai document yg akan dicari.

Sementara untuk socket.io hanya berfokus pada apapun method yang akan dipakai, dengan /socket.io/ sebagai tujuan utama service yg berjalan sebagai perantara komunikasi.

## Scaffolding folder
Penjelasan folder dan file yang berjalan di project ini adalah sebagai berikut:
1. config: berfungsi untuk menambahkan konfigurasi yg biasa dipanggil di package yang lain.
Di sini saya menggunakan viper sebagai binding nilai yang akan diambil melalui config.yml, terutama ketika nama variable yang akan digunakan mungkin saja berbeda, tergantung dari kondisi build type yang akan dijalankan, semisal adanya environment lain (devs, master), atau microservices yang berbeda beda, dll.

2. controller: berfungsi untuk mengeksekusi fungsi fungsi yang akan digunakan di project ini.
Untuk simpelnya saya masih hanya menggunakan auth dan product saja. Ke depannya, jika diperlukan untuk lebih kompleks lagi dalam satu modul/beberapa modul yg berjalan di project ini, saya bisa memecah kategori dari controller nya berdasarkan modul yg berjalan di project tsb.

3. lib: berfungsi untuk helper jika ada fungsi yang dibutuhkan di controller, dimana fungsi tersebut bisa dipanggil berulang kali tanpa harus membuat fungsi yg sama di setiap controller nya.

4. middleware: berfungsi untuk membuat fungsi middleware yang kemungkinan bisa dimodifikasi di project ini. 
Biasanya framework Echo Golang memiliki dependency middleware yang bisa di install di project tertentu. Namun jika memungkinkan untuk menggunakan middleware untuk keperluan tertentu, biasanya saya menyimpannya di folder ini sebagai package khusus.

5. model: berfungsi sebagai tempat memanggil Data Transfer Object (DTO) yang nantinya bisa dipanggil di controller untuk struct tertentu.
Pada umumnya saya di sini juga memisahkan folder yang ada di dalamnya berdasarkan kategori modul atau sub modul dari sebuah program yang berjalan di sana.

6. router: berfungsi untuk menentukan REST API dengan URI apa yang akan dibuat, dengan pengelompokkan berdasarkan method serta fungsi yg mana yang akan dipanggil dari controller yang sudah ada.
Selain model dan controller, di sini juga saya biasa memisahkan file Go yang berjalan berdasarkan submodul/fitur/referensi tertentu.

7. services: berfungsi untuk menjalankan layanan apa saja yang akan di install di project ini, dengan mempertimbangkan fitur yang akan berjalan di dalamnya.
Di folder ini juga saya menambahkan scheduler sebagai tambahan, jika nantinya ke depannya ada fitur yang berfungsi sebagai scheduler dengan goroutine di dalamnya.

8. config.yml: berfungsi sebagai parameter dan nilai yang akan dipanggil di config atau di package yang lain.
Beberapa ada yang ditentukan nilainya langsung dari dalam file tersebut, beberapa ada yang memungkinkan untuk diubah di platform yg berjalan di server seperti Docker yang nantinya dapat diubah sesuai dengan spesifikasi dari server sendiri (domain, host, port, dll.)

9. go.mod & go.sum: berfungsi untuk menyimpan dependencies yang terinstal di dalam project ini.

10. main.go: berfungsi untuk menjalankan project ini, dengan cara menjalankan syntax "go run main.go" sebagai perintah di dalamnya.