# BandHub 開発タスク一覧

> 先輩エンジニアとして、学習効果を最大化するタスク分割を設計した。
> 各タスクには「なぜこの順番か」「何を学ぶか」を明記している。

---

## タスクの進め方（ルール）

1. **必ず上から順番に進める**（依存関係があるため飛ばさない）
2. **1タスクごとにコミットする**（`feat:`, `fix:` 等のプレフィックス付き）
3. **タスク完了時に「動くもの」を確認してからコミットする**（curl, ブラウザ, テスト等）
4. **わからないことがあったら、まず自分で調べて15分考える → それでも解決しなければ質問する**
5. **コードを書く前に「このコードはどの層に属するか」を意識する**

---

## Phase 0: 開発環境の構築と動作確認

> **学習テーマ**: Docker, PostgreSQL, Go, Next.js の基本動作を確認する
> **ゴール**: `docker compose up` で3コンテナが起動し、ヘルスチェックが通ること

### Task 0-1: Docker Compose で全コンテナを起動する

- [ ] `docker compose up --build` で3コンテナが起動することを確認
- [ ] `curl http://localhost:8080/health` で `{"status":"ok"}` が返ること
- [ ] `http://localhost:3000` で Next.js のデフォルトページが表示されること
- [ ] `docker compose exec db psql -U bandhub -d bandhub_dev` で PostgreSQL に接続できること

**確認観点**: 各コンテナのログにエラーが出ていないか？ポートの競合はないか？

### Task 0-2: DBマイグレーションの仕組みを作る

- [ ] `backend/migrations/001_create_users.sql` を作成（CREATE TABLE 文）
- [ ] `backend/migrations/002_create_events.sql` を作成（CREATE TABLE 文）
- [ ] Docker Compose の起動時に自動でマイグレーションが走る仕組みを作る（または手動実行のスクリプト）
- [ ] テーブルが作成されたことを psql で確認

**なぜ最初にやるか**: テーブルがないと後続のすべてのタスクが進められないため。
**学ぶこと**: SQL の DDL（CREATE TABLE）、UUID型、配列型（TEXT[]）、外部キー制約。

**ヒント**: マイグレーションツール（golang-migrate 等）を使うか、シンプルに SQL ファイルを直接実行するかは自分で判断すること。最初はシンプルな方がよい。

---

## Phase 1: ドメイン層を作る（バックエンド最重要フェーズ）

> **学習テーマ**: Clean Architecture の核心 ── ドメイン層は外部に依存しない
> **ゴール**: `internal/domain/` 配下に、GORM にも Echo にも依存しないピュアな Go コードができること

### なぜ Phase 1 でドメイン層から始めるのか？

Clean Architecture で最も重要な原則は**依存性逆転の原則（DIP）**。
ドメイン層は他のどの層にも依存せず、ビジネスルールだけを表現する。
先にドメイン層を固めることで、後続の usecase / infrastructure / handler 層が
「ドメイン層に合わせて作る」という正しい依存方向で開発できる。

**逆にやると何が起きるか？**
handler から作り始めると、DB のカラムや HTTP のリクエスト形式に引きずられて、
ドメインモデルが「技術都合の型」になってしまう。これが技術的負債の典型。

### Task 1-1: 値オブジェクト（Part）を定義する

**ファイル**: `internal/domain/value/part.go`

- [ ] `Part` 型を定義する（string のラッパー型）
- [ ] 有効な値を定数で定義する（Vo, Gt, Ba, Dr, Key, etc.）
- [ ] `NewPart(s string) (Part, error)` で不正な値を弾くバリデーションを実装
- [ ] このファイルに `import` が1行もないこと（標準ライブラリの `fmt` や `errors` は OK）を確認

**学ぶこと**:
- 値オブジェクトとは何か？ → 「同一性（ID）を持たず、値そのもので比較される型」
- なぜ `string` をそのまま使わないのか？ → 不正な値（"ギター" など）が混入するのを型レベルで防ぐ
- Go での値オブジェクトの表現方法（型エイリアス vs 独自型）

### Task 1-2: 値オブジェクト（Role, EventType）を定義する

**ファイル**: `internal/domain/value/role.go`, `internal/domain/value/event_type.go`

- [ ] `Role` 型を定義（member, admin）
- [ ] `EventType` 型を定義（live, camp, other）
- [ ] それぞれ `New` 関数でバリデーション付きコンストラクタを実装
- [ ] `Role` に `IsAdmin() bool` メソッドを追加（後で権限チェックに使う）

