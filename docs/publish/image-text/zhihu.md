# 知乎图文发布参数 (Zhihu Image-Text)

本平台图文发布通过 `contentPublishForm` 承载以下参数。

## 1. contentPublishForm 参数定义

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| `formType` | `string` | **是** | 固定值: `task` | `task` |
| `title` | `string` | **是** | 标题 | - |
| `description` | `string` | **是** | 图文内容，支持 HTML (`<p>`, `<topic>`) | - |
| `images` | `Array` | **是** | 图片数组 (`OldImage[]`) | - |

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
  "platforms": ["知乎"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "ZH_ACC_ID",
        "contentPublishForm": {
          "formType": "task",
          "title": "知乎图文标题",
          "description": "<p>知乎想法内容 <topic text='知识' raw='{\"id\":\"xxx\",\"name\":\"知识\"}'>#知识</topic></p>",
          "images": [
            { "key": "img_zh_01", "size": 1024, "width": 1080, "height": 1440, "format": "jpg" }
          ]
        }
      }
    ]
  }
}
```

## 5. DTO 参考
- 后端类: `ZhiHuDynamicForm`
- 文件路径: `apps/server-api/packages/yxr-open-platform/src/models/platform/zhihu.dto.ts`
