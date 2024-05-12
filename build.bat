cd .\web
call npm install
call npm run build
cd ..\
rmdir /s/q .\md\web\
xcopy /s/e .\web\dist\ .\md\web\
cd .\md
call go build
echo md build finished
pause