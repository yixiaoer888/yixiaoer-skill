# 快手图文发布参数 (Kuaishou Image-Text)

本平台图文发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `description` | `string` | **是** | 图文描述，支持 HTML (`<p>`, `<topic>`)。最多 1000 字符。 | - |
| `images` | `Array` | **是** | 图片数组 (`OldImage[]`) | - |
| `location` | `Object` | 否 | 位置信息 (`PlatformDataItem`) | - |
| `music` | `Object` | 否 | 音乐信息 (`PlatformDataItem`) | - |
| `scheduledTime` | `number` | 否 | 定时发布时间 (Unix 时间戳) | - |
| `collection` | `Object` | 否 | 合集信息 (`Category`) | - |
| `visibleType` | `number` | **是** | 可见类型: 0-公开, 1-私密, 3-好友可见 | 0 |

## 2. 复杂对象结构说明

### 复杂对象：OldImage
| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `width` | `number` | **是** | 图片宽度 | - |
| `height` | `number` | **是** | 图片高度 | - |
| `size` | `number` | **是** | 文件大小 (Bytes) | - |
| `key` | `string` | **是** | 资源 Key (通过上传接口获取) | - |
| `format` | `string` | **是** | 文件格式 (e.g., `jpg`, `png`) | - |

### 复杂对象：Category
| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `yixiaoerId` | `string` | **是** | 蚁小二内部 ID | - |
| `yixiaoerName` | `string` | **是** | 显示名称 | - |
| `raw` | `object` | 否 | 原始数据对象 | - |

## 3. 依赖接口说明

若字段值需通过查询获得，需注明：
- **位置 (location)**: 需通过 `[获取位置](../../get-location.md)` 获得对应的 `PlatformDataItem`。
- **音乐 (music)**: 需通过 `[获取音乐](../../get-music.md)` 获得对应的 `PlatformDataItem`。
- **话题 (topic)**: 在 `description` 中使用 `<topic>` 标签时，话题数据需通过 `[获取话题](../../get-topics.md)` 获得。
- **合集 (collection)**: 需通过 `[获取合集](../../get-collections.md)` 获得对应的 `Category`。

## 4. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "imageText",
  "platforms": ["快手"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "KS_ACC_ID",
        "contentPublishForm": {
          "formType": "task",
          "description": "<p>快手动态内容 <topic text='快手' raw='{\"id\":\"xxx\",\"name\":\"快手\"}'>#快手</topic></p>",
          "images": [
            { "key": "img_ks_01", "size": 1024, "width": 1080, "height": 1440, "format": "jpg" }
          ],
          "visibleType": 0
        }
      }
    ]
  }
}
```

## 5. DTO 参考
- 后端类: `KuaiShouDynamicForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/kuaishou.dto.ts`
