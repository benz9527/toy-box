package b

// 贪心策略

// 手上有一副扑克牌，每张牌按牌面数字记分（J=11,Q=12,K=13，没有大小王)，出牌时按照以下规则记分：
// 出单张，记牌面分数，例如出一张2，得分为2
// 出对或3张，记牌面分数总和再x2，例如出3张3，得分为(3+3+3)x2=18
// 出5张顺，记牌面分数总和再x2，例如出34567顺，得分为(3+4+5+6+7)x2=50
// 出4张炸弹，记牌面分数总和再x3，例如出4张4，得分为4x4x3=48
// 求出一副牌最高的得分数

func GetMaxScoreOfPokers(pokerStr string) int {
	pokers := make([]int, 14)
	for _, poker := range pokerStr {
		switch poker {
		case '0':
			pokers[10] += 1
		case 'J':
			pokers[11] += 1
		case 'Q':
			pokers[12] += 1
		case 'K':
			pokers[13] += 1
		default:
			pokers[poker-'0'] += 1
		}
	}
	_getPokersScore := func(poker, count int) int {
		score := 0
		switch count {
		case 1:
			score = poker
		case 2, 3:
			score = poker * count * 2
		case 4:
			score = poker * count * 3
		default:

		}
		return score
	}
	_getPokerProfit := func(poker, count int) int {
		profit := 0
		switch count {
		case 1:
			//  poker * 2 - poker
			profit = poker
		case 2:
			// a = poker * 2 * 2
			// b = poker + poker * 2
			profit = -poker
		case 3:
			// a = poker * 3 * 2
			// b = poker * 2 + poker * 2 * 2
			profit = 0
		case 4:
			// a = poker * 4  * 3
			// b = poker * 2 + poker * 3 * 2
			profit = -poker * 4
		default:
			profit = -100000
		}
		return profit
	}
	_getStraightPokersProfit := func(startPoker int) int {
		profit := 0
		for i := startPoker; i <= startPoker+4; i++ {
			profit += _getPokerProfit(i, pokers[i])
		}
		return profit
	}
	// 12345, 90JQK
	maxProfit := 0
	maxProfitStartPoker := 0
	ans := 0
	for {
		// 12345,67890
		for i := 1; i <= 9; i++ {
			p := _getStraightPokersProfit(i)
			if p > maxProfit {
				maxProfit = p
				maxProfitStartPoker = i
			}
		}
		if maxProfitStartPoker == 0 {
			break
		}
		for i := maxProfitStartPoker; i <= maxProfitStartPoker+4; i++ {
			ans += i * 2
			pokers[i]--
		}
	}
	for i := 1; i <= 13; i++ {
		ans += _getPokersScore(i, pokers[i])
	}
	return ans
}
