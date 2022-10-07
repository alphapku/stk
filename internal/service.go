package internal

import (
	cfg "StakeBackendGoTest/configs"
	adt "StakeBackendGoTest/internal/adapter"
	def "StakeBackendGoTest/pkg/const"

	"github.com/gin-gonic/gin"
)

// func GetUserEquityPositions(adapterType pkg.AdapterType, _ string) []resp.StakePosition {
// 	_ = ReadMockPositionData() // mockPositions
// 	_ = ReadMockPriceData()    // mockPrices

// 	/* TODO: Transforms, conversions and calculations
// 	 *  - transform mockdata/mockpositions to mockdata/mockprices
// 	 *  - calculate all four profitOrLoss values
// 	 *  - return response/positions
// 	 */

// 	return nil
// }

func NewAdapter(cfg *cfg.Adapter) gin.HandlerFunc {
	switch cfg.AdapterType {
	case def.MockAdapter:
		mock := &adt.MockAdapter{}
		return mock.Do
	default:
		mock := &adt.MockAdapter{}
		return mock.Do
	}
}
