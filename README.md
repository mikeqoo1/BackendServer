# BackendServer
一個通用的Server 雛型

## RESTful API

```
GET /user ：列出所有user
POST /user：新建一個user
GET /user/ID：獲取某個指定user的信息
PUT /user/ID：更新某個指定user的信息（提供該user的全部信息）
PATCH /user/ID：更新某個指定user的信息（提供該user的部分信息）
DELETE /user/ID：刪除某個user
```

## 狀態碼

```
2xx = Success（成功）
3xx = Redirect（重定向）
4xx = User error（客戶端錯誤）
5xx = Server error（伺服器端錯誤）
```

- 200 OK：成功通常我用在查詢（GET）的部分
- 201 Created：資源新增成功通常我用在新增（POST）的部分
- 202 Accepted：請求已接受, 但還在處理中（換句話說可能失敗）
- 204 No Content：請求成功, 沒有返回內容通常我用在刪除（DELETE）或修改（PUT）的部分, 也有人會把修改成功放在 201
- 400 Bad Request：使用者做錯的通用狀態通常我用在有必填欄位未填或資料錯誤的狀況
- 401 Unauthorized：使用者沒有進行驗證
- 403 Forbidden：使用者已經登入, 但沒有權限進行操作
- 404 Not Found：請求的資源不存在
- 410 Gone：資源已經過期
- 500 Internal Server Error：伺服器端錯誤
- 502 Bad Gateway：通常是伺服器的某個服務沒有正確執行
- 503 Service Unavailable：伺服器臨時維護或是快掛了, 暫時無法處理請求
- 504 Gateway Timeout：伺服器上的服務沒有回應