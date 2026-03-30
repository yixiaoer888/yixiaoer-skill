# 百家号图文发布参数 (Baijiahao Image-Text)

本平台图文发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 标题 | - |
| `description` | `string` | **是** | 内容描述，支持 HTML (`<p>`, `<topic>`) | - |
| `cover` | `Object` | **是** | 封面图 (`OldCover`) | - |
| `declaration` | `number` | **是** | 创作声明: 0-不声明, 1-内容由AI生成 | - |
| `location` | `Object` | 否 | 位置信息 | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳) | - |

> [!NOTE]
> 百家号图文 (动态) 后端 DTO 暂未直接支持 `images` 数组，而是通过 `cover` 承载主要图片，更多内容建议放入 `description`。

## 2. Payload 完整示例

```json
{
  "publishType": "imageText",
  "platforms": ["百家号"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "BJH_ACC_ID",
        "contentPublishForm": {
          "formType": "task",
          "title": "百家号动态标题",
          "description": "<p>这是动态内容 <topic text='百度'>#百度</topic></p>",
          "cover": { "key": "img_key", "size": 1024, "width": 800, "height": 600 },
          "declaration": 0
        }
      }
    ]
  }
}
```

## 3. DTO 参考
- 后端类: `BaiJiaHaoDynamicForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/baijiahao.dto.ts`
