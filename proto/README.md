# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [proto/health/v1/health.proto](#proto_health_v1_health-proto)
    - [CheckRequest](#health-v1-CheckRequest)
    - [CheckResponse](#health-v1-CheckResponse)
  
    - [HealthService](#health-v1-HealthService)
  
- [proto/user/v1/user.proto](#proto_user_v1_user-proto)
    - [CreateRequest](#user-v1-CreateRequest)
    - [CreateResponse](#user-v1-CreateResponse)
    - [UpdateRequest](#user-v1-UpdateRequest)
    - [UpdateResponse](#user-v1-UpdateResponse)
    - [User](#user-v1-User)
  
    - [UserService](#user-v1-UserService)
  
- [proto/article/v1/article.proto](#proto_article_v1_article-proto)
    - [Article](#article-v1-Article)
    - [DeleteRequest](#article-v1-DeleteRequest)
    - [DeleteResponse](#article-v1-DeleteResponse)
    - [ListRequest](#article-v1-ListRequest)
    - [ListResponse](#article-v1-ListResponse)
    - [ShareRequest](#article-v1-ShareRequest)
    - [ShareResponse](#article-v1-ShareResponse)
  
    - [ArticleService](#article-v1-ArticleService)
  
- [Scalar Value Types](#scalar-value-types)



<a name="proto_health_v1_health-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## proto/health/v1/health.proto



<a name="health-v1-CheckRequest"></a>

### CheckRequest
チェックリクエスト






<a name="health-v1-CheckResponse"></a>

### CheckResponse
チェックレスポンス





 

 

 


<a name="health-v1-HealthService"></a>

### HealthService
ヘルスサービス

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Check | [CheckRequest](#health-v1-CheckRequest) | [CheckResponse](#health-v1-CheckResponse) | チェック Need X-Api-Key Header |

 



<a name="proto_user_v1_user-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## proto/user/v1/user.proto



<a name="user-v1-CreateRequest"></a>

### CreateRequest
作成リクエスト






<a name="user-v1-CreateResponse"></a>

### CreateResponse
作成レスポンス


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [User](#user-v1-User) |  |  |






<a name="user-v1-UpdateRequest"></a>

### UpdateRequest
更新リクエスト


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [string](#string) |  |  |






<a name="user-v1-UpdateResponse"></a>

### UpdateResponse
更新レスポンス


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user | [User](#user-v1-User) |  |  |






<a name="user-v1-User"></a>

### User
ユーザーモデル


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| user_id | [string](#string) |  |  |





 

 

 


<a name="user-v1-UserService"></a>

### UserService
ユーザーサービス

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Create | [CreateRequest](#user-v1-CreateRequest) | [CreateResponse](#user-v1-CreateResponse) | 作成 |
| Update | [UpdateRequest](#user-v1-UpdateRequest) | [UpdateResponse](#user-v1-UpdateResponse) | 更新 |

 



<a name="proto_article_v1_article-proto"></a>
<p align="right"><a href="#top">Top</a></p>

## proto/article/v1/article.proto



<a name="article-v1-Article"></a>

### Article
記事モデル


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| article_id | [string](#string) |  |  |
| url | [string](#string) |  |  |
| title | [string](#string) |  |  |
| description | [string](#string) |  |  |
| thumbnail | [string](#string) |  |  |
| tags | [string](#string) | repeated |  |






<a name="article-v1-DeleteRequest"></a>

### DeleteRequest
削除リクエスト


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| article_id | [string](#string) |  |  |






<a name="article-v1-DeleteResponse"></a>

### DeleteResponse
削除レスポンス






<a name="article-v1-ListRequest"></a>

### ListRequest
一覧リクエスト


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| page_token | [string](#string) |  |  |
| max_page_size | [uint32](#uint32) |  |  |






<a name="article-v1-ListResponse"></a>

### ListResponse
一覧レスポンス


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| articles | [Article](#article-v1-Article) | repeated |  |
| next_page_token | [string](#string) |  |  |






<a name="article-v1-ShareRequest"></a>

### ShareRequest
共有リクエスト


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| url | [string](#string) |  |  |
| title | [string](#string) |  |  |
| description | [string](#string) |  |  |
| thumbnail | [string](#string) |  |  |






<a name="article-v1-ShareResponse"></a>

### ShareResponse
共有レスポンス


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| article | [Article](#article-v1-Article) |  |  |





 

 

 


<a name="article-v1-ArticleService"></a>

### ArticleService
記事サービス

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Share | [ShareRequest](#article-v1-ShareRequest) | [ShareResponse](#article-v1-ShareResponse) | 共有 Need X-Api-Key Header |
| List | [ListRequest](#article-v1-ListRequest) | [ListResponse](#article-v1-ListResponse) | 一覧 |
| Delete | [DeleteRequest](#article-v1-DeleteRequest) | [DeleteResponse](#article-v1-DeleteResponse) | 削除 |

 



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

