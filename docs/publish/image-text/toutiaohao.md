# 头条号图文发布参数 (Toutiaohao Image-Text)

本平台图文发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `description` | `string` | **是** | 图文描述，支持 HTML (`<p>`, `<topic>`) | - |
| `images` | `Array` | **是** | 图片 OSS 列表 (`OldImage[]`) | - |
| `pubType` | `number` | **是** | 发布类型: 0-草稿, 1-直接发布 | - |
| `declaration` | `number` | 否 | 创作类型 (1-自行拍摄, 3-AI生成等) | - |

## 2. Payload 完整示例

```json
{
  "publishType": "imageText",
  "platforms": ["头条号"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "TT_ACC_ID",
        "contentPublishForm": {
          "formType": "task",
          "description": "<p>今日头条动态测试</p>",
          "images": [
            { "key": "img_key_1", "size": 1024, "width": 720, "height": 1280 }
          ],
          "pubType": 1
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `TouTiaoHaoDynamicForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/toutiaohao.dto.ts`
