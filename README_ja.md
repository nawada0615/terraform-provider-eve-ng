# EVE-NG用Terraformプロバイダー

EVE-NG（Emulated Virtual Environment Next Generation）リソースを管理するためのTerraformプロバイダーです。信頼性の向上と包括的なエラーハンドリングを備えています。

## 機能

このプロバイダーは、堅牢なエラーハンドリングとフォールバック機能を備えた以下のEVE-NGリソースをサポートしています：

### リソース
- **eve_folder** - EVE-NGフォルダの管理
- **eve_lab** - EVE-NGラボの管理
- **eve_network** - ラボネットワークの管理（可視性制御対応）
- **eve_node** - ラボノードの管理（QEMU、Dynamips、IOL、Docker、VPCS）
- **eve_interface_attachment** - ノードインターフェース接続の管理

### データソース
- **eve_templates** - 利用可能なノードテンプレートの取得
- **eve_network_types** - 利用可能なネットワークタイプの取得
- **eve_icons** - 利用可能なアイコンの取得
- **eve_status** - システムステータスの取得

## 最近の改善

### ✅ 信頼性の向上
- **フォールバック機能**: ネットワーク個別読み取り失敗時に自動的にネットワークリストにフォールバック
- **包括的エラーハンドリング**: 詳細なエラーメッセージと適切なステータスコード検証
- **型安全性**: APIレスポンス型の変動に対する堅牢な処理
- **デバッグログ**: トラブルシューティングとモニタリングのための広範なログ

### 🛡️ 堅牢なエラーハンドリング
- **APIレスポンス検証**: すべてのAPIレスポンスの適切な検証
- **ステート管理**: リソースが見つからないシナリオの正しい処理
- **認証**: 詳細なメッセージを含むログインエラーハンドリングの強化
- **ネットワーク管理**: フォールバックサポートを備えた回復力のあるネットワーク操作

### 🎨 ネットワーク可視性制御
- **非表示ネットワーク**: `visibility = "0"`を設定してEVE-NG GUIからネットワークを非表示
- **クリーンなラボインターフェース**: 視覚的な乱雑さなしに接続専用ネットワークを作成
- **プロフェッショナルな外観**: プレゼンテーション用のクリーンなラボ図を維持

## インストール

### ローカルインストール

1. リポジトリをクローン
2. プロバイダーをビルド:
   ```bash
   make build
   ```
3. ローカルにインストール:
   ```bash
   make install
   ```

### 使用方法

```hcl
terraform {
  required_providers {
    eve-ng = {
      source = "local/eve-ng"
      version = "1.0.0"
    }
  }
}

provider "eve-ng" {
  endpoint = "http://192.168.0.101"
  username = "admin"
  password = "eve"
  insecure_skip_verify = true
}

# ラボを作成
resource "eve_lab" "test_lab" {
  path        = "/"
  name        = "test-lab"
  author      = "Terraform"
  description = "Terraformで作成されたテストラボ"
  version     = "1"
}

# ネットワークを作成
resource "eve_network" "net1" {
  lab_file = eve_lab.test_lab.file
  name     = "net1"
  type     = "bridge"
}

# カスタムオプションでQEMUノードを作成
resource "eve_node" "linux_node" {
  lab_file       = eve_lab.test_lab.file
  name           = "linux-node"
  type           = "qemu"
  template       = "linux"
  cpu            = 2
  ram            = 2048
  qemu_version   = "2.4.0"
  qemu_arch      = "x86_64"
  qemu_nic       = "virtio-net-pci"
  qemu_options   = "-enable-kvm"
  cpulimit       = true
  desired_state  = "started"
}

# ノードをネットワークに接続
resource "eve_interface_attachment" "net1_connection" {
  lab_file        = eve_lab.test_lab.file
  node_id         = tonumber(split(":node:", eve_node.linux_node.id)[1])
  interface_index = 0
  target          = "network:${tonumber(split(":network:", eve_network.net1.id)[1])}"
}

# 2番目のノードを作成
resource "eve_node" "linux_node2" {
  lab_file       = eve_lab.test_lab.file
  name           = "linux-node2"
  type           = "qemu"
  template       = "linux"
  cpu            = 1
  ram            = 1024
  qemu_version   = "2.4.0"
  qemu_arch      = "x86_64"
  qemu_nic       = "virtio-net-pci"
  qemu_options   = "-enable-kvm"
  cpulimit       = true
  desired_state  = "started"
}

# ノード間接続用の2番目のネットワークを作成
resource "eve_network" "connection_net" {
  lab_file   = eve_lab.test_lab.file
  name       = "connection-net"
  type       = "bridge"
  visibility = "0"  # 非表示ネットワーク（EVE-NG GUIで表示されない）
}

# 両方のノードを同じネットワークに接続して直接通信を可能にする
resource "eve_interface_attachment" "node1_to_net" {
  lab_file        = eve_lab.test_lab.file
  node_id         = tonumber(split(":node:", eve_node.linux_node.id)[1])
  interface_index = 0
  target          = "network:${tonumber(split(":network:", eve_network.connection_net.id)[1])}"
}

resource "eve_interface_attachment" "node2_to_net" {
  lab_file        = eve_lab.test_lab.file
  node_id         = tonumber(split(":node:", eve_node.linux_node2.id)[1])
  interface_index = 0
  target          = "network:${tonumber(split(":network:", eve_network.connection_net.id)[1])}"
}
```

## 開発

### 前提条件
- Go 1.21以降
- Terraform 1.0以降

### ビルド
```bash
make build
```

### テスト
```bash
# すべてのテストを実行
make test

# レース検出でテストを実行
make test-race

# 特定のテストスイートを実行
make test-unit
make test-client
make test-mock
```

### コード品質
```bash
# コードをフォーマット
make fmt

# リンターを実行
make lint

# vetを実行
make vet
```

## 貢献

1. リポジトリをフォーク
2. 機能ブランチを作成
3. 変更を加える
4. 変更のテストを追加
5. テストスイートを実行
6. プルリクエストを送信

## ライセンス

このプロジェクトはMITライセンスの下でライセンスされています。詳細はLICENSEファイルを参照してください。

## 謝辞

- 優れたネットワークエミュレーションプラットフォームを提供するEVE-NGコミュニティ
- TerraformプラグインSDKを提供するHashiCorp
