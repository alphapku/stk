package internal

func GetUserEquityPositions(_ string) []StakePosition {
	_ = ReadMockPositionData() // mockPositions
	_ = ReadMockPriceData()    // mockPrices

	/* TODO: Transforms, conversions and calculations
	 *  - transform mockdata/mockpositions to mockdata/mockprices
	 *  - calculate all four profitOrLoss values
	 *  - return response/positions
	 */

	return nil
}
