package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/google/uuid"
	articlev1 "github.com/morning-night-guild/platform-app/pkg/connect/proto/article/v1"
	"github.com/morning-night-guild/platform-app/pkg/log"
	"github.com/morning-night-guild/platform-app/pkg/openapi"
)

func (api *API) V1ListArticles(w http.ResponseWriter, r *http.Request, params openapi.V1ListArticlesParams) {
	ctx := log.SetLogCtx(r.Context())

	pageToken := ""
	if params.PageToken != nil {
		pageToken = *params.PageToken
	}

	req := &articlev1.ListRequest{
		PageToken:   pageToken,
		MaxPageSize: uint32(params.MaxPageSize),
	}

	res, err := api.connect.Article.List(ctx, connect.NewRequest(req))
	if err != nil {
		log.GetLogCtx(ctx).Error("failed to list articles", log.ErrorField(err))

		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))

		return
	}

	articles := make([]openapi.Article, len(res.Msg.Articles))

	for i, article := range res.Msg.Articles {
		uid, _ := uuid.Parse(article.Id)

		articles[i] = openapi.Article{
			Id:          &uid,
			Title:       &article.Title,
			Url:         &article.Url,
			Description: &article.Description,
			Thumbnail:   &article.Thumbnail,
			Tags:        &article.Tags,
		}
	}

	rs := openapi.ListArticleResponse{
		Articles:      &articles,
		NextPageToken: &res.Msg.NextPageToken,
	}

	if err := json.NewEncoder(w).Encode(rs); err != nil {
		log.GetLogCtx(ctx).Error("failed to encode response", log.ErrorField(err))

		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (api *API) V1ShareArticle(w http.ResponseWriter, r *http.Request) {
	ctx := log.SetLogCtx(r.Context())

	log.GetLogCtx(ctx).Info(fmt.Sprintf("%+v", w.Header()))

	key := r.Header.Get("Api-Key")
	if key != api.key {
		log.GetLogCtx(ctx).Warn(fmt.Sprintf("invalid api key. api key = %s", key))

		w.WriteHeader(http.StatusUnauthorized)
		_, _ = w.Write([]byte("unauthorized"))

		return
	}

	var body openapi.V1ShareArticleRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.GetLogCtx(ctx).Error("failed to decode request body", log.ErrorField(err))

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	req := &articlev1.ShareRequest{
		Url:         body.Url,
		Title:       *body.Title,
		Description: *body.Description,
		Thumbnail:   *body.Thumbnail,
	}

	res, err := api.connect.Article.Share(ctx, connect.NewRequest(req))
	if err != nil {
		w.WriteHeader(api.HandleConnectError(ctx, err))

		return
	}

	log.GetLogCtx(ctx).Debug(fmt.Sprintf("article shared id = %s", res.Msg.Article.Id))
}
