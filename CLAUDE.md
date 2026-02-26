# CLAUDE.md - BandHub プロジェクトガイド

## プロジェクト概要

軽音サークル（約230名）専用のメンバーポータルサイト。
メンバー情報の検索・閲覧と、サークルイベント（ライブ・合宿等）の共有を行うWebアプリケーション。

### 開発者の意図

- 就活ポートフォリオとして、設計力・技術選定の説明力を重視している
- Clean Architecture / DDD を学習中であり、実装を通じて定着させることが目的の一つ
- バイブコーディングに頼りきらず、自分で理解・実装することを重視している
- **したがって、コードを生成する際は「なぜこう書くのか」の説明をコメントや会話で補足すること**
- 将来的な機能追加（出演費計算、収支管理など）を見据えた拡張性のある設計にする

---

## 技術スタック

| レイヤー | 技術 | バージョン目安 |
|---------|------|-------------|
| フロントエンド | Next.js (App Router) + TypeScript | 14.x |
| UI | Tailwind CSS + shadcn/ui | - |
| API通信 | TanStack Query | v5 |
| バックエンド | Go + Echo | Go 1.22+ / Echo v4 |
| ORM | GORM | v2 |
| データベース | PostgreSQL | 16 |
| 認証 | JWT（アクセストークン + リフレッシュトークン） | - |
| ホスティング(FE) | Vercel | - |
| ホスティング(BE) | Render or Railway | - |
| コンテナ | Docker Compose（ローカル開発） | - |

---

## アーキテクチャ方針

### バックエンド: Clean Architecture

依存性逆転の原則（DIP）を守る。内側の層は外側の層に依存しない。

```
依存の方向: handler → usecase → domain ← infrastructure
```

- **domain（ドメイン層）**: エンティティ、値オブジェクト、リポジトリインターフェース
  - 他のどの層にも依存しない
  - ビジネスルールをここに集約する
- **usecase（ユースケース層）**: アプリケーションのビジネスロジック
  - domain層のインターフェースに依存する（具体的な実装には依存しない）
- **infrastructure（インフラ層）**: DBアクセス（GORM）、外部サービス連携
  - domain層で定義されたインターフェースを実装する
- **handler（プレゼンテーション層）**: HTTPハンドラー、リクエスト/レスポンスの変換
  - usecase層に依存する

### フロントエンド: Feature-Sliced Design（簡易版）

機能ごとにディレクトリを分割し、関心の分離を保つ。

---

## ディレクトリ構成

### バックエンド（Go）

```
backend/
├── cmd/
│   └── server/
│       └── main.go              # エントリーポイント、DI（依存性注入）
├── internal/
│   ├── domain/                  # ドメイン層（最も内側、依存なし）
│   │   ├── entity/              # エンティティ（User, Event）
│   │   ├── value/               # 値オブジェクト（Part, Role, EventType など）
│   │   └── repository/          # リポジトリインターフェース定義
│   ├── usecase/                 # ユースケース層
│   │   ├── user_usecase.go
│   │   └── event_usecase.go
│   ├── infrastructure/          # インフラ層（domain のインターフェースを実装）
│   │   ├── persistence/         # GORM によるリポジトリ実装
│   │   └── auth/                # JWT トークン生成・検証
│   └── handler/                 # プレゼンテーション層（Echo ハンドラー）
│       ├── user_handler.go
│       ├── event_handler.go
│       ├── auth_handler.go
│       ├── middleware/           # 認証・権限チェックミドルウェア
│       ├── request/             # リクエスト DTO
│       └── response/            # レスポンス DTO
├── migrations/                  # DBマイグレーションSQL
├── go.mod
└── go.sum
```

### フロントエンド（Next.js）

```
frontend/
├── src/
│   ├── app/                     # App Router（ページ定義のみ、ロジックは薄く）
│   │   ├── (auth)/              # 未認証ページグループ
│   │   │   ├── login/
│   │   │   └── signup/
│   │   ├── (main)/              # 認証済みページグループ
│   │   │   ├── dashboard/
│   │   │   ├── members/
│   │   │   ├── calendar/
│   │   │   ├── events/
│   │   │   ├── profile/
│   │   │   └── admin/
│   │   └── layout.tsx
│   ├── features/                # 機能スライス（機能ごとに完結）
│   │   ├── auth/
│   │   │   ├── components/      # ログインフォーム等
│   │   │   ├── hooks/           # useLogin, useSignup
│   │   │   ├── api/             # 認証APIクライアント
│   │   │   └── types/           # 認証関連の型
│   │   ├── members/
│   │   │   ├── components/      # MemberCard, MemberSearchForm
│   │   │   ├── hooks/           # useMembers, useMemberSearch
│   │   │   ├── api/             # メンバーAPIクライアント
│   │   │   └── types/
│   │   ├── events/
│   │   │   ├── components/      # Calendar, EventCard, EventForm
│   │   │   ├── hooks/           # useEvents, useCreateEvent
│   │   │   ├── api/
│   │   │   └── types/
│   │   └── profile/
│   │       ├── components/
│   │       ├── hooks/
│   │       ├── api/
│   │       └── types/
│   ├── shared/                  # 共通部品
│   │   ├── components/          # Button, Card, Modal 等（shadcn/ui ベース）
│   │   ├── hooks/               # useAuth（認証状態）等
│   │   ├── lib/                 # API クライアント基盤、ユーティリティ
│   │   └── types/               # 共通型定義
│   └── styles/
├── public/
├── package.json
└── tailwind.config.ts
```

---

## データベース設計

### users テーブル

