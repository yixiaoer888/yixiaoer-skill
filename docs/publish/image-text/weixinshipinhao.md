# 视频号图文发布参数 (WeChat Video Account Image-Text)

本平台图文发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | 否 | 标题 | - |
| `description` | `string` | 否 | 图文描述，支持 HTML (`<p>`, `<topic>`)。最多 1000 字符。 | - |
| `images` | `Array` | **是** | 图片数组 (`OldImage[]`) | - |
| `location` | `Object` | 否 | 位置信息 (`PlatformDataItem`) | - |
| `music` | `Object` | 否 | 音乐信息 (`PlatformDataItem`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳) | - |
| `collection` | `Object` | 否 | 合集信息 | - |
| `pubType` | `number` | 否 | 发布类型: 0-草稿, 1-直接发布 | 1 |

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
- **音乐 (music)**: 需通过 `[获取音乐](../../get-music.md)` 获得对应的 `PlatformDataItem`。
- **话题 (topic)**: 在 `description` 中使用 `<topic>` 标签时，话题数据需通过 `[获取话题](../../get-topics.md)` 获得。

## 4. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "imageText",
  "platforms": ["视频号"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "SPH_ACC_ID",
        "contentPublishForm": {
          "formType": "task",
          "title": "视频号动态标题",
          "description": "<p>视频号图文内容 <topic text='微信' raw='{\"id\":\"xxx\",\"name\":\"微信\"}'>#微信</topic></p>",
          "images": [
            { "key": "img_sph_01", "size": 1024, "width": 1080, "height": 1440, "format": "jpg" }
          ]
        }
      }
    ]
  }
}
```

## 5. DTO 参考
- 后端类: `ShiPingHaoDynamicForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/shipinghao.dto.ts`
