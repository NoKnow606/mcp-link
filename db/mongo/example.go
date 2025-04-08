package mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// ExampleUser 是一个示例用户模型
type ExampleUser struct {
	BaseModel        // 嵌入基础模型提供ID和时间戳字段
	Username  string `bson:"username" json:"username"`
	Email     string `bson:"email" json:"email"`
	Password  string `bson:"password" json:"password,omitempty"`
	Age       int    `bson:"age" json:"age"`
	Active    bool   `bson:"active" json:"active"`
}

// ExampleUserRepository 是用户模型的仓库
type ExampleUserRepository struct {
	*BaseRepository[*ExampleUser]
}

// NewExampleUserRepository 创建新的用户仓库
func NewExampleUserRepository(client *Client) *ExampleUserRepository {
	return &ExampleUserRepository{
		BaseRepository: NewRepository[*ExampleUser](client, "users"),
	}
}

// FindByUsername 根据用户名查找用户
func (r *ExampleUserRepository) FindByUsername(ctx context.Context, username string) (*ExampleUser, error) {
	filter := bson.M{"username": username}
	return r.FindOne(ctx, filter)
}

// FindActiveUsers 查找所有活跃用户
func (r *ExampleUserRepository) FindActiveUsers(ctx context.Context) ([]*ExampleUser, error) {
	filter := bson.M{"active": true}
	return r.Find(ctx, filter)
}

// ExampleUsage demonstrates how to use the MongoDB components
func ExampleUsage() {
	// 创建上下文
	ctx := context.Background()

	// 从环境变量加载配置
	config := LoadConfigFromEnv()

	// 也可以手动创建配置
	// config := &Config{
	//     URI:      "mongodb://localhost:27017",
	//     Database: "example_db",
	// }

	// 初始化客户端
	client := NewClient(config)
	if err := client.Connect(ctx); err != nil {
		fmt.Printf("Failed to connect to MongoDB: %v\n", err)
		return
	}
	defer client.Disconnect(ctx)

	// 测试连接
	if err := client.Ping(ctx); err != nil {
		fmt.Printf("Failed to ping MongoDB: %v\n", err)
		return
	}
	fmt.Println("Successfully connected to MongoDB")

	// 创建用户仓库
	userRepo := NewExampleUserRepository(client)

	// 创建示例用户
	user := &ExampleUser{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "securepassword",
		Age:      30,
		Active:   true,
	}

	// 保存用户
	if err := userRepo.Create(ctx, user); err != nil {
		fmt.Printf("Failed to create user: %v\n", err)
		return
	}
	fmt.Printf("Created user with ID: %s\n", user.ID.Hex())

	// 通过ID查找用户
	foundUser, err := userRepo.FindByID(ctx, user.ID.Hex())
	if err != nil {
		fmt.Printf("Failed to find user by ID: %v\n", err)
		return
	}
	fmt.Printf("Found user: %s (%s)\n", foundUser.Username, foundUser.Email)

	// 通过用户名查找用户
	foundByUsername, err := userRepo.FindByUsername(ctx, "testuser")
	if err != nil {
		fmt.Printf("Failed to find user by username: %v\n", err)
		return
	}
	fmt.Printf("Found user by username: %s\n", foundByUsername.Email)

	// 更新用户
	foundUser.Age = 31
	if err := userRepo.Update(ctx, foundUser); err != nil {
		fmt.Printf("Failed to update user: %v\n", err)
		return
	}
	fmt.Println("User updated successfully")

	// 查找所有活跃用户
	activeUsers, err := userRepo.FindActiveUsers(ctx)
	if err != nil {
		fmt.Printf("Failed to find active users: %v\n", err)
		return
	}
	fmt.Printf("Found %d active users\n", len(activeUsers))

	// 统计用户数量
	count, err := userRepo.Count(ctx, bson.M{})
	if err != nil {
		fmt.Printf("Failed to count users: %v\n", err)
		return
	}
	fmt.Printf("Total users: %d\n", count)

	// 删除用户
	if err := userRepo.Delete(ctx, user.ID.Hex()); err != nil {
		fmt.Printf("Failed to delete user: %v\n", err)
		return
	}
	fmt.Println("User deleted successfully")
}

// 以下是如何在main函数中使用MongoDB组件的示例
/*
func main() {
	ctx := context.Background()

	// 从环境变量加载配置
	config := mongo.LoadConfigFromEnv()

	// 初始化默认客户端
	if err := mongo.InitDefaultClient(ctx, config); err != nil {
		log.Fatalf("Failed to initialize MongoDB client: %v", err)
	}
	defer mongo.CloseDefaultClient(ctx)

	// 获取默认客户端
	client, err := mongo.GetDefaultClient()
	if err != nil {
		log.Fatalf("Failed to get MongoDB client: %v", err)
	}

	// 创建用户仓库
	userRepo := mongo.NewExampleUserRepository(client)

	// 使用仓库进行操作...
}
*/
