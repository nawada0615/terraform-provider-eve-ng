# EVE-NG API 仕様書

## 概要
このドキュメントは、EVE-NG（EVE-NG Community Edition 6.2.0-4）のREST APIの仕様をまとめたものです。実機での検証結果に基づいて作成されています。

## 基本情報
- **ベースURL**: `http://192.168.0.101/api/`
- **認証方式**: Cookie-based認証（`unetlab_session`）
- **Content-Type**: `application/json`
- **APIバージョン**: 6.2.0-4

## 認証

### ログイン
```http
POST /api/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "eve",
  "html5": 1
}
```

**レスポンス**:
```json
{
  "code": 200,
  "status": "success",
  "message": "User logged in (90013)."
}
```

**注意**: `html5`パラメータは必須です。

### ログアウト
```http
GET /api/auth/logout
Cookie: unetlab_session=<session_id>
```

### 認証確認
```http
GET /api/auth
Cookie: unetlab_session=<session_id>
```

## システム情報

### ステータス取得
```http
GET /api/status
Cookie: unetlab_session=<session_id>
```

**レスポンス**:
```json
{
  "code": 200,
  "status": "success",
  "message": "Fetched system status (60001).",
  "data": {
    "version": "6.2.0-4",
    "qemu_version": "2.4.0",
    "uksm": "unsupported",
    "ksm": "enabled",
    "cpulimit": "enabled",
    "cpu": 0,
    "disk": 35,
    "cached": 98,
    "mem": 3,
    "swap": 0,
    "iol": 0,
    "dynamips": 0,
    "qemu": 0,
    "docker": 0,
    "vpcs": 0
  }
}
```

## フォルダ管理

### フォルダ一覧取得
```http
GET /api/folders/
Cookie: unetlab_session=<session_id>
```

**レスポンス**:
```json
{
  "code": 200,
  "status": "success",
  "message": "Successfully listed path (60007).",
  "data": {
    "folders": [
      {
        "name": "projects",
        "path": "/projects"
      }
    ],
    "labs": [
      {
        "file": "lab.unl",
        "path": "/lab.unl",
        "umtime": 1760274335,
        "mtime": "12 Oct 2025 15:05"
      }
    ]
  }
}
```

### フォルダ作成
```http
POST /api/folders
Content-Type: application/json
Cookie: unetlab_session=<session_id>

{
  "name": "folder_name",
  "path": "/parent_path"
}
```

### フォルダ編集
```http
PUT /api/folders/<path>
Content-Type: application/json
Cookie: unetlab_session=<session_id>

{
  "path": "/new_path"
}
```

### フォルダ削除
```http
DELETE /api/folders/<path>
Cookie: unetlab_session=<session_id>
```

## ラボ管理

### ラボ作成
```http
POST /api/labs
Content-Type: application/json
Cookie: unetlab_session=<session_id>

{
  "path": "/",                    // Required: ラボのパス
  "name": "lab_name",             // Required: ラボ名
  "author": "author_name",        // Optional: 作成者名 (デフォルト: "")
  "description": "lab_description", // Optional: 説明 (デフォルト: "")
  "body": "lab_body",             // Optional: 本文 (デフォルト: "")
  "version": "1",                 // Optional: バージョン (デフォルト: "1")
  "scripttimeout": 600            // Optional: スクリプトタイムアウト秒 (デフォルト: 300)
}
```

**レスポンス**:
```json
{
  "code": 200,
  "status": "success",
  "message": "Lab has been created (60019)."
}
```

### ラボ取得
```http
GET /api/labs/<lab_file>
Cookie: unetlab_session=<session_id>
```

**レスポンス**:
```json
{
  "code": 200,
  "status": "success",
  "message": "Lab has been loaded (60020).",
  "data": {
    "author": "test",
    "description": "test lab",
    "body": "",
    "filename": "test-lab.unl",
    "id": "613a38ee-4dd8-46d6-bae7-e3ebb6107ff1",
    "name": "test-lab",
    "version": 0,
    "scripttimeout": 600,
    "lock": 0
  }
}
```

