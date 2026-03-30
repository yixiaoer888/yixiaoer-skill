# 知乎图文发布参数 (Zhihu Image-Text)

本平台图文发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 想法标题 | - |
| `description` | `string` | **是** | 想法内容，支持 HTML (`<p>`, `<topic>`) | - |
| `images` | `Array` | **是** | 图片 OSS 列表 (`OldImage[]`) | - |

## 2. Payload 完整示例

```json
{
  "publishType": "imageText",
  "platforms": ["知乎"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "ZH_ACC_ID",
        "contentPublishForm": {
          "formType": "task",
          "title": "知乎想法标题",
          "description": "<p>这是想法内容 <topic text='职场'>#职场</topic></p>",
          "images": [
            { "key": "img_key_1", "size": 1024, "width": 1000, "height": 1000 }
          ]
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `ZhiHuDynamicForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/zhihu.dto.ts`
