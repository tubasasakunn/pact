# Pact テスト項目書

## テスト構成

```
pact/
├── internal/
│   ├── domain/
│   │   ├── ast/
│   │   │   └── visitor_test.go
│   │   ├── config/
│   │   │   └── config_test.go
│   │   ├── diagram/
│   │   │   ├── common/
│   │   │   │   └── types_test.go
│   │   │   ├── class/
│   │   │   │   └── model_test.go
│   │   │   ├── sequence/
│   │   │   │   └── model_test.go
│   │   │   ├── state/
│   │   │   │   └── model_test.go
│   │   │   └── flow/
│   │   │       └── model_test.go
│   │   └── errors/
│   │       └── errors_test.go
│   │
│   ├── application/
│   │   ├── parser/
│   │   │   └── service_test.go
│   │   ├── project/
│   │   │   └── service_test.go
│   │   ├── transformer/
│   │   │   ├── class_test.go
│   │   │   ├── sequence_test.go
│   │   │   ├── state_test.go
│   │   │   └── flow_test.go
│   │   └── renderer/
│   │       └── service_test.go
│   │
│   ├── infrastructure/
│   │   ├── config/
│   │   │   └── loader_test.go
│   │   ├── parser/
│   │   │   ├── lexer_test.go
│   │   │   ├── parser_test.go
│   │   │   └── adapter_test.go
│   │   ├── resolver/
│   │   │   └── import_test.go
│   │   ├── renderer/
│   │   │   ├── svg/
│   │   │   │   ├── renderer_test.go
│   │   │   │   ├── class_test.go
│   │   │   │   ├── sequence_test.go
│   │   │   │   ├── state_test.go
│   │   │   │   └── flow_test.go
│   │   │   └── canvas/
│   │   │       ├── canvas_test.go
│   │   │       ├── shapes_test.go
│   │   │       └── text_test.go
│   │   └── filesystem/
│   │       ├── reader_test.go
│   │       └── writer_test.go
│   │
│   └── interfaces/
│       └── cli/
│           ├── root_test.go
│           ├── init_test.go
│           ├── generate_test.go
│           ├── validate_test.go
│           ├── check_test.go
│           └── watch_test.go
│
├── pkg/
│   └── pact/
│       └── api_test.go
│
└── test/
    ├── integration/
    │   ├── parse_transform_test.go
    │   ├── import_resolution_test.go
    │   └── full_pipeline_test.go
    └── e2e/
        └── cli_test.go
```

---

# 1. 単体テスト（Unit Tests）

## 1.1 Lexer テスト

### ファイル: `internal/infrastructure/parser/lexer_test.go`

```go
package parser

import "testing"
```

### 1.1.1 基本トークン認識

| ID | テスト名 | 入力 | 期待出力 | 目的 |
|----|---------|------|----------|------|
| L001 | TestLexer_EOF | `""` | `TOKEN_EOF` | 空入力でEOF |
| L002 | TestLexer_Identifier_Simple | `"foo"` | `TOKEN_IDENT("foo")` | 単純な識別子 |
| L003 | TestLexer_Identifier_WithUnderscore | `"user_id"` | `TOKEN_IDENT("user_id")` | アンダースコア含む識別子 |
| L004 | TestLexer_Identifier_WithNumbers | `"token2"` | `TOKEN_IDENT("token2")` | 数字含む識別子 |
| L005 | TestLexer_Identifier_StartWithNumber | `"2token"` | `TOKEN_INT(2), TOKEN_IDENT("token")` | 数字始まりは識別子でない |

### 1.1.2 キーワード認識

| ID | テスト名 | 入力 | 期待出力 | 目的 |
|----|---------|------|----------|------|
| L010 | TestLexer_Keyword_Component | `"component"` | `TOKEN_COMPONENT` | componentキーワード |
| L011 | TestLexer_Keyword_Import | `"import"` | `TOKEN_IMPORT` | importキーワード |
| L012 | TestLexer_Keyword_Type | `"type"` | `TOKEN_TYPE` | typeキーワード |
| L013 | TestLexer_Keyword_Enum | `"enum"` | `TOKEN_ENUM` | enumキーワード |
| L014 | TestLexer_Keyword_Depends | `"depends"` | `TOKEN_DEPENDS` | dependsキーワード |
| L015 | TestLexer_Keyword_On | `"on"` | `TOKEN_ON` | onキーワード |
| L016 | TestLexer_Keyword_Extends | `"extends"` | `TOKEN_EXTENDS` | extendsキーワード |
| L017 | TestLexer_Keyword_Implements | `"implements"` | `TOKEN_IMPLEMENTS` | implementsキーワード |
| L018 | TestLexer_Keyword_Contains | `"contains"` | `TOKEN_CONTAINS` | containsキーワード |
| L019 | TestLexer_Keyword_Aggregates | `"aggregates"` | `TOKEN_AGGREGATES` | aggregatesキーワード |
| L020 | TestLexer_Keyword_Provides | `"provides"` | `TOKEN_PROVIDES` | providesキーワード |
| L021 | TestLexer_Keyword_Requires | `"requires"` | `TOKEN_REQUIRES` | requiresキーワード |
| L022 | TestLexer_Keyword_Flow | `"flow"` | `TOKEN_FLOW` | flowキーワード |
| L023 | TestLexer_Keyword_States | `"states"` | `TOKEN_STATES` | statesキーワード |
| L024 | TestLexer_Keyword_State | `"state"` | `TOKEN_STATE` | stateキーワード |
| L025 | TestLexer_Keyword_Parallel | `"parallel"` | `TOKEN_PARALLEL` | parallelキーワード |
| L026 | TestLexer_Keyword_Region | `"region"` | `TOKEN_REGION` | regionキーワード |
| L027 | TestLexer_Keyword_Initial | `"initial"` | `TOKEN_INITIAL` | initialキーワード |
| L028 | TestLexer_Keyword_Final | `"final"` | `TOKEN_FINAL` | finalキーワード |
| L029 | TestLexer_Keyword_Entry | `"entry"` | `TOKEN_ENTRY` | entryキーワード |
| L030 | TestLexer_Keyword_Exit | `"exit"` | `TOKEN_EXIT` | exitキーワード |
| L031 | TestLexer_Keyword_If | `"if"` | `TOKEN_IF` | ifキーワード |
| L032 | TestLexer_Keyword_Else | `"else"` | `TOKEN_ELSE` | elseキーワード |
| L033 | TestLexer_Keyword_For | `"for"` | `TOKEN_FOR` | forキーワード |
| L034 | TestLexer_Keyword_In | `"in"` | `TOKEN_IN` | inキーワード |
| L035 | TestLexer_Keyword_While | `"while"` | `TOKEN_WHILE` | whileキーワード |
| L036 | TestLexer_Keyword_Return | `"return"` | `TOKEN_RETURN` | returnキーワード |
| L037 | TestLexer_Keyword_Throw | `"throw"` | `TOKEN_THROW` | throwキーワード |
| L038 | TestLexer_Keyword_Await | `"await"` | `TOKEN_AWAIT` | awaitキーワード |
| L039 | TestLexer_Keyword_Async | `"async"` | `TOKEN_ASYNC` | asyncキーワード |
| L040 | TestLexer_Keyword_Throws | `"throws"` | `TOKEN_THROWS` | throwsキーワード |
| L041 | TestLexer_Keyword_When | `"when"` | `TOKEN_WHEN` | whenキーワード |
| L042 | TestLexer_Keyword_After | `"after"` | `TOKEN_AFTER` | afterキーワード |
| L043 | TestLexer_Keyword_Do | `"do"` | `TOKEN_DO` | doキーワード |
| L044 | TestLexer_Keyword_True | `"true"` | `TOKEN_TRUE` | trueキーワード |
| L045 | TestLexer_Keyword_False | `"false"` | `TOKEN_FALSE` | falseキーワード |
| L046 | TestLexer_Keyword_Null | `"null"` | `TOKEN_NULL` | nullキーワード |
| L047 | TestLexer_Keyword_As | `"as"` | `TOKEN_AS` | asキーワード |

### 1.1.3 リテラル

| ID | テスト名 | 入力 | 期待出力 | 目的 |
|----|---------|------|----------|------|
| L050 | TestLexer_String_Simple | `"\"hello\""` | `TOKEN_STRING("hello")` | 単純な文字列 |
| L051 | TestLexer_String_Empty | `"\"\""` | `TOKEN_STRING("")` | 空文字列 |
| L052 | TestLexer_String_Escape_Quote | `"\"say \\\"hi\\\"\""` | `TOKEN_STRING("say \"hi\"")` | エスケープ（引用符） |
| L053 | TestLexer_String_Escape_Backslash | `"\"path\\\\to\""` | `TOKEN_STRING("path\\to")` | エスケープ（バックスラッシュ） |
| L054 | TestLexer_String_Escape_Newline | `"\"line1\\nline2\""` | `TOKEN_STRING("line1\nline2")` | エスケープ（改行） |
| L055 | TestLexer_String_Escape_Tab | `"\"col1\\tcol2\""` | `TOKEN_STRING("col1\tcol2")` | エスケープ（タブ） |
| L056 | TestLexer_String_Unterminated | `"\"hello"` | `TOKEN_STRING("hello"), EOF` | 終端なし文字列（EOFで終了） |
| L060 | TestLexer_Int_Zero | `"0"` | `TOKEN_INT(0)` | ゼロ |
| L061 | TestLexer_Int_Positive | `"42"` | `TOKEN_INT(42)` | 正の整数 |
| L062 | TestLexer_Int_Large | `"1234567890"` | `TOKEN_INT(1234567890)` | 大きな整数 |
| L065 | TestLexer_Float_Simple | `"3.14"` | `TOKEN_FLOAT(3.14)` | 単純な浮動小数点 |
| L066 | TestLexer_Float_LeadingZero | `"0.5"` | `TOKEN_FLOAT(0.5)` | 先頭ゼロ |
| L067 | TestLexer_Float_TrailingDot | `"1."` | `TOKEN_INT(1), TOKEN_DOT` | 末尾ドット（floatではない） |
| L068 | TestLexer_Float_LeadingDot | `".5"` | `TOKEN_DOT, TOKEN_INT(5)` | 先頭ドット（floatではない） |
| L070 | TestLexer_Duration_Milliseconds | `"500ms"` | `TOKEN_DURATION(500, "ms")` | ミリ秒 |
| L071 | TestLexer_Duration_Seconds | `"30s"` | `TOKEN_DURATION(30, "s")` | 秒 |
| L072 | TestLexer_Duration_Minutes | `"5m"` | `TOKEN_DURATION(5, "m")` | 分 |
| L073 | TestLexer_Duration_Hours | `"24h"` | `TOKEN_DURATION(24, "h")` | 時間 |
| L074 | TestLexer_Duration_Days | `"7d"` | `TOKEN_DURATION(7, "d")` | 日 |

