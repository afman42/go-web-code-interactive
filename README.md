Sama seperti ada di CodeWars,LeetCode,dll

# CODE INTERACTIVE WEB

### Kebutuhan

- Install NodeJS terbaru
- Install Golang terbaru
- Install CMAKE
- Install PHP
- Install UPX (compress file)

### Development

- Masuk root folder
- Menjalankan Golang `make run/api`
- Menjalankan Web `make run/web`
- Install Multiple lib di package.json `make npmi i="lib lib"` || Install Semua `make npmi i=""` || Install Satu `make npmi i=lib`
- Uninstall Multiple lib di package.json `make npmu u="lib lib"` || Uninstall satu `make npmu u=lib`

### Preview Setelah Build

- Menjalankan di linux: `make run/preview_linux`
- Menjalankan di windows: `not testing`

### Production

- buat file `.env.prod` di golang dan file `.env.production` di web
- buat variabel `export` di cmd `deploy:` di `.bashrc` lalu `source path/.bashrc` di linux
- windows `not testing`
- Menjalankan di os ubuntu `make deploy/prod` atau Menjalankan server di os alpine `make deploy/prod_alpine`

### Kendala

- PATH NodeJS executable di windows dan linux ketika production
- Atur path di urutan nomor 65 di file `main.go`. global: `node`. relative: `path node.exe`.
- Atur Cors di bagian file `main.go` dan atur port di bagian file `.env` folder web

### Todo

- [x] Buat flag parse di `main.go`
- [x] Embed file html dan folder dist di file `main.go`
- [x] File Docker dan Masukin Web
- [x] Log all routes di web
- [x] Make Web Responsive mobile and dekstop
- Need change code editor like vscode ? Ribet

![IMG_PROD](images/WEB.PNG "Title")

![IMG_PROD](images/IMG_PROD.PNG "Title")

![IMG_PROD](images/after-prod-img.jpg "Title")
