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

<h1 id="morning-night-guild-app-api-article">article</h1>

記事

## v1ListArticles

<a id="opIdv1ListArticles"></a>

`GET /v1/articles`

*記事一覧*

記事一覧を取得する

<h3 id="v1listarticles-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|pageToken|query|string|false|トークン|
|maxPageSize|query|integer|true|ページサイズ|

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

<h3 id="v1listarticles-responses">Responses</h3>

|Status|Meaning|Description|Schema|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|[ListArticleResponse](#schemalistarticleresponse)|
|500|[Internal Server Error](https://tools.ietf.org/html/rfc7231#section-6.6.1)|Internal Server Error|None|

<aside class="success">
This operation does not require authentication
</aside>

## v1ShareArticle

<a id="opIdv1ShareArticle"></a>

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

<h3 id="v1sharearticle-parameters">Parameters</h3>

|Name|In|Type|Required|Description|
|---|---|---|---|---|
|body|body|[V1ShareArticleRequest](#schemav1sharearticlerequest)|true|記事共有リクエストボディ|

<h3 id="v1sharearticle-responses">Responses</h3>

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

<h2 id="tocS_Article">Article</h2>
<!-- backwards compatibility -->
<a id="schemaarticle"></a>
<a id="schema_Article"></a>
<a id="tocSarticle"></a>
<a id="tocsarticle"></a>

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

<h2 id="tocS_ListArticleResponse">ListArticleResponse</h2>
<!-- backwards compatibility -->
<a id="schemalistarticleresponse"></a>
<a id="schema_ListArticleResponse"></a>
<a id="tocSlistarticleresponse"></a>
<a id="tocslistarticleresponse"></a>

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
|articles|[[Article](#schemaarticle)]|false|none|none|
|nextPageToken|string|false|none|次回リクエスト時に指定するページトークン|

<h2 id="tocS_V1ShareArticleRequest">V1ShareArticleRequest</h2>
<!-- backwards compatibility -->
<a id="schemav1sharearticlerequest"></a>
<a id="schema_V1ShareArticleRequest"></a>
<a id="tocSv1sharearticlerequest"></a>
<a id="tocsv1sharearticlerequest"></a>

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