### 1.1.4 演算子・記号

| ID | テスト名 | 入力 | 期待出力 | 目的 |
|----|---------|------|----------|------|
| L080 | TestLexer_Symbol_LBrace | `"{"` | `TOKEN_LBRACE` | 左波括弧 |
| L081 | TestLexer_Symbol_RBrace | `"}"` | `TOKEN_RBRACE` | 右波括弧 |
| L082 | TestLexer_Symbol_LParen | `"("` | `TOKEN_LPAREN` | 左丸括弧 |
| L083 | TestLexer_Symbol_RParen | `")"` | `TOKEN_RPAREN` | 右丸括弧 |
| L084 | TestLexer_Symbol_LBracket | `"["` | `TOKEN_LBRACKET` | 左角括弧 |
| L085 | TestLexer_Symbol_RBracket | `"]"` | `TOKEN_RBRACKET` | 右角括弧 |
| L086 | TestLexer_Symbol_Colon | `":"` | `TOKEN_COLON` | コロン |
| L087 | TestLexer_Symbol_Comma | `","` | `TOKEN_COMMA` | カンマ |
| L088 | TestLexer_Symbol_Dot | `"."` | `TOKEN_DOT` | ドット |
| L089 | TestLexer_Symbol_Arrow | `"->"` | `TOKEN_ARROW` | 矢印 |
| L090 | TestLexer_Symbol_At | `"@"` | `TOKEN_AT` | アットマーク |
| L091 | TestLexer_Symbol_Question | `"?"` | `TOKEN_QUESTION` | クエスチョン |
| L092 | TestLexer_Symbol_Plus | `"+"` | `TOKEN_PLUS` | プラス |
| L093 | TestLexer_Symbol_Minus | `"-"` | `TOKEN_MINUS` | マイナス |
| L094 | TestLexer_Symbol_Star | `"*"` | `TOKEN_STAR` | アスタリスク |
| L095 | TestLexer_Symbol_Slash | `"/"` | `TOKEN_SLASH` | スラッシュ |
| L096 | TestLexer_Symbol_Percent | `"%"` | `TOKEN_PERCENT` | パーセント |
| L097 | TestLexer_Symbol_Eq | `"=="` | `TOKEN_EQ` | 等価 |
| L098 | TestLexer_Symbol_Ne | `"!="` | `TOKEN_NE` | 非等価 |
| L099 | TestLexer_Symbol_Lt | `"<"` | `TOKEN_LT` | 小なり |
| L100 | TestLexer_Symbol_Gt | `">"` | `TOKEN_GT` | 大なり |
| L101 | TestLexer_Symbol_Le | `"<="` | `TOKEN_LE` | 以下 |
| L102 | TestLexer_Symbol_Ge | `">="` | `TOKEN_GE` | 以上 |
| L103 | TestLexer_Symbol_And | `"&&"` | `TOKEN_AND` | 論理積 |
| L104 | TestLexer_Symbol_Or | `"\|\|"` | `TOKEN_OR` | 論理和 |
| L105 | TestLexer_Symbol_Not | `"!"` | `TOKEN_NOT` | 否定 |
| L106 | TestLexer_Symbol_Assign | `"="` | `TOKEN_ASSIGN` | 代入 |
| L107 | TestLexer_Symbol_Nullish | `"??"` | `TOKEN_NULLISH` | null合体 |
| L108 | TestLexer_Symbol_Hash | `"#"` | `TOKEN_HASH` | ハッシュ |
| L109 | TestLexer_Symbol_Tilde | `"~"` | `TOKEN_TILDE` | チルダ |

### 1.1.5 コメント

| ID | テスト名 | 入力 | 期待出力 | 目的 |
|----|---------|------|----------|------|
| L120 | TestLexer_Comment_Line | `"// comment\nfoo"` | `TOKEN_IDENT("foo")` | 行コメントスキップ |
| L121 | TestLexer_Comment_Line_EOF | `"// comment"` | `TOKEN_EOF` | 行コメントで終了 |
| L122 | TestLexer_Comment_Block | `"/* comment */foo"` | `TOKEN_IDENT("foo")` | ブロックコメントスキップ |
| L123 | TestLexer_Comment_Block_Multiline | `"/* line1\nline2 */foo"` | `TOKEN_IDENT("foo")` | 複数行ブロックコメント |
| L124 | TestLexer_Comment_Block_NoNesting | `"/* /* inner */ outer"` | `TOKEN_IDENT("outer")` | ネスト非対応（最初の*/で終了） |
| L125 | TestLexer_Comment_Block_Unterminated | `"/* comment"` | `TOKEN_EOF` | 終端なしブロックコメント（EOFで終了） |

### 1.1.6 空白・位置情報

| ID | テスト名 | 入力 | 期待出力 | 目的 |
|----|---------|------|----------|------|
| L130 | TestLexer_Whitespace_Spaces | `"  foo  "` | `TOKEN_IDENT("foo")` | スペーススキップ |
| L131 | TestLexer_Whitespace_Tabs | `"\t\tfoo\t"` | `TOKEN_IDENT("foo")` | タブスキップ |
| L132 | TestLexer_Whitespace_Newlines | `"\n\nfoo\n"` | `TOKEN_IDENT("foo")` | 改行スキップ |
| L133 | TestLexer_Whitespace_Mixed | `"  \t\n  foo"` | `TOKEN_IDENT("foo")` | 混合空白スキップ |
| L140 | TestLexer_Position_FirstToken | `"foo"` | `Line:1, Column:1` | 最初のトークン位置 |
| L141 | TestLexer_Position_AfterWhitespace | `"  foo"` | `Line:1, Column:3` | 空白後の位置 |
| L142 | TestLexer_Position_AfterNewline | `"\nfoo"` | `Line:2, Column:1` | 改行後の位置 |
| L143 | TestLexer_Position_MultipleTokens | `"foo bar"` | `{1,1}, {1,5}` | 複数トークンの位置 |
| L144 | TestLexer_Position_MultipleLines | `"foo\nbar\nbaz"` | `{1,1}, {2,1}, {3,1}` | 複数行の位置 |

### 1.1.7 複合トークン列

| ID | テスト名 | 入力 | 期待出力 | 目的 |
|----|---------|------|----------|------|
| L150 | TestLexer_Sequence_TypeDecl | `"type Foo { }"` | `TYPE, IDENT, LBRACE, RBRACE` | 型宣言 |
| L151 | TestLexer_Sequence_MethodDecl | `"Login(email: string) -> Token"` | 適切なトークン列 | メソッド宣言 |
| L152 | TestLexer_Sequence_Annotation | `"@description(\"hello\")"` | `AT, IDENT, LPAREN, STRING, RPAREN` | アノテーション |
| L153 | TestLexer_Sequence_Transition | `"Idle -> Active on Start"` | 適切なトークン列 | 状態遷移 |

---

## 1.2 Parser テスト

### ファイル: `internal/infrastructure/parser/parser_test.go`

### 1.2.1 インポート文

| ID | テスト名 | 入力 | 期待AST | 目的 |
|----|---------|------|---------|------|
| P001 | TestParser_Import_Simple | `import "./foo.pact"` | `ImportDecl{Path:"./foo.pact", Alias:nil}` | 単純インポート |
| P002 | TestParser_Import_WithAlias | `import "./foo.pact" as Foo` | `ImportDecl{Path:"./foo.pact", Alias:"Foo"}` | エイリアス付きインポート |
| P003 | TestParser_Import_Multiple | 複数import | 複数ImportDecl | 複数インポート |
| P004 | TestParser_Import_MissingPath | `import` | ParseError | パスなしエラー |
| P005 | TestParser_Import_InvalidPath | `import 123` | ParseError | 不正なパス |
| P006 | TestParser_Import_Position | `import "./foo.pact"` | `Pos{Line:1, Column:1}` | 位置情報 |

### 1.2.2 コンポーネント宣言

| ID | テスト名 | 入力 | 期待AST | 目的 |
|----|---------|------|---------|------|
| P010 | TestParser_Component_Empty | `component Foo { }` | `ComponentDecl{Name:"Foo"}` | 空コンポーネント |
| P011 | TestParser_Component_WithAnnotation | `@description("...") component Foo { }` | アノテーション付き | アノテーション |
| P012 | TestParser_Component_MissingName | `component { }` | ParseError | 名前なしエラー |
| P013 | TestParser_Component_MissingBrace | `component Foo` | ParseError | 波括弧なしエラー |
| P014 | TestParser_Component_Position | `component Foo { }` | `Pos{Line:1, Column:1}` | 位置情報 |

### 1.2.3 型定義