### ラボ削除
```http
DELETE /api/labs/<lab_file>
Cookie: unetlab_session=<session_id>
```

## ノード管理

### ノード作成
```http
POST /api/labs/<lab_file>/nodes
Content-Type: application/json
Cookie: unetlab_session=<session_id>

{
  "name": "node_name",              // Required: ノード名
  "type": "qemu",                   // Required: ノードタイプ (qemu, dynamips, iol, docker, vpcs)
  "template": "linux",              // Required: テンプレート名
  "image": "linux",                 // Optional: イメージ名 (デフォルト: "")
  "icon": "Router-2D-Gen-White-S.svg", // Optional: アイコン (デフォルト: "")
  "top": 100,                       // Optional: Y座標 (デフォルト: 0)
  "left": 100,                      // Optional: X座標 (デフォルト: 0)
  "delay": 0,                       // Optional: 起動遅延秒 (デフォルト: 0)
  "config": "",                     // Optional: 設定 (デフォルト: "")
  "cpu": 1,                         // Optional: CPU数 (デフォルト: 0 = テンプレートのデフォルト)
  "ram": 1024,                      // Optional: RAM (MB) (デフォルト: 0 = テンプレートのデフォルト)
  "ethernet": 4,                    // Optional: Ethernetポート数 (デフォルト: 0 = テンプレートのデフォルト)
  "serial": 0,                      // Optional: シリアルポート数 (デフォルト: 0)
  "cpulimit": false,                // Optional: CPU制限 (デフォルト: false)
  "uuid": "",                       // Optional: UUID (デフォルト: 自動生成)
  "qemu_version": "2.4.0",          // Optional: QEMUバージョン (デフォルト: "")
  "qemu_arch": "x86_64",            // Optional: QEMUアーキテクチャ (デフォルト: "")
  "qemu_nic": "virtio-net-pci",     // Optional: QEMU NICタイプ (デフォルト: "")
  "qemu_options": "",               // Optional: QEMUオプション (デフォルト: "")
  "firstmac": "",                   // Optional: 最初のMACアドレス (デフォルト: 自動生成)
  "management_address": ""          // Optional: 管理アドレス (デフォルト: "")
}
```

**レスポンス**:
```json
{
  "code": 201,
  "status": "success",
  "message": "Lab has been saved (60023).",
  "data": {
    "id": 1
  }
}
```

### ノード取得
```http
GET /api/labs/<lab_file>/nodes/<node_id>
Cookie: unetlab_session=<session_id>
```

**レスポンス**:
```json
{
  "code": 200,
  "status": "success",
  "message": "Successfully listed node (60025).",
  "data": {
    "console": "telnet",
    "config": 0,
    "delay": 0,
    "left": 100,
    "icon": "Router-2D-Gen-White-S.svg",
    "image": "linux-ubuntu-22.04",
    "name": "test-node",
    "status": 0,
    "template": "linux",
    "type": "qemu",
    "top": 100,
    "url": "/html5/#/client/...",
    "cpulimit": 1,
    "cpu": 1,
    "ethernet": 4,
    "ram": 1024,
    "uuid": "9de1d69f-5f2d-4e47-bb3f-312125576dc8",
    "firstmac": "00:50:00:00:01:00",
    "qemu_options": "-machine type=pc,accel=kvm -vga std -usbdevice tablet -boot order=cd -cpu host",
    "qemu_version": "2.4.0",
    "qemu_arch": "x86_64",
    "qemu_nic": "virtio-net-pci"
  }
}
```

### ノード編集
```http
PUT /api/labs/<lab_file>/nodes/<node_id>
Content-Type: application/json
Cookie: unetlab_session=<session_id>

{
  "id": 1,
  "name": "updated_node_name",
  "top": 150,
  "left": 150
}
```

### ノード削除
```http
DELETE /api/labs/<lab_file>/nodes/<node_id>
Cookie: unetlab_session=<session_id>
```

### ノード操作
```http
GET /api/labs/<lab_file>/nodes/<node_id>/start
GET /api/labs/<lab_file>/nodes/<node_id>/stop
GET /api/labs/<lab_file>/nodes/<node_id>/wipe
Cookie: unetlab_session=<session_id>
```