**学ぶこと**: 値オブジェクトにビジネスロジック（IsAdmin）を持たせるパターン。

### Task 1-3: カスタムエラー型を定義する

**ファイル**: `internal/domain/error.go`（または `internal/domain/errors/` ディレクトリ）

- [ ] `ErrNotFound` を定義（ユーザーやイベントが見つからない場合）
- [ ] `ErrDuplicateEmail` を定義（メールアドレスが重複している場合）
- [ ] `ErrInvalidCredentials` を定義（ログイン失敗時）
- [ ] `ErrUnauthorized` を定義（認証されていない場合）
- [ ] `ErrForbidden` を定義（権限がない場合）

**なぜ domain 層でエラーを定義するか？**
usecase 層が「何が起きたか」を表現するとき、HTTP ステータスコード（404, 401 等）に
依存したくない。domain 層のエラーを使えば、handler 層が適切なステータスコードに
変換する責務を持てる。

**学ぶこと**: Go のカスタムエラー型、errors パッケージの使い方。

### Task 1-4: User エンティティを定義する

**ファイル**: `internal/domain/entity/user.go`

- [ ] `User` 構造体を定義（フィールドは CLAUDE.md の users テーブルに対応）
- [ ] `Part` や `Role` など、値オブジェクトで定義した型を使う
- [ ] GORM のタグ（`gorm:"..."`）は **絶対に付けない**（これが Clean Architecture のルール）
- [ ] フィールドはエクスポート（大文字始まり）にする

**考えてほしいポイント**:
- `PasswordHash` フィールドを `User` 構造体に持たせるべきか？
  → レスポンスに含めてはいけないが、ドメインモデルとしては必要
  → どう設計するか考えてみよう

### Task 1-5: Event エンティティを定義する

**ファイル**: `internal/domain/entity/event.go`

- [ ] `Event` 構造体を定義（CLAUDE.md の events テーブルに対応）
- [ ] `EventType` 値オブジェクトを使う
- [ ] `CreatedBy` は `uuid.UUID` 型を使う（google/uuid パッケージ）

### Task 1-6: リポジトリインターフェースを定義する

**ファイル**: `internal/domain/repository/user_repository.go`, `internal/domain/repository/event_repository.go`

- [ ] `UserRepository` インターフェースを定義
  - `FindByID(ctx context.Context, id uuid.UUID) (*entity.User, error)`
  - `FindByEmail(ctx context.Context, email string) (*entity.User, error)`
  - `FindAll(ctx context.Context, filter UserFilter) ([]*entity.User, error)`
  - `Create(ctx context.Context, user *entity.User) error`
  - `Update(ctx context.Context, user *entity.User) error`
- [ ] `UserFilter` 構造体を定義（検索条件: Part, Year, Genre, Query）
- [ ] `EventRepository` インターフェースを定義
  - `FindByID`, `FindAll`, `Create`, `Update`, `Delete`
- [ ] インターフェースが `context.Context` を第1引数に取ること（Go の慣習）

**なぜインターフェースを domain 層に置くか？**（これが DIP の核心）
usecase 層は「データの保存・取得」が必要だが、「GORM を使ってどう保存するか」は
知る必要がない。インターフェースを domain 層に定義し、infrastructure 層が実装する。
これにより usecase 層は具体的なDB実装に依存せず、テスト時にはモックに差し替えられる。

---

## Phase 2: 認証機能（バックエンド ── 全層を縦に貫通する最初の機能）

> **学習テーマ**: Clean Architecture で1つの機能を全層にわたって実装する流れを体験する
> **ゴール**: `POST /api/v1/auth/signup` と `POST /api/v1/auth/login` が動くこと

### なぜ最初の機能が「認証」なのか？

1. 他のすべての機能が「ログイン済みユーザー」を前提にするため
2. 全層（domain → usecase → infrastructure → handler）を縦断するので、
   Clean Architecture の「依存の方向」を一番最初に体で覚えられる
3. JWT, bcrypt, ミドルウェアなど、Web 開発で必須の技術に触れられる

### Task 2-1: パスワードハッシュ化のユーティリティを作る

**ファイル**: `internal/infrastructure/auth/password.go`

