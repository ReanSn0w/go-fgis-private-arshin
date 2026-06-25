package registries

import arshin "github.com/ReanSn0w/go-fgis-private-arshin/public/client"

type Spec struct {
	RegistryID arshin.RegistryID
	Title      string
	ItemType   string
}

const (
	NDRegistryID   arshin.RegistryID = "1"
	SCMRegistryID  arshin.RegistryID = "2"
	GSIRegistryID  arshin.RegistryID = "3"
	CMM2RegistryID arshin.RegistryID = "6"
	MDGRegistryID  arshin.RegistryID = "7"
	CMM3RegistryID arshin.RegistryID = "8"
	ICRegistryID   arshin.RegistryID = "9"
	SSDRegistryID  arshin.RegistryID = "10"
	SURegistryID   arshin.RegistryID = "11"
	GPSRegistryID  arshin.RegistryID = "12"
	MDRegistryID   arshin.RegistryID = "14"
	CMM1RegistryID arshin.RegistryID = "16"
	TSSIRegistryID arshin.RegistryID = "17"
	EPIRegistryID  arshin.RegistryID = "18"
	UTSORegistryID arshin.RegistryID = "19"
	P1WFRegistryID arshin.RegistryID = "47"
)

var (
	ND   = Spec{RegistryID: NDRegistryID, Title: "Нормативные правовые акты Российской Федерации", ItemType: "foei:ND_type"}
	SCM  = Spec{RegistryID: SCMRegistryID, Title: "Шифры калибровочных клейм", ItemType: "foei:SCM_type"}
	GSI  = Spec{RegistryID: GSIRegistryID, Title: "Стандарты государственной системы обеспечения единства измерений", ItemType: "foei:GSI_type"}
	CMM2 = Spec{RegistryID: CMM2RegistryID, Title: "Первичные референтные методики (методы) измерений", ItemType: "foei:CMM2_type"}
	MDG  = Spec{RegistryID: MDGRegistryID, Title: "Международные договоры", ItemType: "foei:MDG_type"}
	CMM3 = Spec{RegistryID: CMM3RegistryID, Title: "Референтные методики (методы) измерений", ItemType: "foei:CMM3_type"}
	IC   = Spec{RegistryID: ICRegistryID, Title: "Международные сличения", ItemType: "foei:IC_type"}
	SSD  = Spec{RegistryID: SSDRegistryID, Title: "Информация и данные ГСССД", ItemType: "foei:SSD_type"}
	SU   = Spec{RegistryID: SURegistryID, Title: "Эталоны единиц величин", ItemType: "foei:SU_type"}
	GPS  = Spec{RegistryID: GPSRegistryID, Title: "Государственные первичные эталоны Российской Федерации", ItemType: "foei:GPS_type"}
	MD   = Spec{RegistryID: MDRegistryID, Title: "Международные документы", ItemType: "foei:MD_type"}
	CMM1 = Spec{RegistryID: CMM1RegistryID, Title: "Аттестованные методики (методы) измерений", ItemType: "foei:CMM1_type"}
	TSSI = Spec{RegistryID: TSSIRegistryID, Title: "Сведения об отнесении технических средств к средствам измерений", ItemType: "foei:TSSI_type"}
	EPI  = Spec{RegistryID: EPIRegistryID, Title: "Перечень измерений, относящихся к сфере государственного регулирования", ItemType: "foei:EPI_type"}
	UTSO = Spec{RegistryID: UTSORegistryID, Title: "Утвержденные типы стандартных образцов", ItemType: "foei:UTSO_type"}
	P1WF = Spec{RegistryID: P1WFRegistryID, Title: "Уведомления об осуществлении деятельности по производству эталонов единиц величин, стандартных образцов и средств измерений", ItemType: "gost:p1wfRequestType4"}
)

var known = []Spec{
	ND,
	SCM,
	GSI,
	CMM2,
	MDG,
	CMM3,
	IC,
	SSD,
	SU,
	GPS,
	MD,
	CMM1,
	TSSI,
	EPI,
	UTSO,
	P1WF,
}

func Known() []Spec {
	specs := make([]Spec, len(known))
	copy(specs, known)
	return specs
}

func SpecForRegistry(registryID arshin.RegistryID) (Spec, bool) {
	for _, spec := range known {
		if spec.RegistryID == registryID {
			return spec, true
		}
	}
	return Spec{}, false
}

func SpecForItemType(itemType string) (Spec, bool) {
	for _, spec := range known {
		if spec.ItemType == itemType {
			return spec, true
		}
	}
	return Spec{}, false
}