| ID | テスト名 | 入力 | 期待AST | 目的 |
|----|---------|------|---------|------|
| P020 | TestParser_Type_Empty | `type Foo { }` | `TypeDecl{Name:"Foo", Kind:struct, Fields:[]}` | 空の型 |
| P021 | TestParser_Type_SingleField | `type Foo { name: string }` | 1フィールド | 単一フィールド |
| P022 | TestParser_Type_MultipleFields | `type Foo { a: int b: string }` | 複数フィールド | 複数フィールド |
| P023 | TestParser_Type_NullableField | `type Foo { name: string? }` | nullable型 | nullable |
| P024 | TestParser_Type_ArrayField | `type Foo { items: string[] }` | 配列型 | 配列 |
| P025 | TestParser_Type_NullableArrayField | `type Foo { items: string?[] }` | nullable配列型 | nullable配列 |
| P026 | TestParser_Type_VisibilityPublic | `type Foo { +name: string }` | `Visibility:public` | public |
| P027 | TestParser_Type_VisibilityPrivate | `type Foo { -name: string }` | `Visibility:private` | private |
| P028 | TestParser_Type_VisibilityProtected | `type Foo { #name: string }` | `Visibility:protected` | protected |
| P029 | TestParser_Type_VisibilityPackage | `type Foo { ~name: string }` | `Visibility:package` | package |
| P030 | TestParser_Type_FieldAnnotation | `type Foo { @desc("x") name: string }` | フィールドにアノテーション | フィールドアノテーション |
| P031 | TestParser_Type_FieldPosition | `type Foo { name: string }` | フィールドに位置情報 | フィールド位置 |
| P035 | TestParser_Enum_Simple | `enum Status { A B C }` | `TypeDecl{Kind:enum, Values:[A,B,C]}` | 単純enum |
| P036 | TestParser_Enum_Empty | `enum Status { }` | 空のenum | 空enum |
| P037 | TestParser_Enum_WithAnnotation | `@desc("x") enum Status { A }` | アノテーション付き | enumアノテーション |

### 1.2.4 関係定義

| ID | テスト名 | 入力 | 期待AST | 目的 |
|----|---------|------|---------|------|
| P040 | TestParser_Relation_DependsOn | `depends on Foo` | `RelationDecl{Kind:depends, Target:Foo}` | depends on |
| P041 | TestParser_Relation_DependsOn_WithType | `depends on Foo: database` | `TargetType:database` | 型指定 |
| P042 | TestParser_Relation_DependsOn_WithAlias | `depends on Foo as F` | `Alias:F` | エイリアス |
| P043 | TestParser_Relation_DependsOn_Full | `depends on Foo: external as F` | 全指定 | フル指定 |
| P044 | TestParser_Relation_Extends | `extends Base` | `RelationDecl{Kind:extends}` | extends |
| P045 | TestParser_Relation_Implements | `implements Iface` | `RelationDecl{Kind:implements}` | implements |
| P046 | TestParser_Relation_Contains | `contains Cache` | `RelationDecl{Kind:contains}` | contains |
| P047 | TestParser_Relation_Aggregates | `aggregates Items` | `RelationDecl{Kind:aggregates}` | aggregates |
| P048 | TestParser_Relation_WithAnnotation | `@desc("x") depends on Foo` | アノテーション付き | アノテーション |
| P049 | TestParser_Relation_AllTargetTypes | 各TargetType | 正しく認識 | 全TargetType |
| P04A | TestParser_Relation_Position | `depends on Foo` | 位置情報付き | 位置情報 |

### 1.2.5 インターフェース定義

| ID | テスト名 | 入力 | 期待AST | 目的 |
|----|---------|------|---------|------|
| P050 | TestParser_Interface_Provides_Empty | `provides API { }` | 空Interface | 空provides |
| P051 | TestParser_Interface_Requires_Empty | `requires Query { }` | 空Interface | 空requires |
| P052 | TestParser_Interface_SingleMethod | `provides API { Get() -> Item }` | 1メソッド | 単一メソッド |
| P053 | TestParser_Interface_MethodWithParams | `provides API { Get(id: string) -> Item }` | パラメータ付き | パラメータ |
| P054 | TestParser_Interface_MethodMultiParams | `provides API { Get(a: int, b: string) -> Item }` | 複数パラメータ | 複数パラメータ |
| P055 | TestParser_Interface_MethodThrows | `provides API { Get() -> Item throws NotFound }` | throws付き | throws |
| P056 | TestParser_Interface_MethodMultiThrows | `provides API { Get() -> Item throws A, B }` | 複数throws | 複数throws |
| P057 | TestParser_Interface_MethodAsync | `provides API { async Send() -> void }` | async付き | async |
| P058 | TestParser_Interface_MethodVoid | `provides API { Delete() -> void }` | void戻り値 | void |
| P059 | TestParser_Interface_MethodAnnotation | `provides API { @desc("x") Get() -> Item }` | アノテーション付き | メソッドアノテーション |
| P060 | TestParser_Interface_MultipleMethod | 複数メソッド | 複数MethodDecl | 複数メソッド |
| P061 | TestParser_Interface_MethodPosition | メソッド | 位置情報付き | 位置情報 |

### 1.2.6 フロー定義

| ID | テスト名 | 入力 | 期待AST | 目的 |
|----|---------|------|---------|------|
| P070 | TestParser_Flow_Empty | `flow Foo { }` | 空Flow | 空フロー |
| P071 | TestParser_Flow_Assign_Simple | `flow F { x = y }` | AssignStep | 単純代入 |
| P072 | TestParser_Flow_Assign_Call | `flow F { x = A.B() }` | メソッド呼び出し代入 | 呼び出し代入 |
| P073 | TestParser_Flow_Assign_CallWithArgs | `flow F { x = A.B(c, d) }` | 引数付き呼び出し | 引数付き |
| P074 | TestParser_Flow_Assign_ThrowOnNull | `flow F { x = A.B() ?? throw E }` | null合体throw | ??throw |
| P075 | TestParser_Flow_Call | `flow F { A.B() }` | CallStep | 単純呼び出し |
| P076 | TestParser_Flow_Call_Await | `flow F { await A.B() }` | async呼び出し | await |
| P077 | TestParser_Flow_Return | `flow F { return x }` | ReturnStep | return |
| P078 | TestParser_Flow_Return_Empty | `flow F { return }` | 値なしreturn | 空return |
| P079 | TestParser_Flow_Throw | `flow F { throw E }` | ThrowStep | throw |
| P080 | TestParser_Flow_If_Simple | `flow F { if x { } }` | IfStep(thenのみ) | 単純if |
| P081 | TestParser_Flow_If_Else | `flow F { if x { } else { } }` | IfStep(else付き) | if-else |
| P082 | TestParser_Flow_If_Nested | `flow F { if x { if y { } } }` | ネストif | ネストif |
| P083 | TestParser_Flow_For | `flow F { for x in items { } }` | ForStep | for |
| P084 | TestParser_Flow_For_Nested | `flow F { for x in a { for y in b { } } }` | ネストfor | ネストfor |
| P085 | TestParser_Flow_While | `flow F { while cond { } }` | WhileStep | while |
| P086 | TestParser_Flow_Complex | 複合フロー | 複合AST | 複合フロー |
| P087 | TestParser_Flow_StepAnnotation | `flow F { @desc("x") a = b }` | ステップにアノテーション | ステップアノテーション |
| P088 | TestParser_Flow_StepPosition | フローステップ | 位置情報付き | 位置情報 |

### 1.2.7 式

| ID | テスト名 | 入力 | 期待AST | 目的 |
|----|---------|------|---------|------|
| P100 | TestParser_Expr_Literal_String | `"hello"` | `LiteralExpr{Value:"hello"}` | 文字列リテラル |
| P101 | TestParser_Expr_Literal_Int | `42` | `LiteralExpr{Value:42}` | 整数リテラル |
| P102 | TestParser_Expr_Literal_Float | `3.14` | `LiteralExpr{Value:3.14}` | 浮動小数点リテラル |
| P103 | TestParser_Expr_Literal_True | `true` | `LiteralExpr{Value:true}` | trueリテラル |
| P104 | TestParser_Expr_Literal_False | `false` | `LiteralExpr{Value:false}` | falseリテラル |
| P105 | TestParser_Expr_Literal_Null | `null` | `LiteralExpr{Value:nil}` | nullリテラル |
| P106 | TestParser_Expr_Variable | `foo` | `VariableExpr{Name:"foo"}` | 変数 |
| P107 | TestParser_Expr_Field | `a.b` | `FieldExpr{Object:a, Field:b}` | フィールドアクセス |
| P108 | TestParser_Expr_Field_Chain | `a.b.c` | ネストFieldExpr | チェーン |
| P109 | TestParser_Expr_Call | `A.B()` | `CallExpr` | 呼び出し |
| P110 | TestParser_Expr_Binary_Add | `a + b` | `BinaryExpr{Op:"+"}` | 加算 |
| P111 | TestParser_Expr_Binary_Sub | `a - b` | `BinaryExpr{Op:"-"}` | 減算 |
| P112 | TestParser_Expr_Binary_Mul | `a * b` | `BinaryExpr{Op:"*"}` | 乗算 |
| P113 | TestParser_Expr_Binary_Div | `a / b` | `BinaryExpr{Op:"/"}` | 除算 |
| P114 | TestParser_Expr_Binary_Mod | `a % b` | `BinaryExpr{Op:"%"}` | 剰余 |
| P115 | TestParser_Expr_Binary_Eq | `a == b` | `BinaryExpr{Op:"=="}` | 等価 |
| P116 | TestParser_Expr_Binary_Ne | `a != b` | `BinaryExpr{Op:"!="}` | 非等価 |
| P117 | TestParser_Expr_Binary_Lt | `a < b` | `BinaryExpr{Op:"<"}` | 小なり |
| P118 | TestParser_Expr_Binary_Gt | `a > b` | `BinaryExpr{Op:">"}` | 大なり |
| P119 | TestParser_Expr_Binary_Le | `a <= b` | `BinaryExpr{Op:"<="}` | 以下 |
| P120 | TestParser_Expr_Binary_Ge | `a >= b` | `BinaryExpr{Op:">="}` | 以上 |
| P121 | TestParser_Expr_Binary_And | `a && b` | `BinaryExpr{Op:"&&"}` | 論理積 |
| P122 | TestParser_Expr_Binary_Or | `a \|\| b` | `BinaryExpr{Op:"\|\|"}` | 論理和 |
| P123 | TestParser_Expr_Unary_Not | `!a` | `UnaryExpr{Op:"!"}` | 否定 |
| P124 | TestParser_Expr_Unary_Neg | `-a` | `UnaryExpr{Op:"-"}` | 負数 |
| P125 | TestParser_Expr_Ternary | `a ? b : c` | `TernaryExpr` | 三項演算子 |
| P126 | TestParser_Expr_Paren | `(a + b)` | 括弧内の式 | 括弧 |
| P127 | TestParser_Expr_Precedence_MulAdd | `a + b * c` | `+(a, *(b,c))` | 優先順位（乗算優先） |
| P128 | TestParser_Expr_Precedence_AndOr | `a \|\| b && c` | `\|\|(a, &&(b,c))` | 優先順位（AND優先） |
| P129 | TestParser_Expr_Precedence_Compare | `a == b && c < d` | 正しい優先順位 | 比較と論理 |
| P130 | TestParser_Expr_Complex | `a.b + c.d() * 2` | 複合式 | 複合式 |
| P131 | TestParser_Expr_Position | 各式 | 位置情報付き | 位置情報 |

