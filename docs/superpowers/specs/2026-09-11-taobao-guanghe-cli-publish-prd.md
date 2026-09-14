# CLI 支持淘宝光合发布与商品挂车产品需求文档

## 1. 文档信息

| 项目 | 内容 |
| --- | --- |
| 产品 | `yxer` CLI |
| 需求 | 淘宝光合视频、图文发布与商品挂车 |
| 文档状态 | 最终方案，待实施 |
| 日期 | 2026-09-11 |
| 主要改动仓库 | `yixiaoer-skill` |
| 配套验证仓库 | `yixiaoer-universal` |

## 2. 背景

蚁小二 Web 端已经支持淘宝光合视频发布、图文发布和关联商品。`yxer` CLI 目前没有注册淘宝光合平台，也没有封装淘宝光合专属的商品分类、商品列表和挂车流程，导致人类用户和 AI Agent 无法通过 CLI 完成与 Web 一致的淘宝光合发布。

Web 与 CLI 应复用同一套网关和平台发布能力。CLI 不直接访问淘宝光合，不管理淘宝登录态、Cookie 或 `publishSession`，也不复制 Web 页面组件和平台运行时逻辑。

现有网关已具备主体能力：

- API Key 可解析为固定的 `userId`、`teamId` 和 `authSource=cli_apikey`；
- 淘宝光合商品分类接口已存在；
- 淘宝光合商品列表接口已存在；
- 商品会话缓存按 `teamId + platformAccountId + pageType` 隔离；
- `POST /taskSets/v2` 已支持淘宝光合视频和图文发布；
- 服务端可将完整 `shopping_cart` 商品对象透传至淘宝光合发布运行时。

因此，本期以 CLI 能力补齐为主，不修改网关生产功能代码；网关仅补验证和自动化测试。若验证发现现有链路不符合契约，只针对失败点另行评估最小修复。

## 3. 产品目标

本期完成后，用户和 AI Agent 可以通过 `yxer`：

1. 查询当前团队可用的淘宝光合账号；
2. 按视频或图文场景查询淘宝光合商品分类；
3. 按商品来源、筛选条件、关键词和分页游标查询商品；
4. 从真实查询结果中选择最多 6 件商品；
5. 发布不挂商品或挂 1 至 6 件商品的淘宝光合视频；
6. 发布不挂商品或挂 1 至 6 件商品的淘宝光合图文；
7. 在正式发布前完成账号确认、资源上传、Payload 校验和 `--dry-run`；
8. 获得稳定、机器可读的 JSON 成功结果和结构化错误。

## 4. 本期范围

### 4.1 包含范围

#### `yixiaoer-skill` CLI

- 注册淘宝光合平台及常用别名；
- 支持淘宝光合 `video` 和 `imageText` 两种发布类型；
- 增加商品分类和商品列表查询命令；
- 支持商品筛选、关键词搜索和游标翻页；
- 支持将查询返回的完整商品对象写入账号级 `shopping_cart`；
- 增加淘宝光合视频和图文 JSON Schema；
- 增加淘宝光合专属 preflight、来源校验和错误提示；
- 确保淘宝光合商品不进入抖音购物车归一化；
- 确保淘宝光合不调用通用购物车 `entitlements` 接口；
- 支持 `prepare`、`schema fields`、`schema get` 和 `publish form` 发现淘宝光合字段；
- 支持 `validate`、`publish --dry-run` 和正式发布；
- 更新 Skill、平台文档和命令参考；
- 补充 CLI 自动化测试和 JSON 输出契约测试。

#### `yixiaoer-universal` 网关

- 验证 API Key 可以访问淘宝光合两个商品接口；
- 补充 API Key 与团队隔离的集成测试；
- 验证跨团队请求不会调用淘宝上游，也不会创建商品会话缓存；
- 验证商品会话不会跨团队、跨账号或跨 `pageType` 复用；
- 不修改网关生产功能代码。

### 4.2 不包含范围

- 淘宝光合文章发布；
- CLI 直接登录淘宝或直接调用淘宝接口；
- CLI 管理 Cookie、`publishSession` 或淘宝授权凭据；
- 修改 Web 页面或 Web 现有请求结构；
- 商品签名令牌 `goodsSelectionToken`；
- 发布请求幂等键；
- 网关新增淘宝专属稳定业务错误码；
- 网关商品响应 DTO 重构；
- 网关将超过 6 件商品的静默截断改为业务错误；
- 同团队内按负责人或运营人限制账号商品访问；
- 自动选择多个账号或多个商品候选；
- 自动重试正式发布。

