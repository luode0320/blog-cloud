cd .\web
npm install
npm run build
cd ..\
rmdir /s/q .\md\web\
xcopy /s/e .\web\dist\ .\md\web\
cd .\md
go build
echo md build finished
pause