### 1.2.8 ステートマシン定義

| ID | テスト名 | 入力 | 期待AST | 目的 |
|----|---------|------|---------|------|
| P140 | TestParser_States_Empty | `states S { initial I }` | 最小StatesDecl | 最小ステートマシン |
| P141 | TestParser_States_WithFinal | `states S { initial I final F }` | final付き | final |
| P142 | TestParser_States_MultipleFinal | `states S { initial I final F1 final F2 }` | 複数final | 複数final |
| P143 | TestParser_States_Transition_OnEvent | `... I -> A on E` | `Trigger:EventTrigger{E}` | イベントトリガー |
| P144 | TestParser_States_Transition_After | `... I -> A after 3s` | `Trigger:AfterTrigger{3,s}` | 時間トリガー |
| P145 | TestParser_States_Transition_When | `... I -> A when cond` | `Trigger:WhenTrigger{cond}` | 条件トリガー |
| P146 | TestParser_States_Transition_Guard | `... I -> A on E when g` | `Guard:g` | ガード条件 |
| P147 | TestParser_States_Transition_Actions | `... I -> A on E do [a, b]` | `Actions:[a,b]` | アクション |
| P148 | TestParser_States_Transition_Full | `... I -> A on E when g do [a]` | フル指定 | フル遷移 |
| P149 | TestParser_States_State_Empty | `states S { initial I state I { } }` | 空状態定義 | 空状態 |
| P150 | TestParser_States_State_Entry | `... state I { entry [a] }` | entry付き | entry |
| P151 | TestParser_States_State_Exit | `... state I { exit [a] }` | exit付き | exit |
| P152 | TestParser_States_State_EntryExit | `... state I { entry [a] exit [b] }` | 両方付き | entry/exit |
| P153 | TestParser_States_Compound | ネストstate | 階層StateDecl | 階層状態 |
| P154 | TestParser_States_Compound_Initial | 内部initial | 内部initial | 階層initial |
| P155 | TestParser_States_Parallel | `parallel P { region R1 { } region R2 { } }` | ParallelStateDecl | 並行状態 |
| P156 | TestParser_States_Region | `region R { initial I I -> A on E }` | RegionDecl | リージョン |
| P157 | TestParser_States_Complex | 複合ステートマシン | 複合AST | 複合 |
| P158 | TestParser_States_Position | ステートマシン | 位置情報付き | 位置情報 |

### 1.2.9 アノテーション

| ID | テスト名 | 入力 | 期待AST | 目的 |
|----|---------|------|---------|------|
| P160 | TestParser_Annotation_Flag | `@deprecated` | `AnnotationDecl{Name:deprecated, Args:[]}` | フラグアノテーション |
| P161 | TestParser_Annotation_SingleValue | `@description("hello")` | `Args:[{Value:"hello"}]` | 単一値 |
| P162 | TestParser_Annotation_MultiValue | `@throws("A", "B")` | `Args:[{Value:"A"},{Value:"B"}]` | 複数値 |
| P163 | TestParser_Annotation_Named | `@source(file: "a.go")` | `Args:[{Key:"file",Value:"a.go"}]` | 名前付き |
| P164 | TestParser_Annotation_Mixed | `@foo("a", b: "c")` | 混合 | 混合 |
| P165 | TestParser_Annotation_Multiple | 複数アノテーション | 複数 | 複数アノテーション |
| P166 | TestParser_Annotation_Position | アノテーション | 位置情報付き | 位置情報 |

### 1.2.10 エラーケース

| ID | テスト名 | 入力 | 期待エラー | 目的 |
|----|---------|------|------------|------|
| P200 | TestParser_Error_UnexpectedToken | `component { }` | "expected identifier" | 予期しないトークン |
| P201 | TestParser_Error_UnterminatedBrace | `component Foo {` | "expected }" | 括弧閉じ忘れ |
| P202 | TestParser_Error_InvalidExpression | `flow F { x = }` | "expected expression" | 不正な式 |
| P203 | TestParser_Error_Position | 各エラー | 正確な位置情報 | エラー位置 |
| P204 | TestParser_Error_Type | 各エラー | `*errors.ParseError` | エラー型 |

---

## 1.3 Parser Adapter テスト

### ファイル: `internal/infrastructure/parser/adapter_test.go`

| ID | テスト名 | 操作 | 期待出力 | 目的 |
|----|---------|------|----------|------|
| PA001 | TestFileParser_ParseFile_Success | 有効なファイル | SpecFile | ファイルパース成功 |
| PA002 | TestFileParser_ParseFile_NotFound | 存在しないファイル | エラー | ファイル不在 |
| PA003 | TestFileParser_ParseFile_Invalid | 不正な構文 | ParseError | 構文エラー |
| PA004 | TestFileParser_ParseString_Success | 有効な文字列 | SpecFile | 文字列パース成功 |
| PA005 | TestFileParser_ParseString_Invalid | 不正な文字列 | ParseError | 構文エラー |
| PA006 | TestFileParser_ParseFile_Position | パース結果 | ファイルパス設定 | ファイルパス伝播 |

---

## 1.4 Config テスト

### ファイル: `internal/domain/config/config_test.go`

| ID | テスト名 | 操作 | 期待出力 | 目的 |
|----|---------|------|----------|------|
| DC001 | TestConfig_Default | Default() | デフォルト値 | デフォルト設定 |
| DC002 | TestConfig_Default_SourceRoot | Default() | "./src" | デフォルトSourceRoot |
| DC003 | TestConfig_Default_PactRoot | Default() | "./.pact" | デフォルトPactRoot |
| DC004 | TestConfig_Default_OutputDir | Default() | "./diagrams" | デフォルト出力先 |
| DC005 | TestConfig_Default_Diagrams | Default() | 全4種類 | デフォルト図種類 |
| DC006 | TestConfig_DiagramEnabled_Exists | DiagramEnabled("class") | true | 有効な図 |
| DC007 | TestConfig_DiagramEnabled_NotExists | DiagramEnabled("foo") | false | 無効な図 |
| DC008 | TestConfig_DiagramEnabled_All | diagrams=["all"] | 常にtrue | all指定 |
| DC009 | TestConfig_IsExcluded | パターンマッチ | 正しい結果 | 除外判定 |

### ファイル: `internal/infrastructure/config/loader_test.go`

| ID | テスト名 | 操作 | 期待出力 | 目的 |
|----|---------|------|----------|------|
| CL001 | TestLoader_Load_NotFound | 存在しないパス | Default Config | ファイル不在時デフォルト |
| CL002 | TestLoader_Load_Valid | 有効なYAML | パースされたConfig | 有効ファイル |
| CL003 | TestLoader_Load_Invalid | 不正なYAML | ConfigError | 不正ファイル |
| CL004 | TestLoader_Load_Partial | 一部設定のみ | デフォルトとマージ | 部分設定 |
| CL005 | TestLoader_Save_Success | Config保存 | ファイル作成 | 保存成功 |
| CL006 | TestLoader_Save_ReadBack | 保存→読込 | 同一Config | ラウンドトリップ |
| CL007 | TestLoader_FindProjectRoot_Found | .pactconfigあり | ルートパス | ルート発見 |
| CL008 | TestLoader_FindProjectRoot_NotFound | .pactconfigなし | ConfigError | ルート未発見 |
| CL009 | TestLoader_FindProjectRoot_Nested | サブディレクトリ | 親のルート | 親ディレクトリ検索 |

---

## 1.5 Import Resolver テスト

### ファイル: `internal/infrastructure/resolver/import_test.go`

| ID | テスト名 | 操作 | 期待出力 | 目的 |
|----|---------|------|----------|------|
| IR001 | TestResolver_NoImports | importなし | 同じファイルリスト | インポートなし |
| IR002 | TestResolver_SingleImport | 1件import | 依存順ソート | 単一インポート |
| IR003 | TestResolver_MultipleImports | 複数import | 依存順ソート | 複数インポート |
| IR004 | TestResolver_TransitiveImports | A→B→C | C,B,A順 | 推移的インポート |
| IR005 | TestResolver_CycleDetection | A→B→A | CycleError | サイクル検出 |
| IR006 | TestResolver_SelfImport | A→A | CycleError | 自己参照 |
| IR007 | TestResolver_DiamondImports | A→B,C B,C→D | 重複なし | ダイヤモンド依存 |
| IR008 | TestResolver_RelativePath | "./sub/b.pact" | 正しく解決 | 相対パス |
| IR009 | TestResolver_ImportNotFound | 存在しないパス | ImportError | インポート未発見 |
| IR010 | TestResolver_ImportParseError | 構文エラーファイル | ImportError(内部にParseError) | インポート先エラー |

