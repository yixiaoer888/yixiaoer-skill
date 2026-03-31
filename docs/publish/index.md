# 通用内容发布基础结构 (Universal Publish Base)

所有通过 `api.ts`（指定 `action: "publish"`）执行的内容发布任务均遵循此通用结构。

> [!IMPORTANT]
> **发布合规性要求**:
> 所有的封面 (`cover`)、视频 (`video`)、图文图片 (`images`) 均**必须**使用通过[资源上传接口](../upload-resource.md)获得的资源 `key`。
> **严禁**直接填写外部网络 URL 或在该填入 Key 的地方留空。

## 1. 基础结构 (Base Structure)

接口要求传入 `CloudTaskPushRequest` 结构。

| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `publishType` | `string` | **是** | 发布类型: `article` (文章), `imageText` (图文), `video` (视频) |
| `platforms` | `string[]` | **是** | 目标平台枚举数组 (如 `["抖音", "微信公众号"]`)，详见 [平台列表](../platform.md) |
| `publishArgs` | `Object` | **是** | 发布参数核心容器 |
| `taskSetId` | `string` | 否 | 任务集唯一标识 (草稿发布时必填) |
| `desc` | `string` | 否 | 任务描述/摘要 |
| `publishChannel` | `string` | 否 | `cloud` (云端) 或 `local` (本机)，默认 `local` |
| `clientId` | `string` | 否 | 客户端连接 ID (`local` 发布时必填) |
| `isDraft` | `boolean` | 否 | 是否仅保存为草稿，默认 `false` |
| `intervalConfig` | `Object` | 否 | 间隔发布配置 (包含 `enable`, `interval`, `timeUnit` 等) |
| `coverKey` | `string` | 否 | 封面资源 Key (上传至 OSS 的路径) |


## 2. 发布参数 (publishArgs)

| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `accountForms` | `Array` | **是** | 账号发布表单列表 |
| `content` | `string` | 条件 | **内容正文**: 文章或图文类型下为**必填** (文章为 HTML，图文为纯文本) |

## 3. 账号表单项 (accountForms Item)

| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| `platformAccountId` | `string` | **是** | 蚁小二平台账号唯一 ID |
| `cover` | `Object` | **是** | **ImageFormItem**: 主封面对象 (`key`, `width`, `height`, `size`) |
| `contentPublishForm`| `Object` | **是** | **透传层**: `{formType: "task", ...}`，其他更多属性需要参考对应平台文档 |
| `video` | `Object` | 条件 | **VideoFormItem**: 视频对象 (**视频发布必填**) |
| `images` | `Array` | 条件 | **ImageFormItem[]**: 图文图片列表 (**图文发布必填**) |
| `mediaId` | `string` | 否 | 第三方库素材 ID |
| `fps` | `number` | 否 | 视频发布帧率 (海外平台使用) |
| `coverKey` | `string` | 否 | 账号级封面 Key (通常与 `cover.key` 一致) |


## 4. 调用指令 (Execution)

所有发布操作均通过 `scripts/api.ts` 脚本执行。调用时需通过 `--payload` 参数传入符合上述 JSON 结构且包含 `"action": "publish"` 的字符串。

```bash
node scripts/api.ts --payload='{"action": "publish", ...}'
```

## 5. 演示数据 (Payload Examples)

以下为几种主要发布类型的完整 Payload 示例。

### 示例 A: 文章发布 (Article)
```json
{
  "publishType": "article",
  "platforms": ["微信公众号"],
  "coverKey": "article_cover_key",
  "publishArgs": {
    "content": "<h1>演示文章标题</h1><p>这是一个演示文章的正文内容，支持 HTML 格式。</p>",
    "accountForms": [
      {
        "platformAccountId": "acc_art_001",
        "coverKey": "article_cover_key",
        "cover": {
          "key": "article_cover_key",
          "width": 900,
          "height": 500,
          "size": 150000
        },
        "contentPublishForm": {
          "formType": "task",
          // 更多属性请阅读平台文档
        }
      }
    ]
  }
}
```

### 示例 B: 图文发布 (Image-Text)
```json
{
  "publishType": "imageText",
  "platforms": ["小红书"],
  "publishArgs": {
    "content": "这是一个图文发布的描述内容，通常为纯文本。 #演示 #Demo",
    "accountForms": [
      {
        "platformAccountId": "acc_img_002",
        "images": [
          { "key": "img_key_1", "width": 1080, "height": 1440, "size": 200000 },
          { "key": "img_key_2", "width": 1080, "height": 1440, "size": 200000 }
        ],
        "coverKey": "img_key_1",
        "cover": { "key": "img_key_1", "width": 1080, "height": 1440, "size": 200000 },
        "contentPublishForm": {
          "formType": "task",
          // 更多属性请阅读平台文档
        }
      }
    ]
  }
}
```

### 示例 C: 视频发布 (Video)
```json
{
  "publishType": "video",
  "platforms": ["抖音"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "acc_vid_003",
        "video": {
          "key": "video_oss_key",
          "width": 1080,
          "height": 1920,
          "size": 52428800
        },
        "coverKey": "video_cover_key",
        "cover": {
          "key": "video_cover_key",
          "width": 1080,
          "height": 1920,
          "size": 307200
        },
        "contentPublishForm": {
          "formType": "task",
          // 更多属性请阅读平台文档
        }
      }
    ]
  }
}
```

---

### [文章发布平台列表 (Article)](./article/index.md)
### [图文发布平台列表 (Image-Text)](./image-text/index.md)
### [视频发布平台列表 (Video)](./video/index.md)

> [!TIP]
> 平台特定的 DTO 字段请查阅对应模态子目录下的各平台文档，将其放入 `contentPublishForm` 中。
> contentPublishForm 中如果存在不清楚的选填字段，请不要随意填写。
> 批量发布，不要多次调用脚本，而是在accountForms里增加多个平台的表单。
> 发布时，请务必关注表单必填项，严格按照表单规范来构建payload。