### ノードインターフェース管理

#### ノードインターフェース取得
```http
GET /api/labs/<lab_file>/nodes/<node_id>/interfaces
Cookie: unetlab_session=<session_id>
```

**レスポンス**:
```json
{
  "code": 200,
  "status": "success",
  "message": "Successfully listed node interfaces (60030).",
  "data": {
    "id": 1,
    "sort": "qemu",
    "ethernet": [
      {
        "name": "e0",
        "network_id": 0
      },
      {
        "name": "e1",
        "network_id": 0
      },
      {
        "name": "e2",
        "network_id": 0
      },
      {
        "name": "e3",
        "network_id": 0
      }
    ],
    "serial": []
  }
}
```

#### ノードインターフェース更新
```http
PUT /api/labs/<lab_file>/nodes/<node_id>/interfaces
Content-Type: application/json
Cookie: unetlab_session=<session_id>

{
  "0": 1,
  "1": 2,
  "2": 3
}
```

**パラメータ**:
- `{interface_id}`: インターフェースID（数値）
- `{network_id}`: 接続先ネットワークID（数値）

**レスポンス**:
```json
{
  "code": 201,
  "status": "success",
  "message": "Lab has been saved (60023)."
}
```

**説明**:
- ノードのインターフェースとネットワークの接続を設定します
- キーはインターフェースID、値は接続先のネットワークIDを指定します
- 例：`{"0":1}` は、インターフェース0をネットワーク1に接続することを意味します

## ネットワーク管理

### ネットワーク作成
```http
POST /api/labs/<lab_file>/networks
Content-Type: application/json
Cookie: unetlab_session=<session_id>

{
  "name": "network_name",    // Required: ネットワーク名
  "type": "bridge",          // Required: ネットワークタイプ (bridge, pnet0-9)
  "icon": "cloud.png",       // Optional: アイコン (デフォルト: "")
  "top": 200,                // Optional: Y座標 (デフォルト: 0)
  "left": 200,               // Optional: X座標 (デフォルト: 0)
  "visibility": "1"          // Optional: 可視性 (デフォルト: "1")
}
```

**レスポンス**:
```json
{
  "code": 201,
  "status": "success",
  "message": "Network has been added to the lab (60006).",
  "data": {
    "id": 1
  }
}
```


### ネットワーク一覧取得
```http
GET /api/labs/<lab_file>/networks
Cookie: unetlab_session=<session_id>
```

### ネットワーク編集
```http
PUT /api/labs/<lab_file>/networks/<network_id>
Content-Type: application/json
Cookie: unetlab_session=<session_id>

{
  "visibility": 0,
  "name": "updated_network_name",
  "top": 200,
  "left": 200
}
```

**レスポンス**:
```json
{
  "code": 201,
  "status": "success",
  "message": "Lab has been saved (60023)."
}
```

### ネットワーク削除
```http
DELETE /api/labs/<lab_file>/networks/<network_id>
Cookie: unetlab_session=<session_id>
```

## テンプレート管理

### テンプレート一覧取得
```http
GET /api/list/templates/
Cookie: unetlab_session=<session_id>
```

**レスポンス**:
```json
{
  "code": 200,
  "status": "success",
  "message": "Successfully listed node templates (60003).",
  "data": {
    "a10": "A10 vThunder.missing",
    "alienvault": "AlienVault Cybersecurity.missing",
    "android": "Android VM.missing",
    "linux": "Linux",
    "vios": "Cisco vIOS Router.missing",
    "viosl2": "Cisco vIOS Switch.missing",
    "vpcs": "Virtual PC (VPCS)",
    "vyos": "VyOS.missing"
  }
}
```

### ネットワークタイプ取得
```http
GET /api/list/networks
Cookie: unetlab_session=<session_id>
```