---

## 1.6 Project Service テスト

### ファイル: `internal/application/project/service_test.go`

| ID | テスト名 | 操作 | 期待出力 | 目的 |
|----|---------|------|----------|------|
| PS001 | TestProjectService_LoadProject_Empty | 空ディレクトリ | 空リスト | 空プロジェクト |
| PS002 | TestProjectService_LoadProject_Single | 1ファイル | 1 SpecFile | 単一ファイル |
| PS003 | TestProjectService_LoadProject_Multiple | 複数ファイル | 複数 SpecFile | 複数ファイル |
| PS004 | TestProjectService_LoadProject_WithImports | import付き | 依存順ソート | インポートあり |
| PS005 | TestProjectService_LoadProject_Exclude | exclude設定 | 除外される | 除外パターン |
| PS006 | TestProjectService_CheckMissing_None | 全対応あり | 空リスト | 欠落なし |
| PS007 | TestProjectService_CheckMissing_Some | 一部欠落 | 欠落リスト | 一部欠落 |
| PS008 | TestProjectService_CheckMissing_All | 全欠落 | 全ソースリスト | 全欠落 |
| PS009 | TestProjectService_SourceForPact_Go | language=go | .go拡張子 | Go対応 |
| PS010 | TestProjectService_SourceForPact_TypeScript | language=ts | .ts拡張子 | TypeScript対応 |
| PS011 | TestProjectService_SourceForPact_Python | language=py | .py拡張子 | Python対応 |

---

## 1.7 Transformer テスト

### 1.7.1 ClassTransformer

#### ファイル: `internal/application/transformer/class_test.go`

| ID | テスト名 | 入力AST | 期待出力 | 目的 |
|----|---------|--------|----------|------|
| TC001 | TestClassTransformer_EmptyComponent | 空Component | 1 Node, 0 Edges | 空コンポーネント |
| TC002 | TestClassTransformer_ComponentWithType | Component + Type | 2 Nodes | 型がノードになる |
| TC003 | TestClassTransformer_ComponentWithMethods | provides付き | methods filled | メソッドがコンパートメントに |
| TC004 | TestClassTransformer_TypeFields | Type + Fields | attributes filled | フィールドが属性に |
| TC005 | TestClassTransformer_FieldVisibility | 各Visibility | 正しいVisibility | 可視性マッピング |
| TC006 | TestClassTransformer_Enum | enum定義 | stereotype: enum | enum |
| TC007 | TestClassTransformer_DependsOn | depends on | Edge(dependency) | 依存エッジ |
| TC008 | TestClassTransformer_Extends | extends | Edge(inheritance) | 継承エッジ |
| TC009 | TestClassTransformer_Implements | implements | Edge(implementation) | 実装エッジ |
| TC010 | TestClassTransformer_Contains | contains | Edge(composition) | コンポジションエッジ |
| TC011 | TestClassTransformer_Aggregates | aggregates | Edge(aggregation) | 集約エッジ |
| TC012 | TestClassTransformer_RequiresInterface | requires | 別ノード(interface) | 要求インターフェース |
| TC013 | TestClassTransformer_MultipleFiles | 複数SpecFile | 統合された図 | 複数ファイル |
| TC014 | TestClassTransformer_FilterComponents | opts指定 | フィルタ適用 | コンポーネントフィルタ |
| TC015 | TestClassTransformer_Annotations | アノテーション付き | Annotations保持 | アノテーション伝播 |
| TC016 | TestClassTransformer_EdgeDecorations | 各関係種類 | 正しいDecoration | エッジ装飾 |
| TC017 | TestClassTransformer_EdgeLineStyle | 各関係種類 | 正しいLineStyle | 線スタイル |

### 1.7.2 SequenceTransformer

#### ファイル: `internal/application/transformer/sequence_test.go`

| ID | テスト名 | 入力AST | 期待出力 | 目的 |
|----|---------|--------|----------|------|
| TS001 | TestSeqTransformer_EmptyFlow | 空Flow | 1 Participant, 0 Events | 空フロー |
| TS002 | TestSeqTransformer_SingleCall | 1 CallStep | 2 Participants, 1 Message | 単一呼び出し |
| TS003 | TestSeqTransformer_CallWithReturn | AssignStep(call) | Message + Return | 戻りメッセージ |
| TS004 | TestSeqTransformer_AsyncCall | await付き | MessageType:async | 非同期 |
| TS005 | TestSeqTransformer_SyncCall | awaitなし | MessageType:sync | 同期 |
| TS006 | TestSeqTransformer_MultipleParticipants | 複数depends | 複数Participant | 複数参加者 |
| TS007 | TestSeqTransformer_ParticipantType_Database | depends: database | Type:database | 参加者タイプ |
| TS008 | TestSeqTransformer_ParticipantType_External | depends: external | Type:external | 外部システム |
| TS009 | TestSeqTransformer_ParticipantType_Queue | depends: queue | Type:queue | キュー |
| TS010 | TestSeqTransformer_ParticipantType_Actor | depends: actor | Type:actor | アクター |
| TS011 | TestSeqTransformer_IfStep | IfStep | FragmentEvent(alt) | 条件分岐 |
| TS012 | TestSeqTransformer_IfElse | If-Else | alt + else | if-else |
| TS013 | TestSeqTransformer_ForStep | ForStep | FragmentEvent(loop) | forループ |
| TS014 | TestSeqTransformer_WhileStep | WhileStep | FragmentEvent(loop) | whileループ |
| TS015 | TestSeqTransformer_NestedIf | ネストIf | ネストFragment | ネストif |
| TS016 | TestSeqTransformer_OptionReturnMessages | opts.IncludeReturn=false | returnなし | オプション |
| TS017 | TestSeqTransformer_FlowNotFound | 存在しないflowName | TransformError | フロー未発見 |
| TS018 | TestSeqTransformer_MessageOrder | 複数ステップ | 正しい順序 | メッセージ順序 |
| TS019 | TestSeqTransformer_Activation | 呼び出し | Activation events | アクティベーション |

### 1.7.3 StateTransformer

#### ファイル: `internal/application/transformer/state_test.go`

| ID | テスト名 | 入力AST | 期待出力 | 目的 |
|----|---------|--------|----------|------|
| TST001 | TestStateTransformer_MinimalStates | initial のみ | Initial + 1 State | 最小 |
| TST002 | TestStateTransformer_InitialState | initial指定 | Type:initial のState | 初期状態 |
| TST003 | TestStateTransformer_FinalState | final指定 | Type:final のState | 最終状態 |
| TST004 | TestStateTransformer_AtomicStates | 遷移から収集 | Type:atomic | 原子状態 |
| TST005 | TestStateTransformer_Transition_Event | on E | Trigger:EventTrigger | イベントトリガー |
| TST006 | TestStateTransformer_Transition_After | after 3s | Trigger:AfterTrigger | 時間トリガー |
| TST007 | TestStateTransformer_Transition_When | when cond | Trigger:WhenTrigger | 条件トリガー |
| TST008 | TestStateTransformer_Transition_Guard | when付き | Guard設定 | ガード |
| TST009 | TestStateTransformer_Transition_Actions | do [a,b] | Actions設定 | アクション |
| TST010 | TestStateTransformer_State_Entry | entry [a] | Entry設定 | entry |
| TST011 | TestStateTransformer_State_Exit | exit [a] | Exit設定 | exit |
| TST012 | TestStateTransformer_CompoundState | ネストstate | Type:compound, Children | 階層状態 |
| TST013 | TestStateTransformer_CompoundInitial | 内部initial | 内部Initial設定 | 階層initial |
| TST014 | TestStateTransformer_ParallelState | parallel | Type:parallel, Regions | 並行状態 |
| TST015 | TestStateTransformer_Region | region定義 | Region構造 | リージョン |
| TST016 | TestStateTransformer_StatesNotFound | 存在しないstatesName | TransformError | 未発見 |
| TST017 | TestStateTransformer_Annotations | アノテーション付き | Annotations保持 | アノテーション |
| TST018 | TestStateTransformer_DurationUnits | 各単位(ms,s,m,h,d) | 正しい変換 | 時間単位 |

### 1.7.4 FlowTransformer

#### ファイル: `internal/application/transformer/flow_test.go`

| ID | テスト名 | 入力AST | 期待出力 | 目的 |
|----|---------|--------|----------|------|
| TF001 | TestFlowTransformer_EmptyFlow | 空Flow | Start + End | 空フロー |
| TF002 | TestFlowTransformer_StartEndNodes | 任意Flow | terminal nodes | 端子ノード |
| TF003 | TestFlowTransformer_AssignStep | AssignStep | process node | 代入→処理 |
| TF004 | TestFlowTransformer_CallStep | CallStep | process node | 呼び出し→処理 |
| TF005 | TestFlowTransformer_ReturnStep | ReturnStep | terminal + no outgoing | return |
| TF006 | TestFlowTransformer_ThrowStep | ThrowStep | terminal + no outgoing | throw |
| TF007 | TestFlowTransformer_IfStep | IfStep | decision node + edges | 判断 |
| TF008 | TestFlowTransformer_IfStep_Labels | If | edges with Yes/No | 分岐ラベル |
| TF009 | TestFlowTransformer_IfElse | If-Else | 両分岐エッジ | if-else |
| TF010 | TestFlowTransformer_ForStep | ForStep | decision + loop edge | forループ |
| TF011 | TestFlowTransformer_WhileStep | WhileStep | decision + loop edge | whileループ |
| TF012 | TestFlowTransformer_NestedIf | ネストIf | ネスト構造 | ネストif |
| TF013 | TestFlowTransformer_Sequential | 複数ステップ | 連続エッジ | 順次 |
| TF014 | TestFlowTransformer_Swimlane_Infer | 複数target | 推論されたSwimlanes | スイムレーン推論 |
| TF015 | TestFlowTransformer_OptionSwimlanes | opts.IncludeSwimlanes | スイムレーン有無 | オプション |
| TF016 | TestFlowTransformer_FlowNotFound | 存在しないflowName | TransformError | 未発見 |
| TF017 | TestFlowTransformer_NodeShapes | 各ステップ種類 | 正しいShape | 形状マッピング |
| TF018 | TestFlowTransformer_EdgeConnectivity | 複合フロー | 全ノード接続 | 接続性 |
| TF019 | TestFlowTransformer_MergeNode | if-else後 | connector node | マージノード |