- [ ] `HashPassword(raw string) (string, error)` を実装（bcrypt）
- [ ] `CheckPassword(hash, raw string) error` を実装
- [ ] bcrypt のコストパラメータを理解する（デフォルトの 10 でOK）

**学ぶこと**: なぜパスワードをハッシュ化するのか？ bcrypt vs SHA-256 の違い。

### Task 2-2: JWT トークンの生成・検証を実装する

**ファイル**: `internal/infrastructure/auth/jwt.go`

- [ ] `GenerateAccessToken(userID uuid.UUID, role string) (string, error)` を実装
- [ ] `GenerateRefreshToken(userID uuid.UUID) (string, error)` を実装
- [ ] `ValidateToken(tokenString string) (*Claims, error)` を実装
- [ ] アクセストークンの有効期限: 15分、リフレッシュトークン: 7日
- [ ] JWT のシークレットキーは環境変数から取得する

**学ぶこと**:
- JWT の構造（Header.Payload.Signature）
- アクセストークンとリフレッシュトークンの役割の違い
- なぜアクセストークンの有効期限を短くするのか？

### Task 2-3: 認証ユースケースを実装する

**ファイル**: `internal/usecase/auth_usecase.go`

- [ ] `AuthUsecase` 構造体を定義
  - 依存: `UserRepository`（インターフェース）、パスワードハッシャー、トークン生成器
- [ ] `Signup(ctx, email, password, displayName) (*TokenPair, error)` を実装
  - メール重複チェック → パスワードハッシュ化 → ユーザー作成 → トークン生成
- [ ] `Login(ctx, email, password) (*TokenPair, error)` を実装
  - メールでユーザー検索 → パスワード照合 → トークン生成
- [ ] `RefreshToken(ctx, refreshToken) (*TokenPair, error)` を実装

**考えてほしいポイント**:
- `AuthUsecase` は `bcrypt` パッケージを直接 import すべきか？
  → infrastructure の詳細に依存してしまう。インターフェースで抽象化すべきか？
  → ただし過度な抽象化は学習段階では避けたい。どこでバランスを取るか考える

### Task 2-4: UserRepository の GORM 実装を作る

**ファイル**: `internal/infrastructure/persistence/user_repository.go`

- [ ] `userRepository` 構造体を定義（小文字始まり = unexported）
- [ ] `NewUserRepository(db *gorm.DB) repository.UserRepository` コンストラクタ
- [ ] GORM 用のモデル構造体を **このファイル内に** 定義する（`gorm:"..."` タグ付き）
  - domain の `entity.User` とは別の構造体
  - 変換メソッド `toEntity()` / `fromEntity()` を用意する
- [ ] `FindByID`, `FindByEmail`, `Create` を実装（他は後でOK）
- [ ] GORM のエラーを domain 層のカスタムエラーに変換する

**なぜ GORM モデルと domain エンティティを分けるか？**
domain 層の `User` に `gorm:"..."` タグを付けると、domain 層が GORM に依存する。
infrastructure 層に別のモデルを用意し、変換するのが Clean Architecture のやり方。
「面倒じゃないか？」と思うだろうが、この分離こそがテスト容易性と保守性を生む。

### Task 2-5: 認証ハンドラーを実装する

**ファイル**: `internal/handler/auth_handler.go`

- [ ] リクエスト DTO を定義（`request/auth_request.go`）
  - `SignupRequest { Email, Password, DisplayName }`
  - `LoginRequest { Email, Password }`
- [ ] レスポンス DTO を定義（`response/auth_response.go`）
  - `TokenResponse { AccessToken, RefreshToken }`
- [ ] `AuthHandler` 構造体を定義（依存: `AuthUsecase`）
- [ ] `Signup(c echo.Context) error` を実装
  - リクエストのバインドとバリデーション → usecase 呼び出し → レスポンス返却
- [ ] `Login(c echo.Context) error` を実装
- [ ] domain 層のエラーを HTTP ステータスコードに変換するロジックを実装
  - `ErrDuplicateEmail` → 409 Conflict
  - `ErrInvalidCredentials` → 401 Unauthorized

**学ぶこと**: handler 層の責務は「HTTPの世界」と「ビジネスロジックの世界」の橋渡し。
リクエストの変換、usecase の呼び出し、レスポンスの整形のみ。ビジネスロジックは書かない。

### Task 2-6: 認証ミドルウェアを実装する

**ファイル**: `internal/handler/middleware/auth_middleware.go`

