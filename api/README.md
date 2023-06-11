<!-- Generator: Widdershins v4.0.1 -->

<h1 id="morning-night-guild-app-api">Morning Night Guild - App API v0.0.1</h1>

> Scroll down for code samples, example requests and responses. Select a language for code samples from the tabs above or the mobile navigation menu.

This is the AppAPI API documentation.

Base URLs:

* <a href="http://localhost:8080/api">http://localhost:8080/api</a>

<a href="https://example.com">Terms of service</a>
Email: <a href="mailto:morning.night.guild@example.com">Support</a> 
 License: MIT

# Authentication

* API Key (apiKey)
    - Parameter Name: **api-key**, in: header. 

* API Key (authTokenCookie)
    - Parameter Name: **auth-token**, in: cookie. 

* API Key (sessionTokenCookie)
    - Parameter Name: **session-token**, in: cookie. 

<h1 id="morning-night-guild-app-api-auth">auth</h1>

認証

## v1AuthInvite

<a id="opIdv1AuthInvite"></a>

`POST /v1/auth/invite`

*招待*

ユーザーを招待する
招待コードはメールアドレスに送信される

> Body parameter

```json
{
  "email": "morning.night.guild@example.com"
}
```

<h3 id="v1authinvite-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[V1AuthInviteRequestSchema](#schemav1authinviterequestschema)|true|招待リクエストボディ|

> Example responses

> 200 Response

```json
{
  "code": "xxxxxxxx"
}
```

<h3 id="v1authinvite-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[V1AuthInviteResponseSchema](#schemav1authinviteresponseschema)|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|None|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthorized|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
apiKey
</aside>

## v1AuthJoin

<a id="opIdv1AuthJoin"></a>

`POST /v1/auth/join`

*参加*

招待コードを用いてサインアップする

> Body parameter

```json
{
  "code": "xxxxxxxx",
  "password": "password"
}
```

<h3 id="v1authjoin-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[V1AuthJoinRequestSchema](#schemav1authjoinrequestschema)|true|参加リクエストボディ|

> Example responses

> 200 Response

```json
{
  "articles": [
    {
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "url": "https://example.com",
      "title": "title",
      "description": "description",
      "thumbnail": "https://example.com",
      "tags": [
        "tag"
      ]
    }
  ],
  "nextPageToken": "string"
}
```

<h3 id="v1authjoin-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|[V1ArticleListResponseSchema](#schemav1articlelistresponseschema)|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not Found|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="success">
This operation does not require authentication
</aside>

## v1AuthSignUp

<a id="opIdv1AuthSignUp"></a>

`POST /v1/auth/signup`

*サインアップ(テスト用)*

ユーザーを登録する

> Body parameter

```json
{
  "email": "morning.night.guild@example.com",
  "password": "password"
}
```

<h3 id="v1authsignup-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[V1AuthSignUpRequestSchema](#schemav1authsignuprequestschema)|true|サインアップリクエストボディ|

<h3 id="v1authsignup-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|None|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|None|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthorized|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
apiKey
</aside>

## v1AuthSignIn

<a id="opIdv1AuthSignIn"></a>

`POST /v1/auth/signin`

*サインイン*

ユーザーを認証する

> Body parameter

```json
{
  "email": "morning.night.guild@example.com",
  "password": "password",
  "publicKey": "string",
  "expiresIn": 3600
}
```

<h3 id="v1authsignin-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[V1AuthSignInRequestSchema](#schemav1authsigninrequestschema)|true|サインインリクエストボディ|

<h3 id="v1authsignin-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|None|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="success">
This operation does not require authentication
</aside>

## v1AuthVerify

<a id="opIdv1AuthVerify"></a>

`GET /v1/auth/verify`

*検証*

検証を行う

> Example responses

> 401 Response

```json
{
  "code": "00000000-0000-0000-0000-000000000000"
}
```

<h3 id="v1authverify-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|None|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthorized|[V1AuthVerifyUnauthorizedResponseSchema](#schemav1authverifyunauthorizedresponseschema)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
authTokenCookie, sessionTokenCookie
</aside>

## v1AuthRefresh

<a id="opIdv1AuthRefresh"></a>

`GET /v1/auth/refresh`

*リフレッシュ*

セッショントークンとクライアント署名により認証トークンを再発行する

