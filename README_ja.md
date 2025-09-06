# Terraform Provider for Multipass

[Canonical Multipass](https://multipass.run/)を使用してUbuntu仮想マシンを管理するためのTerraformプロバイダーです。

Multipassは、Linux、Windows、macOS向けの軽量VMマネージャーで、Ubuntuインスタンスを迅速に作成・管理できます。

## ⚠️ 重要な注意事項

**このプロジェクトは主にClaude Code AIの支援により生成されました。メンテナー（sh05）は生成されたすべてのコードについて十分なレビューを実施できていません。このプロバイダーをご利用の際は自己責任でお願いします。本番環境で使用する前に、ご自身の環境で十分にテストを実施してください。メンテナーは、このプロバイダーの使用により発生する可能性のある問題や課題について、責任を負うことはできません。**

## 機能

✅ **実装済み**
- 🚀 Multipassインスタンスの作成、読み込み、更新、削除
- 📊 既存インスタンスの照会とすべてのインスタンスの一覧表示
- ⚙️ CPU、メモリ、ディスク、Ubuntuイメージバージョンの設定
- 🔧 インスタンスのカスタマイズ用のCloud-initサポート
- 📋 既存インスタンスのTerraform状態へのインポート
- 🧪 包括的テストスイート

🔄 **今後の拡張機能** ([ロードマップ](#ロードマップ)を参照)
- ネットワーク設定とポートフォワーディング
- ボリュームマウントサポート
- スナップショット管理
- マルチプラットフォームバイナリリリース

## インストール

### Terraform Registryから（近日公開予定）

Terraform Registryに公開されたら、Terraform設定で直接プロバイダーを使用できます：

```hcl
terraform {
  required_providers {
    multipass = {
      source  = "registry.terraform.io/sh05/multipass"
      version = "~> 0.1.0"
    }
  }
}
```

### ローカル開発

1. リポジトリをクローン：
```bash
git clone https://github.com/sh05/terraform-provider-multipass
cd terraform-provider-multipass
```

2. プロバイダーをローカルでビルド・インストール：
```bash
make install-local
```

これにより、ローカル開発用にプロバイダーが `~/.terraform.d/plugins/` にインストールされます。

## 使用方法

### 基本的な例

```hcl
terraform {
  required_providers {
    multipass = {
      source  = "registry.terraform.io/sh05/multipass"
      version = "~> 0.1.0"
    }
  }
}

provider "multipass" {
  # オプション: multipassバイナリがPATHにない場合はパスを指定
  # binary_path = "/usr/local/bin/multipass"
}

# 基本的なUbuntuインスタンスの作成
resource "multipass_instance" "example" {
  name   = "my-ubuntu-vm"
  image  = "22.04"
  cpu    = "2"
  memory = "2G"
  disk   = "10G"
  
  # オプション: タイムアウトの設定
  timeouts {
    create = "20m"  # インスタンス作成に最大20分を許可
    delete = "5m"   # 削除に最大5分を許可
  }
}

# インスタンスの情報を取得
data "multipass_instance" "example" {
  name = multipass_instance.example.name
}

# すべてのインスタンスを一覧表示
data "multipass_instance" "all" {}

# インスタンス情報の出力
output "instance_ip" {
  value = data.multipass_instance.example.instance.ipv4
}

output "instance_state" {
  value = data.multipass_instance.example.instance.state
}
```

### Cloud-Initを使用した高度な例

```hcl
resource "multipass_instance" "web_server" {
  name       = "web-server"
  image      = "22.04"
  cpu        = "2"
  memory     = "4G"
  disk       = "20G"
  cloud_init = "./cloud-init.yaml"
}
```

`cloud-init.yaml` ファイルを作成：
```yaml
#cloud-config
packages:
  - nginx
  - git

runcmd:
  - systemctl enable nginx
  - systemctl start nginx
```

## リソースとデータソース

### リソース

#### `multipass_instance`

Multipass Ubuntuインスタンスを管理します。

**引数：**
- `name`（必須） - インスタンス名
- `image`（オプション） - Ubuntuイメージ（デフォルト：最新LTS）
- `cpu`（オプション） - CPU数
- `memory`（オプション） - メモリ割り当て（例："1G"、"512M"）
- `disk`（オプション） - ディスク容量（例："5G"、"10G"）
- `cloud_init`（オプション） - Cloud-init設定ファイルのパス
- `timeouts`（オプション） - タイムアウト設定ブロック
  - `create`（オプション） - インスタンス作成のタイムアウト（デフォルト：15分）
  - `read`（オプション） - インスタンス読み込みのタイムアウト（デフォルト：5分）
  - `update`（オプション） - インスタンス更新のタイムアウト（デフォルト：10分）
  - `delete`（オプション） - インスタンス削除のタイムアウト（デフォルト：10分）

**属性：**
- `id` - インスタンス識別子（名前と同じ）
- `state` - 現在のインスタンス状態
- `ipv4` - インスタンスに割り当てられたIPv4アドレスのリスト

### データソース

#### `multipass_instance`

Multipassインスタンスの情報を照会します。

**引数：**
- `name`（オプション） - 特定のインスタンス名。指定されていない場合は、すべてのインスタンスを一覧表示します。

**属性：**
- `instance` - 単一インスタンス情報（`name`が指定された場合）
- `instances` - すべてのインスタンスのリスト（`name`が指定されていない場合）

## 開発

### 前提条件

- Go 1.23以上
- Terraform 1.0以上
- MultipassがインストールされてPATHでアクセス可能

### ビルド

```bash
# 依存関係のインストール
make deps

# プロバイダーのビルド
make build

# テストの実行
make test

# 受け入れテストの実行（Multipassが必要）
make testacc

# 開発用にローカルインストール
make install-local
```

### テスト

```bash
# ユニットテスト
go test ./...

# 受け入れテスト（TF_ACC=1とMultipassが必要）
TF_ACC=1 go test ./... -v

# 特定のインスタンスでテスト
cd examples/complete-examples/vm-info-output
terraform init
terraform plan
terraform apply
```

## アーキテクチャ

このプロバイダーは、モダンな[Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)を使用して構築され、以下のパターンに従います：

- **Provider**: メインプロバイダー設定とクライアント初期化
- **Resources**: MultipassインスタンスのCRUD操作
- **Data Sources**: 既存インスタンスの読み取り専用クエリ
- **Client**: Multipass CLIのGoラッパー
- **Testing**: 包括的なユニットおよび受け入れテスト

## ロードマップ

### バージョン 0.2.0（次回リリース）
- [ ] ネットワーク設定サポート
- [ ] ポートフォワーディング管理
- [ ] ボリュームマウント（`multipass mount`統合）
- [ ] より良いエラーハンドリングと検証
- [ ] 自動リリース用のCI/CDパイプライン

### バージョン 0.3.0
- [ ] スナップショット管理（`multipass snapshot`コマンド）
- [ ] インスタンスライフサイクル操作（開始、停止、サスペンド、再起動）
- [ ] コマンド実行用のExec操作
- [ ] マルチプラットフォームバイナリリリース（Linux、macOS、Windows）

### バージョン 1.0.0
- [ ] 本番環境対応
- [ ] 完全なMultipass APIカバレッジ
- [ ] 拡張されたドキュメントと例
- [ ] パフォーマンス最適化
- [ ] セキュリティ強化

## 貢献

1. リポジトリをフォーク
2. 機能ブランチを作成（`git checkout -b feature/amazing-feature`）
3. 変更を行い、テストを追加
4. テストスイートを実行（`make test`）
5. [CHANGELOG.md](CHANGELOG.md)に変更内容を更新
6. 変更をコミット（`git commit -m 'Add amazing feature'`）
7. ブランチにプッシュ（`git push origin feature/amazing-feature`）
8. プルリクエストを開く

### リリースプロセス

このプロジェクトは自動リリース用に[GoReleaser](https://goreleaser.com/)を使用します：

```bash
# 新しいタグを作成
git tag v0.1.0
git push origin v0.1.0

# GoReleaserが自動的にバイナリ付きのリリースを作成
```

## ライセンス

このプロジェクトはApache License 2.0の下でライセンスされています - 詳細は[LICENSE](LICENSE)ファイルを参照してください。

## 謝辞

- 優れたVM管理ツールを提供する[Canonical Multipass](https://multipass.run/)
- TerraformとPlugin Frameworkを提供する[HashiCorp](https://www.hashicorp.com/)
- ガイダンスとベストプラクティスを提供するTerraformコミュニティ