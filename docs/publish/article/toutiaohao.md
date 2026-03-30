# 头条号文章发布参数 (TouTiaoHao Article)

本平台文章发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 文章标题 | - |
| `content` | `string` | **是** | 文章 HTML 正文 | - |
| `covers` | `Array` | **是** | 封面图 OSS 列表 (`OldCover[]`) | - |
| `isFirst` | `boolean` | 否 | 是否在头条首发 | - |
| `location` | `Object` | 否 | 位置字段 (`PlatformDataItem`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (时间戳) | - |
| `advertisement` | `number` | 否 | 广告投放赚取收益 (2-是, 3-否) | `3` |
| `declaration` | `number` | 否 | 创作声明 (1:自行拍摄, 2:取自站外, 3:AI生成, 6:虚构演绎, 7:投资观点, 8:健康医疗) | - |

## 2. Payload 完整示例

```json
{
  "publishType": "article",
  "platforms": ["TouTiaoHao"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "YOUR_ACCOUNT_ID",
        "contentPublishForm": {
          "formType": "task",
          "title": "今日头条发布测试",
          "content": "<p>正文内容...</p>",
          "covers": [
            { "key": "cover_key", "size": 100, "width": 800, "height": 600 }
          ],
          "advertisement": 3
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `TouTiaoHaoArticleForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/toutiaohao.dto.ts`
