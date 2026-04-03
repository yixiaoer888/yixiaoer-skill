# 百家号图文发布参数 (BaiJiaHao Image-Text)

> [!IMPORTANT]
> **前提条件 (Prerequisite)**:
> 在使用本平台的特定参数之前，你 **必须** 已经阅读并理解了 [图文发布首页 (Index)](./index.md) 中定义的 Payload 根结构。本页仅描述 `contentPublishForm` 内部的平台差异化字段。

本平台图文发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 标题 (1-20 字符) | - |
| `description` | `string` | **是** | 图文描述，支持 HTML (`<p>`, `<topic>`)。 (1-1000 字符) | - |
| `cover` | `Object` | **是** | 封面图对象 (`OldCover`) | - |
| `pubType` | `number` | **是** | 发布类型: 0-草稿, 1-直接发布 | - |
| `declaration` | `number` | **是** | 创作声明: 0-不声明, 1-内容由 AI 生成 | - |
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

### 复杂对象：PlatformDataItem
| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `id` | `string` | **是** | 位置 ID | - |
| `text` | `string` | **是** | 项名称/文本 | - |
| `raw` | `object` | 否 | 原始数据对象 | - |

## 3. 依赖接口说明

若字段值需通过查询获得，需注明：
- **位置 (location)**: 需通过 `[获取位置](../../get-locations.md)` 获得对应的 `PlatformDataItem`。
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
        "platformAccountId": "acc_bjh_it_001",
        "contentPublishForm": {
          "formType": "task",
          "title": "百家号图文标题",
          "description": "<p>内容摘要 <topic text='科技' raw='{\"yixiaoerId\":\"123\",\"yixiaoerName\":\"科技\",\"raw\":{\"id\":\"xxx\",\"topic\":\"科技\"}}'>#科技</topic></p>",
          "cover": { 
            "key": "img_key_123", 
            "size": 102400, 
            "width": 800, 
            "height": 600 
          },
          "pubType": 1,
          "declaration": 0
        }
      }
    ]
  }
}
```

## 相关接口

| 目标数据 | 对应 Action | 相关文档 |
| :--- | :--- | :--- |
| `images.key` | `upload` | [资源上传](../../upload-resource.md) |
