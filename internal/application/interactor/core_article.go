package interactor

import (
	"context"
	"fmt"

	"github.com/morning-night-guild/platform-app/internal/application/usecase"
	"github.com/morning-night-guild/platform-app/internal/domain/model"
	"github.com/morning-night-guild/platform-app/internal/domain/model/article"
	"github.com/morning-night-guild/platform-app/internal/domain/repository"
	"github.com/morning-night-guild/platform-app/pkg/log"
)

var _ usecase.CoreArticle = (*CoreArticle)(nil)

// CoreArticle.
type CoreArticle struct {
	articleRepository repository.Article
	userRepository    repository.User
}

func NewCoreArticle(
	articleRepository repository.Article,
	userRepository repository.User,
) *CoreArticle {
	return &CoreArticle{
		articleRepository: articleRepository,
		userRepository:    userRepository,
	}
}

// Share.
func (itr *CoreArticle) Share(
	ctx context.Context,
	input usecase.CoreArticleShareInput,
) (usecase.CoreArticleShareOutput, error) {
	item := model.CreateArticle(input.URL, input.Title, input.Description, input.Thumbnail, []article.Tag{})

	if err := itr.articleRepository.Save(ctx, item); err != nil {
		return usecase.CoreArticleShareOutput{}, err
	}

	return usecase.CoreArticleShareOutput{
		Article: item,
	}, nil
}

// List.
func (itr *CoreArticle) List(
	ctx context.Context,
	input usecase.CoreArticleListInput,
) (usecase.CoreArticleListOutput, error) {
	articles, err := itr.articleRepository.FindAll(ctx, input.Index, input.Size, input.Filter...)
	if err != nil {
		return usecase.CoreArticleListOutput{}, err
	}

	return usecase.CoreArticleListOutput{
		Articles: articles,
	}, nil
}

// Delete.
func (itr *CoreArticle) Delete(
	ctx context.Context,
	input usecase.CoreArticleDeleteInput,
) (usecase.CoreArticleDeleteOutput, error) {
	if _, err := itr.articleRepository.Find(ctx, input.ArticleID); err != nil {
		return usecase.CoreArticleDeleteOutput{}, err
	}

	if err := itr.articleRepository.Delete(ctx, input.ArticleID); err != nil {
		return usecase.CoreArticleDeleteOutput{}, err
	}

	return usecase.CoreArticleDeleteOutput{}, nil
}

func (itr *CoreArticle) AddToUser(
	ctx context.Context,
	input usecase.CoreArticleAddToUserInput,
) (usecase.CoreArticleAddToUserOutput, error) {
	// NOTE:
	// user_articleテーブルでuserとarticleに外部キーを設定しているため
	// UserとArticleのFindは不要だが(存在しなければ追加に失敗&複数無駄な通信が発生)
	// UseCaseのわかりやすさ優先で以下の実装とする
	// パフォーマンス問題が発生したら修正する
	if _, err := itr.userRepository.Find(ctx, input.UserID); err != nil {
		log.GetLogCtx(ctx).Warn(fmt.Sprintf("user not found. id=%s", input.UserID), log.ErrorField(err))

		return usecase.CoreArticleAddToUserOutput{}, err
	}

	if _, err := itr.articleRepository.Find(ctx, input.ArticleID); err != nil {
		log.GetLogCtx(ctx).Warn(fmt.Sprintf("article not found. id=%s", input.ArticleID), log.ErrorField(err))

		return usecase.CoreArticleAddToUserOutput{}, err
	}

	if err := itr.articleRepository.AddToUser(ctx, input.ArticleID, input.UserID); err != nil {
		msg := fmt.Sprintf("failed to add article to user. article_id=%s, user_id=%s", input.ArticleID, input.UserID)

		log.GetLogCtx(ctx).Warn(msg, log.ErrorField(err))

		return usecase.CoreArticleAddToUserOutput{}, err
	}

	return usecase.CoreArticleAddToUserOutput{}, nil
}