| カラム | 型 | 説明 |
|--------|-----|------|
| id | UUID, PK | |
| email | VARCHAR, UNIQUE, NOT NULL | ログイン用 |
| password_hash | VARCHAR, NOT NULL | bcrypt ハッシュ |
| display_name | VARCHAR, NOT NULL | 表示名 |
| avatar_url | VARCHAR, nullable | プロフィール画像 |
| parts | TEXT[] | 担当パート（例: ["Gt", "Vo"]） |
| year | INTEGER | 学年（1〜4） |
| faculty | VARCHAR | 学部 |
| genres | TEXT[] | 好きなジャンル |
| bio | TEXT | 自己紹介 |
| role | VARCHAR, DEFAULT 'member' | member / admin |
| created_at | TIMESTAMP | |
| updated_at | TIMESTAMP | |

### events テーブル

| カラム | 型 | 説明 |
|--------|-----|------|
| id | UUID, PK | |
| title | VARCHAR, NOT NULL | イベント名 |
| description | TEXT | 説明 |
| event_type | VARCHAR, NOT NULL | live / camp / other |
| venue | VARCHAR | 会場 |
| start_datetime | TIMESTAMP, NOT NULL | 開始日時 |
| end_datetime | TIMESTAMP | 終了日時 |
| created_by | UUID, FK → users.id | 作成者 |
| created_at | TIMESTAMP | |
| updated_at | TIMESTAMP | |

---

## 権限モデル

| ロール | 説明 | できること |
|--------|------|-----------|
| member | 一般メンバー | プロフィール編集（自分のみ）、メンバー検索・閲覧、カレンダー閲覧 |
| admin | 幹部（会長・副会長・会計） | 上記すべて ＋ イベント CRUD、メンバーのロール変更 |

- 権限チェックは Echo のミドルウェアで行う
- 幹部のみのエンドポイントには `AdminOnly` ミドルウェアを適用

---

## API設計

### 認証

```
POST /api/v1/auth/signup    サインアップ
POST /api/v1/auth/login     ログイン
POST /api/v1/auth/refresh   トークンリフレッシュ
POST /api/v1/auth/logout    ログアウト
```

### ユーザー

```
GET  /api/v1/users/me       自分のプロフィール
PUT  /api/v1/users/me       プロフィール更新
GET  /api/v1/users          メンバー一覧（クエリ: ?part=Gt&year=2&genre=ロック&q=田中）
GET  /api/v1/users/:id      メンバー詳細
```

### イベント

```
GET    /api/v1/events          一覧（クエリ: ?month=2026-03）
POST   /api/v1/events          作成（幹部のみ）
GET    /api/v1/events/:id      詳細
PUT    /api/v1/events/:id      更新（幹部のみ）
DELETE /api/v1/events/:id      削除（幹部のみ）
```

### 管理

```
PUT /api/v1/admin/users/:id/role  ロール変更（幹部のみ）
```

---

## ドメイン知識（軽音サークル固有の用語・概念）

- **パート**: 楽器の担当。Vo（ボーカル）、Gt（ギター）、Ba（ベース）、Dr（ドラム）、Key（キーボード）など。1人が複数パートを持てる
- **ライブ**: サークル主催のライブイベント。月1回程度、ライブハウスで開催。ライブごとにバンドを組む（固定バンドではない）
- **箱代**: ライブハウスのレンタル費用。出演者で割り勘する
- **出演費**: 箱代を出演者間で分担した金額（＝出演者が支払うもの）
- **幹部**: 会長・副会長・会計。イベント管理やメンバー管理の権限を持つ
- **合宿**: サークル全体の宿泊イベント

---

## コーディング規約

### 共通

- コメントは日本語で書く
- コミットメッセージは日本語で書く
- コミットメッセージのプレフィックス: `feat:`, `fix:`, `refactor:`, `docs:`, `test:`, `chore:`

### Go

- フォーマッター: `gofmt`
- リンター: `golangci-lint`
- 命名規則: Go の標準に従う（キャメルケース、エクスポートは大文字始まり）
- エラーハンドリング: カスタムエラー型を domain 層で定義し、usecase 層で使う
- テスト: usecase 層を中心にユニットテストを書く。リポジトリはインターフェースでモック可能にする

### TypeScript / React

- フォーマッター: Prettier
- リンター: ESLint
- コンポーネント: アロー関数 + `export const`
- 状態管理: TanStack Query でサーバー状態を管理。ローカル状態は useState / useReducer
- スタイリング: Tailwind CSS のユーティリティクラスのみ（CSS ファイルは原則書かない）

### Git 運用（GitHub Flow）

- `main`: 常にデプロイ可能な状態
- `feature/xxx`: 機能開発ブランチ（例: `feature/メンバー検索`, `feature/カレンダー`）
- `fix/xxx`: バグ修正ブランチ
- 作業完了後、main にマージ

---

## Claude Code への指示

### やってほしいこと

- Clean Architecture の層の分離を厳守すること
- コードを書く際、なぜその設計にしたのか簡潔に説明すること
- domain 層に外部ライブラリ（GORM, Echo など）の import が入らないようにすること
- 新しいファイルを作成する際、どの層に属するか明記すること
- テストを書く際、リポジトリはインターフェースを使ってモックすること

### やらないでほしいこと

- 一度に大量のコードを生成しないこと（段階的に進める）
- 設計の判断を勝手に行わないこと（迷ったら確認する）
- domain 層に infrastructure の詳細（GORM のタグなど）を混ぜないこと
- 過度に複雑なパターンを導入しないこと（学習段階であることを考慮）