# log_Server

# 起動方法
```
docker-compose up --build
```

# 使用方法
HTTPで通信します
/logsにGETメゾットで送信すると、ログをすべて取得\
/logsにPOSTメゾットでJSONを送信すると、ログを登録
JSONの書き方は、Test.jsonの通り


# 改善点
- ログをすべて取得するため、多分重くなる
- 暗号化