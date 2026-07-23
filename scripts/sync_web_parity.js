const fs = require('fs');
const path = require('path');

const root = path.resolve(__dirname, '..', 'schemas', 'platforms');
const optionalize = (schema, names) => { schema.required = schema.required.filter((name) => !names.includes(name)); };
const changes = {
  'acfun.article.schema.json': (s) => {
    if (s.properties.type) { s.properties.createType = s.properties.type; delete s.properties.type; }
    s.properties.createType.enum = [1, 2];
    s.properties.createType.default = 1;
    if (s.properties.description) s.properties.description.maxLength = 50;
    if (s.properties.tags) s.properties.tags.maxItems = 6;
  },
  'douyin.imageText.schema.json': (s) => {
    optionalize(s, ['shopping_cart']);
    s.properties.description.maxLength = 1000;
    s.properties.visibleType = { type: 'number', enum: [0, 1, 2], default: 0 };
    s.properties.group_shopping = { type: 'object' };
  },
  'chejiahao.article.schema.json': (s) => {
    if (s.properties.type) { s.properties.createType = s.properties.type; delete s.properties.type; }
    if (s.properties.createType) s.properties.createType.enum = [1, 2, 3];
  },
  'kuaichuanhao.article.schema.json': (s) => {
    if (s.properties.type) { s.properties.createType = s.properties.type; delete s.properties.type; }
    s.properties.declaration = { type: 'number', enum: [0, 1, 2, 3, 4, 5], default: 0 };
  },
  'wangyihao.article.schema.json': (s) => {
    if (s.properties.type) { s.properties.createType = s.properties.type; delete s.properties.type; }
    if (s.properties.declaration) s.properties.declaration.enum = [0, 1, 2, 3, 4];
  },
  'yidianhao.article.schema.json': (s) => {
    if (s.properties.declaration) s.properties.declaration.enum = [0, 3, 4, 5];
  },
  'kuaishou.imageText.schema.json': (s) => {
    optionalize(s, ['declaration']);
    s.properties.description.maxLength = 495;
    s.properties.visibleType.enum = [0, 1, 2];
    s.properties.declaration = { type: 'number', enum: [0, 1, 2, 3], default: 0 };
  },
  'shipinhao.imageText.schema.json': (s) => {
    optionalize(s, ['title', 'description', 'declaration']);
    s.properties.title.maxLength = 22;
    s.properties.declaration = { type: 'number', enum: [0, 1, 2, 3, 7, 8], default: 0 };
  },
  'toutiaohao.imageText.schema.json': (s) => {
    optionalize(s, ['declaration']);
    s.properties.description.maxLength = 2000;
    s.properties.declaration.enum = [0, 1, 2, 3, 6, 7, 8];
  },
  'xhs.imageText.schema.json': (s) => {
    optionalize(s, ['title', 'declaration', 'createType', 'shopping_cart']);
    s.properties.visibleType.enum = [0, 1, 2];
    s.properties.group = { type: 'object' };
    s.properties.shopping_cart = { type: 'array', minItems: 1, items: { type: 'object' } };
    s.properties.declaration = { type: 'number', enum: [0, 1, 2], default: 0 };
  },
  'xinlang.imageText.schema.json': (s) => {
    s.required = [...new Set([...s.required, 'visibleType', 'declaration'])];
    s.properties.description.maxLength = 5000;
    s.properties.visibleType = { type: 'number', enum: [0, 1], default: 0 };
    s.properties.declaration = { type: 'number', enum: [0, 1, 2, 3, 4], default: 0 };
  },
  'xiaohongshu.video.schema.json': (s) => {
    s.properties.visibleType.enum = [0, 1, 2];
    s.properties.declaration.enum = [0, 1, 2];
  },
  'souhushipin.video.schema.json': (s) => { s.properties.declaration.enum = [0, 1, 2, 3]; },
  'tengxunshipin.video.schema.json': (s) => { s.properties.declaration.enum = [0, 1, 2, 3, 4]; },
  'toutiaohao.video.schema.json': (s) => {
    s.properties.title.maxLength = 30;
    s.properties.declaration.enum = [0, 1, 2, 3, 6, 7, 8];
  },
  'wangyihao.video.schema.json': (s) => { s.properties.declaration.enum = [0, 1, 2, 3, 4]; },
  'yidianhao.video.schema.json': (s) => {
    if (s.properties.type) { s.properties.createType = s.properties.type; delete s.properties.type; }
    s.properties.createType.enum = [1, 2];
    s.properties.declaration.enum = [0, 3, 4, 5];
  },
  'xinlang.video.schema.json': (s) => {
    if (s.properties.type) { s.properties.createType = s.properties.type; delete s.properties.type; }
    s.properties.createType.enum = [1, 2, 3];
  },
};

for (const [name, update] of Object.entries(changes)) {
  const file = path.join(root, name);
  const schema = JSON.parse(fs.readFileSync(file, 'utf8'));
  update(schema);
  fs.writeFileSync(file, `${JSON.stringify(schema, null, 2)}\n`);
}
