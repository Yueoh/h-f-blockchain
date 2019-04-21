

package main



import (
	"fmt"
	"strconv"
	"bytes"
	"time"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
	 "encoding/json"
	
)


type SimpleChaincode struct {
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("BlockChain:ex02 Init")
	_, args := stub.GetFunctionAndParameters()
	t.addTransferMesg(stub,args)
	
	return shim.Success([]byte("Init Success!"))
}

func (t *SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	fmt.Println("BlockChain:ex02 Invoke")
	function, args := stub.GetFunctionAndParameters()
	if function == "addPatientMesg" {
		// 医院增加PatientMesg数据
		return t.addPatientMesg(stub, args)
	} else if function=="queryPatientMesg" {
		//工厂查询patientMesg数据
		return t.queryPatientMesg(stub,args)
	}else if function=="addMedicineMesg" {
		//工厂增加MedicineMesg数据
		return t.addMedicineMesg(stub,args)
	}else if function=="queryMedicineMesg" {
		//医院查看MedicineMesg数据
		return t.queryMedicineMesg(stub,args)
	}else if function=="addTransferMesg" {
		//工厂上传TransferMesg数据
		return t.addTransferMesg(stub,args)
	}else if function=="queryTransferMesg" {
		//查询TransferMesg数据
		return t.queryTransferMesg(stub,args)
	}

	
	return shim.Error("Invalid invoke function name.")
}


//	增加patientMesg数据
//	args:医院id、疾病id、基本信息、病症状况、医生诊断、疗程、药物id、剂量、每日用药次数、给药途径、开始服药时间、结束服药时间、合同id、是否紧急
func (t *SimpleChaincode) addPatientMesg(stub shim.ChaincodeStubInterface,args []string) pb.Response {

	if len(args)!=14{
		return shim.Error("Parameter count must 14")
	}
	if len(args[10])!=14{
		return shim .Error("Parameter TimeOfBegin length must be 14")
	}
	if len(args[11])!=14{
		return shim .Error("Parameter TimeOfEnd length must be 14")
	}
	var HospitalId ,DiseaseId string 
	HospitalId=args[0]
	DiseaseId=args[1]
	var patientMesg PatientMesg 
	patientMesg.HospitalId=HospitalId
	patientMesg.DiseaseId=DiseaseId
	patientMesg.BasicMesg=args[2]
	patientMesg.DiseaseCondition=args[3]
	patientMesg.Diagnose=args[4]
	patientMesg.Treatment=args[5]
	patientMesg.MedicineId=args[6]
	patientMesg.Dosage=args[7]
	patientMesg.TimesOfMedicineUsePerDay=args[8]
	patientMesg.MedicineRoute=args[9]
	patientMesg.TimeOfBegin=args[10]
	patientMesg.TimeOfEnd=args[11]
	patientMesg.DiseaseContractId=args[12]
	patientMesg.IsUrgent=args[13]
	patientMesg.Timestamp=time.Now().Unix()

	//json序列化
	patientMesgJsonBytes,err:=json.Marshal(patientMesg)
	if err!=nil{
		return shim .Error("Json serialize patientMesg fail,HospitalId= "+HospitalId)
	}
	//生成联合主键
	Patientkey,err :=stub.CreateCompositeKey("PatientMesg",[]string{HospitalId,DiseaseId})
	if err!=nil{
		return shim .Error("Fail to CreateCompositeKey")
	}
	//保存数据
	err=stub.PutState(Patientkey,patientMesgJsonBytes)
	if err!=nil{
		return shim .Error("Fail to Putstate,HospitalId= "+HospitalId)
	}
	return shim.Success(patientMesgJsonBytes)
	//return shim.Success(hospitalMesgJsonBytes)
	// return shim.Success([]byte("Success to addHospitalMesg!"))
}


//制药厂查询医院上传数据:DiseaseId,HospitalId
func (t *SimpleChaincode) queryPatientMesg(stub shim.ChaincodeStubInterface,args []string) pb.Response {

	if len(args)!=2{
		return shim .Error("Parameter count must 2")
	}

	patientMesgResult,err:=stub.GetStateByPartialCompositeKey("PatientMesg",[]string{args[0],args[1]})	//联合查询
	if err!=nil{
		return shim .Error("Fail to GetStateByPartialCompositeKey")
	}
	
	defer patientMesgResult.Close()
	
	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf(`{"PatientMesg_Key": {"HospitalId": "%s", "DiseaseId": "%s"},`, args[0], args[1]))
	buff.WriteString(`"PatientResult":`)
	bArrayMemberAlreadyWritten := false
	for patientMesgResult.HasNext() {
		queryResponse, err := patientMesgResult.Next()
		if err != nil {
			return shim .Error("Fail queryResponse")
		}

		if bArrayMemberAlreadyWritten == true {
		   buff.WriteString(",")
		}
		buff.WriteString("{\"Key\":")
		buff.WriteString("\"")
		buff.WriteString(queryResponse.Key)
		buff.WriteString("\"")
  
		buff.WriteString(", \"Record\":")
		buff.WriteString(string(queryResponse.Value))
		buff.WriteString("}")
		bArrayMemberAlreadyWritten = true
	 }
	
	buff.WriteString("}")

	return shim.Success(buff.Bytes())
}


