package tool

import (
	"reflect"
	"testing"
	"time"
)

type SessionCache struct {
	UserID                          uint64    `gorm:"column:UserID;primary_key;unique"`                         //条目TID
	CreatedAt                       time.Time `gorm:"column:CreatedAt;index:CreatedAt_Index" json:",omitempty"` //登陆日期
	UpdatedAt                       time.Time `gorm:"column:UpdatedAt" json:",omitempty"`                       //修改日期
	IdleCoinAward                   string    `gorm:"column:IdleCoinAward"`                                     //离线金币，用于认领
	WheelPrizeCount                 int       `gorm:"column:WheelPrizeCount"`                                   //转盘次数
	WheelNotStealNum                int64     `gorm:"column:WheelNotStealNum"`                                  //转盘，没有抽到盗贼的次数
	WheelPrizeDataTableID           uint64    `gorm:"column:WheelPrizeDataTableID"`                             //
	WheelCoinRand                   int       `gorm:"column:WheelCoinRand"`                                     //轮盘金币随机项
	WheelFunctionItemRand           int       `gorm:"column:WheelFunctionItemRand"`                             //轮盘道具随机项
	SystemUnlockOver                int       `gorm:"column:SystemUnlockOver"`                                  //所有功能是否开启 -1=false,0=?,1=true
	NextRestoreVITTime              time.Time `gorm:"column:NextRestoreVITTime"`                                //下一次恢复体力的时间
	NextRestoreStockTradeNumberTime time.Time `gorm:"column:NextRestoreStockTradeNumberTime"`                   //下一次恢复购买次数的时间
	RestVITTime                     time.Time `gorm:"column:RestVITTime"`                                       //重置体力时间
	RestBuyStockTrendNumberTime     time.Time `gorm:"column:RestBuyStockTrendNumberTime"`                       //
	RestBuyNewStockNumTime          time.Time `gorm:"column:RestBuyNewStockNumTime"`                            //
	BuyStockTrendNumber             int64     `gorm:"column:BuyStockTrendNumber"`                               //已经购买股票次数
	UseStockTrendNumber             int64     `gorm:"column:UseStockTrendNumber"`                               //已经使用股票次数
	BuyedNewStockStarCoinNum        string    `gorm:"column:BuyedNewStockStarCoinNum;default:'0'"`              //已经购买新股的星币数量
	StockWeek                       uint64    `gorm:"column:StockWeek;default:0"`                               //股票一年的第几周
	StockWeekGainAwards             uint64    `gorm:"column:StockWeekGainAwards;default:0"`                     //
	LastDaySCoinRank                uint64    `gorm:"column:LastDaySCoinRank;default:0"`                        //
	LastDaySCoinRankTime            time.Time `gorm:"column:LastDaySCoinRankTime"`                              //
	CanNewStockTradeTime            time.Time `gorm:"column:CanNewStockTradeTime"`                              //
	_oldSessionCache                *SessionCache
}

func BenchmarkCopyAndChange(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CopyAndChange(&SessionCache{UserID:545,IdleCoinAward:"241545"}, &SessionCache{})
	}
}

func TestCopyAndChange(t *testing.T) {
	type args struct {
		source interface{}
		target interface{}
	}

	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		{name:"sdfsd",args:args{source:&SessionCache{UserID:545,IdleCoinAward:"241545"},target:&SessionCache{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CopyAndChange(tt.args.source, tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CopyAndChange() = %v, want %v", got, tt.want)
			}
		})
	}
}