## 5. 系统边界

```text
用户 / AI Agent
      |
      v
yxer CLI（yixiaoer-skill）
  - 命令和 JSON 输出
  - Schema 与本地预检
  - 资源上传
  - 商品查询与来源记录
      |
      | Authorization: API Key
      v
Gateway（yixiaoer-universal）
  - API Key 鉴权
  - 团队和账号隔离
  - publishSession 缓存
  - 任务集创建
  - 最终业务校验
      |
      v
淘宝光合发布运行时
```

边界原则：

- CLI 是网关客户端，不是新的平台适配层；
- 网关是身份、团队、账号和最终发布结果的权威；
- `platformAccountId` 是商品查询和账号级发布的边界；
- 商品查询结果是挂车 Payload 的唯一数据来源；
- `publishSession` 只在服务端保存，不返回 CLI；
- API Key 只发送给蚁小二网关，不写入 Payload、日志或错误详情。

## 6. 用户流程

### 6.1 标准流程

```text
doctor
  -> 查询并确认淘宝光合账号
  -> prepare / schema fields
  -> 上传视频、封面或图片
  -> 查询商品分类
  -> 查询并选择商品（可选）
  -> 组装或导出 payload.json
  -> validate
  -> publish --dry-run
  -> 用户授权
  -> publish
```

### 6.2 候选确认

- 只有一个可用账号时，可以自动选择，并向用户说明选择依据；
- 多个账号时必须展示账号稳定 ID、名称、平台和状态，由用户确认；
- 多个商品候选时必须通过商品 ID 或索引明确选择；
- CLI 不得默认选择第一个商品；
- 没有 `status=1` 的淘宝光合账号时停止流程；
- 商品查询为空时允许不挂车发布，不得编造商品对象。

### 6.3 发布门禁

正式发布前必须使用同一份 `payload.json` 和同一套发布通道参数依次完成：

```text
yxer validate
yxer publish --dry-run
用户授权
yxer publish
```

修改账号、资源、商品、内容、发布时间或发布通道后，必须重新执行 `validate` 和 `publish --dry-run`。

## 7. 平台与类型

### 7.1 平台标识

CLI canonical key：

```text
taobaoguanghe
```

至少支持以下输入别名：

```text
淘宝光合
taobaoguanghe
taobao-guanghe
TaoBaoGuangHe
```

所有发布请求中的平台名称统一规范化为：

```json
"platforms": ["淘宝光合"]
```

### 7.2 内容类型映射

| CLI 输入 | CLI 标准发布类型 | 商品接口 `pageType` |
| --- | --- | --- |
| `video` | `video` | `video` |
| `imageText` | `imageText` | `photo` |

不支持 `article`。传入不支持的类型时，CLI 在请求远端前返回结构化使用错误。

## 8. 命令需求

命令名称沿用现有 `yxer query` 风格。

### 8.1 查询账号

```bash
yxer accounts list 淘宝光合 --status 1 --json
```

账号必须来自 CLI 实际响应，不能手写或从其他团队复用。

### 8.2 查询商品分类

```bash
yxer query taobao-guanghe-goods-tabs <account_id> --type video --json
yxer query taobao-guanghe-goods-tabs <account_id> --type imageText --json
```

请求接口：

```text
GET /platform-accounts/{accountId}/taobao-guanghe/goods-tabs?pageType={video|photo}
```

要求：

- `--type` 只接受 `video` 或 `imageText`；
- CLI 将 `imageText` 映射成 `photo`；
- 返回值仅执行通用响应信封解包，不改写分类对象；
- `source`、`filterValue` 和二级筛选值均以接口返回为准；
- CLI 不根据分类名称推导 `source`。

### 8.3 查询商品列表

```bash
yxer query taobao-guanghe-goods <account_id> \
  --type video \
  --source <source> \
  --filter-value <filterValue> \
  --second-filter-value <secondLevelFilterValue> \
  --query <keyword> \
  --next-page <nextPage> \
  --json
```

请求接口：

```text
GET /platform-accounts/{accountId}/taobao-guanghe/goods
```

查询参数：

| CLI 参数 | 网关参数 | 必填 | 说明 |
| --- | --- | --- | --- |
| `--type` | `pageType` | 是 | `video -> video`，`imageText -> photo` |
| `--query` | `keyword` | 否 | 商品关键词 |
| `--keyword` | `keyword` | 否 | `--query` 的兼容别名 |
| `--next-page` | `nextPage` | 否 | 不透明分页游标 |
| `--filter-value` | `filterValue` | 否 | 一级筛选值 |
| `--second-filter-value` | `secondLevelFilterValue` | 否 | 二级筛选值 |
| `--source` | `source` | 否 | 商品来源，必须来自分类响应 |

