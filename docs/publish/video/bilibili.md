# 哔哩哔哩 视频发布

## 1. contentPublishForm 数据结构

| 字段名 | 类型 | 必填 | 说明 | 默认值 |
| :--- | :--- | :--- | :--- | :--- |
| formType | string | 是 | 固定为 `task` | `task` |
| title | string | 是 | 视频标题 | - |
| description | string | 是 | 视频描述 | - |
| tags | string[] | 是 | 视频标签 | - |
| category | object[] | 是 | 视频分类，使用 `CascadingPlatformDataItem[]` 结构 | - |
| declaration | number | 是 | 创作者申明：0-不申明, 1-AI合成, 2-危险行为, 3-仅供娱乐, 4-引人不适, 5-理性适度消费, 6-个人观点 | 0 |
| type | number | 是 | 类型：1-自制, 2-转载 | 1 |
| scheduledTime | number | 否 | 定时发布时间戳（单位：秒） | - |
| contentSourceUrl | string | 否 | 原文 URL 链接，当 `type` 为转载时必填 | - |
| collection | object | 否 | 合集信息，使用 `Collection` 结构 | - |

## 2. 复杂对象结构

### OldCover
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| width | number | 是 | 宽度 |
| height | number | 是 | 高度 |
| size | number | 是 | 文件大小 (Byte) |
| path | string | 否 | 文件绝对路径 (与 `key` 二选一) |
| key | string | 否 | OSS 对象存储 Key (与 `path` 二选一) |

### CascadingPlatformDataItem (多级分类)
Bilibili 的分类通常是一个数组，包含一级分区和二级分区的 `PlatformDataItem`。

| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| id | string | 是 | ID |
| text | string | 是 | 文本内容 |
| raw | object | 是 | 平台原始数据 |

### Collection
| 字段名 | 类型 | 必填 | 说明 |
| :--- | :--- | :--- | :--- |
| yixiaoerId | string | 是 | 合集 ID |
| yixiaoerName | string | 是 | 合集名称 |

### 数据获取途径

| 目标字段 | 对应 Action | 文档参考 |
| :--- | :--- | :--- |
| `category` | `categories` | [获取账号分类](../../get-publish-categories.md) |
| `collection` | `collections` | [获取合集列表](../../get-collections.md) |

## 3. JSON 示例

```json
{
  "publishType": "video",
  "platforms": ["Bilibili"],
  "publishArgs": {
    "accountForms": [
      {
        "platformAccountId": "BILIBILI_ACC_ID",
        "video": {
          "key": "v_key",
          "size": 1024000,
          "width": 1920,
          "height": 1080,
          "duration": 60
        },
        "contentPublishForm": {
          "formType": "task",
          "title": "B站视频标题示例",
          "description": "这是关于B站视频发布的详细描述。",
          "tags": ["数码", "生活"],
          "category": [
            {
              "id": "1",
              "text": "生活",
              "raw": {}
            }
          ],
          "declaration": 0,
          "type": 1
        }
      }
    ]
  }
}
```