---

## 1.8 Renderer テスト

### 1.8.1 Canvas テスト

#### ファイル: `internal/infrastructure/renderer/canvas/canvas_test.go`

| ID | テスト名 | 操作 | 期待出力 | 目的 |
|----|---------|------|----------|------|
| RC001 | TestCanvas_Empty | New() | 空SVG | 空キャンバス |
| RC002 | TestCanvas_SetSize | SetSize(800,600) | viewBox設定 | サイズ設定 |
| RC003 | TestCanvas_Rect | Rect(10,20,100,50) | rect要素 | 矩形 |
| RC004 | TestCanvas_RoundRect | RoundRect(...,rx,ry) | rx,ry付きrect | 角丸矩形 |
| RC005 | TestCanvas_Circle | Circle(50,50,25) | circle要素 | 円 |
| RC006 | TestCanvas_Ellipse | Ellipse(50,50,30,20) | ellipse要素 | 楕円 |
| RC007 | TestCanvas_Line | Line(0,0,100,100) | line要素 | 線 |
| RC008 | TestCanvas_Path | Path("M0,0 L100,100") | path要素 | パス |
| RC009 | TestCanvas_Polygon | Polygon(points) | polygon要素 | 多角形 |
| RC010 | TestCanvas_Text | Text(10,20,"hello") | text要素 | テキスト |
| RC011 | TestCanvas_Option_Fill | Fill("#ff0000") | fill属性 | 塗り |
| RC012 | TestCanvas_Option_Stroke | Stroke("#000000") | stroke属性 | 線色 |
| RC013 | TestCanvas_Option_StrokeWidth | StrokeWidth(2) | stroke-width属性 | 線幅 |
| RC014 | TestCanvas_Option_Class | Class("node") | class属性 | クラス |
| RC015 | TestCanvas_Option_Multiple | 複数Option | 複数属性 | 複数オプション |
| RC016 | TestCanvas_Defs | AddDef(marker) | defs内に追加 | 定義 |
| RC017 | TestCanvas_WriteTo | WriteTo(w) | 有効なSVG | 出力 |

#### ファイル: `internal/infrastructure/renderer/canvas/shapes_test.go`

| ID | テスト名 | 操作 | 期待出力 | 目的 |
|----|---------|------|----------|------|
| RS001 | TestShapes_Diamond | Diamond(...) | ひし形polygon | ひし形 |
| RS002 | TestShapes_Arrow | Arrow(...) | 線+矢印 | 矢印 |
| RS003 | TestShapes_Stadium | Stadium(...) | 角丸長方形 | 端子形状 |
| RS004 | TestShapes_Cylinder | Cylinder(...) | 円柱 | DB形状 |
| RS005 | TestShapes_Parallelogram | Parallelogram(...) | 平行四辺形 | IO形状 |

#### ファイル: `internal/infrastructure/renderer/canvas/text_test.go`

| ID | テスト名 | 操作 | 期待出力 | 目的 |
|----|---------|------|----------|------|
| RT001 | TestText_MeasureText | MeasureText("hello", 12) | 概算サイズ | テキスト測定 |
| RT002 | TestText_WrapText | WrapText(長文, 100, 12) | 折り返し行 | テキスト折り返し |

### 1.8.2 ClassRenderer テスト

#### ファイル: `internal/infrastructure/renderer/svg/class_test.go`

| ID | テスト名 | 入力 | 検証内容 | 目的 |
|----|---------|------|----------|------|
| RCL001 | TestClassRenderer_EmptyDiagram | 空Diagram | 有効なSVG | 空図 |
| RCL002 | TestClassRenderer_SingleNode | 1 Node | rect + text | 単一ノード |
| RCL003 | TestClassRenderer_NodeWithAttributes | attributes付き | コンパートメント | 属性表示 |
| RCL004 | TestClassRenderer_NodeWithMethods | methods付き | コンパートメント | メソッド表示 |
| RCL005 | TestClassRenderer_Stereotype | stereotype付き | ≪stereotype≫ | ステレオタイプ |
| RCL006 | TestClassRenderer_VisibilitySymbols | 各Visibility | +/-/#/~ | 可視性記号 |
| RCL007 | TestClassRenderer_Edge_Dependency | dependency | 点線矢印 | 依存 |
| RCL008 | TestClassRenderer_Edge_Inheritance | inheritance | 実線白三角 | 継承 |
| RCL009 | TestClassRenderer_Edge_Implementation | implementation | 点線白三角 | 実装 |
| RCL010 | TestClassRenderer_Edge_Composition | composition | 黒ひし形 | コンポジション |
| RCL011 | TestClassRenderer_Edge_Aggregation | aggregation | 白ひし形 | 集約 |
| RCL012 | TestClassRenderer_Layout | 複数ノード | 重ならない | レイアウト |
| RCL013 | TestClassRenderer_EdgeRouting | 複数エッジ | 交差最小化 | エッジルーティング |

### 1.8.3 SequenceRenderer テスト

#### ファイル: `internal/infrastructure/renderer/svg/sequence_test.go`

| ID | テスト名 | 入力 | 検証内容 | 目的 |
|----|---------|------|----------|------|
| RSQ001 | TestSeqRenderer_EmptyDiagram | 空Diagram | 有効なSVG | 空図 |
| RSQ002 | TestSeqRenderer_SingleParticipant | 1 Participant | rect + lifeline | 単一参加者 |
| RSQ003 | TestSeqRenderer_MultipleParticipants | 複数Participant | 水平配置 | 複数参加者 |
| RSQ004 | TestSeqRenderer_ParticipantType_Actor | Type:actor | 人型 | アクター |
| RSQ005 | TestSeqRenderer_ParticipantType_Database | Type:database | 円柱 | DB |
| RSQ006 | TestSeqRenderer_ParticipantType_Queue | Type:queue | 特殊形状 | キュー |
| RSQ007 | TestSeqRenderer_Message_Sync | sync | 実線閉じ矢印 | 同期 |
| RSQ008 | TestSeqRenderer_Message_Async | async | 実線開き矢印 | 非同期 |
| RSQ009 | TestSeqRenderer_Message_Return | return | 点線矢印 | 戻り |
| RSQ010 | TestSeqRenderer_Message_Label | label付き | テキスト表示 | ラベル |
| RSQ011 | TestSeqRenderer_Fragment_Alt | alt | 枠 + "alt" | 条件分岐 |
| RSQ012 | TestSeqRenderer_Fragment_Loop | loop | 枠 + "loop" | ループ |
| RSQ013 | TestSeqRenderer_Fragment_Nested | ネストFragment | ネスト枠 | ネスト |
| RSQ014 | TestSeqRenderer_Activation | activation | 細長い矩形 | アクティベーション |
| RSQ015 | TestSeqRenderer_MessageOrder | 複数Message | 上から順 | 順序 |

### 1.8.4 StateRenderer テスト

#### ファイル: `internal/infrastructure/renderer/svg/state_test.go`

| ID | テスト名 | 入力 | 検証内容 | 目的 |
|----|---------|------|----------|------|
| RST001 | TestStateRenderer_EmptyDiagram | 空Diagram | 有効なSVG | 空図 |
| RST002 | TestStateRenderer_InitialState | initial | 黒丸 | 初期状態 |
| RST003 | TestStateRenderer_FinalState | final | 二重丸 | 最終状態 |
| RST004 | TestStateRenderer_AtomicState | atomic | 角丸矩形 | 原子状態 |
| RST005 | TestStateRenderer_StateWithName | name付き | テキスト表示 | 状態名 |
| RST006 | TestStateRenderer_StateWithEntry | entry付き | entry/ 表示 | entry |
| RST007 | TestStateRenderer_StateWithExit | exit付き | exit/ 表示 | exit |
| RST008 | TestStateRenderer_Transition | 遷移 | 矢印 | 遷移 |
| RST009 | TestStateRenderer_Transition_Label | label付き | event[guard]/action | ラベル |
| RST010 | TestStateRenderer_CompoundState | compound | 外枠 + 内部 | 階層状態 |
| RST011 | TestStateRenderer_ParallelState | parallel | 点線区切り | 並行状態 |
| RST012 | TestStateRenderer_Region | region | 区切り線 | リージョン |
| RST013 | TestStateRenderer_Layout | 複数状態 | 重ならない | レイアウト |
| RST014 | TestStateRenderer_SelfTransition | 自己遷移 | ループ矢印 | 自己遷移 |

### 1.8.5 FlowRenderer テスト

#### ファイル: `internal/infrastructure/renderer/svg/flow_test.go`