- [ ] `AuthMiddleware` を実装
  - Authorization ヘッダーから Bearer トークンを取得
  - JWT を検証し、ユーザー情報を Echo の Context にセット
  - トークンが無効なら 401 を返す
- [ ] `AdminOnly` ミドルウェアを実装
  - Context からロールを取得し、admin でなければ 403 を返す

### Task 2-7: main.go で DI とルーティングを組み立てる

**ファイル**: `cmd/server/main.go`

- [ ] GORM で PostgreSQL に接続
- [ ] `UserRepository` → `AuthUsecase` → `AuthHandler` の順に DI
- [ ] ルーティングを設定（`/api/v1/auth/signup`, `/api/v1/auth/login`）
- [ ] curl で実際にサインアップ → ログインできることを確認

**確認手順**:
```bash
# サインアップ
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","display_name":"テスト太郎"}'

# ログイン
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'
```

**ここまで来たらお祝い**: Clean Architecture で1つの機能が全層を通して動いた。
これが今後のすべての機能開発のテンプレートになる。

---

## Phase 3: ユーザー機能（バックエンド）

> **学習テーマ**: Phase 2 で覚えたパターンを反復して定着させる
> **ゴール**: メンバー一覧・検索・詳細・プロフィール更新の API が動くこと

### Task 3-1: ユーザーユースケースを実装する

**ファイル**: `internal/usecase/user_usecase.go`

- [ ] `GetProfile(ctx, userID) (*entity.User, error)` ── 自分のプロフィール取得
- [ ] `UpdateProfile(ctx, userID, input) (*entity.User, error)` ── プロフィール更新
- [ ] `ListUsers(ctx, filter) ([]*entity.User, error)` ── メンバー一覧・検索
- [ ] `GetUser(ctx, targetID) (*entity.User, error)` ── メンバー詳細

**考えてほしいポイント**:
- `UpdateProfile` の `input` はどんな型にする？ entity.User をそのまま渡す？専用の DTO を作る？
- 「自分以外のプロフィールは編集できない」というルールはどの層で守る？

### Task 3-2: UserRepository の残りのメソッドを実装する

- [ ] `FindAll(ctx, filter)` を実装（GORM の Where 句を動的に組み立てる）
- [ ] `Update(ctx, user)` を実装
- [ ] 検索条件（パート、学年、ジャンル、フリーワード）が正しくフィルタリングされることを確認

### Task 3-3: ユーザーハンドラーを実装する

**ファイル**: `internal/handler/user_handler.go`

- [ ] `GET /api/v1/users/me` ── 認証ミドルウェア経由で userID を取得し、プロフィールを返す
- [ ] `PUT /api/v1/users/me` ── プロフィール更新
- [ ] `GET /api/v1/users` ── クエリパラメータで検索条件を受け取る
- [ ] `GET /api/v1/users/:id` ── メンバー詳細
- [ ] レスポンス DTO で `PasswordHash` を **絶対に返さない** ことを確認

### Task 3-4: main.go にユーザールーティングを追加する

- [ ] `UserUsecase` → `UserHandler` の DI を追加
- [ ] 認証が必要なエンドポイントに `AuthMiddleware` を適用
- [ ] curl でメンバー一覧・検索が動くことを確認

---

## Phase 4: イベント機能（バックエンド）

> **学習テーマ**: 権限制御（admin のみ作成可能）を Clean Architecture で実現する
> **ゴール**: イベントの CRUD API が動き、admin のみ作成・更新・削除できること

### Task 4-1: EventRepository の GORM 実装を作る

**ファイル**: `internal/infrastructure/persistence/event_repository.go`

- [ ] Phase 2-4 と同様に、GORM モデルと domain エンティティの変換を実装
- [ ] `FindAll` に月での絞り込み（`?month=2026-03`）を実装

### Task 4-2: イベントユースケースを実装する

**ファイル**: `internal/usecase/event_usecase.go`

- [ ] `ListEvents(ctx, month) ([]*entity.Event, error)`
- [ ] `GetEvent(ctx, id) (*entity.Event, error)`
- [ ] `CreateEvent(ctx, userID, input) (*entity.Event, error)`
- [ ] `UpdateEvent(ctx, id, input) (*entity.Event, error)`
- [ ] `DeleteEvent(ctx, id) error`

