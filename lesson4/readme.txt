
Аплоад файла:
curl -F 'file=@testfile.txt' http://localhost:8001/upload

Список загруженных файлов:
curl -i http://localhost:8001/

Список загруженных файлов с фильтром по расширению:
curl -i http://localhost:8001?ext=txt