| ID | テスト名 | 入力 | 検証内容 | 目的 |
|----|---------|------|----------|------|
| RFL001 | TestFlowRenderer_EmptyDiagram | 空Diagram | 有効なSVG | 空図 |
| RFL002 | TestFlowRenderer_TerminalNode | terminal | 角丸長方形 | 端子 |
| RFL003 | TestFlowRenderer_ProcessNode | process | 長方形 | 処理 |
| RFL004 | TestFlowRenderer_DecisionNode | decision | ひし形 | 判断 |
| RFL005 | TestFlowRenderer_IONode | io | 平行四辺形 | 入出力 |
| RFL006 | TestFlowRenderer_DatabaseNode | database | 円柱 | DB |
| RFL007 | TestFlowRenderer_NodeLabel | label付き | テキスト表示 | ラベル |
| RFL008 | TestFlowRenderer_Edge | エッジ | 矢印線 | エッジ |
| RFL009 | TestFlowRenderer_EdgeLabel | label付き | テキスト表示 | エッジラベル |
| RFL010 | TestFlowRenderer_Swimlane | swimlanes | 縦区切り | スイムレーン |
| RFL011 | TestFlowRenderer_SwimlaneLabel | name付き | ヘッダ表示 | スイムレーン名 |
| RFL012 | TestFlowRenderer_Layout | 複数ノード | 上→下配置 | レイアウト |
| RFL013 | TestFlowRenderer_BranchMerge | 分岐・合流 | 正しい接続 | 分岐合流 |

---

## 1.9 Filesystem テスト

### ファイル: `internal/infrastructure/filesystem/reader_test.go`

| ID | テスト名 | 操作 | 期待出力 | 目的 |
|----|---------|------|----------|------|
| FR001 | TestReader_ReadFile_Success | 存在するファイル | 内容 | 読み込み成功 |
| FR002 | TestReader_ReadFile_NotFound | 存在しないファイル | エラー | ファイル不在 |
| FR003 | TestReader_FindPactFiles_Empty | 空ディレクトリ | 空リスト | 空ディレクトリ |
| FR004 | TestReader_FindPactFiles_Found | .pactあり | パスリスト | ファイル発見 |
| FR005 | TestReader_FindPactFiles_Nested | サブディレクトリ | 再帰的に発見 | ネスト検索 |
| FR006 | TestReader_FindPactFiles_Exclude | excludeパターン | 除外される | 除外 |
| FR007 | TestReader_FindSourceFiles | ソースあり | パスリスト | ソース検索 |

### ファイル: `internal/infrastructure/filesystem/writer_test.go`

| ID | テスト名 | 操作 | 期待出力 | 目的 |
|----|---------|------|----------|------|
| FW001 | TestWriter_WriteFile_Success | ファイル書き込み | ファイル作成 | 書き込み成功 |
| FW002 | TestWriter_WriteFile_CreateDir | 存在しないディレクトリ | ディレクトリ作成 | ディレクトリ自動作成 |
| FW003 | TestWriter_WriteFile_Overwrite | 既存ファイル | 上書き | 上書き |
| FW004 | TestWriter_EnsureDir | ディレクトリ作成 | ディレクトリ存在 | ディレクトリ作成 |

---

## 1.10 Domain Model テスト

### ファイル: `internal/domain/ast/visitor_test.go`

| ID | テスト名 | 操作 | 期待出力 | 目的 |
|----|---------|------|----------|------|
| AV001 | TestVisitor_Walk_SpecFile | Walk(SpecFile) | VisitSpecFile呼出 | SpecFile訪問 |
| AV002 | TestVisitor_Walk_AllNodes | 複合AST | 全ノード訪問 | 全ノード走査 |
| AV003 | TestVisitor_Walk_Order | AST | 深さ優先順 | 走査順序 |
| AV004 | TestVisitor_BaseVisitor | BaseVisitor | エラーなし | デフォルト実装 |

### ファイル: `internal/domain/errors/errors_test.go`

| ID | テスト名 | 操作 | 期待出力 | 目的 |
|----|---------|------|----------|------|
| DE001 | TestParseError_Error | Error() | 位置+メッセージ | エラー文字列 |
| DE002 | TestSemanticError_Error | Error() | 位置+メッセージ | エラー文字列 |
| DE003 | TestImportError_Error | Error() | 位置+パス+メッセージ | エラー文字列 |
| DE004 | TestCycleError_Error | Error() | サイクル情報 | エラー文字列 |
| DE005 | TestTransformError_Error | Error() | ソース+ターゲット+メッセージ | エラー文字列 |
| DE006 | TestConfigError_Error | Error() | パス+メッセージ | エラー文字列 |

---

# 2. 結合テスト（Integration Tests）

### ファイル: `test/integration/parse_transform_test.go`

### 2.1 パース→変換

| ID | テスト名 | 入力 | 検証内容 | 目的 |
|----|---------|------|----------|------|
| I001 | TestIntegration_ParseToClassDiagram | .pact文字列 | パース→Class変換成功 | パース→クラス図 |
| I002 | TestIntegration_ParseToSequenceDiagram | .pact文字列 | パース→Seq変換成功 | パース→シーケンス図 |
| I003 | TestIntegration_ParseToStateDiagram | .pact文字列 | パース→State変換成功 | パース→状態図 |
| I004 | TestIntegration_ParseToFlowchart | .pact文字列 | パース→Flow変換成功 | パース→フローチャート |
| I005 | TestIntegration_ParseMultipleFiles | 複数.pact | 複数パース成功 | 複数ファイル |

### ファイル: `test/integration/import_resolution_test.go`

### 2.2 インポート解決

| ID | テスト名 | 入力 | 検証内容 | 目的 |
|----|---------|------|----------|------|
| I006 | TestIntegration_ImportResolution_Single | A imports B | B→A順で解決 | 単一インポート |
| I007 | TestIntegration_ImportResolution_Chain | A→B→C | C→B→A順 | チェーン |
| I008 | TestIntegration_ImportResolution_Diamond | A→B,C→D | 正しい順序 | ダイヤモンド |
| I009 | TestIntegration_ImportResolution_Cycle | A↔B | CycleError | サイクル |
| I010 | TestIntegration_ImportResolution_NotFound | A→存在しない | ImportError | 未発見 |

### ファイル: `test/integration/full_pipeline_test.go`

### 2.3 フルパイプライン

