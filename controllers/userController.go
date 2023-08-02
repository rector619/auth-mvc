package controllers

import (
	"context"
	"fmt"
	"log"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/aremxyplug-be/database"
	"github.com/aremxyplug-be/types/dto"

	helper "github.com/aremxyplug-be/helpers"
	"github.com/aremxyplug-be/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

//HashPassword is used to encrypt the password before it is stored in the DB
func HashPassword(password string) string {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
    if err != nil {
        log.Panic(err)
    }

    return string(bytes)
}

//VerifyPassword checks the input password while verifying it with the password in the DB.
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
    err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
    check := true
    msg := ""

    if err != nil {
        msg = fmt.Sprintf("login or passowrd is incorrect")
        check = false
    }

    return check, msg
}

//CreateUser is the api used to tget a single user
func SignUp() gin.HandlerFunc {
    return func(c *gin.Context) {
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        defer cancel()
        var user models.User
        if err := c.BindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        validationErr := validate.Struct(user)
        if validationErr != nil {
            fmt.Print(validationErr)
            c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
            return
        }
        count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
        fmt.Println("count results", count)
       
        if err != nil {
            log.Print(err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
            return
        }
        if count > 0 {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "this email already exists"})
            return
        }

    
        count, err = userCollection.CountDocuments(ctx, bson.M{"phonenumber": user.Phonenumber})
        defer cancel()
        if err != nil {
            log.Panic(err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone number"})
            return
        }

        if count > 0 {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "this phone number already exists"})
            return
        }

        count, err = userCollection.CountDocuments(ctx, bson.M{"username": user.Username})
        defer cancel()
        if err != nil {
            log.Panic(err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the username"})
            return
        }

        if count > 0 {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "this username already exists"})
            return
        }

        password := HashPassword(*user.Password)
        user.Password = &password

        user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
        user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
        user.ID = primitive.NewObjectID()
        user.User_id = user.ID.Hex()
        token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.Fullname, *user.Username, user.User_id)
        user.Token = &token
        user.Refresh_token = &refreshToken

        _, insertErr := userCollection.InsertOne(ctx, user,)
        if insertErr != nil {
            msg := fmt.Sprintf("User item was not created")
            c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
            return
        }
        defer cancel()

        signUpResponse := dto.SignUpResponse{
            Id:            user.ID,
            Fullname: *user.Fullname,
            Username: *user.Username,
            Email: *user.Email,
            Token: *user.Token,
            Refresh_token: *user.Refresh_token,
        }

        c.JSON(http.StatusOK, gin.H{"message": "user created successfully", "status": http.StatusOK, "data": signUpResponse})

    }
}

//Login is the api used to tget a single user
func Login() gin.HandlerFunc {
    return func(c *gin.Context) {
        var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
        var user models.User
        var foundUser models.User

        if err := c.BindJSON(&user); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        filter := bson.M{
            "$or": []bson.M{
                bson.M{"email": user.Email},
                bson.M{"username": user.Username},
            },
        }

        err := userCollection.FindOne(ctx, filter).Decode(&foundUser)
        defer cancel()
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "email or username or passowrd is incorrect"})
            return
        }

        passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
        defer cancel()
        if passwordIsValid != true {
            c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
            return
        }

        token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.Fullname, *foundUser.Username, foundUser.User_id)

        helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)

        loginResponse := dto.LoginResponse{
            Id:            foundUser.ID,
            Email: *foundUser.Email,
            Token: *foundUser.Token,
            Refresh_token: *foundUser.Refresh_token,
        }

        c.JSON(http.StatusOK, gin.H{"message": "login successful", "status": http.StatusOK, "data": loginResponse})

    }
}