要求：

- `nextPage` 不解析、不拼接、不改写；
- 所有查询参数只进行一次标准 URL 编码；
- 分页结束以响应中 `nextPage` 缺失、为 `null` 或为空字符串为准；
- 首版不提供自动拉取全部页面的 `--all`；
- 返回商品对象时完整保留 `raw`；
- CLI 日志不得打印 `raw`、API Key 或服务端会话信息。

### 8.4 页面式表单选择

淘宝光合 `shopping_cart` 必须注册为查询型动态字段，路径固定为：

```text
publishArgs.accountForms[].contentPublishForm.shopping_cart
```

推荐流程：

```bash
yxer publish form start 淘宝光合 video --output publish-form.json
yxer query taobao-guanghe-goods-tabs <account_id> --type video --json
yxer query taobao-guanghe-goods <account_id> --type video --json
yxer publish form choose publish-form.json shopping_cart \
  --value-file goods.json \
  --id <yixiaoerId> \
  --source-command "yxer query taobao-guanghe-goods <account_id> --type video --json"
```

`choose`、`verify`、`review` 和 `export` 必须校验：

- 来源命令是淘宝光合商品命令；
- 查询账号与目标账号一致；
- 查询类型与发布类型一致；
- 当前商品对象与来源记录的哈希一致；
- 多账号表单分别拥有各自的来源记录；
- 手工修改商品对象后来源校验失败。

## 9. 商品数据契约

### 9.1 商品对象

发布使用商品接口返回的完整对象：

```json
{
  "yixiaoerId": "goods-123",
  "yixiaoerName": "商品名称",
  "yixiaoerImageUrl": "https://example.com/image.jpg",
  "price": 129,
  "yixiaoerDesc": "商品描述",
  "raw": {}
}
```

必填字段：

- `yixiaoerId`：非空字符串；
- `yixiaoerName`：非空字符串；
- `raw`：对象。

可选字段：

- `yixiaoerImageUrl`；
- `price`；
- `yixiaoerDesc`；
- 服务端返回的其他字段。

### 9.2 数据使用规则

- 商品对象必须来自目标账号最近一次 CLI 查询；
- CLI 不裁剪、重构或补造 `raw`；
- `yixiaoerId` 仅用于候选选择和去重，不能替代完整商品对象；
- 不使用抖音的 `sale_title/images/data` 包装结构；
- 不将淘宝商品转换成通用 `ShoppingCartItem`；
- 一个账号最多关联 6 件商品；
- 同一账号表单内 `yixiaoerId` 不得重复；
- 第 7 件商品必须在 CLI 请求网关前失败，不能静默截断；
- 多账号发布时，账号 A 的商品不得用于账号 B。

## 10. 发布 Payload

### 10.1 视频

```json
{
  "action": "publish",
  "publishType": "video",
  "platforms": ["淘宝光合"],
  "publishChannel": "cloud",
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "<account_id>",
        "video": {
          "key": "<uploaded_video_key>",
          "size": 12345678,
          "width": 1080,
          "height": 1920,
          "duration": 20,
          "format": "mp4"
        },
        "cover": {
          "key": "<uploaded_cover_key>",
          "size": 345678,
          "width": 1080,
          "height": 1920,
          "format": "jpg"
        },
        "coverKey": "<uploaded_cover_key>",
        "contentPublishForm": {
          "formType": "task",
          "short_title": "视频短标题",
          "desc": "视频描述",
          "tags": [],
          "allow_download": false,
          "statement": { "type": 0 },
          "shopping_cart": []
        }
      }
    ]
  }
}
```

### 10.2 图文

```json
{
  "action": "publish",
  "publishType": "imageText",
  "platforms": ["淘宝光合"],
  "publishChannel": "cloud",
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "<account_id>",
        "images": [
          {
            "key": "<uploaded_image_key>",
            "size": 456789,
            "width": 1080,
            "height": 1440,
            "format": "jpg"
          }
        ],
        "contentPublishForm": {
          "formType": "task",
          "title": "图文标题",
          "desc": "图文正文",
          "tags": [],
          "statement": { "type": 0 },
          "shopping_cart": []
        }
      }
    ]
  }
}
```

