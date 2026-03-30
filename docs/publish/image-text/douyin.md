# 抖音图文发布参数 (Douyin Image-Text)

本平台图文发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | 否 | 抖音图文标题 | - |
| `description` | `string` | **是** | 图文描述，支持 HTML 格式（`<p>` 标签及 `<topic>` 标签） | - |
| `images` | `Array` | **是** | 图片 OSS 列表 (`OldImage[]`) | - |
| `location` | `Object` | 否 | 物理地址信息 | - |
| `musice` | `Object` | 否 | 音乐素材信息 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间戳 | - |
| `collection` | `Object` | 否 | 合集信息 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "imageText",
  "platforms": ["抖音"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "YOUR_ACCOUNT_ID",
        "images": [
          { "key": "img1_oss_key", "size": 100, "width": 800, "height": 600 }
        ],
        "contentPublishForm": {
          "formType": "task",
          "title": "图文标题",
          "description": "<p>这是一个关于 <topic text='搞笑' raw='{\\\"yixiaoerId\\\":\\\"123\\\"}'>#搞笑</topic> 的示例</p>"
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `DouYinDynamicForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/douyin.dto.ts`
