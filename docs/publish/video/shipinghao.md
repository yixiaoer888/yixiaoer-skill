# 视频号 视频发布

## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| title | string | 否 | 视频标题 | - |
| short_title | string | 否 | 视频短标题 | - |
| description | string | 否 | 视频描述，支持 HTML 格式和 `@` 好友/话题标签 | - |
| horizontalCover | object | 否 | 视频横板封面，使用 `OldCover` 结构 | - |
| type | number | 是 | 视频原创类型：1-非原创，2-原创 | 1 |
| createType | number | 是 | 创建类型：1-草稿，2-直接发布 | 2 |
| pubType | number | 否 | 发布类型：0-草稿，1-直接发布 | 1 |
| location | object | 否 | 视频位置，使用 `PlatformDataItem` 结构 | - |
| scheduledTime | number | 否 | 定时发布时间戳（单位：秒） | - |
| shoppingCart | object | 否 | 关联商品信息 | - |
| collection | object | 否 | 合集信息 | - |
| activity | object | 否 | 活动信息 | - |

## 2. 复杂对象结构

### OldCover
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| width | number | 是 | 宽度 |
| height | number | 是 | 高度 |
| size | number | 是 | 文件大小 (Byte) |
| path | string | 否 | 文件绝对路径 (与 `key` 二选一) |
| key | string | 否 | OSS 对象存储 Key (与 `path` 二选一) |

### PlatformDataItem
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| id | string | 是 | ID |
| text | string | 是 | 文本内容 |
| raw | object | 是 | 平台原始数据 |

## 3. JSON 示例

```json
{
  "publishType": "video",
  "platforms": ["Shipinghao"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "SPH_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1080,
          "height": 1920,
          "duration": 15
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "视频号视频标题",
          "description": "<p>这是视频号内容的精彩描述 #生活 #记录</p>",
          "type": 2,
          "createType": 2,
          "pubType": 1,
          "location": {
            "id": "sph_loc_1",
            "text": "广州市",
            "raw": {}
          }
        }
      }
    ]
  }
}
```
