cd .\web
npm install
npm run build
cd ..\
Remove-Item -Path "md\web\" -Recurse -Force
xcopy /s/e .\web\dist\ .\md\web\
