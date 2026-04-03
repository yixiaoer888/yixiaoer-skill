# 新浪微博图文发布参数 (Sina Weibo Image-Text)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [图文发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。


本平台图文发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `description` | `string` | **是** | 微博博文内容，支持 HTML (`<p>`, `<topic>`)。最多 1000 字符。 | - |
| `images` | `Array` | **是** | 图片数组 (`OldImage[]`) | - |
| `location` | `Object` | 否 | 位置信息 (`PlatformDataItem`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳) | - |

## 2. 复杂对象结构说明

### 复杂对象：OldImage
| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `width` | `number` | **是** | 图片宽度 | - |
| `height` | `number` | **是** | 图片高度 | - |
| `size` | `number` | **是** | 文件大小 (Bytes) | - |
| `key` | `string` | **是** | 资源 Key (通过上传接口获取) | - |
| `format` | `string` | **是** | 文件格式 (e.g., `jpg`, `png`) | - |

## 3. 依赖接口说明

若字段值需通过查询获得，需注明：
- **位置 (location)**: 需通过 `[获取位置](../../get-location.md)` 获得对应的 `PlatformDataItem`。
- **话题 (topic)**: 在 `description` 中使用 `<topic>` 标签时，话题数据需通过 `[获取话题](../../get-topics.md)` 获得。

## 4. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "imageText",
  "platforms": ["新浪微博"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "WB_ACC_ID",
        "contentPublishForm": {
          "formType": "task",
          "description": "<p>微博动态内容 <topic text='微博' raw='{\"id\":\"xxx\",\"name\":\"微博\"}'>#微博</topic></p>",
          "images": [
            { "key": "img_wb_01", "size": 1024, "width": 1080, "height": 1440, "format": "jpg" }
          ]
        }
      }
    ]
  }
}
```

## 5. DTO 参考
- 后端类: `XinLangWeiBoDynamicForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/xinlangweibo.dto.ts`
