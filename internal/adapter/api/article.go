package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
	"github.com/morning-night-guild/platform-app/internal/usecase/port"
	"github.com/morning-night-guild/platform-app/pkg/log"
	"github.com/morning-night-guild/platform-app/pkg/openapi"
)

func (api *API) V1ListArticles(w http.ResponseWriter, r *http.Request, params openapi.V1ListArticlesParams) {
	ctx := r.Context()

	pageToken := ""
	if params.PageToken != nil {
		pageToken = *params.PageToken
	}

	token := value.NewNextToken(pageToken)

	size, err := value.NewSize(params.MaxPageSize)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to list articles", log.ErrorField(err))

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	input := port.APIArticleListInput{
		Index: token.ToIndex(),
		Size:  size,
	}

	res, err := api.article.list.Execute(ctx, input)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to list articles", log.ErrorField(err))

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	articles := make([]openapi.Article, len(res.Articles))

	for i, article := range res.Articles {
		id := uuid.MustParse(article.ID.String())
		tags := article.TagList.StringSlice()
		articles[i] = openapi.Article{
			Id:          &id,
			Title:       api.StringToPointer(article.Title.String()),
			Url:         api.StringToPointer(article.URL.String()),
			Description: api.StringToPointer(article.Description.String()),
			Thumbnail:   api.StringToPointer(article.Thumbnail.String()),
			Tags:        &tags,
		}
	}

	next := token.CreateNextToken(size).String()

	rs := openapi.ListArticleResponse{
		Articles:      &articles,
		NextPageToken: &next,
	}

	if err := json.NewEncoder(w).Encode(rs); err != nil {
		log.GetLogCtx(ctx).Warn("failed to encode response", log.ErrorField(err))

		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (api *API) V1ShareArticle(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	key := r.Header.Get("Api-Key")
	if key != api.key {
		log.GetLogCtx(ctx).Warn(fmt.Sprintf("invalid api key. api key = %s", key))

		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	var body openapi.V1ShareArticleRequest

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.GetLogCtx(ctx).Warn("failed to decode request body", log.ErrorField(err))

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	input := port.APIArticleShareInput{
		URL:         article.URL(body.Url),
		Title:       article.Title(api.PointerToString(body.Title)),
		Description: article.Description(api.PointerToString(body.Description)),
		Thumbnail:   article.Thumbnail(api.PointerToString(body.Thumbnail)),
	}

	if _, err := api.article.share.Execute(ctx, input); err != nil {
		w.WriteHeader(api.HandleConnectError(ctx, err))

		return
	}
}