//	增加MedicineMesg数据
//	args:工厂id、药物id、药物名称、药理、剂量、不良反映
func (t *SimpleChaincode) addMedicineMesg(stub shim.ChaincodeStubInterface,args []string) pb.Response {

	if len(args)!=6{
		return shim.Error("Parameter count must 6")
	}
	
	var MedicineId ,FactoryId string 
	FactoryId =args[0]
	MedicineId=args[1]
	var medicineMesg MedicineMesg 
	medicineMesg.FactoryId=FactoryId
	medicineMesg.MedicineId=MedicineId
	medicineMesg.MedicineName=args[2]
	medicineMesg.Pharmacology=args[3]
	medicineMesg.Dosage=args[4]
	medicineMesg.Adv_reaction=args[5]
	medicineMesg.Timestamp=time.Now().Unix()

	//json序列化
	medicineMesgJsonBytes,err:=json.Marshal(medicineMesg)
	if err!=nil{
		return shim .Error("Json serialize medicineMesg fail,factoryId= "+FactoryId)
	}
	//生成联合主键
	Medicinekey,err :=stub.CreateCompositeKey("MedicineMesg",[]string{FactoryId,MedicineId})
	if err!=nil{
		return shim .Error("Fail to CreateCompositeKey")
	}
	//保存数据
	err=stub.PutState(Medicinekey,medicineMesgJsonBytes)
	if err!=nil{
		return shim .Error("Fail to Putstate,FactoryId= "+FactoryId)
	}
	return shim.Success(medicineMesgJsonBytes)
}

//医院查询制药厂上传数据:FactoryId、Medicineid
func (t *SimpleChaincode) queryMedicineMesg(stub shim.ChaincodeStubInterface,args []string) pb.Response {

	if len(args)!=2{
		return shim .Error("Parameter count must 2")
	}

	medicineMesgResult,err:=stub.GetStateByPartialCompositeKey("MedicineMesg",[]string{args[0],args[1]})	//联合查询
	if err!=nil{
		return shim .Error("Fail to GetStateByPartialCompositeKey")
	}
	
	defer medicineMesgResult.Close()
	
	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf(`{"MedicineMesg_Key": {"FactoryId": "%s", "MedicineId": "%s"},`, args[0], args[1]))
	buff.WriteString(`"MedicineResult":`)
	bArrayMemberAlreadyWritten := false
	for medicineMesgResult.HasNext() {
		queryResponse, err := medicineMesgResult.Next()
		if err != nil {
			return shim .Error("Fail queryResponse")
		}

		if bArrayMemberAlreadyWritten == true {
		   buff.WriteString(",")
		}
		buff.WriteString("{\"Key\":")
		buff.WriteString("\"")
		buff.WriteString(queryResponse.Key)
		buff.WriteString("\"")
  
		buff.WriteString(", \"Record\":")
		buff.WriteString(string(queryResponse.Value))
		buff.WriteString("}")
		bArrayMemberAlreadyWritten = true
	 }
	
	buff.WriteString("}")

	return shim.Success(buff.Bytes())
}

//	工厂增TransferMesg数据
//	args:转账凭证id、药物id、医院id、工厂id、转账金额
func (t *SimpleChaincode) addTransferMesg(stub shim.ChaincodeStubInterface,args []string) pb.Response {

	if len(args)!=5{
		return shim.Error("Parameter count must 5")
	}
	
	var transferMesg TransferMesg
	var MedicineId ,FactoryId	string
	var err error
	MedicineId=args[1]
	FactoryId=args[3]
	transferMesg.TransferId=args[0]
	transferMesg.MedicineId=MedicineId
	transferMesg.HospitalId=args[2]
	transferMesg.FactoryId=FactoryId
	transferMesg.TransferAccountValue,err =strconv.Atoi(args[4])
	transferMesg.Timestamp=time.Now().Unix()

	//json序列化
	transferMesgJsonBytes,err:=json.Marshal(transferMesg)
	if err!=nil{
		return shim .Error("Json serialize medicineMesg fail,factoryId= "+FactoryId)
	}
	//生成联合主键
	Transferkey,err :=stub.CreateCompositeKey("TransferMesg",[]string{FactoryId,MedicineId})
	if err!=nil{
		return shim .Error("Fail to CreateCompositeKey")
	}
	//保存数据
	err=stub.PutState(Transferkey,transferMesgJsonBytes)
	if err!=nil{
		return shim .Error("Fail to Putstate,FactoryId= "+FactoryId)
	}
	return shim.Success(transferMesgJsonBytes)
}

//查询制药厂上传转账数据:FactoryId、Medicineid
func (t *SimpleChaincode) queryTransferMesg(stub shim.ChaincodeStubInterface,args []string) pb.Response {

	if len(args)!=2{
		return shim .Error("Parameter count must 2")
	}

	transferMesgResult,err:=stub.GetStateByPartialCompositeKey("TransferMesg",[]string{args[0],args[1]})	//联合查询
	if err!=nil{
		return shim .Error("Fail to GetStateByPartialCompositeKey")
	}
	
	defer transferMesgResult.Close()
	
	var buff bytes.Buffer
	buff.WriteString(fmt.Sprintf(`{"TransferMesg_Key": {"FactoryId": "%s", "MedicineId": "%s"},`, args[0], args[1]))
	buff.WriteString(`"TransferResult":`)
	bArrayMemberAlreadyWritten := false
	for transferMesgResult.HasNext() {
		queryResponse, err := transferMesgResult.Next()
		if err != nil {
			return shim .Error("Fail queryResponse")
		}

		if bArrayMemberAlreadyWritten == true {
		   buff.WriteString(",")
		}
		buff.WriteString("{\"Key\":")
		buff.WriteString("\"")
		buff.WriteString(queryResponse.Key)
		buff.WriteString("\"")
  
		buff.WriteString(", \"Record\":")
		buff.WriteString(string(queryResponse.Value))
		buff.WriteString("}")
		bArrayMemberAlreadyWritten = true
	 }
	
	buff.WriteString("}")

	return shim.Success(buff.Bytes())
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