**レスポンス**:
```json
{
  "code": 200,
  "status": "success",
  "message": "Successfully listed network types (60002).",
  "data": {
    "bridge": "bridge",
    "ovs": "ovs",
    "pnet0": "pnet0",
    "pnet1": "pnet1",
    "pnet2": "pnet2",
    "pnet3": "pnet3",
    "pnet4": "pnet4",
    "pnet5": "pnet5",
    "pnet6": "pnet6",
    "pnet7": "pnet7",
    "pnet8": "pnet8",
    "pnet9": "pnet9"
  },
  "icons": {
    "01-Cloud-Default.svg": "01-Cloud-Default.svg",
    "cloud.png": "cloud.png",
    "Router-2D-Gen-White-S.svg": "Router-2D-Gen-White-S.svg"
  }
}
```

## エラーレスポンス

### 一般的なエラー形式
```json
{
  "code": 404,
  "status": "fail",
  "message": "Resource not found (60038)."
}
```

### 認証エラー
```json
{
  "code": 412,
  "status": "unauthorized",
  "message": "User is not authenticated or session timed out (90001)."
}
```

### 権限エラー
```json
{
  "code": 401,
  "status": "forbidden",
  "message": "User is not authenticated or session timed out (90001)."
}
```

## パラメータ仕様一覧

### ラボリソース (eve_lab)

| パラメータ | 必須/任意 | デフォルト値 | 説明 |
|-----------|---------|------------|------|
| path | Required | - | ラボのパス（`/`で始まる必要がある） |
| name | Required | - | ラボ名 |
| author | Optional | "" | 作成者名 |
| description | Optional | "" | 説明 |
| body | Optional | "" | 本文 |
| version | Optional | "1" | バージョン |
| scripttimeout | Optional | 300 | スクリプトタイムアウト（秒） |
| lock | Optional | false | ロック状態 |

### ノードリソース (eve_node)

| パラメータ | 必須/任意 | デフォルト値 | 説明 |
|-----------|---------|------------|------|
| lab_file | Required | - | ラボファイルパス |
| name | Required | - | ノード名 |
| type | Required | - | ノードタイプ (qemu, dynamips, iol, docker, vpcs) |
| template | Required | - | テンプレート名 |
| image | Optional | "" | イメージ名 |
| icon | Optional | "" | アイコン |
| top | Optional | 0 | Y座標 |
| left | Optional | 0 | X座標 |
| delay | Optional | 0 | 起動遅延（秒） |
| config | Optional | "" | 設定 |
| cpu | Optional | 0 | CPU数（0はテンプレートのデフォルト） |
| ram | Optional | 0 | RAM (MB)（0はテンプレートのデフォルト） |
| ethernet | Optional | 0 | Ethernetポート数（0はテンプレートのデフォルト） |
| serial | Optional | 0 | シリアルポート数 |
| cpulimit | Optional | false | CPU制限 |
| uuid | Optional | "" | UUID（空文字列の場合は自動生成） |
| qemu_version | Optional | "" | QEMUバージョン |
| qemu_arch | Optional | "" | QEMUアーキテクチャ |
| qemu_nic | Optional | "" | QEMU NICタイプ |
| qemu_options | Optional | "" | QEMUオプション |
| firstmac | Optional | "" | 最初のMACアドレス（空文字列の場合は自動生成） |
| timos_line | Optional | "" | TiMOS line設定 |
| timos_license | Optional | "" | TiMOSライセンス |
| management_address | Optional | "" | 管理アドレス |
| desired_state | Optional | "stopped" | 希望する電源状態 (stopped, started) |
| reboot_on_change | Optional | false | 変更時に再起動 |
| wipe_on_destroy | Optional | false | 削除時にデータをワイプ |

### ネットワークリソース (eve_network)

| パラメータ | 必須/任意 | デフォルト値 | 説明 |
|-----------|---------|------------|------|
| lab_file | Required | - | ラボファイルパス |
| name | Required | - | ネットワーク名 |
| type | Required | - | ネットワークタイプ (bridge, pnet0-9) |
| top | Optional | 0 | Y座標 |
| left | Optional | 0 | X座標 |
| icon | Optional | "" | アイコン |
| visibility | Optional | "1" | 可視性 |

## 実用例：ノード間接続

GUIでノードAとノードBを直接接続する場合のAPI呼び出しの流れを示します。

