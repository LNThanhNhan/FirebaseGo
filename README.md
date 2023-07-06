# Webserver tích hợp các hợp các dịch vụ của Firebase

## Mục đích
&nbsp; &nbsp; Nhằm tiết kiệm thời gian mỗi lần xây dựng backend thì phải xây dựng thêm các dịch vụ cơ bản như Firesbase vào. Thì giờ đây,chỉ cần viết API gọi đến web server này.

## Hiện tại có
+ ## Storage

## Lưu ý!
+ Cần phải tạo các file .env và serviceAccountKey.json để chứa các key từ FirebaseAdmin và setting cho web server.
+ Web server này mặc định chạy trên localhost với cổng 8080.

## Danh mục API


### Tải ảnh lên  
Request
```http
  POST /img/create
```

| Tham số | Kiểu     | Mô tả                |
| :-------- | :------- | :------------------------- |
| `image` | `file` | **(Bắt buộc)** Ảnh cần gửi lên Firebase Storage  |

Response
| Tham số | Kiểu     | Mô tả                |
| :-------- | :------- | :------------------------- |
| `success` | `bool` | Chỉ trạng thái Request thành công hay thất bại |
| `message` | `string` | Tin nhắn trả về từ Response |
| `data` | `object` | Chứa ID ảnh trên Storage |


### Cập nhật ảnh

```http
  PUT /img/update
```
Request
| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `int` |  **(Bắt buộc)** ID ảnh trên Storage|
| `image` | `file` | **(Bắt buộc)** Upload ảnh lên Firebase Storage  |

Response
| Tham số | Kiểu     | Mô tả                |
| :-------- | :------- | :------------------------- |
| `success` | `bool` | Chỉ trạng thái Request thành công hay thất bại |
| `message` | `string` | Tin nhắn trả về từ Response |
| `data` | `object` | Chứa ID ảnh cập nhật trên Storage |


### Xóa ảnh

```http
  DELETE /img/delete
```
Request
| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `id`      | `int` |  **(Bắt buộc)** ID ảnh trên Storage|


Response
| Tham số | Kiểu     | Mô tả                |
| :-------- | :------- | :------------------------- |
| `success` | `bool` | Chỉ trạng thái Request thành công hay thất bại |
| `message` | `string` | Tin nhắn trả về từ Response |
| `data` | `object` | Chứa thông báo thành công  |
