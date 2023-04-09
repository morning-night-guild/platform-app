package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
	"github.com/morning-night-guild/platform-app/pkg/log"
	"github.com/morning-night-guild/platform-app/pkg/openapi"
)

// 記事一覧
// (GET /v1/articles).
func (hdl *Handler) V1ArticleList(w http.ResponseWriter, r *http.Request, params openapi.V1ArticleListParams) {
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

	uid, err := hdl.ExtractUserID(ctx, r)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to extract user id", log.ErrorField(err))

		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	input := usecase.APIArticleListInput{
		UserID: uid,
		Index:  token.ToIndex(),
		Size:   size,
	}

	output, err := hdl.article.List(ctx, input)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to list articles", log.ErrorField(err))

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	articles := make([]openapi.ArticleSchema, len(output.Articles))

	for i, article := range output.Articles {
		id := uuid.MustParse(article.ID.String())
		tags := article.TagList.StringSlice()
		articles[i] = openapi.ArticleSchema{
			Id:          &id,
			Title:       hdl.StringToPointer(article.Title.String()),
			Url:         hdl.StringToPointer(article.URL.String()),
			Description: hdl.StringToPointer(article.Description.String()),
			Thumbnail:   hdl.StringToPointer(article.Thumbnail.String()),
			Tags:        &tags,
		}
	}

	next := token.CreateNextToken(size).String()

	res := openapi.V1ArticleListResponseSchema{
		Articles:      &articles,
		NextPageToken: &next,
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.GetLogCtx(ctx).Warn("failed to encode outputponse", log.ErrorField(err))

		w.WriteHeader(http.StatusInternalServerError)
	}
}

// 記事共有
// (POST /v1/articles).
func (hdl *Handler) V1ArticleShare(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	key := r.Header.Get("Api-Key")
	if key != hdl.key {
		log.GetLogCtx(ctx).Warn(fmt.Sprintf("invalid api key. api key = %s", key))

		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	var body openapi.V1ArticleShareRequestSchema

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.GetLogCtx(ctx).Warn("failed to decode request body", log.ErrorField(err))

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	input := usecase.APIArticleShareInput{
		URL:         article.URL(body.Url),
		Title:       article.Title(hdl.PointerToString(body.Title)),
		Description: article.Description(hdl.PointerToString(body.Description)),
		Thumbnail:   article.Thumbnail(hdl.PointerToString(body.Thumbnail)),
	}

	if _, err := hdl.article.Share(ctx, input); err != nil {
		w.WriteHeader(hdl.HandleConnectError(ctx, err))

		return
	}
}