<h3 id="v1authrefresh-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|code|query|string|true|署名付きコード|
|signature|query|string|true|署名|
|expiresIn|query|integer|false|none|

<h3 id="v1authrefresh-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|None|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthorized|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
sessionTokenCookie
</aside>

## v1AuthSignOut

<a id="opIdv1AuthSignOut"></a>

`GET /v1/auth/signout`

*サインアウト*

サインアウトする

<h3 id="v1authsignout-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|None|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthorized|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
authTokenCookie, sessionTokenCookie
</aside>

## v1AuthSignOutAll

<a id="opIdv1AuthSignOutAll"></a>

`GET /v1/auth/signout/all`

*サインアウトオール*

該当ユーザーのすべてのセッションにてサインアウトする

<h3 id="v1authsignoutall-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|None|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthorized|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
authTokenCookie, sessionTokenCookie
</aside>

## v1AuthChangePassword

<a id="opIdv1AuthChangePassword"></a>

`PUT /v1/auth/password`

*パスワード変更*

パスワードを変更する。認証トークンとセッショントークンは新たに発行される。（これまで発行されているトークンは無効化される）

> Body parameter

```json
{
  "oldPassword": "OldPassword",
  "newPassword": "NewPassword",
  "publicKey": "string",
  "expiresIn": 3600
}
```

<h3 id="v1authchangepassword-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[V1AuthChangePasswordRequestSchema](#schemav1authchangepasswordrequestschema)|true|パスワード更新リクエストボディ|

