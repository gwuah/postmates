package utils

import (
	"crypto/rand"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/electra-systems/core-api/database/models"
	"github.com/electra-systems/core-api/shared"
)

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func GeneratePhoneNumber(phone string) string {
	return fmt.Sprintf("233%s", phone[1:])
}

func GenerateOTP() string {
	max := 4
	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}

func StringifyLngLat(props shared.Coord) string {
	return fmt.Sprintf("%f,%f", props.Longitude, props.Latitude)
}

func ConvertToUint64(num string) uint64 {
	id64, _ := strconv.ParseUint(num, 10, 64)
	return id64
}

func ConvertToVehicleType(id string) models.VehicleType {
	switch strings.ToLower(id) {
	case "motor":
		return models.Motor
	case "car":
		return models.Car
	default:
		return models.Motor
	}
}