示例中的占位符仅用于说明。正式 Payload 必须使用 `accounts list`、`upload` 和商品查询的真实结果。

### 10.3 发布时间

- 立即发布：不传 `scheduledTime`；
- 定时发布：传 13 位 Unix 毫秒时间戳 `scheduledTime`；
- CLI 在平台归一化阶段将 `scheduledTime` 映射为服务端兼容的 `prePubTime`；
- `scheduledTime: 0` 不作为 CLI 的立即发布写法，避免与通用 13 位时间戳校验冲突；
- CLI 可执行稳定的格式检查，易变化的时间窗口由服务端最终校验。

### 10.4 发布通道

本期产品文档和验收以云发布为准：

```text
publishChannel=cloud
```

已确认本机客户端具备淘宝光合运行时，淘宝光合视频和图文纳入本机发布支持列表。CLI 使用与其他平台一致的 `--publish-channel local --client-id` 通道参数，不自动切换通道。

## 11. 本地校验规则

### 11.1 视频

- 必须有账号级 `video`；
- 必须有账号级 `cover` 和与其一致的 `coverKey`；
- `short_title` 最多 30 个字符；
- `desc` 最多 1000 个字符；
- `allow_download` 为布尔值；
- `statement.type` 只允许 `0、1、2、4、5、6`；
- `statement.type=4` 时允许非空 `forwardFrom`；
- `shopping_cart` 为数组，最多 6 件。

视频标题和描述允许同时为空，与当前 Web 视频校验保持一致。

### 11.2 图文

- 图片数量为 1 至 9 张；
- 每张图片的宽和高均不得小于 720 像素；
- 标题和正文至少一个非空；
- 标题最多 30 个字符；
- 正文最多 1000 个字符；
- `tags` 最多 4 个；
- `statement` 与视频规则一致；
- `shopping_cart` 为数组，最多 6 件。

### 11.3 商品

- 商品必须含非空 `yixiaoerId`、非空 `yixiaoerName` 和对象类型 `raw`；
- 同一账号内商品 ID 不得重复；
- 商品数量超过 6 时返回结构化错误；
- 商品图片等商品元数据 URL 不作为待上传媒体外链拦截；
- `raw` 中的 URL 不参与通用外链拒绝；
- 商品查询结果不得跨账号或跨内容类型复用。

## 12. 兼容处理

### 12.1 通用购物车权限

淘宝光合不使用当前通用 `/platform-accounts/{id}/entitlements` 检查。`validate`、`publish --dry-run` 和正式发布遇到淘宝光合 `shopping_cart` 时，应跳过该通用接口，由淘宝光合商品查询和服务端发布链路完成最终资格判断。

### 12.2 购物车归一化

淘宝光合商品必须保持扁平完整对象，不得进入以下抖音兼容逻辑：

- 包装为 `sale_title/images/data`；
- 从 `raw` 派生顶层图片；
- 自动补充购买标题；
- 将完整对象移动到 `data`。

### 12.3 最终请求

`publish --dry-run` 展示的 `request.publishArgs.accountForms[].contentPublishForm.shopping_cart` 与正式发布发往 `/taskSets/v2` 的结构必须一致。正式发布边界不得再次裁剪、改名或包装淘宝商品。

## 13. 输出合同

### 13.1 成功输出

- stdout 只能输出 JSON；
- 查询成功沿用统一成功信封；
- `action` 使用稳定值：
  - `taobao-guanghe-goods-tabs`
  - `taobao-guanghe-goods`
- 商品响应完整放在 `data` 中；
- 不向 stdout 混入分页说明、提示或诊断文本。

示例：

```json
{
  "ok": true,
  "action": "taobao-guanghe-goods",
  "version": "<cli-version>",
  "data": {
    "dataList": [],
    "nextPage": null
  }
}
```

### 13.2 错误输出

- stdout 保持为空；
- stderr 输出结构化 JSON；
- Agent 可处理错误必须使用 `yxerrors.Error`；
- 错误至少包含 `code`、`message`、`category`；
- 有修复路径时包含 `hint`；
- 有明确下一步时包含 `nextCommand`；
- 仅在安全重试时设置 `retryable=true`。

建议的 CLI 错误：

