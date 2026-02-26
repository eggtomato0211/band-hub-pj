package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()

	// ミドルウェアの設定
	// Logger: リクエストのログを出力する（開発時のデバッグに便利）
	// Recover: パニック発生時にサーバーを落とさず500を返す
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ヘルスチェック用エンドポイント
	// デプロイ先（Render等）がサーバーの生存確認に使う
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// TODO: ここで各層の依存性注入（DI）を行い、ルーティングを設定する
	// 例:
	//   userRepo := persistence.NewUserRepository(db)
	//   userUsecase := usecase.NewUserUsecase(userRepo)
	//   userHandler := handler.NewUserHandler(userUsecase)
	//   api := e.Group("/api/v1")
	//   api.GET("/users", userHandler.List)

	// 環境変数からポートを取得（デフォルト: 8080）
	// Render 等のホスティングサービスは PORT 環境変数でポートを指定する
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(e.Start(":" + port))
}