**考えてほしいポイント**:
- 「admin のみ作成可能」のチェックは usecase でやる？ handler（ミドルウェア）でやる？
- → ミドルウェアでやるのが Echo の一般的なパターン。usecase にはロールの概念を持ち込まない方がシンプル

### Task 4-3: イベントハンドラーとルーティングを実装する

- [ ] `event_handler.go` を実装
- [ ] 作成・更新・削除に `AdminOnly` ミドルウェアを適用
- [ ] main.go にルーティング追加
- [ ] curl で admin ユーザーでイベントが作成できること、member ユーザーでは 403 が返ることを確認

---

## Phase 5: 管理者機能（バックエンド）

> **学習テーマ**: 小さな機能を手際よく追加する
> **ゴール**: admin がメンバーのロールを変更できること

### Task 5-1: ロール変更のユースケースとハンドラーを実装する

- [ ] `PUT /api/v1/admin/users/:id/role`
- [ ] usecase でロール変更を実装
- [ ] handler + AdminOnly ミドルウェア
- [ ] curl で確認

---

## Phase 6: バックエンドのユニットテスト

> **学習テーマ**: インターフェースを使ったモック、テスタブルな設計の価値を実感する
> **ゴール**: usecase 層のテストが通ること

### なぜ Phase 6 でテストを書くのか？

テストを最初に書く（TDD）のが理想だが、Clean Architecture のパターンを理解する前に
テストを書こうとすると、「何をテストすべきか」がわからない。
Phase 2-5 で実装を経験した後なら、「usecase 層のロジックをテストしたい」
「repository はモックに差し替えたい」という動機が自然に生まれる。

### Task 6-1: モック用のリポジトリを作る

**ファイル**: `internal/usecase/mock_test.go`（または testutil パッケージ）

- [ ] `UserRepository` のモック実装を作る（手動モック or testify/mock）
- [ ] `EventRepository` のモック実装を作る

**学ぶこと**: インターフェースを domain 層で定義した恩恵がここで活きる。
GORM やデータベースなしで usecase のテストが書ける。

### Task 6-2: AuthUsecase のテストを書く

- [ ] サインアップ成功ケース
- [ ] サインアップ時にメール重複でエラーになるケース
- [ ] ログイン成功ケース
- [ ] ログイン時にパスワード不一致でエラーになるケース

### Task 6-3: UserUsecase のテストを書く

- [ ] プロフィール取得テスト
- [ ] メンバー一覧取得テスト
- [ ] プロフィール更新テスト

### Task 6-4: EventUsecase のテストを書く

- [ ] イベント作成テスト
- [ ] イベント一覧取得テスト

---

## Phase 7: フロントエンド基盤を作る

> **学習テーマ**: Next.js App Router, TanStack Query, shadcn/ui のセットアップ
> **ゴール**: フロントエンドの「土台」を整えて、機能開発にすぐ入れる状態にする

### Task 7-1: shadcn/ui をセットアップする

- [ ] shadcn/ui を初期化（`npx shadcn-ui@latest init`）
- [ ] 基本コンポーネントを追加: Button, Card, Input, Label, Form

**学ぶこと**: shadcn/ui は「コピーして使う」UI ライブラリ。npm パッケージとは違い、
ソースコードがプロジェクトに直接入る。カスタマイズしやすい反面、自分で管理する必要がある。

### Task 7-2: API クライアント基盤を作る

**ファイル**: `src/shared/lib/api-client.ts`

- [ ] fetch のラッパー関数を作る（ベース URL, JSON ヘッダー, エラーハンドリング）
- [ ] アクセストークンを自動付与する仕組みを作る
- [ ] 401 レスポンス時にリフレッシュトークンで自動再取得する仕組みを作る（ここは難しいので後回しでもOK）

**考えてほしいポイント**:
- axios を使うか、fetch API を使うか？
- トークンの保存先は localStorage? cookie? → それぞれのメリデメを調べてみよう

### Task 7-3: TanStack Query をセットアップする

**ファイル**: `src/shared/lib/query-client.ts`, `src/app/providers.tsx`

- [ ] `npm install @tanstack/react-query`
- [ ] `QueryClientProvider` を App Router のレイアウトに組み込む
- [ ] React Query DevTools を追加（開発時のデバッグ用）

**学ぶこと**: TanStack Query がサーバー状態管理を担う。
キャッシュ、再フェッチ、ローディング/エラー状態を自動管理してくれる。
これにより useState + useEffect で API コールするパターンが不要になる。

