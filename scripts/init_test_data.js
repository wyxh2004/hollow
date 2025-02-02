// 连接到 hollow 数据库
db = db.getSiblingDB("hollow");

// 清空现有数据
db.users.drop();
db.boxes.drop();
db.messages.drop();

// 创建测试用户
const testUsers = [
  {
    _id: ObjectId(),
    email: "test1@example.com",
    // 密码是 "password123"
    password: "$2a$10$1qcjIeVnv.L0Y2bfJMpFS.jaSB0jDigLyP1CBJ4Nd36KuTenWNski",
    created_at: new Date(Date.now() + 8 * 60 * 60 * 1000),
    updated_at: new Date(Date.now() + 8 * 60 * 60 * 1000),
  },
  {
    _id: ObjectId(),
    email: "test2@example.com",
    // 密码是 "password123"
    password: "$2a$10$1qcjIeVnv.L0Y2bfJMpFS.jaSB0jDigLyP1CBJ4Nd36KuTenWNski",
    created_at: new Date(),
    updated_at: new Date(),
  },
];

// 插入测试用户
db.users.insertMany(testUsers);

// 创建测试盒子
const testBoxes = [
  {
    _id: ObjectId(),
    name: "心情分享盒子",
    description: "分享你今天的心情和感受",
    owner_id: testUsers[0]._id,
    created_at: new Date(Date.now() + 8 * 60 * 60 * 1000),
    updated_at: new Date(Date.now() + 8 * 60 * 60 * 1000),
  },
  {
    _id: ObjectId(),
    name: "美食推荐盒子",
    description: "分享你最近吃到的美食",
    owner_id: testUsers[0]._id,
    created_at: new Date(),
    updated_at: new Date(),
  },
  {
    _id: ObjectId(),
    name: "学习交流盒子",
    description: "分享你的学习经验和心得",
    owner_id: testUsers[1]._id,
    created_at: new Date(),
    updated_at: new Date(),
  },
];

// 插入测试盒子
db.boxes.insertMany(testBoxes);

// 创建测试留言
const testMessages = [
  {
    _id: ObjectId(),
    box_id: testBoxes[0]._id,
    sender_id: testUsers[1]._id,
    content: "今天心情很好，完成了一个重要的项目！",
    is_anonymous: false,
    like_count: 2,
    liked_by: [testUsers[0]._id],
    created_at: new Date(),
  },
  {
    _id: ObjectId(),
    box_id: testBoxes[0]._id,
    content: "最近压力有点大，但是在努力坚持。",
    is_anonymous: true,
    like_count: 1,
    liked_by: [testUsers[1]._id],
    created_at: new Date(),
  },
  {
    _id: ObjectId(),
    box_id: testBoxes[1]._id,
    sender_id: testUsers[0]._id,
    content: "推荐一家新开的火锅店，味道非常不错！",
    is_anonymous: false,
    like_count: 0,
    liked_by: [],
    created_at: new Date(),
  },
];

// 插入测试留言
db.messages.insertMany(testMessages);

// 打印插入结果
print("测试数据已创建:");
print("用户数量:", db.users.count());
print("盒子数量:", db.boxes.count());
print("留言数量:", db.messages.count());
