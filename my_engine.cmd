@echo off
setlocal

:: Путь к bee2evp
set BEE2_PATH=.\bpki\demo\
set OPENSSL_CONF=%BEE2_PATH%\openssl.cnf
set OPENSSL_BIN=%BEE2_PATH%

:: Перейти в каталог скрипта
cd /d %~dp0

:: Генерация ключа
"%OPENSSL_BIN%\openssl.exe" genpkey -engine bee2evp -algorithm bign -pkeyopt params:bign-curve256v1 -out private_key.pem

:: Генерация публичного ключа
"%OPENSSL_BIN%\openssl.exe" pkey -engine bee2evp -in private_key.pem -pubout -out public_key.pem  

:: Генерация сертификата
"%OPENSSL_BIN%\openssl.exe" req -engine bee2evp -new -x509 -key private_key.pem -out certificate.pem -config cert.conf

:: CMS-сообщение
"%OPENSSL_BIN%\openssl.exe" cms -engine bee2evp -sign -in message.txt -signer certificate.pem -inkey private_key.pem  -out message.cms -outform PEM

echo ==== Готово ====
pause
