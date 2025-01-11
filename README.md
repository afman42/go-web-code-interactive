Sama seperti ada di CodeWars,LeetCode,dll

# CODE INTERACTIVE WEB

### Kebutuhan

- Install NodeJS terbaru
- Install Golang terbaru
- Install CMAKE

### Development

- Masuk root folder
- Menjalankan Golang `make run/api`
- Menjalankan Web `make run/web`

### Preview Setelah Build

- Menjalankan di linux: `preview\api_linux`
- Menjalankan web: `preview\dist_web`
- Menjalankan di windows: `not testing`

### Production

- Menjalankan `make build`

### Kendala

- PATH NodeJS executable di windows ketika production
- Atur path di urutan nomor 65 di file `main.go`. global: `node`. relative: `path node.exe`.
- Atur Cors di bagian file `main.go` dan atur port di bagian file `.env` folder web

### Todo

- [] Buat flag parse di `main.go`
- [] Embed file html dan folder dist di file `main.go`

![IMG_PROD](WEB.PNG "Title")

![IMG_PROD](IMG_PROD.PNG "Title")