### 前提条件
- ラボ `test-connection-lab-v2.unl` が既に存在
- ノード1（ノードA）とノード2（ノードB）が既に作成済み

### ステップ1: ノードのインターフェース情報を確認

#### ノード1のインターフェース取得
```http
GET /api/labs/test-connection-lab-v2.unl/nodes/1/interfaces
Cookie: unetlab_session=<session_id>
```

**レスポンス**:
```json
{
  "code": 200,
  "status": "success",
  "message": "Successfully listed node interfaces (60030).",
  "data": {
    "id": 1,
    "sort": "qemu",
    "ethernet": [
      {
        "name": "e0",
        "network_id": 0
      },
      {
        "name": "e1",
        "network_id": 0
      },
      {
        "name": "e2",
        "network_id": 0
      },
      {
        "name": "e3",
        "network_id": 0
      }
    ],
    "serial": []
  }
}
```

#### ノード2のインターフェース取得
```http
GET /api/labs/test-connection-lab-v2.unl/nodes/2/interfaces
Cookie: unetlab_session=<session_id>
```

**レスポンス**:
```json
{
  "code": 200,
  "status": "success",
  "message": "Successfully listed node interfaces (60030).",
  "data": {
    "id": 2,
    "sort": "qemu",
    "ethernet": [
      {
        "name": "e0",
        "network_id": 0
      },
      {
        "name": "e1",
        "network_id": 0
      },
      {
        "name": "e2",
        "network_id": 0
      },
      {
        "name": "e3",
        "network_id": 0
      }
    ],
    "serial": []
  }
}
```

### ステップ2: 接続用ネットワークを作成

```http
POST /api/labs/test-connection-lab-v2.unl/networks
Content-Type: application/json
Cookie: unetlab_session=<session_id>

{
  "count": 1,
  "name": "Net-esxiiface_0",
  "type": "bridge",
  "left": 458,
  "top": 154,
  "visibility": 1,
  "postfix": 0
}
```

**レスポンス**:
```json
{
  "code": 201,
  "status": "success",
  "message": "Network has been added to the lab (60006).",
  "data": {
    "id": 1
  }
}
```

### ステップ3: ノード1のインターフェースをネットワークに接続

```http
PUT /api/labs/test-connection-lab-v2.unl/nodes/1/interfaces
Content-Type: application/json
Cookie: unetlab_session=<session_id>

{
  "0": 1
}
```

**レスポンス**:
```json
{
  "code": 201,
  "status": "success",
  "message": "Lab has been saved (60023)."
}
```

### ステップ4: ノード2のインターフェースを同じネットワークに接続

```http
PUT /api/labs/test-connection-lab-v2.unl/nodes/2/interfaces
Content-Type: application/json
Cookie: unetlab_session=<session_id>

{
  "0": 1
}
```

**レスポンス**:
```json
{
  "code": 201,
  "status": "success",
  "message": "Lab has been saved (60023)."
}
```

### ステップ5: ネットワークの可視性を非表示に設定（オプション）

```http
PUT /api/labs/test-connection-lab-v2.unl/networks/1
Content-Type: application/json
Cookie: unetlab_session=<session_id>

{
  "visibility": 0
}
```

**レスポンス**:
```json
{
  "code": 201,
  "status": "success",
  "message": "Lab has been saved (60023)."
}
```

### 結果
この一連のAPI呼び出しにより、ノード1のe0インターフェースとノード2のe0インターフェースが同じネットワーク（ID: 1）に接続され、ノード間で直接通信が可能になります。

## 注意事項

1. **パス形式**: すべてのラボパスは`/`で始まる必要があります
2. **認証**: ログイン時に`html5`パラメータが必須です
3. **デフォルト値の動作**: 数値の`0`や空文字列`""`をデフォルト値として設定した場合、EVE-NGはテンプレートの設定を使用します
4. **オプショナルパラメータ**: 上記の表でOptionalとマークされたパラメータは省略可能で、省略した場合はデフォルト値が使用されます
5. **ネットワーク接続**: ノード間を接続するには、両方のノードのインターフェースを同じネットワークIDに設定する必要があります
