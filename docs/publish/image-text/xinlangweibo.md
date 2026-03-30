# 新浪微博图文发布参数 (Weibo Image-Text)

本平台图文发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `description` | `string` | **是** | 图文内容，支持 HTML (`<p>`, `<topic>`) | - |
| `images` | `Array` | **是** | 图片 OSS 列表 (`OldImage[]`) | - |
| `location` | `Object` | 否 | 位置信息 (`PlatformDataItem`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳) | - |

## 2. Payload 完整示例

```json
{
  "publishType": "imageText",
  "platforms": ["新浪微博"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "WEIBO_ACC_ID",
        "contentPublishForm": {
          "formType": "task",
          "description": "<p>微博图文测试 <topic text='热点'>#热点</topic></p>",
          "images": [
            { "key": "img_key_1", "size": 1024, "width": 800, "height": 600 }
          ]
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `XinLangWeiBoDynamicForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/xinlangweibo.dto.ts`
