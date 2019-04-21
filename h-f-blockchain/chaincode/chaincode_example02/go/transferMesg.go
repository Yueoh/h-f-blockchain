

package main




//转账信息
type TransferMesg struct {
	Timestamp  int64	`json:"timestamp"`  	//时间戳
	TransferId string 	`json:"transferId"`		//转账凭证id
	MedicineId string	`json:"medicineId"`	//药物id
	HospitalId	string 	`json:"hospitalId"`	//医院id
	FactoryId	string 	`json:"factoryId"`	//工厂id
	TransferAccountValue	int	`json:"transferAccountValue"`	//转账金额
}


