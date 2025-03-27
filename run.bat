@echo off
echo Haber Sitesi Başlatılıyor...

REM Gerekli bağımlılıkları kontrol et
echo Bağımlılıklar kontrol ediliyor...
go mod tidy

REM Ana klasörler oluşturuluyor
if not exist "web\uploads" mkdir web\uploads

REM Uygulamayı çalıştır
echo Uygulama başlatılıyor...
go run cmd/server/main.go

pause 