# Pact

**AIとの契約 — 仕様から図へ、図からコードへ**

Pact は、ソフトウェアの仕様を記述するためのDSLです。1つの `.pact` ファイルから、クラス図・シーケンス図・ステートマシン図・フローチャートを自動生成します。

---

## なぜ Pact が必要か

AIが生成したコードのレビューは難しい。

- AIは大量のコードを瞬時に生成する
- しかしコードは詳細すぎて、全体像が見えない
- レビューに時間がかかり、見落としが発生する

Pact は発想を逆転させます。

**人間が書くのは「仕様」だけ。コードはAIに任せる。**

```
人間 → Pact仕様 → 図（視覚的検証）→ コード（AI生成）
```

---

## 思想

### 仕様は1箇所に

設計書、UML図、コードコメント、テストケース。これらは容易に乖離します。Pact はすべてを1つのファイルに統合します。

### 複数の視点、同一の真実

1つの `.pact` ファイルから4種類の図を生成します。

- **クラス図** — 何が何に依存しているか
- **シーケンス図** — 誰が誰と、どの順序で話すか
- **ステートマシン図** — どんな状態があり、何が遷移を起こすか
- **フローチャート** — どんな手順で処理が進むか

図の間に矛盾があれば、それは仕様の問題です。

### コードと仕様の1対1対応

```
src/auth/service.go      ←→  .pact/auth/service.pact
src/order/handler.go     ←→  .pact/order/handler.pact
```

どのコードがどの仕様に対応するか一目瞭然。仕様の漏れも検出できます。

### AIとの契約

Pact（協定・契約）という名前は、この言語の本質を表しています。

- 人間は Pact ファイルで「何を作るか」を定義する
- AIは Pact ファイルに従ってコードを生成する

レビューすべきはコードではなく、契約です。

---

## プロジェクト構造

```
my-project/
├── .pactconfig              # 設定ファイル
├── src/                     # コード本体
│   ├── auth/
│   │   └── service.go
│   └── order/
│       └── handler.go
└── .pact/                   # Pact 仕様（コードと同じ階層構造）
    ├── auth/
    │   └── service.pact
    └── order/
        └── handler.pact
```

`.pactconfig` でコード本体のルートを指定します。

```yaml
source_root: ./src
pact_root: ./.pact
```

---

## 例

```pact
@description("ユーザー認証サービス")
component AuthService {

  type Token {
    +value: string
    +expiresAt: datetime
  }

  depends on UserRepository: database
  depends on TokenGenerator

  provides AuthAPI {
    Login(email: string, password: string) -> Token
      throws InvalidCredentials, UserNotFound
  }

  flow Login {
    user = UserRepository.FindByEmail(email) 
      ?? throw UserNotFound
    
    if !verify(password, user.passwordHash) {
      throw InvalidCredentials
    }
    
    token = TokenGenerator.Generate(user.id)
    return token
  }

  states AuthState {
    initial Idle
    
    Idle -> Authenticating on LoginRequest
    Authenticating -> Authenticated on Success
    Authenticating -> Failed on Failure
    Failed -> Idle after 3s
  }
}
```

---

## 使い方

```bash
# 初期化
pact init

# 図を生成
pact generate

# 構文チェック
pact validate

# 仕様がないコードを検出
pact check --missing

# ファイル監視
pact watch
```

---

## インストール

```bash
go install github.com/example/pact@latest
```

---

## ドキュメント

- [サンプルギャラリー](https://tubasasakunn.github.io/pact/sample/) — 生成されたSVGダイアグラムの一覧

---

## ライセンス

MIT License