| 场景 | code | category | retryable |
| --- | --- | --- | --- |
| 内容类型不支持 | `taobao_guanghe_invalid_content_type` | `taobao_guanghe_goods` | `false` |
| 商品超过 6 件 | `taobao_guanghe_goods_limit` | `taobao_guanghe_goods` | `false` |
| 商品结构不完整 | `taobao_guanghe_goods_invalid` | `taobao_guanghe_goods` | `false` |
| 商品来源账号不一致 | `taobao_guanghe_goods_account_mismatch` | `taobao_guanghe_goods_source` | `false` |
| 商品来源类型不一致 | `taobao_guanghe_goods_type_mismatch` | `taobao_guanghe_goods_source` | `false` |
| 商品查询失败 | 沿用远程错误 code | `taobao_guanghe_goods_query` | 取决于远端错误 |

账号状态需要更新时，CLI 应提示用户在蚁小二打开该账号的“创作者后台”更新状态，并给出重新查询命令。

## 14. 安全与隐私

- API Key 不得出现在 stdout、stderr、日志、Payload、测试快照和请求哈希明文中；
- `publishSession` 不得返回 CLI；
- 普通人工摘要不打印商品 `raw`；
- 结构化 JSON 查询结果可以包含服务端明确返回的完整商品对象，以便后续发布；
- 错误详情不得回显完整请求头；
- CLI 不接受客户端传入 `teamId` 或 `userId` 来覆盖 API Key 身份；
- 跨团队账号请求应返回 `404`，不暴露目标账号是否真实存在；
- 正式发布超时不得自动重试，用户应先查询发布记录确认任务是否已创建。

## 15. 网关验证与测试补充

本期不修改网关生产功能代码，但必须完成以下验证。

### 15.1 API Key 集成测试

至少准备团队 A、团队 B，各自拥有独立 API Key 和淘宝光合账号。

| API Key | 请求账号 | 接口 | 预期 |
| --- | --- | --- | --- |
| 团队 A Key | 团队 A 账号 | `goods-tabs` | `200`，响应结构正确 |
| 团队 A Key | 团队 A 账号 | `goods` | `200`，响应结构正确 |
| 团队 B Key | 团队 B 账号 | `goods-tabs` | `200`，响应结构正确 |
| 团队 B Key | 团队 B 账号 | `goods` | `200`，响应结构正确 |
| 团队 A Key | 团队 B 账号 | 两个接口 | `404` |
| 团队 B Key | 团队 A 账号 | 两个接口 | `404` |
| 无 Key | 任意账号 | 两个接口 | `401` |
| 过期 Key | 任意账号 | 两个接口 | `401` |
| 已退出团队的 Key | 原团队账号 | 两个接口 | `401` |
| 冻结成员 Key | 本团队账号 | 两个接口 | `401` |

额外断言：

- API Key 请求的 `authSource` 为 `cli_apikey`；
- 跨团队请求不调用淘宝开放平台服务；
- 跨团队请求不创建或读取目标账号的商品会话缓存；
- 跨团队错误不包含 `openAccountId`、商品或团队详情；
- 只携带 API Key、不携带 Web Cookie 时也能成功访问本团队账号。

### 15.2 缓存隔离测试

缓存维度必须为：

```text
teamId + platformAccountId + pageType
```

验证：

- 同一账号的 `video` 与 `photo` 使用不同缓存；
- 不同账号的相同 `pageType` 使用不同缓存；
- 不同团队不会读取彼此缓存；
- 缺失缓存时，网关可先获取 tabs 建立会话，再查询商品；
- 跨团队请求在建立缓存前失败。

### 15.3 验证失败处理

若网关验证失败：

1. 停止淘宝光合 CLI 上线；
2. 记录失败用例及真实响应；
3. 单独评估最小网关修复；
4. 修复通过全部隔离测试后再恢复上线；
5. 不允许在 CLI 中绕过或模拟服务端权限检查。

## 16. CLI 自动化测试

实施采用 TDD，至少覆盖：

1. 平台中文名、canonical key 和全部别名正确归一化；
2. `accounts list 淘宝光合` 使用标准中文平台名请求账号；
3. tabs 请求路径准确，`video` 和 `imageText` 映射正确；
4. goods 请求完整透传所有可选查询参数；
5. 中文关键词正确编码，`nextPage` 不被解析或二次编码；
6. tabs、goods 成功输出的 JSON 信封和 `action` 稳定；
7. 查询失败时 stdout 为空，stderr 为结构化 JSON；
8. 视频和图文 schema 可被 `schema list/get/fields` 和 `prepare` 发现；
9. 视频接受合规资源和表单，拒绝超过长度和非法声明；
10. 图文接受 1 至 9 张合规图片，拒绝 0 张、10 张和小于 720 像素的图片；
11. 图文拒绝标题与正文同时为空；
12. 商品 0、1、6 件通过，7 件失败；
13. 商品缺少 ID、名称或 `raw` 时失败；
14. 同一账号存在重复商品 ID 时失败；
15. 淘宝商品完整对象和 `raw` 在 validate、dry-run、云发布请求中保持不变；
16. 淘宝商品不会被包装成抖音购物车结构；
17. 淘宝挂车不会调用通用 `entitlements`；
18. `publish form choose` 拒绝错误命令来源、跨账号来源、跨类型来源和对象漂移；
19. 立即发布不需要 `scheduledTime`，定时发布接受 13 位毫秒时间戳；
20. 淘宝光合文章发布失败；
21. `--publish-channel local` 在客户端运行时支持时可通过校验并进入本机发布链路；
22. 现有平台、通用商品查询和购物车发布测试保持通过；
23. 全量 Go 测试通过，stdout/stderr 合同不回归。