### Task 7-4: 認証状態管理の仕組みを作る

**ファイル**: `src/shared/hooks/useAuth.ts`, `src/features/auth/api/auth-api.ts`

- [ ] トークンの保存・取得・削除のユーティリティ
- [ ] `useAuth` フック（ログイン状態の判定、ログアウト処理）
- [ ] 認証ガード（未認証なら /login にリダイレクト）

### Task 7-5: レイアウトを作る

**ファイル**: `src/app/(main)/layout.tsx`

- [ ] サイドバーまたはヘッダーナビゲーションを作る
- [ ] ナビゲーション項目: ダッシュボード, メンバー, カレンダー, プロフィール
- [ ] 認証ガードを適用（未ログインなら /login にリダイレクト）
- [ ] まずは見た目だけ作る（リンクは後で接続する）

---

## Phase 8: ログイン・サインアップ画面（フロントエンド）

> **学習テーマ**: Feature-Sliced Design で機能を完結させるパターン
> **ゴール**: ブラウザからログイン・サインアップができ、認証後にダッシュボードに遷移すること

### Task 8-1: ログインページを実装する

**ファイル**: `src/features/auth/` 配下

- [ ] `src/features/auth/types/auth.ts` ── 型定義（LoginRequest, TokenResponse 等）
- [ ] `src/features/auth/api/auth-api.ts` ── API 呼び出し関数
- [ ] `src/features/auth/hooks/useLogin.ts` ── TanStack Query の useMutation
- [ ] `src/features/auth/components/LoginForm.tsx` ── ログインフォーム UI
- [ ] `src/app/(auth)/login/page.tsx` ── ページ（LoginForm をインポートするだけ）

**学ぶこと**: Feature-Sliced Design のファイル配置。
ページファイル（app/）は薄く、ロジックとUIは features/ に置く。

### Task 8-2: サインアップページを実装する

- [ ] Task 8-1 と同様の構造で実装
- [ ] フォームにバリデーションを追加（メール形式、パスワード最低文字数）

### Task 8-3: 認証フローを通しで動かして確認する

- [ ] サインアップ → 自動ログイン → ダッシュボードに遷移
- [ ] ログアウト → ログインページに遷移
- [ ] 未認証で /dashboard にアクセス → /login にリダイレクト

---

## Phase 9: メンバー機能（フロントエンド）

> **学習テーマ**: TanStack Query を使った一覧取得・検索、コンポーネント設計
> **ゴール**: メンバー一覧の表示と検索ができること

### Task 9-1: メンバー一覧ページを実装する

**ファイル**: `src/features/members/` 配下

- [ ] `types/member.ts` ── User 型定義
- [ ] `api/members-api.ts` ── メンバー一覧取得 API
- [ ] `hooks/useMembers.ts` ── useQuery でメンバー一覧取得
- [ ] `components/MemberCard.tsx` ── メンバーカード（名前、パート、学年）
- [ ] `components/MemberList.tsx` ── カードのグリッド表示
- [ ] `src/app/(main)/members/page.tsx` ── ページ

### Task 9-2: メンバー検索機能を実装する

- [ ] `components/MemberSearchForm.tsx` ── 検索フォーム（パート、学年、ジャンル、フリーワード）
- [ ] `hooks/useMemberSearch.ts` ── 検索条件を useQuery のキーに含める
- [ ] 検索条件を変更したら自動で再取得されることを確認（TanStack Query のキャッシュキー）

**学ぶこと**: TanStack Query のクエリキーによるキャッシュ管理。
検索条件が変わるとクエリキーが変わり、自動で再フェッチが走る仕組み。

### Task 9-3: メンバー詳細ページを実装する

- [ ] `src/app/(main)/members/[id]/page.tsx` ── 動的ルート
- [ ] `components/MemberDetail.tsx` ── 詳細表示（プロフィール、自己紹介、好きなジャンル等）

---

## Phase 10: プロフィール編集（フロントエンド）

> **学習テーマ**: フォームの状態管理、楽観的更新
> **ゴール**: 自分のプロフィールを編集・保存できること

### Task 10-1: プロフィールページを実装する

- [ ] `src/features/profile/` 配下に API, hooks, components を実装
- [ ] プロフィール表示と編集フォームを作る
- [ ] 保存後に TanStack Query のキャッシュを無効化（invalidateQueries）して最新データを反映

---

## Phase 11: カレンダー・イベント機能（フロントエンド）

