package tool

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

type BaseModel struct {
	ID        uint64    `gorm:"column:ID;primary_key;unique" json:",omitempty"`           //条目TID
	CreatedAt time.Time `gorm:"column:CreatedAt;index:CreatedAt_Index" json:",omitempty"` //登陆日期
	UpdatedAt time.Time `gorm:"column:UpdatedAt" json:",omitempty"`                       //修改日期
	//DeletedAt *time.Time `gorm:"column:DeletedAt" json:",omitempty"` //删除日期
	//Delete    int        `gorm:"column:Delete"`                //0=无，1=删除，
}
type SessionCache struct {
	BaseModel
	//AccountID                   uint64             `gorm:"column:AccountID;index:AccountID;unique_index:AccountIDWorldKeyAreaID_Unique_Index"`
	WorldKey                    string             `gorm:"column:WorldKey;index:WorldKey_Index;unique_index:WorldKeyAreaIDAccount_Unique_Index;default:''"`
	AreaID                      uint64             `gorm:"column:AreaID;index:AreaID_Index;unique_index:WorldKeyAreaIDAccount_Unique_Index;default:'0'"`
	Account                     string             `gorm:"column:Account;index:Account_Index;unique_index:WorldKeyAreaIDAccount_Unique_Index"` //`gorm:"column:Account;unique;NOT NULL;unique_index:Account_Unique_Index"` //
	NickName                    string             `gorm:"column:NickName"`                                                //
	Avatar                      string             `gorm:"column:Avatar"`                                                  //
	SuperiorUserID              uint64             `gorm:"column:SuperiorUserID"`                                          //
	PlayerTitleDataTableID      uint64             `gorm:"column:PlayerTitleDataTableID"`                                  //
	LastLoginFactoryDataTableID uint64 `gorm:"column:LastLoginFactoryDataTableID"`                             //最后登陆的工厂
	OffLineTime                 time.Time          `gorm:"column:OffLineTime;index:OnOffTime_Index"`                       //离线时间
	OnLineTime                  time.Time          `gorm:"column:OnLineTime;index:OnOffTime_Index"`                        //上线时间
	Star                        uint32             `gorm:"column:Star;index:Star_Index"`                                   //
	StarTime                    time.Time          `gorm:"column:StarTime;index:StarTime_Index"`                           //添加星星的时间
	Robot                       uint64             `gorm:"column:Robot;default:'0'"`                                       //000000000  0=没有开启，1=已经开户
	Guide                       string             `gorm:"column:Guide;type:TEXT"`                              //引导列表
	SystemUnlock                string             `gorm:"column:SystemUnlock;type:TEXT"`                       //功能开启列表
	Task                        string             `gorm:"column:Task;type:TEXT"`                               //任务成就
	PKPoint                     int64              `gorm:"column:PKPoint;default:'0'"`                                     //PKPoint
	InviteAward                 uint64             `gorm:"column:InviteAward;default:'0'"`                                 //二进制：000000 获取奖励,位表示，分享次数
	LogicIndex                  int                `gorm:"column:LogicIndex;index:LogicIndex_Index;default:-1"`                                   //
	Authorize int `gorm:"column:Authorize;index:Authorize_Index;default:0"`

}

func BenchmarkCopyAndChange(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CopyAndChange(&SessionCache{AreaID:545+uint64(i),NickName:"241545"+strconv.Itoa(i)}, &SessionCache{})
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
		{name:"sdfsd",args:args{source:&SessionCache{AreaID:5455,NickName:"241545二妹dg"},target:&SessionCache{}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CopyAndChange(tt.args.source, tt.args.target); !reflect.DeepEqual(got, tt.want) {
				t.Logf("CopyAndChange() = %v, want %v", got, tt.want)
			}
		})
	}
}