| ID | テスト名 | 入力 | 検証内容 | 目的 |
|----|---------|------|----------|------|
| I020 | TestPipeline_ClassDiagram_Full | .pact → SVG | 有効なSVG出力 | クラス図パイプライン |
| I021 | TestPipeline_SequenceDiagram_Full | .pact → SVG | 有効なSVG出力 | シーケンス図パイプライン |
| I022 | TestPipeline_StateDiagram_Full | .pact → SVG | 有効なSVG出力 | 状態図パイプライン |
| I023 | TestPipeline_Flowchart_Full | .pact → SVG | 有効なSVG出力 | フローチャートパイプライン |
| I024 | TestPipeline_AllDiagrams | .pact → 4種SVG | 全種類生成成功 | 全図生成 |
| I025 | TestPipeline_ComplexSpec | 複雑な.pact | 全図正常生成 | 複雑ケース |
| I026 | TestPipeline_RealWorldExample | testdata/*.pact | 実例で動作 | 実例 |

### 2.4 エラーハンドリング

| ID | テスト名 | 入力 | 検証内容 | 目的 |
|----|---------|------|----------|------|
| I030 | TestIntegration_ParseError_Propagation | 不正.pact | エラー伝播 | パースエラー伝播 |
| I031 | TestIntegration_TransformError_Propagation | 不整合AST | エラー伝播 | 変換エラー伝播 |
| I032 | TestIntegration_ErrorPosition | エラー | 正確な位置 | エラー位置 |
| I033 | TestIntegration_ErrorType | 各エラー | 正しいエラー型 | エラー型 |

---

# 3. E2Eテスト（End-to-End Tests）

### ファイル: `test/e2e/cli_test.go`

### 3.1 CLI init コマンド

| ID | テスト名 | コマンド | 検証内容 | 目的 |
|----|---------|---------|----------|------|
| E001 | TestCLI_Init_Success | `pact init` | .pactconfig作成 | 初期化成功 |
| E002 | TestCLI_Init_Exists | `pact init`（既存あり） | 上書きまたはエラー | 既存ファイル |
| E003 | TestCLI_Init_Content | `pact init` | 正しい内容 | 初期内容 |

### 3.2 CLI generate コマンド

| ID | テスト名 | コマンド | 検証内容 | 目的 |
|----|---------|---------|----------|------|
| E010 | TestCLI_Generate_SingleFile | `pact generate` | SVG生成 | 単一ファイル |
| E011 | TestCLI_Generate_Directory | `pact generate` | 複数SVG生成 | ディレクトリ |
| E012 | TestCLI_Generate_Output | `pact generate -o ./out` | 出力先指定 | 出力先 |
| E013 | TestCLI_Generate_TypeClass | `pact generate -t class` | クラス図のみ | 種類指定 |
| E014 | TestCLI_Generate_TypeSequence | `pact generate -t sequence` | シーケンス図のみ | 種類指定 |
| E015 | TestCLI_Generate_TypeState | `pact generate -t state` | 状態図のみ | 種類指定 |
| E016 | TestCLI_Generate_TypeFlow | `pact generate -t flow` | フローチャートのみ | 種類指定 |
| E017 | TestCLI_Generate_TypeMultiple | `pact generate -t class,sequence` | 複数種類 | 複数種類 |
| E018 | TestCLI_Generate_TypeAll | `pact generate -t all` | 全種類 | 全種類 |
| E019 | TestCLI_Generate_NotFound | `pact generate`（ファイルなし） | エラー | ファイル不在 |
| E01A | TestCLI_Generate_InvalidSpec | `pact generate`（不正ファイル） | エラー | 不正ファイル |

### 3.3 CLI validate コマンド

| ID | テスト名 | コマンド | 検証内容 | 目的 |
|----|---------|---------|----------|------|
| E020 | TestCLI_Validate_Valid | `pact validate` | 成功(exit 0) | 有効ファイル |
| E021 | TestCLI_Validate_Invalid | `pact validate`（不正） | 失敗(exit 1) + エラー出力 | 無効ファイル |
| E022 | TestCLI_Validate_Directory | `pact validate` | 全ファイル検証 | ディレクトリ |
| E023 | TestCLI_Validate_ErrorOutput | `pact validate`（不正） | エラー詳細出力 | エラー出力 |

### 3.4 CLI check コマンド

| ID | テスト名 | コマンド | 検証内容 | 目的 |
|----|---------|---------|----------|------|
| E030 | TestCLI_Check_Success | `pact check` | 成功(exit 0) | チェック成功 |
| E031 | TestCLI_Check_Missing_None | `pact check --missing`（全対応） | "All source files have .pact" | 欠落なし |
| E032 | TestCLI_Check_Missing_Some | `pact check --missing`（一部欠落） | 欠落リスト出力 | 一部欠落 |
| E033 | TestCLI_Check_Missing_Exit | `pact check --missing`（欠落あり） | 失敗(exit 1) | 終了コード |
| E034 | TestCLI_Check_ParseError | `pact check`（構文エラー） | エラー出力 | 構文エラー検出 |

### 3.5 CLI watch コマンド

| ID | テスト名 | 操作 | 検証内容 | 目的 |
|----|---------|------|----------|------|
| E040 | TestCLI_Watch_FileChange | ファイル変更 | 自動再生成 | ファイル監視 |
| E041 | TestCLI_Watch_NewFile | 新規ファイル | 自動生成 | 新規ファイル |
| E042 | TestCLI_Watch_DeleteFile | ファイル削除 | 適切な処理 | ファイル削除 |

---

# 4. 公開API テスト

### ファイル: `pkg/pact/api_test.go`

| ID | テスト名 | 操作 | 検証内容 | 目的 |
|----|---------|------|----------|------|
| A001 | TestAPI_New | New() | Clientインスタンス | 初期化 |
| A002 | TestAPI_ParseFile | ParseFile(path) | SpecFile取得 | ファイルパース |
| A003 | TestAPI_ParseFile_NotFound | ParseFile(存在しない) | エラー | ファイル不在 |
| A004 | TestAPI_ParseString | ParseString(content) | SpecFile取得 | 文字列パース |
| A005 | TestAPI_ParseString_Invalid | ParseString(不正) | ParseError | 構文エラー |
| A006 | TestAPI_ToClassDiagram | ToClassDiagram(...) | Diagram取得 | クラス図変換 |
| A007 | TestAPI_ToSequenceDiagram | ToSequenceDiagram(...) | Diagram取得 | シーケンス図変換 |
| A008 | TestAPI_ToSequenceDiagram_NotFound | ToSequenceDiagram(存在しないflow) | TransformError | フロー未発見 |
| A009 | TestAPI_ToStateDiagram | ToStateDiagram(...) | Diagram取得 | 状態図変換 |
| A010 | TestAPI_ToStateDiagram_NotFound | ToStateDiagram(存在しないstates) | TransformError | ステート未発見 |
| A011 | TestAPI_ToFlowchart | ToFlowchart(...) | Diagram取得 | フローチャート変換 |
| A012 | TestAPI_RenderClassDiagram | RenderClassDiagram(...) | SVG出力 | クラス図レンダー |
| A013 | TestAPI_RenderSequenceDiagram | RenderSequenceDiagram(...) | SVG出力 | シーケンス図レンダー |
| A014 | TestAPI_RenderStateDiagram | RenderStateDiagram(...) | SVG出力 | 状態図レンダー |
| A015 | TestAPI_RenderFlowchart | RenderFlowchart(...) | SVG出力 | フローチャートレンダー |

---

# 5. テストデータ

### ファイル構成

```
testdata/
├── valid/
│   ├── minimal.pact           # 最小限のコンポーネント
│   ├── with_types.pact        # 型定義あり
│   ├── with_relations.pact    # 関係定義あり
│   ├── with_interfaces.pact   # インターフェースあり
│   ├── with_flow.pact         # フローあり
│   ├── with_states.pact       # ステートマシンあり
│   ├── with_annotations.pact  # アノテーションあり
│   ├── complex.pact           # 複合例
│   ├── auth_service.pact      # 実例：認証サービス
│   ├── order_service.pact     # 実例：注文サービス
│   └── multi/
│       ├── service_a.pact     # import "./shared.pact"
│       ├── service_b.pact     # import "./shared.pact"
│       └── shared.pact        # 共通定義
│
├── invalid/
│   ├── syntax/
│   │   ├── missing_brace.pact
│   │   ├── missing_name.pact
│   │   ├── invalid_token.pact
│   │   └── unterminated_string.pact
│   │
│   ├── semantic/
│   │   ├── undefined_type.pact
│   │   ├── duplicate_name.pact
│   │   └── invalid_reference.pact
│   │
│   └── import/
│       ├── cycle_a.pact       # import "./cycle_b.pact"
│       ├── cycle_b.pact       # import "./cycle_a.pact"
│       └── not_found.pact     # import "./missing.pact"
│
├── config/
│   ├── default.pactconfig     # デフォルト設定
│   ├── custom.pactconfig      # カスタム設定
│   └── invalid.pactconfig     # 不正な設定
│
└── expected/
    ├── minimal_class.svg
    ├── minimal_sequence.svg
    ├── minimal_state.svg
    └── minimal_flow.svg
```

---

# 6. テスト実行

### Makefile

```makefile
.PHONY: test test-unit test-integration test-e2e test-coverage test-watch lint

# 全テスト実行
test: lint test-unit test-integration test-e2e

# 単体テスト
test-unit:
	go test -v -race ./internal/...

# 結合テスト
test-integration:
	go test -v -race ./test/integration/...

# E2Eテスト
test-e2e:
	go test -v -race ./test/e2e/...

# 公開APIテスト
test-api:
	go test -v -race ./pkg/...

# カバレッジ
test-coverage:
	go test -coverprofile=coverage.out -covermode=atomic ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# 特定パッケージのテスト
test-pkg:
	go test -v -race ./$(PKG)/...

# 特定テストの実行
test-run:
	go test -v -race -run $(RUN) ./...

# 監視モード
test-watch:
	watchexec -e go "make test-unit"

# リント
lint:
	golangci-lint run ./...

# ベンチマーク
bench:
	go test -bench=. -benchmem ./internal/infrastructure/parser/...
	go test -bench=. -benchmem ./internal/application/transformer/...
```

### テスト実行順序（TDD）

1. **Domain テスト**
   - `internal/domain/config/config_test.go`
   - `internal/domain/errors/errors_test.go`
   - `internal/domain/ast/visitor_test.go`

2. **Infrastructure テスト**
   - `internal/infrastructure/filesystem/*_test.go`
   - `internal/infrastructure/parser/lexer_test.go` (L001-L153)
   - `internal/infrastructure/parser/parser_test.go` (P001-P204)
   - `internal/infrastructure/parser/adapter_test.go` (PA001-PA006)
   - `internal/infrastructure/config/loader_test.go` (CL001-CL009)
   - `internal/infrastructure/resolver/import_test.go` (IR001-IR010)

3. **Application テスト**
   - `internal/application/parser/service_test.go`
   - `internal/application/project/service_test.go` (PS001-PS011)
   - `internal/application/transformer/*_test.go` (TC001-TF019)
   - `internal/application/renderer/service_test.go`

4. **Renderer テスト**
   - `internal/infrastructure/renderer/canvas/*_test.go` (RC001-RT002)
   - `internal/infrastructure/renderer/svg/*_test.go` (RCL001-RFL013)

5. **結合テスト** (I001-I033)

6. **E2Eテスト** (E001-E042)

7. **公開APIテスト** (A001-A015)

---

# 7. テストヘルパー

### ファイル: `internal/testutil/helper.go`

```go
package testutil

import (
    "os"
    "path/filepath"
    "testing"

    "pact/internal/domain/ast"
)

// TempDir はテスト用の一時ディレクトリを作成する
func TempDir(t *testing.T) string {
    t.Helper()
    dir, err := os.MkdirTemp("", "pact-test-*")
    if err != nil {
        t.Fatalf("failed to create temp dir: %v", err)
    }
    t.Cleanup(func() { os.RemoveAll(dir) })
    return dir
}

// WriteFile はテスト用ファイルを作成する
func WriteFile(t *testing.T, dir, name, content string) string {
    t.Helper()
    path := filepath.Join(dir, name)
    if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
        t.Fatalf("failed to create dir: %v", err)
    }
    if err := os.WriteFile(path, []byte(content), 0644); err != nil {
        t.Fatalf("failed to write file: %v", err)
    }
    return path
}

// MinimalPact は最小限の.pact内容を返す
func MinimalPact(name string) string {
    return "component " + name + " { }"
}

// AssertNoError はエラーがないことを確認する
func AssertNoError(t *testing.T, err error) {
    t.Helper()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}

// AssertError はエラーがあることを確認する
func AssertError(t *testing.T, err error) {
    t.Helper()
    if err == nil {
        t.Fatal("expected error, got nil")
    }
}

// AssertEqual は値が等しいことを確認する
func AssertEqual[T comparable](t *testing.T, expected, actual T) {
    t.Helper()
    if expected != actual {
        t.Fatalf("expected %v, got %v", expected, actual)
    }
}

// AssertPosition は位置情報を確認する
func AssertPosition(t *testing.T, pos ast.Position, line, column int) {
    t.Helper()
    if pos.Line != line || pos.Column != column {
        t.Fatalf("expected position %d:%d, got %d:%d", line, column, pos.Line, pos.Column)
    }
}
```

---

# 8. カバレッジ目標

| パッケージ | 目標カバレッジ |
|-----------|---------------|
| internal/domain/* | 90%+ |
| internal/application/* | 85%+ |
| internal/infrastructure/parser | 90%+ |
| internal/infrastructure/resolver | 85%+ |
| internal/infrastructure/config | 85%+ |
| internal/infrastructure/renderer | 80%+ |
| internal/infrastructure/filesystem | 80%+ |
| internal/interfaces/cli | 75%+ |
| pkg/pact | 85%+ |
| **全体** | **80%+** |