> **学習テーマ**: カレンダー UI の実装、管理者限定機能の制御
> **ゴール**: カレンダーにイベントが表示され、admin はイベントを作成できること

### Task 11-1: カレンダーページを実装する

- [ ] 月表示のカレンダー UI を作る（ライブラリを使うか自作するか検討）
- [ ] イベントをカレンダー上に表示
- [ ] 月の切り替えで API を再取得

### Task 11-2: イベント詳細ページを実装する

- [ ] `src/app/(main)/events/[id]/page.tsx`
- [ ] イベント情報（タイトル、日時、会場、説明）を表示

### Task 11-3: イベント作成・編集フォームを実装する（admin のみ）

- [ ] admin ユーザーのみに「イベント作成」ボタンを表示
- [ ] イベント作成フォーム
- [ ] イベント編集・削除機能

---

## Phase 12: 管理者画面（フロントエンド）

> **学習テーマ**: ロールベースの UI 制御
> **ゴール**: admin がメンバーのロールを変更できること

### Task 12-1: 管理者ページを実装する

- [ ] `src/app/(main)/admin/page.tsx`
- [ ] メンバー一覧にロール変更ボタンを表示
- [ ] ロール変更の確認ダイアログ
- [ ] admin 以外がアクセスしたらリダイレクトまたは 403 表示

---

## Phase 13: 仕上げと品質向上

> **学習テーマ**: プロダクション品質に仕上げる
> **ゴール**: ポートフォリオとして見せられる品質にすること

### Task 13-1: エラーハンドリングの統一

- [ ] バックエンド: エラーレスポンスの形式を統一する（`{"error": "メッセージ", "code": "エラーコード"}`）
- [ ] フロントエンド: API エラー時のトースト通知を統一する

### Task 13-2: ローディング・空状態の UI を整える

- [ ] スケルトンローディング（メンバーカード等）
- [ ] 検索結果0件のときのメッセージ
- [ ] エラー発生時のリトライボタン

### Task 13-3: レスポンシブ対応

- [ ] モバイル表示のナビゲーション（ハンバーガーメニュー）
- [ ] メンバーカードのグリッド列数をブレイクポイントで調整

### Task 13-4: README を書く

- [ ] プロジェクト概要、技術スタック、アーキテクチャ図
- [ ] ローカル開発の手順（docker compose up だけで動く状態に）
- [ ] スクリーンショット

---

## Phase 14: デプロイ

> **学習テーマ**: インフラの基本（CI/CD, 環境変数管理）
> **ゴール**: 実際に動く URL を面接官に見せられること

### Task 14-1: バックエンドを Render（または Railway）にデプロイする

- [ ] Dockerfile の最適化（マルチステージビルド）
- [ ] 環境変数の設定（DB_HOST, JWT_SECRET 等）
- [ ] PostgreSQL アドオンの設定

### Task 14-2: フロントエンドを Vercel にデプロイする

- [ ] Vercel にリポジトリを接続
- [ ] 環境変数 `NEXT_PUBLIC_API_URL` を本番 API の URL に設定
- [ ] CORS の設定をバックエンドに追加

### Task 14-3: 本番環境で動作確認する

- [ ] サインアップ → ログイン → 各機能の一通りの動作確認
- [ ] エラーが起きないことの確認

---

## おまけ: 将来の拡張タスク（余裕があれば）

ポートフォリオに「将来構想」として書けるネタ。実装しなくても設計だけしておくと面接で話せる。

- **出演費計算機能**: ライブごとの箱代 ÷ 出演者数を自動計算
- **バンド編成管理**: ライブごとにバンドメンバーを登録・管理
- **通知機能**: 新しいイベントが追加されたらメール or プッシュ通知
- **画像アップロード**: プロフィール画像の S3 アップロード
- **CI/CD パイプライン**: GitHub Actions でテスト・デプロイを自動化

---

## 振り返りチェックリスト

各 Phase 完了時に確認すること:

- [ ] コードは動くか？（curl, ブラウザ, テスト）
- [ ] domain 層に外部ライブラリの import が入っていないか？
- [ ] handler 層にビジネスロジックが入っていないか？
- [ ] usecase 層が具体的な実装（GORM 等）に依存していないか？
- [ ] 変更はコミットしたか？コミットメッセージは適切か？
- [ ] 自分の言葉で「なぜこう設計したか」を説明できるか？
