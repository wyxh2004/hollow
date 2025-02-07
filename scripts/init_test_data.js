// 连接到 hollow 数据库
db = db.getSiblingDB("hollow");

// 清空现有数据
db.users.drop();
db.boxes.drop();
db.messages.drop();
db.fs.files.drop();
db.fs.chunks.drop();

// 创建默认头像ID
const defaultAvatarId1 = ObjectId();
const defaultAvatarId2 = ObjectId();

// 创建测试用户
const testUsers = [
  {
    _id: ObjectId(),
    email: "test1@example.com",
    // 密码是 "password123"
    password: "$2a$10$1qcjIeVnv.L0Y2bfJMpFS.jaSB0jDigLyP1CBJ4Nd36KuTenWNski",
    avatar_id: defaultAvatarId1,
    created_at: new Date(Date.now() + 8 * 60 * 60 * 1000),
    updated_at: new Date(Date.now() + 8 * 60 * 60 * 1000),
  },
  {
    _id: ObjectId(),
    email: "test2@example.com",
    // 密码是 "password123"
    password: "$2a$10$1qcjIeVnv.L0Y2bfJMpFS.jaSB0jDigLyP1CBJ4Nd36KuTenWNski",
    avatar_id: defaultAvatarId2,
    created_at: new Date(Date.now() + 8 * 60 * 60 * 1000),
    updated_at: new Date(Date.now() + 8 * 60 * 60 * 1000),
  },
];

// 插入测试用户
db.users.insertMany(testUsers);

// 创建测试盒子
const testBoxes = [
  {
    _id: ObjectId(),
    name: "心情分享",
    description: "分享你今天的心情和感受",
    owner_id: testUsers[0]._id,
    created_at: new Date(Date.now() + 8 * 60 * 60 * 1000),
    updated_at: new Date(Date.now() + 8 * 60 * 60 * 1000),
  },
  {
    _id: ObjectId(),
    name: "美食推荐",
    description: "分享你最近吃到的美食",
    owner_id: testUsers[0]._id,
    created_at: new Date(Date.now() + 8 * 60 * 60 * 1000),
    updated_at: new Date(Date.now() + 8 * 60 * 60 * 1000),
  },
  {
    _id: ObjectId(),
    name: "学习交流",
    description: "分享你的学习经验和心得",
    owner_id: testUsers[1]._id,
    created_at: new Date(Date.now() + 8 * 60 * 60 * 1000),
    updated_at: new Date(Date.now() + 8 * 60 * 60 * 1000),
  },
];

// 插入测试盒子
db.boxes.insertMany(testBoxes);

// 创建测试话题
const testMessages = [
  {
    _id: ObjectId(),
    box_id: testBoxes[0]._id,
    sender_id: testUsers[1]._id,
    content: "今天心情很好，完成了一个重要的项目！",
    is_anonymous: false,
    like_count: 2,
    liked_by: [testUsers[0]._id],
    created_at: new Date(Date.now() + 8 * 60 * 60 * 1000),
  },
  {
    _id: ObjectId(),
    box_id: testBoxes[0]._id,
    content: "最近压力有点大，但是在努力坚持。",
    is_anonymous: true,
    like_count: 1,
    liked_by: [testUsers[1]._id],
    created_at: new Date(Date.now() + 8 * 60 * 60 * 1000),
  },
  {
    _id: ObjectId(),
    box_id: testBoxes[1]._id,
    sender_id: testUsers[0]._id,
    content: "推荐一家新开的火锅店，味道非常不错！",
    is_anonymous: false,
    like_count: 0,
    liked_by: [],
    created_at: new Date(Date.now() + 8 * 60 * 60 * 1000),
  },
];

// 插入测试话题
db.messages.insertMany(testMessages);

// 创建默认头像文件（这里使用一个简单的示例，实际应用中需要真实的图片数据）
const defaultAvatarData = BinData(
  0,
  "R0lGODlhAQABAIAAAAAAAP///yH5BAEAAAAALAAAAAABAAEAAAIBRAA7"
); // 1x1 透明图片

// 插入默认头像到 GridFS
db.fs.files.insertMany([
  {
    _id: defaultAvatarId1,
    filename: "default_avatar_1.png",
    contentType: "image/png",
    length: defaultAvatarData.length(),
    uploadDate: new Date(),
  },
  {
    _id: defaultAvatarId2,
    filename: "default_avatar_2.png",
    contentType: "image/png",
    length: defaultAvatarData.length(),
    uploadDate: new Date(),
  },
]);

// 插入头像数据块
db.fs.chunks.insertMany([
  {
    _id: ObjectId(),
    files_id: defaultAvatarId1,
    n: 0,
    data: defaultAvatarData,
  },
  {
    _id: ObjectId(),
    files_id: defaultAvatarId2,
    n: 0,
    data: defaultAvatarData,
  },
]);

// 打印插入结果
print("测试数据已创建:");
print("用户数量:", db.users.count());
print("盒子数量:", db.boxes.count());
print("话题数量:", db.messages.count());
print("头像文件数量:", db.fs.files.count());
