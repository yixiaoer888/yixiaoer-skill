# 百家号图文发布参数 (Baijiahao Image-Text)

本平台图文发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 标题 | - |
| `description` | `string` | **是** | 图文描述，支持 HTML (`<p>`, `<topic>`)。最多 1000 字符。 | - |
| `cover` | `Object` | **是** | 封面图对象 (`OldCover`) | - |
| `declaration` | `number` | **是** | 创作声明: 0-不声明, 1-内容由AI生成 | - |
| `location` | `Object` | 否 | 位置信息 (`PlatformDataItem`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳) | - |

## 2. 复杂对象结构说明

### 复杂对象：OldCover
| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `width` | `number` | **是** | 图片宽度 | - |
| `height` | `number` | **是** | 图片高度 | - |
| `size` | `number` | **是** | 文件大小 (Bytes) | - |
| `key` | `string` | **是** | 资源 Key (通过上传接口获取) | - |
| `path` | `string` | 否 | 资源路径 (可选) | - |

### 复杂对象：PlatformDataItem
| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `string` | **是** | 外部平台项 ID | - |
| `text` | `string` | **是** | 项名称/文本 | - |
| `raw` | `object` | 否 | 原始数据对象 | - |

## 3. 依赖接口说明

若字段值需通过查询获得，需注明：
- **位置 (location)**: 需通过 `[获取位置](../../get-location.md)` 获得对应的 `PlatformDataItem`。
- **话题 (topic)**: 在 `description` 中使用 `<topic>` 标签时，话题数据需通过 `[获取话题](../../get-topics.md)` 获得。

## 4. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "imageText",
  "platforms": ["百家号"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "BJH_ACC_ID",
        "contentPublishForm": {
          "formType": "task",
          "title": "百家号动态标题",
          "description": "<p>这是内容摘要 <topic text='科技' raw='{\"id\":\"xxx\",\"name\":\"科技\"}'>#科技</topic></p>",
          "cover": { 
            "key": "img_key_123", 
            "size": 102400, 
            "width": 800, 
            "height": 600 
          },
          "declaration": 0
        }
      }
    ]
  }
}
```

## 5. DTO 参考
- 后端类: `BaiJiaHaoDynamicForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/baijiahao.dto.ts`
