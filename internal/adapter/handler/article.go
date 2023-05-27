package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/google/uuid"
	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/model/user"
	"github.com/morning-night-guild/platform-app/internal/domain/value"
	"github.com/morning-night-guild/platform-app/pkg/log"
	"github.com/morning-night-guild/platform-app/pkg/openapi"
)

// 記事一覧
// (GET /v1/articles).
func (hdl *Handler) V1ArticleList( //nolint:cyclop
	w http.ResponseWriter,
	r *http.Request,
	params openapi.V1ArticleListParams,
) {
	ctx := r.Context()

	uid, err := hdl.ExtractUserID(ctx, r)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to extract user id", log.ErrorField(err))

		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	ctx = user.SetUIDCtx(ctx, uid)

	scope := article.All

	if params.Scope != nil {
		s, err := article.NewScope(string(*params.Scope))
		if err != nil {
			log.GetLogCtx(ctx).Warn("failed to list articles", log.ErrorField(err))

			w.WriteHeader(http.StatusBadRequest)

			return
		}

		scope = s
	}

	pageToken := ""
	if params.PageToken != nil {
		pageToken = *params.PageToken
	}

	token := value.NewNextToken(pageToken)

	size := model.DefaultArticleSize

	if params.MaxPageSize != nil {
		s, err := value.NewSize(*params.MaxPageSize)
		if err != nil {
			log.GetLogCtx(ctx).Warn("failed to list articles", log.ErrorField(err))

			w.WriteHeader(http.StatusBadRequest)

			return
		}

		size = s
	}

	input := usecase.APIArticleListInput{
		Scope:  scope,
		UserID: uid,
		Index:  token.ToIndex(),
		Size:   size,
	}

	if params.Title != nil {
		input.Filter = []value.Filter{
			value.NewFilter("title", *params.Title),
		}
	}

	output, err := hdl.article.List(ctx, input)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to list articles", log.ErrorField(err))

		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	articles := make([]openapi.ArticleSchema, len(output.Articles))

	for i, article := range output.Articles {
		id := uuid.MustParse(article.ArticleID.String())
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
func (hdl *Handler) V1ArticleShare(
	w http.ResponseWriter,
	r *http.Request,
) {
	hdl.V1InternalArticleShare(w, r)
}

// 記事追加
// (POST /v1/articles/{articleId}).
func (hdl *Handler) V1ArticleAddOwn(
	w http.ResponseWriter,
	r *http.Request,
	articleID types.UUID,
) {
	ctx := r.Context()

	uid, err := hdl.ExtractUserID(ctx, r)
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to extract user id", log.ErrorField(err))

		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	ctx = user.SetUIDCtx(ctx, uid)

	aid, err := article.NewID(articleID.String())
	if err != nil {
		log.GetLogCtx(ctx).Warn("failed to add article", log.ErrorField(err))

		w.WriteHeader(http.StatusBadRequest)

		return
	}

	input := usecase.APIArticleAddToUserInput{
		ArticleID: aid,
		UserID:    uid,
	}

	if _, err := hdl.article.AddToUser(ctx, input); err != nil {
		log.GetLogCtx(ctx).Warn("failed to add article", log.ErrorField(err))

		w.WriteHeader(hdl.HandleConnectError(ctx, err))

		return
	}
}

// 記事削除
// (DELETE /v1/articles/{articleId}).
func (hdl *Handler) V1ArticleRemoveOwn(
	w http.ResponseWriter,
	r *http.Request,
	articleID types.UUID,
) {
	ctx := r.Context()

	log.GetLogCtx(ctx).Debug(fmt.Sprintf("%v %v %v", w, r, articleID))

	w.WriteHeader(http.StatusNotImplemented)
}

// 記事共有
// (POST /v1/internal/articles).
func (hdl *Handler) V1InternalArticleShare(
	w http.ResponseWriter,
	r *http.Request,
) {
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

// 記事削除
// (DELETE /v1/internal/articles/{articleId}).
func (hdl *Handler) V1InternalArticleDelete(
	w http.ResponseWriter,
	r *http.Request,
	articleID types.UUID,
) {
	ctx := r.Context()

	key := r.Header.Get("Api-Key")
	if key != hdl.key {
		log.GetLogCtx(ctx).Warn(fmt.Sprintf("invalid api key. api key = %s", key))

		w.WriteHeader(http.StatusUnauthorized)

		return
	}

	input := usecase.APIArticleDeleteInput{
		ArticleID: article.ID(articleID),
	}

	if _, err := hdl.article.Delete(ctx, input); err != nil {
		w.WriteHeader(hdl.HandleConnectError(ctx, err))

		return
	}
}