<h3 id="v1authchangepassword-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|None|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|None|
|401|[Unauthorized](https://tools.ietf.org/html/rfc7235#section-3.1)|Unauthorized|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
authTokenCookie, sessionTokenCookie
</aside>

<h1 id="morning-night-guild-app-api-article">article</h1>

記事

## v1ArticleList

<a id="opIdv1ArticleList"></a>

`GET /v1/articles`

*記事一覧*

記事一覧を取得する

<h3 id="v1articlelist-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|scope|query|string|false|取得範囲|
|maxPageSize|query|integer|false|ページサイズ|
|pageToken|query|string|false|トークン|
|title|query|string|false|タイトルによる部分一致検索|

#### Enumerated Values

|Parameter|Value|
|---|---|
|scope|all|
|scope|own|

> Example responses

> 200 Response

```json
{
  "articles": [
    {
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "url": "https://example.com",
      "title": "title",
      "description": "description",
      "thumbnail": "https://example.com",
      "tags": [
        "tag"
      ]
    }
  ],
  "nextPageToken": "string"
}
```

<h3 id="v1articlelist-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|[V1ArticleListResponseSchema](#schemav1articlelistresponseschema)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
authTokenCookie, sessionTokenCookie
</aside>

## v1ArticleShare

<a id="opIdv1ArticleShare"></a>

`POST /v1/articles`

*記事共有*

記事を共有する

> Body parameter

```json
{
  "url": "https://example.com",
  "title": "title",
  "description": "description",
  "thumbnail": "https://example.com"
}
```

<h3 id="v1articleshare-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[V1ArticleShareRequestSchema](#schemav1articlesharerequestschema)|true|記事共有リクエストボディ|

<h3 id="v1articleshare-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|None|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|None|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Forbidden|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
apiKey
</aside>

## v1ArticleAddOwn

<a id="opIdv1ArticleAddOwn"></a>

`POST /v1/articles/{articleId}`

*記事追加*

操作者が管理する記事として追加する

<h3 id="v1articleaddown-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|articleId|path|string(uuid)|true|記事ID|

<h3 id="v1articleaddown-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|None|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|None|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Forbidden|None|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not Found|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
authTokenCookie, sessionTokenCookie
</aside>

## v1ArticleRemoveOwn

<a id="opIdv1ArticleRemoveOwn"></a>

`DELETE /v1/articles/{articleId}`

*記事削除*

操作者が管理する記事から削除する

<h3 id="v1articleremoveown-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|articleId|path|string(uuid)|true|記事ID|

<h3 id="v1articleremoveown-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|None|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|None|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Forbidden|None|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not Found|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
authTokenCookie, sessionTokenCookie
</aside>

## v1InternalArticleShare

<a id="opIdv1InternalArticleShare"></a>

`POST /v1/internal/articles`

*記事共有*

記事を共有する

> Body parameter

```json
{
  "url": "https://example.com",
  "title": "title",
  "description": "description",
  "thumbnail": "https://example.com"
}
```

<h3 id="v1internalarticleshare-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[V1ArticleShareRequestSchema](#schemav1articlesharerequestschema)|true|記事共有リクエストボディ|

<h3 id="v1internalarticleshare-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|None|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|None|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Forbidden|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
apiKey
</aside>

## v1InternalArticleDelete

<a id="opIdv1InternalArticleDelete"></a>

`DELETE /v1/internal/articles/{articleId}`

*記事削除*

記事を削除する

<h3 id="v1internalarticledelete-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|articleId|path|string(uuid)|true|記事ID|

<h3 id="v1internalarticledelete-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功、指定した記事IDが存在しない場合も成功扱いとなる|None|
|400|[Bad Request](https://tools.ietf.org/html/rfc7231#section-6.5.1)|Bad Request|None|
|403|[Forbidden](https://tools.ietf.org/html/rfc7231#section-6.5.3)|Forbidden|None|
|404|[Not Found](https://tools.ietf.org/html/rfc7231#section-6.5.4)|Not Found|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="warning">
To perform this operation, you must be authenticated by means of one of the following methods:
apiKey
</aside>

<h1 id="morning-night-guild-app-api-health">health</h1>

ヘルスチェック

## v1HealthAPI

<a id="opIdv1HealthAPI"></a>

`GET /v1/health/api`

*apiヘルスチェック*

ヘルスチェック

<h3 id="v1healthapi-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="success">
This operation does not require authentication
</aside>

## v1HealthCore

<a id="opIdv1HealthCore"></a>

`GET /v1/health/core`

*coreヘルスチェック*

ヘルスチェック

<h3 id="v1healthcore-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|OK|None|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="success">
This operation does not require authentication
</aside>

# Schemas

<h2 id="tocS_ArticleSchema">ArticleSchema</h2>
<!-- backwards compatibility -->
<a id="schemaarticleschema"></a>
<a id="schema_ArticleSchema"></a>
<a id="tocSarticleschema"></a>
<a id="tocsarticleschema"></a>

```json
{
  "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
  "url": "https://example.com",
  "title": "title",
  "description": "description",
  "thumbnail": "https://example.com",
  "tags": [
    "tag"
  ]
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|id|string(uuid)|false|none|id|
|url|string(uri)|false|none|記事のURL|
|title|string|false|none|タイトル|
|description|string|false|none|description|
|thumbnail|string(uri)|false|none|サムネイルのURL|
|tags|[string]|false|none|タグ|

<h2 id="tocS_V1AuthInviteRequestSchema">V1AuthInviteRequestSchema</h2>
<!-- backwards compatibility -->
<a id="schemav1authinviterequestschema"></a>
<a id="schema_V1AuthInviteRequestSchema"></a>
<a id="tocSv1authinviterequestschema"></a>
<a id="tocsv1authinviterequestschema"></a>

```json
{
  "email": "morning.night.guild@example.com"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|email|string(email)|true|none|メールアドレス|

<h2 id="tocS_V1AuthInviteResponseSchema">V1AuthInviteResponseSchema</h2>
<!-- backwards compatibility -->
<a id="schemav1authinviteresponseschema"></a>
<a id="schema_V1AuthInviteResponseSchema"></a>
<a id="tocSv1authinviteresponseschema"></a>
<a id="tocsv1authinviteresponseschema"></a>

```json
{
  "code": "xxxxxxxx"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|code|string|true|none|招待コード|

<h2 id="tocS_V1AuthJoinRequestSchema">V1AuthJoinRequestSchema</h2>
<!-- backwards compatibility -->
<a id="schemav1authjoinrequestschema"></a>
<a id="schema_V1AuthJoinRequestSchema"></a>
<a id="tocSv1authjoinrequestschema"></a>
<a id="tocsv1authjoinrequestschema"></a>

```json
{
  "code": "xxxxxxxx",
  "password": "password"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|code|string|true|none|招待コード|
|password|string|true|none|パスワード|

<h2 id="tocS_V1AuthSignUpRequestSchema">V1AuthSignUpRequestSchema</h2>
<!-- backwards compatibility -->
<a id="schemav1authsignuprequestschema"></a>
<a id="schema_V1AuthSignUpRequestSchema"></a>
<a id="tocSv1authsignuprequestschema"></a>
<a id="tocsv1authsignuprequestschema"></a>

```json
{
  "email": "morning.night.guild@example.com",
  "password": "password"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|email|string(email)|true|none|メールアドレス|
|password|string|true|none|パスワード|

<h2 id="tocS_V1AuthSignInRequestSchema">V1AuthSignInRequestSchema</h2>
<!-- backwards compatibility -->
<a id="schemav1authsigninrequestschema"></a>
<a id="schema_V1AuthSignInRequestSchema"></a>
<a id="tocSv1authsigninrequestschema"></a>
<a id="tocsv1authsigninrequestschema"></a>

```json
{
  "email": "morning.night.guild@example.com",
  "password": "password",
  "publicKey": "string",
  "expiresIn": 3600
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|email|string(email)|true|none|メールアドレス|
|password|string|true|none|パスワード|
|publicKey|string(base64)|true|none|公開鍵|
|expiresIn|integer|false|none|トークン有効期限(秒)|

<h2 id="tocS_V1AuthVerifyUnauthorizedResponseSchema">V1AuthVerifyUnauthorizedResponseSchema</h2>
<!-- backwards compatibility -->
<a id="schemav1authverifyunauthorizedresponseschema"></a>
<a id="schema_V1AuthVerifyUnauthorizedResponseSchema"></a>
<a id="tocSv1authverifyunauthorizedresponseschema"></a>
<a id="tocsv1authverifyunauthorizedresponseschema"></a>

```json
{
  "code": "00000000-0000-0000-0000-000000000000"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|code|string(uuid)|true|none|リフレッシュコード<br>このコードを使用してトークンを新たに取得することができます。<br>リフレッシュできる見込みがない場合(セッショントークンがない状態でのリクエスト)ではリフレッシュ用コードは払い出しません。|

<h2 id="tocS_V1AuthChangePasswordRequestSchema">V1AuthChangePasswordRequestSchema</h2>
<!-- backwards compatibility -->
<a id="schemav1authchangepasswordrequestschema"></a>
<a id="schema_V1AuthChangePasswordRequestSchema"></a>
<a id="tocSv1authchangepasswordrequestschema"></a>
<a id="tocsv1authchangepasswordrequestschema"></a>

```json
{
  "oldPassword": "OldPassword",
  "newPassword": "NewPassword",
  "publicKey": "string",
  "expiresIn": 3600
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|oldPassword|string|true|none|旧パスワード|
|newPassword|string|true|none|新パスワード|
|publicKey|string(base64)|true|none|公開鍵|
|expiresIn|integer|false|none|トークン有効期限(秒)|

<h2 id="tocS_V1ArticleListResponseSchema">V1ArticleListResponseSchema</h2>
<!-- backwards compatibility -->
<a id="schemav1articlelistresponseschema"></a>
<a id="schema_V1ArticleListResponseSchema"></a>
<a id="tocSv1articlelistresponseschema"></a>
<a id="tocsv1articlelistresponseschema"></a>

```json
{
  "articles": [
    {
      "id": "497f6eca-6276-4993-bfeb-53cbbbba6f08",
      "url": "https://example.com",
      "title": "title",
      "description": "description",
      "thumbnail": "https://example.com",
      "tags": [
        "tag"
      ]
    }
  ],
  "nextPageToken": "string"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|articles|[[ArticleSchema](#schemaarticleschema)]|false|none|none|
|nextPageToken|string|false|none|次回リクエスト時に指定するページトークン|

<h2 id="tocS_V1ArticleShareRequestSchema">V1ArticleShareRequestSchema</h2>
<!-- backwards compatibility -->
<a id="schemav1articlesharerequestschema"></a>
<a id="schema_V1ArticleShareRequestSchema"></a>
<a id="tocSv1articlesharerequestschema"></a>
<a id="tocsv1articlesharerequestschema"></a>

```json
{
  "url": "https://example.com",
  "title": "title",
  "description": "description",
  "thumbnail": "https://example.com"
}

```

### Properties

|Name|Type|Required|Restrictions|Description|
|---|---|---|---|---|
|url|string(uri)|true|none|記事のURL|
|title|string|false|none|タイトル|
|description|string|false|none|description|
|thumbnail|string(uri)|false|none|サムネイルのURL|