## 17. 端到端验收

使用测试环境真实淘宝光合账号完成：

1. CLI 查询视频商品分类和商品列表；
2. CLI 查询图文商品分类和商品列表，确认请求使用 `pageType=photo`；
3. 不挂商品的视频完成 validate、dry-run 和真实发布；
4. 挂 1 件商品的视频发布成功，淘宝侧关联商品正确；
5. 挂 6 件商品的视频发布成功；
6. 第 7 件商品在 CLI 创建任务前被拒绝；
7. 不挂商品的图文发布成功；
8. 挂 1 件和 6 件商品的图文发布成功；
9. 图文图片数量和 720x720 下限生效；
10. 商品会话过期后，网关现有自动恢复逻辑可正常工作；
11. 账号状态过期时，CLI 返回可执行的更新账号提示；
12. 双团队交叉访问均返回 404，且没有上游调用和缓存污染；
13. dry-run 请求中的商品对象与最终 `/taskSets/v2` 请求一致；
14. 发布成功返回有效 `taskSetId`，发布失败不误报成功。

## 18. 上线门槛

以下条件全部满足后才允许发布 CLI 版本：

- CLI 单元测试和命令测试全部通过；
- 全量 Go 测试通过；
- 网关 API Key 集成测试通过；
- 双团队隔离测试通过；
- 视频和图文真实账号端到端验收通过；
- Web 淘宝光合视频、图文和挂车回归通过；
- API Key、`publishSession` 和敏感 `raw` 未出现在日志；
- Skill、平台文档、命令参考和版本记录已更新；
- 新 CLI 包完成安装验证，`yxer doctor` 正常；
- 正式发布失败时不会自动重试创建任务。

## 19. 发布与回滚

### 19.1 发布策略

- 先发布测试版本；
- 仅对内部测试团队开放淘宝光合 CLI 流程；
- 完成视频、图文、不挂车、1 件和 6 件商品验证；
- 观察商品查询失败、账号状态失败和发布失败分布；
- 确认无跨团队访问和缓存串用后再发布稳定版本。

### 19.2 回滚策略

如发现权限、缓存隔离或错误发布问题：

- 下架包含淘宝光合能力的 CLI 版本；
- 保持网关和 Web 原有能力不变；
- 在 CLI 平台能力列表中移除淘宝光合入口；
- 不删除用户账号、商品或既有发布任务；
- 修复并重新完成全部上线门槛后再发布。

## 20. 后续需求

以下事项进入后续版本评估：

- 网关签发短期 `goodsSelectionToken`，绑定团队、账号、内容类型和商品快照哈希；
- `/taskSets/v2` 支持幂等键，避免超时重试产生重复任务；
- 网关显式拒绝超过 6 件商品，不再静默截断；
- 为淘宝商品接口补正式响应 DTO 和 OpenAPI Schema；
- 提供稳定的账号过期、会话过期、商品失效和商品超限错误码；
- 增加同团队内账号负责人或运营人级权限；
- 验证淘宝光合本机发布在客户端运行时中的兼容性；
- 评估受限制的自动翻页能力。

## 21. 最终决策

本期采用“CLI 完整商品快照透传 + 现有统一发布入口”的方案：

- 功能实现集中在 `yixiaoer-skill`；
- `yixiaoer-universal` 不修改生产功能代码，只补 API Key、团队隔离和缓存隔离验证；
- CLI 支持淘宝光合视频和图文的云发布与本机发布；
- 商品必须来自目标账号和目标内容类型的真实查询结果，最多 6 件；
- 商品签名令牌、幂等键和网关契约硬化延期处理；
- 任一网关权限或隔离验证失败，均阻止本功能上线。
