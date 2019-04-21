

package main




//药物信息
type MedicineMesg struct {
	Timestamp  int64	`json:"timestamp"`  	//时间戳
	FactoryId	string  `json:"factoryId"` 		//工厂id
	MedicineId	string `json:"medicineId"` 	//药物id
	MedicineName	string `json:"medicineName"` 	//药物名称
	Pharmacology  string `json:"pharmacology"`        //药理
	Dosage  string	`json:"dosage"`       //剂量
	Adv_reaction  string	`json:"adv_reaction"`    //不良反应
}
