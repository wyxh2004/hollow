package controllers

import (
	"context"
	"hollow/models"
	"hollow/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BoxController struct{}

// CreateBox godoc
// @Summary      创建新的话题盒子
// @Description  创建一个新的话题盒子
// @Tags         boxes
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        box  body      models.Box  true  "盒子信息"
// @Success      201  {object}  CreateBoxResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/boxes [post]
func (bc *BoxController) CreateBox(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var box models.Box
	if err := c.ShouldBindJSON(&box); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	box.OwnerID = userID.(primitive.ObjectID)
	box.CreatedAt = time.Now().Add(time.Hour * 8)
	box.UpdatedAt = time.Now().Add(time.Hour * 8)

	result, err := utils.DB.Collection("boxes").InsertOne(context.Background(), box)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create box"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Box created successfully",
		"id":      result.InsertedID,
	})
}

// ListBoxes godoc
// @Summary      获取所有话题盒子
// @Description  获取所有已创建的话题盒子列表
// @Tags         boxes
// @Produce      json
// @Success      200  {array}   models.Box
// @Failure      500  {object}  ErrorResponse
// @Router       /api/boxes [get]
func (bc *BoxController) ListBoxes(c *gin.Context) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := utils.DB.Collection("boxes").Find(context.Background(), bson.M{}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch boxes"})
		return
	}
	defer cursor.Close(context.Background())

	var boxes []models.Box
	if err := cursor.All(context.Background(), &boxes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode boxes"})
		return
	}

	c.JSON(http.StatusOK, boxes)
}

// GetBox godoc
// @Summary      获取盒子详情及话题
// @Description  获取指定盒子的详细信息和所有话题
// @Tags         boxes
// @Produce      json
// @Param        id   path      string  true  "盒子ID"
// @Success      200  {object}  BoxDetailResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/boxes/{id} [get]
func (bc *BoxController) GetBox(c *gin.Context) {
	boxID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid box ID"})
		return
	}

	// 获取当前用户ID（如果已登录）
	userID, exists := c.Get("user_id")
	var currentUserID primitive.ObjectID
	if exists {
		currentUserID = userID.(primitive.ObjectID)
	}

	var box models.Box
	err = utils.DB.Collection("boxes").FindOne(context.Background(), bson.M{"_id": boxID}).Decode(&box)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Box not found"})
		return
	}

	// 获取话题
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := utils.DB.Collection("messages").Find(context.Background(), bson.M{"box_id": boxID}, opts)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}
	defer cursor.Close(context.Background())

	var messages []models.Message
	if err := cursor.All(context.Background(), &messages); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode messages"})
		return
	}

	// 转换为MessageResponse
	var messageResponses []models.MessageResponse
	for _, msg := range messages {
		response := models.MessageResponse{
			ID:          msg.ID,
			BoxID:       msg.BoxID,
			Content:     msg.Content,
			IsAnonymous: msg.IsAnonymous,
			LikeCount:   msg.LikeCount,
			CreatedAt:   msg.CreatedAt,
		}

		// 如果不是匿名消息，获取发送者邮箱
		if !msg.IsAnonymous && !msg.SenderID.IsZero() {
			var sender models.User
			err := utils.DB.Collection("users").FindOne(context.Background(), bson.M{"_id": msg.SenderID}).Decode(&sender)
			if err == nil {
				response.SenderEmail = sender.Email
			}
		}

		// 检查当前用户是否已点赞
		if exists {
			for _, likedBy := range msg.LikedBy {
				if likedBy == currentUserID {
					response.IsLiked = true
					break
				}
			}
		}

		messageResponses = append(messageResponses, response)
	}

	c.JSON(http.StatusOK, gin.H{
		"box":      box,
		"messages": messageResponses,
	})
}

// CreateMessage godoc
// @Summary      发送话题
// @Description  在指定盒子中创建新话题
// @Tags         messages
// @Accept       json
// @Produce      json
// @Param        id       path      string         true  "盒子ID"
// @Param        message  body      models.Message  true  "话题内容"
// @Success      201      {object}  CreateMessageResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /api/boxes/{id}/messages [post]
func (bc *BoxController) CreateMessage(c *gin.Context) {
	boxID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid box ID"})
		return
	}

	// 获取当前用户ID（如果已登录）
	userID, exists := c.Get("user_id")
	var senderID primitive.ObjectID
	if exists {
		senderID = userID.(primitive.ObjectID)
	}

	var message models.Message
	if err := c.ShouldBindJSON(&message); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	message.BoxID = boxID
	message.CreatedAt = time.Now().Add(time.Hour * 8)

	// 如果用户已登录且选择不匿名，则设置发送者ID
	if !message.IsAnonymous && exists {
		message.SenderID = senderID
	}

	result, err := utils.DB.Collection("messages").InsertOne(context.Background(), message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create message"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Message created successfully",
		"id":      result.InsertedID,
	})
}

// LikeMessage godoc
// @Summary      点赞/取消点赞话题
// @Description  为指定话题添加或取消点赞
// @Tags         messages
// @Produce      json
// @Security     Bearer
// @Param        id   path      string  true  "话题ID"
// @Success      200  {object}  LikeMessageResponse
// @Failure      400  {object}  ErrorResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /api/messages/{id}/like [post]
func (bc *BoxController) LikeMessage(c *gin.Context) {
	messageID, err := primitive.ObjectIDFromHex(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid message ID"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	currentUserID := userID.(primitive.ObjectID)

	// 检查消息是否存在
	var message models.Message
	err = utils.DB.Collection("messages").FindOne(context.Background(), bson.M{"_id": messageID}).Decode(&message)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Message not found"})
		return
	}

	// 检查用户是否已经点赞
	isLiked := false
	for _, likedBy := range message.LikedBy {
		if likedBy == currentUserID {
			isLiked = true
			break
		}
	}

	var update bson.M
	if isLiked {
		// 取消点赞
		update = bson.M{
			"$pull": bson.M{"liked_by": currentUserID},
			"$inc":  bson.M{"like_count": -1},
		}
	} else {
		// 添加点赞
		update = bson.M{
			"$addToSet": bson.M{"liked_by": currentUserID},
			"$inc":      bson.M{"like_count": 1},
		}
	}

	_, err = utils.DB.Collection("messages").UpdateOne(
		context.Background(),
		bson.M{"_id": messageID},
		update,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update like status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Like status updated successfully",
		"isLiked": !isLiked,
	})
}

// CreateBoxResponse represents the response for creating a box
type CreateBoxResponse struct {
	Message string `json:"message" example:"Box created successfully"`
	ID      string `json:"id" example:"507f1f77bcf86cd799439011"`
}

// BoxDetailResponse represents the response for getting box details
type BoxDetailResponse struct {
	Box      models.Box               `json:"box"`
	Messages []models.MessageResponse `json:"messages"`
}

// CreateMessageResponse represents the response for creating a message
type CreateMessageResponse struct {
	Message string `json:"message" example:"Message created successfully"`
	ID      string `json:"id" example:"507f1f77bcf86cd799439011"`
}

// LikeMessageResponse represents the response for liking/unliking a message
type LikeMessageResponse struct {
	Message string `json:"message" example:"Like status updated successfully"`
	IsLiked bool   `json:"isLiked" example:"true"`
}
