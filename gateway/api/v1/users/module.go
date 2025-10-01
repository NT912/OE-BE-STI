package users

import (
	"log"

	"gateway/configs"
	"gateway/models"
	"gateway/rpc"
	"gateway/utils"

	"github.com/gin-gonic/gin"
)

// Export instance
var (
	UserRepo *UserRepository
	UserSvc  *UserService
)

func InitModule(r *gin.Engine, rpcClient *rpc.RPCClient) {
	db := configs.DB
	if configs.Env.AppEnv != "production" {
		db.AutoMigrate(&models.User{})
	}

	// Seed sysadmin nếu chưa có
	var count int64
	db.Model(&models.User{}).Where("role = ?", "sysadmin").Count(&count)
	if count == 0 {
		hashed, _ := utils.HashPassword("SysAdmin@123")
		db.Create(&models.User{
			Name:     "System Admin",
			Email:    "sysadmin@system.local",
			Password: hashed,
			Role:     "sysadmin",
		})
	}

	UserRepo = NewUserRepository(db)
	UserSvc = NewUserService(UserRepo, rpcClient)
	controller := NewUserController(UserSvc)

	api := r.Group("/api/v1")
	controller.RegisterRoutes(api)
}

func InitRPC() *rpc.RPCClient {
	rpcClient, err := rpc.NewRPCClient()
	if err != nil {
		log.Fatalf("Failed to init RPC client: %v", err)
	}

	log.Println("✅ RPC Client initialized")
	return rpcClient
}
