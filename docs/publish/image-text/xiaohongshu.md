# 小红书图文发布参数 (Xiaohongshu Image-Text)

本平台图文发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | 否 | 笔记标题 (最多20字) | - |
| `description` | `string` | **是** | 笔记内容，支持 HTML (`<p>`, `<topic>`) | - |
| `images` | `Array` | **是** | 图片 OSS 列表 (`OldImage[]`) | - |
| `visibleType` | `number` | **是** | 可见性: 0-公开, 1-私密, 3-好友 | - |
| `location` | `Object` | 否 | 位置信息 | - |
| `music` | `Object` | 否 | 音乐挂载 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 | - |
| `collection` | `Object` | 否 | 合集信息 | - |

## 2. Payload 完整示例

```json
{
  "publishType": "imageText",
  "platforms": ["小红书"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "XHS_ACC_ID",
        "images": [{ "key": "img_key", "size": 100, "width": 800, "height": 600 }],
        "contentPublishForm": {
          "formType": "task",
          "title": "小红书图文测试",
          "description": "<p>这是描述 <topic text='测试'>#测试</topic></p>",
          "visibleType": 0
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `XiaoHongShuDynamicForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/xiaohongshu.dto.ts`
