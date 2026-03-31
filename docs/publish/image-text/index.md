# 图文发布 (Image-Text Publish)

所有通过 `api.ts`（指定 `action: "publish"`）执行的图文发布任务均遵循以下数据结构。

> [!IMPORTANT]
> **发布合规性要求**:
> 所有的封面 (`cover`)、图文图片 (`images`) 均**必须**使用通过[资源上传接口](../../upload-resource.md)获得的资源 `key`。
> **严禁**直接填写外部网络 URL 或在该填入 Key 的地方留空。

## 1. 数据结构 (Data Structure)

接口要求传入 `CloudTaskPushRequest` 结构。

### 1.1 基础结构 (Base Structure)

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `publishType` | `string` | **是** | 固定为 `imageText` | - |
| `platforms` | `string[]` | **是** | 目标平台枚举数组，详见下方平台列表 | - |
| `coverKey` | `string` | **是** | 任务封面资源 Key | - |
| `publishArgs` | `Object` | **是** | 发布参数核心容器 | - |
| `taskSetId` | `string` | 否 | 任务集唯一标识 (草稿发布时必填) | - |
| `desc` | `string` | 否 | 任务描述/摘要 | - |
| `publishChannel` | `string` | 否 | `cloud` (云端) 或 `local` (本机) | `local` |
| `clientId` | `string` | 否 | 客户端连接 ID (`local` 发布时必填) | - |
| `isDraft` | `boolean` | 否 | 是否仅保存为草稿 | `false` |

### 1.2 发布参数 (publishArgs)

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `content` | `string` | **是** | **图文描述**: 纯文本格式 | - |
| `accountForms` | `Array` | **是** | 账号发布表单列表 | - |

### 1.3 账号表单项 (accountForms Item)

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `platformAccountId` | `string` | **是** | 蚁小二平台账号唯一 ID | - |
| `images` | `Array` | **是** | **ImageFormItem[]**: 图文图片列表 (`key`, `width`, `height`, `size`) | - |
| `cover` | `Object` | **是** | **ImageFormItem**: 主封面对象 | - |
| `contentPublishForm`| `Object` | **是** | **透传层**: `{formType: "task", ...}`，其他更多属性需要参考对应平台文档 | - |
| `coverKey` | `string` | **是** | 账号级封面 Key (必须与 `cover.key` 一致) | - |

## 2. 发布示例 (Payload Example)

```json
{
  "action": "publish",
  "publishType": "imageText",
  "platforms": ["小红书"],
  "coverKey": "img_key_1",
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
          "formType": "task"
        }
      }
    ]
  }
}
```

## 3. 支持平台列表 (Support Platforms)

以下平台支持通过 `publishType: "imageText"` 进行发布。

| 平台名称 | 标识符 | 文档链接 |
| :--- | :--- | :--- |
| **抖音** | `抖音`, `DouYin` | [douyin.md](./douyin.md) |
| **小红书** | `小红书`, `XiaoHongShu` | [xiaohongshu.md](./xiaohongshu.md) |
| **快手** | `快手`, `KuaiShou` | [kuaishou.md](./kuaishou.md) |
| **新浪微博** | `新浪微博`, `XinLangWeiBo` | [xinlangweibo.md](./xinlangweibo.md) |
| **视频号** | `视频号`, `ShiPinHao` | [weixinshipinhao.md](./weixinshipinhao.md) |
| **百家号** | `百家号`, `BaiJiaHao` | [baijiahao.md](./baijiahao.md) |
| **头条号** | `头条号`, `TouTiaoHao` | [toutiaohao.md](./toutiaohao.md) |
| **知乎** | `知乎`, `ZhiHu` | [zhihu.md](./zhihu.md) |

> [!TIP]
> 持续增加中... 请参考后端 DTO `*DynamicForm` 扩展新平台。

