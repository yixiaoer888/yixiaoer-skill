# 头条号图文发布参数 (Toutiaohao Image-Text)

本平台图文发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `description` | `string` | **是** | 图文描述，支持 HTML (`<p>`, `<topic>`)。最多 1000 字符。 | - |
| `images` | `Array` | **是** | 图片数组 (`OldImage[]`) | - |
| `declaration` | `number` | 否 | 创作类型: 1-自行拍摄, 2-取自站外, 3-AI生成, 6-虚构演绎... | - |
| `pubType` | `number` | **是** | 发布类型: 0-草稿, 1-直接发布 | 1 |

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
- **话题 (topic)**: 在 `description` 中使用 `<topic>` 标签时，话题数据需通过 `[获取话题](../../get-topics.md)` 获得。

## 4. Payload 完整示例

```json
{
  "action": "publish",
  "publishType": "imageText",
  "platforms": ["头条号"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "TTH_ACC_ID",
        "contentPublishForm": {
          "formType": "task",
          "description": "<p>今日头条内容 <topic text='头条' raw='{\"id\":\"xxx\",\"name\":\"头条\"}'>#头条</topic></p>",
          "images": [
            { "key": "img_tth_01", "size": 1024, "width": 1080, "height": 1440, "format": "jpg" }
          ],
          "pubType": 1
        }
      }
    ]
  }
}
```

## 5. DTO 参考
- 后端类: `TouTiaoHaoDynamicForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/toutiaohao.dto.ts`
