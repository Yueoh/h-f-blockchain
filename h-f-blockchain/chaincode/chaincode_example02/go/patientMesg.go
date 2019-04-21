

package main




//医院上传信息
type PatientMesg struct {
	Timestamp  int64	`json:"timestamp"`  	//时间戳
	HospitalId  string 	//医院id
	DiseaseId  string	//疾病id
	BasicMesg string	`json:"basicMesg"`	//基本信息
	DiseaseCondition  string	`json:"diseaseCondition"`	//病症状况
	Diagnose string		`json:"diagnose"`	//医生诊断
	Treatment string	`json:"treatment"`	//疗程
	MedicineId string	`json:"medicineId"`	//药物id
	Dosage string	`json:"dosage"`	//剂量
	TimesOfMedicineUsePerDay string	`json:"timesOfMedicineUsePerDay"`	//每日用药次数
	MedicineRoute string `json:"medicineRoute"`	//给药途径
	TimeOfBegin string `json:"timeOfBegin"`	//开始服药时间
	TimeOfEnd string 	`json:"timeOfEnd"`	//结束服药时间
	DiseaseContractId	string	`json:"diseaseContractId"`	//病人签署合同id
	IsUrgent string 		`json:"isUrgent"`	//是否紧急
}


