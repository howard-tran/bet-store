package controller

import (
	"GoBackend/entity"
	"GoBackend/service"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func GetInfoWallet(ctx *gin.Context) {
	ClaimJwt, _ := ctx.Get("ClaimJwt")
	var data, _ = ClaimJwt.(entity.JwtClaimEntity)
	profile, err := Profileservice.GetProfile(data.ID)

	if err != nil {
		fmt.Printf("[GetInfoWallet] Profile data fail: %s\n", err.Error())
		ctx.JSON(http.StatusOK, service.CreateMsgErrorJsonResponse(http.StatusBadRequest, "Profile data fail"))
		return
	}

	resultWallet, err := Walletservice.GetInfoWallet(profile.ID)

	if err != nil {
		fmt.Printf("[GetInfoWallet] Error: %s\n", err.Error())
		ctx.JSON(http.StatusOK, service.CreateMsgErrorJsonResponse(http.StatusBadRequest, "[GetInfoWallet] Error: "+err.Error()))
		return
	}
	ctx.JSON(200, service.CreateMsgSuccessJsonResponse(resultWallet))
}

func WebhookWalletHandle(ctx *gin.Context) {
	defer ctx.JSON(200, gin.H{})
	var enti entity.HookEntity

	err := ctx.BindJSON(&enti)

	if err != nil {
		fmt.Println("[WebhookWalletHandle] Error:%s" + err.Error())
		return
	}

	err = Transferservice.AddTransfer(enti)

	if err != nil {
		fmt.Println("[WebhookWalletHandle] Error: " + err.Error())
		return
	}

	if enti.Signature != os.Getenv("HookSecret") {
		return
	}

	profile, err := Profileservice.GetProfilebyUsername(enti.Comment)

	if err != nil {
		fmt.Println("[WebhookWalletHandle] Error: " + err.Error())
		return
	}

	err = Walletservice.PayWallet(profile.ID, enti)

	if err != nil {
		fmt.Println("[WebhookWalletHandle] Error: " + err.Error())
		return
	}
	fmt.Printf("[WebhookWalletHandle] Spend success: %s|%d\n", profile.Username, enti.Amount)
}
