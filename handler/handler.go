package handler

import (
	"bytes"
	"doanbackend/config"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type PatientRequest struct {
	Gender          string  `json:"gender" binding:"required,oneof=Male Female Other"`
	Age             float64 `json:"age" binding:"required,gt=0,lt=120"`
	Hypertension    int     `json:"hypertension" binding:"oneof=0 1"`
	HeartDisease    int     `json:"heart_disease" binding:"oneof=0 1"`
	EverMarried     string  `json:"ever_married" binding:"required,oneof=Yes No"`
	WorkType        string  `json:"work_type" binding:"required,oneof=Private Self-employed Govt_job children Never_worked"`
	ResidenceType   string  `json:"Residence_type" binding:"required,oneof=Urban Rural"`
	AvgGlucoseLevel float64 `json:"avg_glucose_level" binding:"required,gt=0"`
	Bmi             float64 `json:"bmi" binding:"required,gt=0"`
	SmokingStatus   string  `json:"smoking_status" binding:"required,oneof='never smoked' 'formerly smoked' smokes Unknown"`
}

type PredictResponse struct {
	Probability float64  `json:"probability"`
	HighRisk    bool     `json:"high_risk"`
	RiskLevel   string   `json:"risk_level"`
	Advice      string   `json:"advice"`
	Factors     []string `json:"factors"`
	Warning     string   `json:"warning"`
}

var mlClient = &http.Client{Timeout: 10 * time.Second}

func HandleData(c *gin.Context) {
	var req PatientRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": "Something went wrong, please try again!",
		})
		return
	}
	body, err := json.Marshal(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Lỗi xử lý dữ liệu."})
		return
	}
	resp, err := mlClient.Post(config.App.PyURL, "application/json", bytes.NewReader(body))
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadGateway, gin.H{"message": "Không kết nối được tới dịch vụ dự đoán."})
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		raw, _ := io.ReadAll(resp.Body)
		fmt.Println("Python trả lỗi:", resp.StatusCode, string(raw))
		c.JSON(http.StatusBadGateway, gin.H{"message": "Dịch vụ dự đoán trả về lỗi."})
		return
	}
	var result PredictResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Không đọc được kết quả dự đoán."})
		return
	}
	if result.Factors == nil {
		result.Factors = []string{}
	}
	c.JSON(http.StatusOK, result)
	c.Abort()
}
